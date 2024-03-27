package main

import "strings"

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
