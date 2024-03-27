package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nabin3/a-simple-web-server/internal/auth"
	"github.com/nabin3/a-simple-web-server/internal/database"
)

// Defining handler func for "POST /api/login"
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Struct for sending data
	type resp struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Struct for recieved data
	type data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decoding recieved json
	recievedData := data{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recievedData)
	if err != nil {
		log.Printf("error decoding from JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	// Creating database connection
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

	// Checking if user exist in database
	user, err := db.GetUserByEmail(recievedData.Email)
	if err != nil {
		log.Printf("error occured in handlerLogin at GetUserByEmail, error: %v", err)
		respondWithError(w, 401, "internal error")
		return
	}

	// Checking if user is authentic
	success := auth.ValidatePassword(user.PasswordHash, recievedData.Password)
	if !success {
		respondWithError(w, 401, "invalid creditional")
		return
	}

	jwtToken, err := auth.JwtAccessTokenGenerator(recievedData.Email, cfg.jwtSecret)
	if err != nil {
		log.Printf("error happend in handlerLogin func for /api/login endpoint, error souce is JwtAccessTokenGenerator func, error: %v", err)
		respondWithError(w, 500, "server faced problem generating token")
		return
	}

	jwtRefreshToken, err := auth.JwtRefreshTokenGenerator(recievedData.Email, cfg.jwtSecret)
	if err != nil {
		log.Printf("error happend in handlerLogin func for /api/login endpoint, error souce is JwtRefreshTokenGenerator func, error: %v", err)
		respondWithError(w, 500, "server faced problem generating refresh token")
		return
	}
	respUserData := resp{
		Token:        jwtToken,
		RefreshToken: jwtRefreshToken,
	}
	respondWithJSON(w, 200, respUserData)
}
