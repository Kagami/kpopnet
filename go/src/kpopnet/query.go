package kpopnet

import (
	"database/sql"
	"encoding/json"
)

// Get all bands.
func getBands(tx *sql.Tx) (bands []Band, bandById map[string]Band, err error) {
	bands = make([]Band, 0)
	bandById = make(map[string]Band)
	rs, err := tx.Stmt(prepared["get_bands"]).Query()
	if err != nil {
		return
	}
	defer rs.Close()
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
		bandById[id] = band
	}
	if err = rs.Err(); err != nil {
		return
	}
	return
}

// Get all idols.
func getIdols(tx *sql.Tx) (idols []Idol, idolById map[string]Idol, err error) {
	idols = make([]Idol, 0)
	idolById = make(map[string]Idol)
	rs, err := tx.Stmt(prepared["get_idols"]).Query()
	if err != nil {
		return
	}
	defer rs.Close()
	for rs.Next() {
		var id string
		var bandId string
		var data []byte
		var idol Idol
		if err = rs.Scan(&id, &bandId, &data); err != nil {
			return
		}
		if err = json.Unmarshal(data, &idol); err != nil {
			return
		}
		idol["id"] = id
		idol["band_id"] = bandId
		idols = append(idols, idol)
		idolById[id] = idol
	}
	if err = rs.Err(); err != nil {
		return
	}
	return
}

// Get and set idol preview property.
func getIdolPreviews(tx *sql.Tx, idolById map[string]Idol) (err error) {
	rs, err := tx.Stmt(prepared["get_idol_previews"]).Query()
	if err != nil {
		return
	}
	defer rs.Close()
	for rs.Next() {
		var idolId string
		var imageId string
		if err = rs.Scan(&idolId, &imageId); err != nil {
			return
		}
		if idol, ok := idolById[idolId]; ok {
			idol["image_id"] = imageId
		}
	}
	if err = rs.Err(); err != nil {
		return
	}
	return
}

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

	bands, _, err := getBands(tx)
	if err != nil {
		return
	}
	idols, idolById, err := getIdols(tx)
	if err != nil {
		return
	}
	err = getIdolPreviews(tx, idolById)
	if err != nil {
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

func getImageInfo(imageId string) (info *ImageInfo, err error) {
	var rectStr string
	var idolId string
	var confirmed bool
	err = prepared["get_face"].QueryRow(imageId).Scan(&rectStr, &idolId, &confirmed)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errNoIdol
		}
		return
	}
	rect := str2rect(rectStr)
	info = &ImageInfo{rect, idolId, confirmed}
	return
}
