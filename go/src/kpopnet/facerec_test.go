package kpopnet

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testConn = "user=meguca password=meguca dbname=meguca sslmode=disable"
)

var (
	testData = map[string]string{
		"elkie.jpg":      "Elkie, CLC",
		"chaeyoung.jpg":  "Chaeyoung, Twice",
		"chaeyoung2.jpg": "Chaeyoung, Twice",
		"sejeong.jpg":    "Sejeong, Gugudan",
		"jimin.jpg":      "Jimin, AOA",
		"jimin2.jpg":     "Jimin, AOA",
		"jimin4.jpg":     "Jimin, AOA",
	}
)

func getMaps() (idolById map[string]Idol, bandById map[string]Band, err error) {
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)
	if _, idolById, err = getIdols(tx); err != nil {
		return
	}
	if _, bandById, err = getBands(tx); err != nil {
		return
	}
	return
}

func getTestFilePath(fname string) string {
	return filepath.Join("testdata", fname)
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
	idolById, bandById, err := getMaps()
	if err != nil {
		t.Fatal(err)
	}
	for fname, expected := range testData {
		t.Run(fname, func(t *testing.T) {
			names := strings.Split(expected, ", ")
			expectedIname := names[0]
			expectedBname := names[1]

			actualIdolId, err := recognizeFile(getTestFilePath(fname))
			if err != nil {
				t.Fatal(err)
			}
			if actualIdolId == nil {
				t.Errorf("%s: expected “%s” but not recognized", fname, expected)
				return
			}

			idol := idolById[*actualIdolId]
			band := bandById[idol["band_id"].(string)]
			actualIname := idol["name"]
			actualBname := band["name"]
			if expectedIname != actualIname || expectedBname != actualBname {
				actual := fmt.Sprintf("%s, %s", actualIname, actualBname)
				t.Errorf("%s: expected “%s” but got “%s”", fname, expected, actual)
			}
		})
	}
}
