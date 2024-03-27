// GetChirps returns all chirps in the database
package database

import "errors"

func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	all_chirps := make([]Chirp, 0, len(data.Chirps))
	for _, single_chirp := range data.Chirps {
		all_chirps = append(all_chirps, single_chirp)
	}

	return all_chirps, nil
}

// GetChirpByID returns a chirp specific to a given ID
func (db *DB) GetChirpByID(id int) (Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if single_chirp, exist := data.Chirps[id]; !exist {
		return Chirp{}, errors.New("chirp with associated id doesn't exist")
	} else {
		return single_chirp, nil
	}
}
