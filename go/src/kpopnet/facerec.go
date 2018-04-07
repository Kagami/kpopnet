package kpopnet

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"mime/multipart"

	"github.com/Kagami/go-dlib"
)

const (
	// Maximum number of recognizer threads executing at the same time.
	numRecWorkers = 1

	minDimension = 300
	maxDimension = 5000
)

var (
	faceRec *dlib.FaceRec
	recJobs = make(chan recRequest)
)

type recRequest struct {
	fh *multipart.FileHeader
	ch chan<- recResult
}

type recResult struct {
	idolId *string
	err    error
}

type TrainData struct {
	labels  []string
	samples []dlib.FaceDescriptor
}

func StartFaceRec(dataDir string) (err error) {
	faceRec, err = dlib.NewFaceRec(getModelsDir(dataDir))
	if err != nil {
		return fmt.Errorf("Error initializing face recognizer: %v", err)
	}
	for i := 0; i < numRecWorkers; i++ {
		go recWorker()
	}
	return
}

// Execute recognizing jobs.
func recWorker() {
	for {
		req := <-recJobs
		idolId, err := RecognizeMultipart(req.fh)
		req.ch <- recResult{idolId, err}
	}
}

// Recognize user-provided image with the specific concurrency level.
// Note that we don't read file beforehand to minimize memory
// consumption.
func RequestRecognizeMultipart(fh *multipart.FileHeader) (idolId *string, err error) {
	ch := make(chan recResult)
	go func() {
		recJobs <- recRequest{fh, ch}
	}()
	res := <-ch
	return res.idolId, res.err
}

// Simple wrapper to work with uploaded files.
// Recognize immediately.
func RecognizeMultipart(fh *multipart.FileHeader) (idolId *string, err error) {
	fd, err := fh.Open()
	if err != nil {
		err = errParseFile
		return
	}
	defer fd.Close()
	imgData, err := ioutil.ReadAll(fd)
	if err != nil {
		err = errParseFile
		return
	}
	idolId, err = Recognize(imgData)
	return
}

// Recognize immediately.
// TODO(Kagami): Search for already recognized idol using imageId.
func Recognize(imgData []byte) (idolId *string, err error) {
	// TODO(Kagami): Invalidate?
	v, err := cached(trainDataCacheKey, func() (interface{}, error) {
		trainData, err := GetTrainData()
		if err == nil {
			// NOTE(Kagami): We don't copy this data to C++ side so need to
			// keep in cache to prevent GC.
			faceRec.SetSamples(trainData.samples)
		}
		return trainData, err
	})
	if err != nil {
		return
	}
	trainData := v.(*TrainData)

	r := bytes.NewReader(imgData)
	c, typ, err := image.DecodeConfig(r)
	if err != nil || typ != "jpeg" ||
		c.Width < minDimension ||
		c.Height < minDimension ||
		c.Width > maxDimension ||
		c.Height > maxDimension {
		err = errBadImage
		return
	}

	face, err := faceRec.RecognizeSingle(imgData)
	if err != nil || face == nil {
		return
	}

	idx := faceRec.Classify(face.Descriptor)
	if idx < 0 {
		err = errNoIdol
		return
	}
	idolId = &trainData.labels[idx]
	return
}
