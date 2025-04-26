package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type errorVals struct {
		Error string `json:"error"`
	}

	type validVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	// Error with JSON
	if err != nil {
		errBody := errorVals{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(errBody)
		if err != nil {
			log.Printf("Error marshalling JSON %s", err)
			w.WriteHeader(500)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	// Chirp is too long
	if len(params.Body) > 140 {
		errBody := errorVals{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errBody)
		if err != nil {
			log.Printf("Error marshalling JSON %s", err)
			w.WriteHeader(500)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	// Valid chirp

	// check for bad words
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned_body := checkProfanity(params.Body, badWords)

	validBody := validVals{
		Cleaned_body: cleaned_body,
	}
	dat, err := json.Marshal(validBody)
	if err != nil {
		log.Printf("Error marshalling JSON %s", err)
		w.WriteHeader(500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}

func checkProfanity(payload string, badWords map[string]struct{}) string {
	words := strings.Split(payload, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if _, ok := badWords[lower]; ok {
			words[i] = "****"
		}
	}
	cleaned_payload := strings.Join(words, " ")
	return cleaned_payload
}
