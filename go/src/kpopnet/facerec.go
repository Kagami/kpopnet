// +build !nodlib

package kpopnet

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
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
	samples []face.Descriptor
	cats    []int32
	labels  map[int]string
}

func StartFaceRec(dataDir string) error {
	return startFaceRec(getModelsDir(dataDir))
}

// Useful for tests.
func startFaceRec(modelsDir string) (err error) {
	faceRec, err = face.NewRecognizer(modelsDir)
	if err != nil {
		return fmt.Errorf("error initializing face recognizer: %v", err)
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
			faceRec.SetSamples(data.samples, data.cats)
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
		c.Height > maxDimension ||
		c.ColorModel != color.YCbCrModel {
		err = errBadImage
		return
	}

	f, err := faceRec.RecognizeSingle(imgData)
	if _, ok := err.(face.ImageLoadError); ok {
		err = errBadImage
	}
	if err != nil || f == nil {
		return
	}

	catIdx := faceRec.Classify(f.Descriptor)
	if catIdx < 0 {
		err = errNoIdol
		return
	}
	id := data.labels[catIdx]
	return &id, nil
}

// Get all confirmed face descriptors.
func getTrainData() (data *trainData, err error) {
	var samples []face.Descriptor
	var cats []int32
	labels := make(map[int]string)

	rs, err := prepared["get_train_data"].Query()
	if err != nil {
		return
	}
	defer rs.Close()
	var catIdx int32
	var prevIdolId string
	catIdx = -1
	for rs.Next() {
		var idolId string
		var descrBytes []byte
		if err = rs.Scan(&idolId, &descrBytes); err != nil {
			return
		}
		descriptor := bytes2descr(descrBytes)
		samples = append(samples, descriptor)
		if idolId != prevIdolId {
			catIdx++
			labels[int(catIdx)] = idolId
		}
		cats = append(cats, catIdx)
		prevIdolId = idolId
	}
	if err = rs.Err(); err != nil {
		return
	}

	data = &trainData{
		samples: samples,
		cats:    cats,
		labels:  labels,
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
