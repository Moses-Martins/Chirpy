package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/Moses-Martins/Chirpy/internal/auth"
	"github.com/Moses-Martins/Chirpy/internal/database"
)


type NewPassEmail struct {
	Password string `json:"password"`
	Email string `json:"email"`
}


func (cfg *apiConfig) UserUpdate(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := NewPassEmail{}
	
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error Hashing Password: %s", err)
		w.WriteHeader(500)
		return
	}


	cfg.DB.UpdateEmailAndPassword(req.Context(), database.UpdateEmailAndPasswordParams{
		Email: params.Email,
		HashedPassword: params.Password,
		ID: ValidatedID,
	}) 


	userDb, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)

	respBody := UserShown{
		ID:       userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:    userDb.Email,
		IsChirpyRed: userDb.IsChirpyRed,
	}

	data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)





}

