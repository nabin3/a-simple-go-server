package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nabin3/a-simple-web-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Token string `json:"token"`
	}

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

	subjectEmail, err := auth.ValidateJWT(unverifiedToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error in handlerRefresh at auth.ValidateJWT: %v", err)
		respondWithError(w, 401, "token couldn't be validated")
		return
	}

	newAccessToken, err := auth.JwtRefreshTokenGenerator(subjectEmail, cfg.jwtSecret)
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
