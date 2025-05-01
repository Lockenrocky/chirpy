package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lockenrocky/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not find token"))
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Could not get user from refresh token"))
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		time.Hour,
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Could not validate token"))
		return
	}

	type resp struct {
		Token string `json:"token"`
	}
	newResponse := resp{
		Token: accessToken,
	}
	dat, err := json.Marshal(newResponse)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
