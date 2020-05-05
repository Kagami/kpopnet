package kpopnet

import (
	"encoding/json"
	"fmt"
	"image"
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

// TODO(Kagami): Add close matches field to simplify confirmation.
type ImageInfo struct {
	Rectangle image.Rectangle
	IdolId    string
	Confirmed bool
}

func (i ImageInfo) MarshalJSON() ([]byte, error) {
	r := i.Rectangle
	s := fmt.Sprintf(
		`{"rect":[%d,%d,%d,%d],"id":"%s","confirmed":"%v"}`,
		r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, i.IdolId, i.Confirmed)
	return []byte(s), nil
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
	return filepath.Join(getProfilesDir(d), bname)
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
	for _, dir := range bandDirs {
		var data []byte
		var band Band
		bname := dir.Name()
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

// Read and update profiles in database.
func ImportProfiles(connStr string, dataDir string) (err error) {
	if err = StartDb(nil, connStr); err != nil {
		return
	}
	ps, err := ReadProfiles(dataDir)
	if err != nil {
		err = fmt.Errorf("error reading profiles: %v", err)
		return
	}
	err = UpdateProfiles(ps)
	if err != nil {
		err = fmt.Errorf("error updating DB profiles: %v", err)
		return
	}
	return
}
