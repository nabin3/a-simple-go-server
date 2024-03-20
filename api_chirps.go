package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/nabin3/a-simple-web-server/internal/database"
)

// Defining handler func for "POST /api/chirps"
func handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	// Struct for recieved data
	type data struct {
		Body string `json:"body"`
	}

	// Decoding recieved json
	recievedData := data{}
	decoder := json.NewDecoder(r.Body)
	err1 := decoder.Decode(&recievedData)
	if err1 != nil {
		log.Printf("error unmarshaling from JSON: %s", err1)
		w.WriteHeader(500)
		return
	}

	// Checking if recieved Chirp is of valid length or not
	if len(recievedData.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	} else {
		db, err := database.NewDB("./database.json")
		if err != nil {
			log.Printf("can't create connection to database, error: %v", err)
			respondWithError(w, 500, "Internal Erro")
			return
		}

		cleaned_Chirp := profaneWordsReplacor(recievedData.Body)
		chirp, err1 := db.CreateChirp(cleaned_Chirp)
		if err1 != nil {
			log.Printf("error: %v", err1)
			respondWithError(w, 500, "internal error")
			return
		}

		respondWithJSON(w, 201, chirp)
	}

}

// profane word filter
func profaneWordsReplacor(msg string) string {
	// Here considered empty struct rather than bool because of memory efficiency
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// Creating slice of string by spliting the msg styring
	msgSlice := strings.Split(msg, " ")

	// Checking if a word from the slice of string(msgSlice) is a badword, if it is then reolacing it with ****
	for index, word := range msgSlice {
		if _, exist := profaneWords[strings.ToLower(word)]; exist {
			msgSlice[index] = "****"
		}
	}

	return strings.Join(msgSlice, " ")
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

	// Sorting all_chirps slice ny ID
	slices.SortStableFunc(all_chirps, func(a, b database.Chirp) int {
		return cmp.Compare(a.ID, b.ID)
	})

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
