package kpopnet

import (
	"crypto/sha1"
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

type Face struct {
	descriptor dlib.FaceDescriptor
	imageId    string
}

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
	return [2]string{bname, iname}, true
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

func recognizeIdolImage(ipath string) (face *Face, err error) {
	d, err := faceRec.GetDescriptor(ipath)
	if err != nil || d == nil {
		return
	}
	hash, err := getImageSha1(ipath)
	if err != nil {
		return
	}
	face = &Face{*d, hash}
	return
}

// TODO(Kagami): Use multiple threads?
func recognizeIdolImages(idir string) (faces []Face, err error) {
	idolImages, err := ioutil.ReadDir(idir)
	if err != nil {
		return
	}
	// No need to validate names/formats because everything was checked by
	// Python spider.
	for _, file := range idolImages {
		var face *Face
		fname := file.Name()
		ipath := filepath.Join(idir, fname)
		face, err = recognizeIdolImage(ipath)
		if err != nil {
			return
		}
		if face == nil {
			log.Printf("Not a single face on %s", ipath)
			continue
		}
		faces = append(faces, *face)
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
		var faces []Face
		iname := dir.Name()
		key := [2]string{bname, iname}
		idol, ok := idolByNames[key]
		if !ok {
			err = fmt.Errorf("Can't find %s (%s)", iname, bname)
			return
		}
		idir := filepath.Join(bdir, iname)
		log.Printf("Importing %s", idir)
		faces, err = recognizeIdolImages(idir)
		if err != nil {
			return
		}
		idolId := idol["id"].(string)
		for _, face := range faces {
			descrSize := unsafe.Sizeof(face.descriptor)
			descrPtr := unsafe.Pointer(&face.descriptor)
			descrBytes := (*[1 << 30]byte)(descrPtr)[:descrSize:descrSize]
			if _, err = st.Exec(descrBytes, face.imageId, idolId); err != nil {
				return
			}
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
