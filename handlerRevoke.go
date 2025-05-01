package main

import (
	"net/http"

	"github.com/Lockenrocky/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not find token"))
		return
	}

	_, err = cfg.db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not revoke Token!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
