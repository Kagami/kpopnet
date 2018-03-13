// Database queries.
package db

import (
	"encoding/json"

	"kpopnet/profile"
)

// Get all profiles.
// FIXME(Kagami): Cache it!
func GetProfiles() (ps *profile.Profiles, err error) {
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
	bands := []profile.Band{}
	for rs.Next() {
		var id string
		var data []byte
		var band profile.Band
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
	idols := []profile.Idol{}
	for rs2.Next() {
		var id string
		var bandId string
		var data []byte
		var idol profile.Idol
		if err = rs2.Scan(&id, &bandId, &data); err != nil {
			return
		}
		if err = json.Unmarshal(data, &idol); err != nil {
			return
		}
		idol["id"] = id
		idol["bandId"] = bandId
		idols = append(idols, idol)
	}
	if err = rs2.Err(); err != nil {
		return
	}

	ps = &profile.Profiles{
		Bands: bands,
		Idols: idols,
	}
	return
}

// Insert/update database profiles.
func UpdateProfiles(ps *profile.Profiles) (err error) {
	return
}
