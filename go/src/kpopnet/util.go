package kpopnet

import (
	"unsafe"

	"github.com/Kagami/go-dlib"
)

// Zero-copy conversions.

func descr2bytes(d dlib.FaceDescriptor) []byte {
	size := unsafe.Sizeof(d)
	return (*[1 << 30]byte)(unsafe.Pointer(&d))[:size:size]
}

func bytes2descr(b []byte) dlib.FaceDescriptor {
	return *(*dlib.FaceDescriptor)(unsafe.Pointer(&b[0]))
}
