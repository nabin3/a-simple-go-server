package database

import "fmt"

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
