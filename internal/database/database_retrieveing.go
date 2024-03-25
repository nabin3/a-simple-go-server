package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	jsonData, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	data := DBStructure{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return DBStructure{}, err
	}

	return data, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	return data.Chirps, nil
}

// GetChirpByID returns a chirp specific to a given ID
func (db *DB) GetChirpByID(id int) (Chirp, error) {
	data, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if id <= len(data.Chirps) && id > 0 {
		return data.Chirps[id-1], nil
	} else {
		return Chirp{}, fmt.Errorf("chirp with id: %d doesn't exist", id)
	}
}

// Get user by email
func (db *DB) GetUserByEmail(email string) (User, error) {
	data, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, exist := data.Users[email]
	if !exist {
		return User{}, fmt.Errorf("user with the given email id doesn't exist")
	}

	return user, nil
}
