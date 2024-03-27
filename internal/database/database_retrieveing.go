package database

import (
	"encoding/json"
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
