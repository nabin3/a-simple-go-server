package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps               []Chirp              `json:"chirps"`
	Users                map[string]User      `json:"users"`
	RevokedRefreshTokens map[string]time.Time `json:"revoked_refreshed_tokens"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID           int    `json:"id"`
	PasswordHash string `json:"password_hash"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	// Intilalizing our database file
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.ReadFile(db.path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			DBStructureBytes, err1 := json.Marshal(DBStructure{make([]Chirp, 0), make(map[string]User), make(map[string]time.Time)})
			if err1 != nil {
				return fmt.Errorf("dataBase file does not exist and in time of creating database file we needed an empty data(json), in process of creating that empty json we faced this error: %v", err1)
			}

			err2 := os.WriteFile(db.path, DBStructureBytes, 0644)
			if err2 != nil {
				return fmt.Errorf("database file does not exist and Can't create database file, error: %v", err2)
			}
		}
	}

	return nil
}
