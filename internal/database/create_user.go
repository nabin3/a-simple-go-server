package database

import (
	"fmt"

	"github.com/nabin3/a-simple-web-server/internal/auth"
)

type RespFormat struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// CreateUsers creates a new user and saves it to disk and return that new user
func (db *DB) CreateUser(givenId int, email, password string) (RespFormat, error) {
	// Retrieving databse data
	dataFromDatabase, err := db.loadDB()
	if err != nil {
		return RespFormat{}, err
	}

	// Checking if user did not give a pass word
	if len(password) == 0 {
		return RespFormat{}, fmt.Errorf("please give a password")
	}

	// Creating a hash for given password
	passHash, err := auth.CreatePasswordHash(password)
	if err != nil {
		return RespFormat{}, err
	}

	// Creating a new user
	var id int
	if givenId == 0 {
		id = len(dataFromDatabase.Users) + 1
	} else {
		id = givenId
	}
	newUser := User{
		ID:           id,
		PasswordHash: passHash,
	}

	// Assigning a new user
	if _, exists := dataFromDatabase.Users[email]; !exists {
		dataFromDatabase.Users[email] = newUser
	} else {
		return RespFormat{}, fmt.Errorf("email alredy in use")
	}

	// Writting to database file
	if err = db.writeDB(dataFromDatabase); err != nil {
		return RespFormat{}, err
	}

	return RespFormat{newUser.ID, email}, nil
}
