// Database queries.
package db

// Band info.
type Band map[string]interface{}

// Idol info.
type Idol map[string]interface{}

type Profiles struct {
	Bands []Band `json:"bands"`
	Idols []Idol `json:"idols"`
}

// Get all profiles.
// FIXME(Kagami): Cache it!
func GetProfiles() (ps *Profiles, err error) {
	ps = &Profiles{
		Bands: []Band{},
		Idols: []Idol{},
	}
	return
}
