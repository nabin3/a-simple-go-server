package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		log.Printf("error decoding from JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

	newUser, err := db.CreateUser(recievedData.Email, recievedData.Password)
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 500, fmt.Sprintf("error: %v", err))
		return
	}

	respondWithJSON(w, 201, newUser)
}

// Defining handler func for "POST /api/login"
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Struct for sending data
	type resp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
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

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Printf("can't create connection to database, error: %v", err)
		respondWithError(w, 500, "Internal Error")
		return
	}

	userId, err := db.CheckUserCreditional(recievedData.Email, recievedData.Password)
	if err != nil {
		log.Printf("error: %v", err)
		respondWithError(w, 500, fmt.Sprintf("error: %v", err))
		return
	}

	if userId == 0 {
		respondWithError(w, 401, "Invalid Creditionals")
	} else {
		respUserData := resp{
			ID:    userId,
			Email: recievedData.Email,
		}
		respondWithJSON(w, 200, respUserData)
	}
}
