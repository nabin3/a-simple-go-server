package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nabin3/a-simple-web-server/internal/auth"
	"github.com/nabin3/a-simple-web-server/internal/database"
)

// Defining handler func for "POST /api/users"
func handlerUsersPost(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("handlerUsersPost: error decoding from JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	// Creating database connection
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("handlerUsersPost_func: can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

	// Creating new user and getting user's unique_id and email id for giving response to client
	newUser, err := db.CreateUser(0, recievedData.Email, recievedData.Password)
	if err != nil {
		log.Printf("handleUsersPost_func: error from CreateUser_func: %v", err)
		respondWithError(w, 500, fmt.Sprintf("error: %v", err))
		return
	}

	respondWithJSON(w, 201, newUser)
}

// Defining handler for PUT /api/users/
func (cfg *apiConfig) handlerUsersPut(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("handlerUsersPost: error decoding from JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	// Creating database connection
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("handlerUsersPost_func: can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

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
		log.Printf("error in handlerUsersPut at RetrieveJwtTokenIssuer: %v", err)
		respondWithError(w, 500, "server-error")
		return
	}
	if tokenIssuer != "chirpy-access" {
		log.Printf("user error in handlerUsersPut at issuer checker: user has not given a access token, issuer: %s", tokenIssuer)
		respondWithError(w, 401, "bad token")
		return
	}

	// Validating recieved token
	emailInToken, err := auth.ValidateJWT(unverifiedToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("error in handlerUsersPut: %v", err)
		respondWithError(w, 401, "token couldn't be validated")
		return
	}

	// Getting a user from Users map of database
	user, err := db.GetUserByEmail(emailInToken)
	if err != nil {
		log.Printf("error in handlerUsersPut: %v", err)
		respondWithError(w, 401, "could't find any user for given token")
		return
	}

	// Deleting a user with given email id
	if err := db.DeleteUser(emailInToken); err != nil {
		log.Printf("error in handlerUsersPut: %v", err)
		respondWithError(w, 500, "server problem")
		return
	}

	// Creating new user
	newUser, err := db.CreateUser(user.ID, recievedData.Email, recievedData.Password)
	if err != nil {
		log.Printf("handleUsersPost_func: error from CreateUser_func: %v", err)
		respondWithError(w, 500, fmt.Sprintf("error: %v", err))
		return
	}

	respondWithJSON(w, 200, newUser)
}
