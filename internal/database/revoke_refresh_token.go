package database

import (
	"fmt"
	"time"
)

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	data, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, exist := data.RevokedRefreshTokens[refreshToken]; !exist {
		data.RevokedRefreshTokens[refreshToken] = time.Now()
		return nil
	}

	return fmt.Errorf("given refreshed token has already been revoked")
}
