package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Lockenrocky/chirpy/internal/auth"
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
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("Something went wrong")
		return
	}

	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Error getting bearer Token"))
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.secret)
	if err != nil {
		w.WriteHeader(401)
		fmt.Println(err)
		w.Write([]byte("Error getting user_id"))
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body, UserID: userID})
	if err != nil {
		log.Fatalf("Something went wrong %s", err)
		w.WriteHeader(400)
		w.Write([]byte("Error getting chirp"))
		return
	}

	created_chirp := resp{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
	}

	dat, err := json.Marshal(created_chirp)
	if err != nil {
		log.Fatal("Error marshalling JSON")
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	author_id := r.URL.Query().Get("author_id")
	sortingOrder := r.URL.Query().Get("sort")

	var chirps = make([]database.Chirp, 0)

	if author_id != "" {
		user_id, err := uuid.Parse(author_id)
		if err != nil {
			log.Fatal(err)
		}
		chirps, err = cfg.db.SelectAllChirpsFromAuthor(r.Context(), user_id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		chirps, _ = cfg.db.SelectAllChirps(r.Context())
	}

	if sortingOrder == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	} else if sortingOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}

	allChirps := []resp{}
	for i := range chirps {
		createdChirp := resp{
			ID:         chirps[i].ID,
			Created_at: chirps[i].CreatedAt,
			Updated_at: chirps[i].UpdatedAt,
			Body:       chirps[i].Body,
			User_id:    chirps[i].UserID,
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Chirp not found!"))
		return
	}

	selectedChirp := resp{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
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

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Could not find access token"))
		return
	}

	user_ID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Could not find user id"))
		return
	}

	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Could not get chirp id"))
		return
	}

	chirp, err := cfg.db.SelectChirp(r.Context(), chirp_id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Chirp not found"))
		return
	}

	if chirp.UserID != user_ID {
		w.WriteHeader(403)
		w.Write([]byte("You dont own the chirp"))
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirp_id)
	if err != nil {
		w.WriteHeader(403)
		w.Write([]byte("Could not delete chirp"))
		return
	}

	w.WriteHeader(204)

}
