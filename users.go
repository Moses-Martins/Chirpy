package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/Moses-Martins/Chirpy/internal/auth"
	"github.com/Moses-Martins/Chirpy/internal/database"
)

type AcceptEmail struct {
	Password string `json:"password"`
	Email string `json:"email"`
}
 
type UserShown struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *apiConfig) CreateUsers(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := AcceptEmail{}
	
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
 
	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error Hashing Password: %s", err)
		w.WriteHeader(500)
		return
	}

	userDb, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: params.Password,
	})
	if err != nil {
    	http.Error(w, "Cannot Create User", http.StatusNotFound)
        return
	}

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
	w.WriteHeader(201)
	w.Write(data)
	
}