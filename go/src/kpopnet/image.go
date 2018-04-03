package kpopnet

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

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

func recognizeIdolImages(idir string) (ds []dlib.FaceDescriptor, err error) {
	idolImages, err := ioutil.ReadDir(idir)
	if err != nil {
		return
	}
	// No need to validate names/formats because everything was checked by
	// Python spider.
	for _, file := range idolImages {
		var d *dlib.FaceDescriptor
		fname := file.Name()
		ipath := filepath.Join(idir, fname)
		d, err = faceRec.GetDescriptor(ipath)
		if err != nil {
			return
		}
		if d == nil {
			log.Printf("Not a single face on %s", ipath)
			continue
		}
		ds = append(ds, *d)
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

	idolDirs, err := ioutil.ReadDir(bdir)
	if err != nil {
		return
	}
	for _, dir := range idolDirs {
		var ds []dlib.FaceDescriptor
		iname := dir.Name()
		key := [2]string{bname, iname}
		idol, ok := idolByNames[key]
		if !ok {
			err = fmt.Errorf("Can't find %s (%s)", iname, bname)
			return
		}
		idir := filepath.Join(bdir, iname)
		log.Printf("Importing %s", idir)
		ds, err = recognizeIdolImages(idir)
		if err != nil {
			return
		}
		fmt.Println(idol["id"], len(ds))
	}
	return
}

// Read and update idol faces in database.
func ImportImages(connStr string, dataDir string) (err error) {
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

	for _, dir := range bandDirs {
		bname := dir.Name()
		bdir := filepath.Join(getImagesDir(dataDir), bname)
		if err = importBandImages(bdir, bname, idolByNames); err != nil {
			err = fmt.Errorf("Error importing %s images: %v", bname, err)
			return
		}
		break
	}
	return
}
