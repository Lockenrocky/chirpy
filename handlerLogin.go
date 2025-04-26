package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Lockenrocky/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decoder.Decode(&params)

	user, err := cfg.db.Login(r.Context(), params.Email)
	if err != nil {
		w.Header().Set("Context-Type", "plain/text")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email or password"))
		return
	}

	if auth.CheckPasswordHash(user.HashedPassword, params.Password) != nil {
		w.Header().Set("Context-Type", "plain/text")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect email or password"))
		return
	}

	type resp struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	loggedin_user := resp{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	dat, err := json.Marshal(loggedin_user)
	if err != nil {
		log.Fatal("Error marshaling json")
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
