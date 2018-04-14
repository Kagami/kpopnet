package kpopnet

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	testConn = "user=meguca password=meguca dbname=meguca sslmode=disable"
)

var (
	testIdols = map[string]string{
		"elkie1": "235a1504-e54b-4aae-bb11-a2e33d3c2ea8",
	}
)

func getTestFilePath(name string) string {
	return filepath.Join("testdata", name+".jpg")
}

func recognizeFile(fpath string) (idolId *string, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return
	}
	imgData, err := ioutil.ReadAll(fd)
	if err != nil {
		return
	}
	return recognize(imgData)
}

func TestIdols(t *testing.T) {
	if err := StartDb(nil, testConn); err != nil {
		t.Fatal(err)
	}
	if err := startFaceRec("testdata"); err != nil {
		t.Fatal(err)
	}
	for name, expectedIdolId := range testIdols {
		actualIdolId, err := recognizeFile(getTestFilePath(name))
		if err != nil {
			t.Fatal(err)
		}
		if actualIdolId == nil {
			t.Errorf("%s: no result", name)
		} else if expectedIdolId != *actualIdolId {
			t.Errorf("%s: expected %s but got %s", name, expectedIdolId, *actualIdolId)
		}
	}
}
