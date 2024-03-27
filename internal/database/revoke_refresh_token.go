package database

import (
	"fmt"
	"time"
)

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	// Reading from database
	data, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, exist := data.RevokedRefreshTokens[refreshToken]; !exist {
		data.RevokedRefreshTokens[refreshToken] = time.Now().UTC()

		// Writting back to database
		err = db.writeDB(data)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("given refreshed token has already been revoked")
}
