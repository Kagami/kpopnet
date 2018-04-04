package kpopnet

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/Kagami/go-dlib"
)

type namesKey [2]string
type idolNamesMap map[namesKey]Idol

var (
	faceRec *dlib.FaceRec
)

func getImagesDir(d string) string {
	return filepath.Join(d, "images")
}

func getModelsDir(d string) string {
	return filepath.Join(d, "models")
}

func getImageSha1(ipath string) (hashHex string, err error) {
	f, err := os.Open(ipath)
	if err != nil {
		return
	}
	defer f.Close()

	h := sha1.New()
	if _, err = io.Copy(h, f); err != nil {
		return
	}
	hash := h.Sum(nil)
	hashHex = hex.EncodeToString(hash[:])
	return
}

func getNamesKey(bname, iname string) namesKey {
	return [2]string{bname, iname}
}

func getIdolNamesKey(idol Idol, bandById map[string]Band) (key namesKey, ok bool) {
	iname, ok := idol["name"].(string)
	if !ok {
		return
	}
	bandId, ok := idol["band_id"].(string)
	if !ok {
		return
	}
	band, ok := bandById[bandId]
	if !ok {
		return
	}
	bname, ok := band["name"].(string)
	if !ok {
		return
	}
	return getNamesKey(bname, iname), true
}

func getIdolNamesMap() (idolByNames idolNamesMap, err error) {
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)
	if err = setReadOnly(tx); err != nil {
		return
	}
	_, bandById, err := getBands(tx)
	if err != nil {
		return
	}
	idols, _, err := getIdols(tx)
	if err != nil {
		return
	}
	idolByNames = make(idolNamesMap)
	for _, idol := range idols {
		if key, ok := getIdolNamesKey(idol, bandById); ok {
			idolByNames[key] = idol
		}
	}
	return
}

func recognizeIdolImage(ipath string) (face *dlib.Face, id string, err error) {
	face, err = faceRec.RecognizeSingle(ipath)
	if err != nil || face == nil {
		return
	}
	id, err = getImageSha1(ipath)
	return
}

// TODO(Kagami): Use multiple threads?
func recognizeIdolImages(idir string) (faces []dlib.Face, ids []string, err error) {
	idolImages, err := ioutil.ReadDir(idir)
	if err != nil {
		return
	}
	// No need to validate names/formats because everything was checked by
	// Python spider.
	for _, file := range idolImages {
		var face *dlib.Face
		var imageId string
		fname := file.Name()
		ipath := filepath.Join(idir, fname)
		face, imageId, err = recognizeIdolImage(ipath)
		if err != nil {
			return
		}
		if face == nil {
			log.Printf("Not a single face on %s", ipath)
			continue
		}
		faces = append(faces, *face)
		ids = append(ids, imageId)
	}
	return
}

func importIdolImages(st *sql.Stmt, idir string, idol Idol) (err error) {
	faces, imageIds, err := recognizeIdolImages(idir)
	if err != nil {
		return
	}
	idolId := idol["id"].(string)
	for i, face := range faces {
		r := face.Rectangle
		rectStr := fmt.Sprintf("((%d,%d),(%d,%d))", r[0], r[1], r[2], r[3])
		descrSize := unsafe.Sizeof(face.Descriptor)
		descrPtr := unsafe.Pointer(&face.Descriptor)
		descrBytes := (*[1 << 30]byte)(descrPtr)[:descrSize:descrSize]
		imageId := imageIds[i]
		if _, err = st.Exec(rectStr, descrBytes, imageId, idolId); err != nil {
			return
		}
	}
	return
}

func importBandImages(bdir, bname string, idolByNames idolNamesMap) (err error) {
	// Use single transaction per band.
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)
	st := tx.Stmt(prepared["upsert_face"])

	idolDirs, err := ioutil.ReadDir(bdir)
	if err != nil {
		return
	}
	for _, dir := range idolDirs {
		iname := dir.Name()
		idir := filepath.Join(bdir, iname)
		key := getNamesKey(bname, iname)
		idol, ok := idolByNames[key]
		if !ok {
			log.Printf("Can't find %s (%s)", iname, bname)
			continue
		}
		log.Printf("Importing %s", idir)
		if err = importIdolImages(st, idir, idol); err != nil {
			return
		}
	}
	return
}

// Read and update idol faces in database.
func ImportImages(connStr string, dataDir string, onlyBands []string) (err error) {
	if err = StartDb(nil, connStr); err != nil {
		return
	}

	idolByNames, err := getIdolNamesMap()
	if err != nil {
		err = fmt.Errorf("Error querying idols: %v", err)
		return
	}

	faceRec, err = dlib.NewFaceRec(getModelsDir(dataDir))
	if err != nil {
		err = fmt.Errorf("Error initializing face recognizer: %v", err)
		return
	}

	bandDirs, err := ioutil.ReadDir(getImagesDir(dataDir))
	if err != nil {
		err = fmt.Errorf("Error reading bands: %v", err)
		return
	}

	bandFilter := make(map[string]bool)
	for _, bname := range onlyBands {
		bandFilter[bname] = true
	}

	for _, dir := range bandDirs {
		bname := dir.Name()
		if len(bandFilter) > 0 && !bandFilter[bname] {
			continue
		}
		bdir := filepath.Join(getImagesDir(dataDir), bname)
		if err = importBandImages(bdir, bname, idolByNames); err != nil {
			err = fmt.Errorf("Error importing %s images: %v", bname, err)
			return
		}
	}
	return
}
