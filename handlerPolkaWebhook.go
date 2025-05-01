package main

import (
	"encoding/json"
	"net/http"

	"github.com/Lockenrocky/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Could find ApiKey in header"))
		return
	}

	if apiKey != cfg.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong ApiKey"))
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			User_id uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.User_id)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Could not update user"))
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
