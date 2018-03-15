package kpopnet

import (
	"encoding/json"
)

// Get all profiles.
func GetProfiles() (ps *Profiles, err error) {
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)
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
	bands := []Band{}
	for rs.Next() {
		var id string
		var data []byte
		var band Band
		if err = rs.Scan(&id, &data); err != nil {
			return
		}
		if err = json.Unmarshal(data, &band); err != nil {
			return
		}
		band["id"] = id
		bands = append(bands, band)
	}
	if err = rs.Err(); err != nil {
		return
	}

	rs2, err := tx.Stmt(prepared["get_idols"]).Query()
	if err != nil {
		return
	}
	defer rs2.Close()
	idols := []Idol{}
	for rs2.Next() {
		var id string
		var bandId string
		var data []byte
		var idol Idol
		if err = rs2.Scan(&id, &bandId, &data); err != nil {
			return
		}
		if err = json.Unmarshal(data, &idol); err != nil {
			return
		}
		idol["id"] = id
		idol["band_id"] = bandId
		idols = append(idols, idol)
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

// Prepare band structure to be stored in DB.
// ID fields are removed to avoid duplication.
func getBandData(band Band) (data []byte, err error) {
	delete(band, "id")
	delete(band, "urls") // Don't need this
	data, err = json.Marshal(band)
	return
}

// Prepare idol structure to be stored in DB.
// ID fields are removed to avoid duplication.
func getIdolData(idol Idol) (data []byte, err error) {
	delete(idol, "id")
	delete(idol, "band_id")
	data, err = json.Marshal(idol)
	return
}

// Insert/update database profiles.
func UpdateProfiles(ps *Profiles) (err error) {
	tx, err := beginTx()
	if err != nil {
		return
	}
	defer endTx(tx, &err)

	st := tx.Stmt(prepared["upsert_band"])
	for _, band := range ps.Bands {
		id := band["id"]
		var data []byte
		data, err = getBandData(band)
		if err != nil {
			return
		}
		if _, err = st.Exec(id, data); err != nil {
			return
		}
	}

	st = tx.Stmt(prepared["upsert_idol"])
	for _, idol := range ps.Idols {
		id := idol["id"]
		bandId := idol["band_id"]
		var data []byte
		data, err = getIdolData(idol)
		if err != nil {
			return
		}
		if _, err = st.Exec(id, bandId, data); err != nil {
			return
		}
	}

	return
}
