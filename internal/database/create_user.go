package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type RespFormat struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// CreateUsers creates a new user and saves it to disk and return that new user
func (db *DB) CreateUser(email, password string) (RespFormat, error) {
	// Retrieving databse data
	dataFromDatabase, err := db.loadDB()
	if err != nil {
		return RespFormat{}, err
	}

	// Checking if user did not give a pass word
	if len(password) == 0 {
		return RespFormat{}, fmt.Errorf("Please give a password")
	}

	// Creating a hash for given password
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return RespFormat{}, err
	}

	// Creating a new user
	newUser := User{
		ID:           len(dataFromDatabase.Users) + 1,
		PasswordHash: string(hashByte),
	}

	// Assigning a new user
	if _, exists := dataFromDatabase.Users[email]; !exists {
		return RespFormat{}, fmt.Errorf("email alredy in use")
	} else {
		dataFromDatabase.Users[email] = newUser
	}

	// Writting to database file
	if err = db.writeDB(dataFromDatabase); err != nil {
		return RespFormat{}, err
	}

	return RespFormat{newUser.ID, email}, nil
}
