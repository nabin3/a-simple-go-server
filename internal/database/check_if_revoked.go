package database

import "fmt"

func (db *DB) AlreadyRevoked(refreshToken string) error {
	data, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, exist := data.RevokedRefreshTokens[refreshToken]; !exist {
		return fmt.Errorf("this refresh token has not been revoked yet")
	}

	return nil
}
