package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/Moses-Martins/Chirpy/internal/auth"
	"github.com/Moses-Martins/Chirpy/internal/database"
)

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(403)
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(403)
		return
	}

	_, err = cfg.DB.DeleteChirpByUserID(req.Context(), database.DeleteChirpByUserIDParams{
		ID: id,
		UserID: ValidatedID,
	})

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)


}