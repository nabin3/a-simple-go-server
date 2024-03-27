package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nabin3/a-simple-web-server/internal/auth"
	"github.com/nabin3/a-simple-web-server/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Defining response body format
	type resp struct {
		Token string `json:"token"`
	}

	// Getting token from header
	unverifiedToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	// Checking the type of token by exmining the issuer
	tokenIssuer, err := auth.RetrieveJwtTokenIssuer(unverifiedToken)
	if err != nil {
		log.Printf("error in handlerRefresh at auth.RetrieveJwtTokenIssuer: %v", err)
		respondWithError(w, 500, "server-error")
		return
	}
	if tokenIssuer != "chirpy-refresh" {
		respondWithError(w, 401, "bad token")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

	// Checking if the refresh token is already revoked
	err = db.AlreadyRevoked(unverifiedToken)
	if err != nil {
		respondWithError(w, 401, "already revoked")
		return
	}

	// Validating recieved refresh token
	subjectEmail, err := auth.ValidateJWT(unverifiedToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error in handlerRefresh at auth.ValidateJWT: %v", err)
		respondWithError(w, 401, "token couldn't be validated")
		return
	}

	// Generating a new access token
	newAccessToken, err := auth.JwtAccessTokenGenerator(subjectEmail, cfg.jwtSecret)
	if err != nil {
		log.Printf("error in handlerRefresh at auth.JwtRefreshTokenGenarator: %v", err)
		respondWithError(w, 500, "couldn't generate access token")
		return
	}

	respondWithJSON(w, 200, resp{
		Token: newAccessToken,
	},
	)
}
