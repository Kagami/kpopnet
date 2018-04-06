package kpopnet

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"unsafe"

	"github.com/Kagami/go-dlib"
)

const (
	minDimension = 300
	maxDimension = 5000
)

var (
	faceRec *dlib.FaceRec
)

type TrainData struct {
	labels  []string
	samples []dlib.FaceDescriptor
}

// Zero-copy conversions.

func descr2bytes(d dlib.FaceDescriptor) []byte {
	size := unsafe.Sizeof(d)
	return (*[1 << 30]byte)(unsafe.Pointer(&d))[:size:size]
}

func bytes2descr(b []byte) dlib.FaceDescriptor {
	return *(*dlib.FaceDescriptor)(unsafe.Pointer(&b[0]))
}

func StartFaceRec(dataDir string) (err error) {
	faceRec, err = dlib.NewFaceRec(getModelsDir(dataDir))
	if err != nil {
		return fmt.Errorf("Error initializing face recognizer: %v", err)
	}
	return
}

// TODO(Kagami): Search for already recognized idol using imageId.
func Recognize(imgData []byte) (idolId *string, err error) {
	v, err := cached(trainDataCacheKey, func() (interface{}, error) {
		return GetTrainData()
	})
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

	idx, err := faceRec.Classify(trainData.samples, face.Descriptor)
	if err != nil {
		return
	}
	if idx < 0 {
		err = errNoIdol
		return
	}
	idolId = &trainData.labels[idx]
	return
}
