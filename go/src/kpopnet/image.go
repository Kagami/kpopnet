package kpopnet

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/Kagami/go-dlib"
)

var (
	faceRec *dlib.FaceRec
)

func getImagesDir(d string) string {
	return filepath.Join(d, "images")
}

func getModelsDir(d string) string {
	return filepath.Join(d, "models")
}

func getIdolNameMap() (idolByName map[string]Idol, err error) {
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)
	if err = setReadOnly(tx); err != nil {
		return
	}
	idols, _, err := getIdols(tx)
	if err != nil {
		return
	}
	idolByName = make(map[string]Idol)
	for _, idol := range idols {
		if name, ok := idol["name"].(string); ok {
			idolByName[name] = idol
		}
	}
	return
}

func recognizeIdolImages(idir string) (faces []*dlib.Face, err error) {
	idolImages, err := ioutil.ReadDir(idir)
	if err != nil {
		return
	}
	// No need to validate names/formats because everything was checked by
	// Python spider.
	for _, file := range idolImages {
		var idolFaces []*dlib.Face
		fname := file.Name()
		ipath := filepath.Join(idir, fname)
		idolFaces, err = faceRec.Recognize(ipath)
		if err != nil {
			return
		}
		// TODO(Kagami): RecognizeSingle.
		if len(idolFaces) != 1 {
			log.Printf("Not a single idol face on %s (%d)", ipath, len(idolFaces))
			continue
		}
		faces = append(faces, idolFaces[0])
	}
	return
}

func importBandImages(bdir string, idolByName map[string]Idol) (err error) {
	// Use single transaction per band.
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)

	idolDirs, err := ioutil.ReadDir(bdir)
	if err != nil {
		return
	}
	for _, dir := range idolDirs {
		var faces []*dlib.Face
		iname := dir.Name()
		idol, ok := idolByName[iname]
		if !ok {
			continue
		}
		idir := filepath.Join(bdir, iname)
		log.Printf("Importing %s", idir)
		faces, err = recognizeIdolImages(idir)
		if err != nil {
			return
		}
		fmt.Println(idol["id"], faces)
	}
	return
}

// Read and update idol faces in database.
func ImportImages(connStr string, dataDir string) (err error) {
	if err = StartDb(nil, connStr); err != nil {
		return
	}

	idolByName, err := getIdolNameMap()
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

	for _, dir := range bandDirs {
		bname := dir.Name()
		bdir := filepath.Join(getImagesDir(dataDir), bname)
		if err = importBandImages(bdir, idolByName); err != nil {
			err = fmt.Errorf("Error importing %s images: %v", bname, err)
			return
		}
		break
	}
	return
}
