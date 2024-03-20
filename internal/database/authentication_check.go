package database

import (
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CheckUserCreditional(email, password string) (int, error) {
	datafromDatabase, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	user, exists := datafromDatabase.Users[email]
	if !exists {
		return 0, nil
	}

	// This folloing function return nil on success and not nil error on unsuccess
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, nil
	}

	return user.ID, nil
}
