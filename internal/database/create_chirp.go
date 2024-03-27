package database

// CreateChirp creates a new chirp and saves it to disk and return that new user
func (db *DB) CreateChirp(userEmail, chirpMSG string) (Chirp, error) {
	dataFromDatabase, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{
		EmailID: userEmail,
		MSG:     chirpMSG,
	}

	// Adding new Chirp to Chirps map of our databse
	dataFromDatabase.Chirps[len(dataFromDatabase.Chirps)+1] = newChirp

	// Writting back to database
	if err = db.writeDB(dataFromDatabase); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}
