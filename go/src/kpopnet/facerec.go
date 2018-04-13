// +build !nodlib

package kpopnet

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"unsafe"

	"github.com/Kagami/go-face"
)

const (
	// Maximum number of recognizer threads executing at the same time.
	numRecWorkers = 1

	minDimension = 300
	maxDimension = 5000
)

var (
	faceRec *face.Recognizer
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

type trainData struct {
	labels  []string
	samples []face.Descriptor
}

func StartFaceRec(dataDir string) (err error) {
	faceRec, err = face.NewRecognizer(getModelsDir(dataDir))
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
		idolId, err := recognizeMultipart(req.fh)
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
func recognizeMultipart(fh *multipart.FileHeader) (idolId *string, err error) {
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
	idolId, err = recognize(imgData)
	return
}

// Recognize immediately.
// TODO(Kagami): Search for already recognized idol using imageId.
func recognize(imgData []byte) (idolId *string, err error) {
	// TODO(Kagami): Invalidate?
	v, err := cached(trainDataCacheKey, func() (interface{}, error) {
		data, err := getTrainData()
		if err == nil {
			// NOTE(Kagami): We don't copy this data to C++ side so need to
			// keep in cache to prevent GC.
			faceRec.SetSamples(data.samples)
		}
		return data, err
	})
	if err != nil {
		return
	}
	data := v.(*trainData)

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
	idolId = &data.labels[idx]
	return
}

// Get all confirmed face descriptors.
func getTrainData() (data *trainData, err error) {
	var labels []string
	var samples []face.Descriptor

	rs, err := prepared["get_train_data"].Query()
	if err != nil {
		return
	}
	defer rs.Close()
	for rs.Next() {
		var idolId string
		var descrBytes []byte
		if err = rs.Scan(&idolId, &descrBytes); err != nil {
			return
		}
		labels = append(labels, idolId)
		descriptor := bytes2descr(descrBytes)
		samples = append(samples, descriptor)
	}
	if err = rs.Err(); err != nil {
		return
	}

	data = &trainData{
		labels:  labels,
		samples: samples,
	}
	return
}

// Zero-copy conversions.

func descr2bytes(d face.Descriptor) []byte {
	size := unsafe.Sizeof(d)
	return (*[1 << 30]byte)(unsafe.Pointer(&d))[:size:size]
}

func bytes2descr(b []byte) face.Descriptor {
	return *(*face.Descriptor)(unsafe.Pointer(&b[0]))
}
