package kpopnet

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	INDEX_NAME = "index"
)

// Band info.
type Band map[string]interface{}

// Idol info.
type Idol map[string]interface{}

// All bands and idols.
type Profiles struct {
	Bands []Band `json:"bands"`
	Idols []Idol `json:"idols"`
}

func checkName(name string) {
	if name == INDEX_NAME {
		panic("Bad name")
	}
}

func getProfilesDir(d string) string {
	return filepath.Join(d, "profiles")
}

func getBandDir(d string, bname string) string {
	checkName(bname)
	return filepath.Join(d, "profiles", bname)
}

func getBandPath(d string, bname string) string {
	return filepath.Join(getBandDir(d, bname), INDEX_NAME+".json")
}

func getIdolPath(d string, bname string, iname string) string {
	checkName(iname)
	return filepath.Join(getBandDir(d, bname), iname+".json")
}

func readBandIdols(d string, bname string) (idols []Idol, err error) {
	idolFiles, err := ioutil.ReadDir(getBandDir(d, bname))
	if err != nil {
		return
	}
	for _, ifile := range idolFiles {
		var data []byte
		var idol Idol
		iname := strings.TrimSuffix(ifile.Name(), ".json")
		if iname == INDEX_NAME {
			continue
		}
		data, err = ioutil.ReadFile(getIdolPath(d, bname, iname))
		if err != nil {
			return
		}
		if err = json.Unmarshal(data, &idol); err != nil {
			return
		}
		idols = append(idols, idol)
	}
	return
}

// Read all profiles from JSON-encoded files in provided directory.
func ReadProfiles(d string) (ps *Profiles, err error) {
	var bands []Band
	var idols []Idol

	bandDirs, err := ioutil.ReadDir(getProfilesDir(d))
	if err != nil {
		return
	}
	for _, bdir := range bandDirs {
		var data []byte
		var band Band
		bname := bdir.Name()
		data, err = ioutil.ReadFile(getBandPath(d, bname))
		if err != nil {
			return
		}
		// NOTE(Kagami): We don't validate decoded structs here (e.g.
		// mandatory id/name fields) because it will be checked by
		// PostgreSQL table constraints.
		if err = json.Unmarshal(data, &band); err != nil {
			return
		}
		bands = append(bands, band)

		var bandIdols []Idol
		bandIdols, err = readBandIdols(d, bname)
		if err != nil {
			return
		}
		idols = append(idols, bandIdols...)
	}

	ps = &Profiles{
		Bands: bands,
		Idols: idols,
	}
	return
}
