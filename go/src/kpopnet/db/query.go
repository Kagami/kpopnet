// Database queries.
package db

import (
	"bytes"
	"encoding/json"
)

type Profiles struct {
	Bands []json.RawMessage `json:"bands"`
	Idols []json.RawMessage `json:"idols"`
}

// Add ID key to encoded JSON data.
func fixBandData(buf []byte, id string) []byte {
	return bytes.Join([][]byte{
		[]byte("{\"id\":\""),
		[]byte(id),
		[]byte("\","),
		buf[1:]}, nil)
}

// Add ID and band ID keys to encoded JSON data.
func fixIdolData(buf []byte, id string, bandId string) []byte {
	return bytes.Join([][]byte{
		[]byte("{\"id\":\""),
		[]byte(id),
		[]byte("\",\"band_id\":\""),
		[]byte(bandId),
		[]byte("\","),
		buf[1:]}, nil)
}

// Get all profiles.
// FIXME(Kagami): Cache it!
func GetProfiles() (ps *Profiles, err error) {
	tx, err := getTx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	if setReadOnly(tx) != nil {
		return
	}
	if setRepeatableRead(tx) != nil {
		return
	}

	r, err := tx.Stmt(prepared["get_bands"]).Query()
	if err != nil {
		return
	}
	defer r.Close()
	bands := []json.RawMessage{}
	for r.Next() {
		var id string
		var data []byte
		if err = r.Scan(&id, &data); err != nil {
			return
		}
		bands = append(bands, fixBandData(data, id))
	}
	if err = r.Err(); err != nil {
		return
	}

	r2, err := tx.Stmt(prepared["get_idols"]).Query()
	if err != nil {
		return
	}
	defer r2.Close()
	idols := []json.RawMessage{}
	for r2.Next() {
		var id string
		var bandId string
		var data []byte
		if err = r2.Scan(&id, &bandId, &data); err != nil {
			return
		}
		idols = append(idols, fixIdolData(data, id, bandId))
	}
	if err = r2.Err(); err != nil {
		return
	}

	ps = &Profiles{
		Bands: bands,
		Idols: idols,
	}
	return
}
