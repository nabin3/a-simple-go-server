package main

import (
	"log"
	"net/http"

	"github.com/nabin3/a-simple-web-server/internal/auth"
	"github.com/nabin3/a-simple-web-server/internal/database"
)

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error in handlerRevokeRefreshToken at auth.GetBearerToken: %v", err)
		respondWithError(w, 401, "bad request")
		return
	}

	_, err = auth.ValidateJWT(refreshToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("issue at handlerRevokeRefreshToken at auth.ValidateJWT: %v", err)
		respondWithError(w, 401, "can't validate token")
		return
	}

	// Creating database connection
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("isuue at handlerRevokeRefreshToken, can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

	err = db.RevokeRefreshToken(refreshToken)
	if err != nil {
		log.Printf("issue in handlerRevokeRefreshToken at database.RevokeRefreshToken: %v", err)
		respondWithError(w, 401, "")
		return
	}

	respondWithJSON(w, 200, "successfully revoked")
}
