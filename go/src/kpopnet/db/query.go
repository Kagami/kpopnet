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
	if err = setReadOnly(tx); err != nil {
		return
	}
	if err = setRepeatableRead(tx); err != nil {
		return
	}

	rs, err := tx.Stmt(prepared["get_bands"]).Query()
	if err != nil {
		return
	}
	defer rs.Close()
	bands := []json.RawMessage{}
	for rs.Next() {
		var id string
		var data []byte
		if err = rs.Scan(&id, &data); err != nil {
			return
		}
		bands = append(bands, fixBandData(data, id))
	}
	if err = rs.Err(); err != nil {
		return
	}

	rs2, err := tx.Stmt(prepared["get_idols"]).Query()
	if err != nil {
		return
	}
	defer rs2.Close()
	idols := []json.RawMessage{}
	for rs2.Next() {
		var id string
		var bandId string
		var data []byte
		if err = rs2.Scan(&id, &bandId, &data); err != nil {
			return
		}
		idols = append(idols, fixIdolData(data, id, bandId))
	}
	if err = rs2.Err(); err != nil {
		return
	}

	ps = &Profiles{
		Bands: bands,
		Idols: idols,
	}
	return
}
