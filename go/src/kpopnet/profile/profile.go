// Working with profile structures and serialization.
package profile

// Band info.
type Band map[string]interface{}

// Idol info.
type Idol map[string]interface{}

// All bands and idols.
type Profiles struct {
	Bands []Band `json:"bands"`
	Idols []Idol `json:"idols"`
}

// Read all profiles from JSON-encoded files in provided directory.
func ReadAll(datadir string) (ps *Profiles, err error) {
	return
}
