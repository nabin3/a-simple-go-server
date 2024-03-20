package database

import (
	"encoding/json"
	"os"
)

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	jsonyFiedData, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, jsonyFiedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
