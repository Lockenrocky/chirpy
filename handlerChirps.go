package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Lockenrocky/chirpy/internal/database"
	"github.com/google/uuid"
)

type resp struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_id    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string        `json:"body"`
		User_id uuid.NullUUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("Something went wrong")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body, UserID: params.User_id})
	if err != nil {
		log.Fatalf("Something went wrong %s", err)
		return
	}

	created_chirp := resp{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID.UUID,
	}

	dat, err := json.Marshal(created_chirp)
	if err != nil {
		log.Fatal("Error marshalling JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.SelectAllChirps(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	allChirps := []resp{}
	for i := range chirps {
		createdChirp := resp{
			ID:         chirps[i].ID,
			Created_at: chirps[i].CreatedAt,
			Updated_at: chirps[i].UpdatedAt,
			Body:       chirps[i].Body,
			User_id:    chirps[i].UserID.UUID,
		}
		allChirps = append(allChirps, createdChirp)
	}

	dat, err := json.Marshal(allChirps)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Fatal()
	}

	chirp, err := cfg.db.SelectChirp(r.Context(), id)
	if err != nil {
		w.Header().Set("Context-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Chirp not found!"))
		log.Fatal(err)
	}

	selectedChirp := resp{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID.UUID,
	}

	dat, err := json.Marshal(selectedChirp)
	if err != nil {
		w.Header().Set("Context-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Chirp not found!"))
		log.Fatal("Somethin went wrong marshalling json")
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
