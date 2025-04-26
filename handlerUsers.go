package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Lockenrocky/chirpy/internal/auth"
	"github.com/Lockenrocky/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type resp struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("Something went wrong")
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Fatal("Error hashing password")
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashed_password})
	if err != nil {
		log.Fatalf("Someting went wrong! %s", err)
		return
	}

	created_user := resp{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	dat, err := json.Marshal(created_user)
	if err != nil {
		log.Fatalf("Error marshalling JSON %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}
