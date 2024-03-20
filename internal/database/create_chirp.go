package database

// CreateChirp creates a new chirp and saves it to disk and return that new user
func (db *DB) CreateChirp(chirp string) (Chirp, error) {
	dataFromDatabase, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{
		ID:   len(dataFromDatabase.Chirps) + 1,
		Body: chirp,
	}

	dataFromDatabase.Chirps = append(dataFromDatabase.Chirps, newChirp)

	if err = db.writeDB(dataFromDatabase); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}
