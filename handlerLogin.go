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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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

	jwtToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		return
	}

	refToken, _ := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not save refresh token"))
	}

	type resp struct {
		ID            uuid.UUID `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Token         string    `json:"token"`
		Refresh_token string    `json:"refresh_token"`
		IsChirpyRed   bool      `json:"is_chirpy_red"`
	}

	loggedin_user := resp{
		ID:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		Token:         jwtToken,
		Refresh_token: refToken,
		IsChirpyRed:   user.IsChirpyRed,
	}

	dat, err := json.Marshal(loggedin_user)
	if err != nil {
		log.Fatal("Error marshaling json")
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
