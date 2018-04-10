// +build nodlib

package kpopnet

import (
	"mime/multipart"
)

func StartFaceRec(dataDir string) (err error) {
	return
}

func RequestRecognizeMultipart(fh *multipart.FileHeader) (idolId *string, err error) {
	err = errNoIdol
	return
}
