package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nabin3/a-simple-web-server/internal/auth"
	"github.com/nabin3/a-simple-web-server/internal/database"
)

// Defining handler func for "POST /api/chirps"
func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	// Struct for recieved data
	type data struct {
		Body string `json:"body"`
	}

	// Check if access token is valid

	// Getting token from request haeader
	unverifiedToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error in handlerUsersPut: %v", err)
		respondWithError(w, 401, "bad request")
		return
	}

	// Checking the type of token by exmining the issuer
	tokenIssuer, err := auth.RetrieveJwtTokenIssuer(unverifiedToken)
	if err != nil {
		log.Printf("error in handlerChirpsPost at RetrieveJwtTokenIssuer: %v", err)
		respondWithError(w, 500, "server-error")
		return
	}
	if tokenIssuer != "chirpy-access" {
		log.Printf("user error in handlerChirpsPost at token-issuer checker: user has not given a access token, issuer: %s", tokenIssuer)
		respondWithError(w, 401, "bad token")
		return
	}

	// Validating recieved token
	emailInToken, err := auth.ValidateJWT(unverifiedToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error in handlerChirpsPost: %v", err)
		respondWithError(w, 401, "token couldn't be validated")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Erro")
		return
	}

	// Getting a user from Users map of database
	_, err = db.GetUserByEmail(emailInToken)
	if err != nil {
		log.Printf("error in handlerUsersPut: %v", err)
		respondWithError(w, 401, "could't find any user for given token")
		return
	}
	// Code for validating token is ended here

	// Decoding recieved json
	recievedData := data{}
	decoder := json.NewDecoder(r.Body)
	err1 := decoder.Decode(&recievedData)
	if err1 != nil {
		log.Printf("error in handlerChirpsPost at decoder.Decode: %s", err1)
		w.WriteHeader(500)
		return
	}

	// Checking if recieved Chirp is of valid length or not
	if len(recievedData.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	} else {
		cleaned_Chirp := profaneWordsReplacor(recievedData.Body)
		chirp, err1 := db.CreateChirp(emailInToken, cleaned_Chirp)
		if err1 != nil {
			log.Printf("error: %v", err1)
			respondWithError(w, 500, "internal error")
			return
		}

		respondWithJSON(w, 201, chirp)
	}
}

// Defining handler func for "GET /api/chirps"
func handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 500, "internal error")
		return
	}

	all_chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 500, "internal error")
		return
	}

	respondWithJSON(w, 201, all_chirps)
}

// Defining handler function for "GET /api/chirps/{chirp_id}"
func handlerChirpGetByID(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 500, "internal error")
		return
	}

	id, err := strconv.Atoi(r.PathValue("chirp_id"))
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 404, "given id is non-integer")
		return
	}

	chirp, err := db.GetChirpByID(id)
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 404, fmt.Sprintf("error: %v", err))
		return
	}

	respondWithJSON(w, 200, chirp)
}
