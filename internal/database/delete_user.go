package database

import "fmt"

func (db *DB) DeleteUser(email string) error {
	data, err := db.loadDB()
	if err != nil {
		return fmt.Errorf("colud not load data from database: %v", err)
	}

	delete(data.Users, email)

	return nil
}
