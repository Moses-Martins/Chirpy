package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"strings"
	"github.com/google/uuid"
	"github.com/Moses-Martins/Chirpy/internal/database"
	"github.com/Moses-Martins/Chirpy/internal/auth"
)


type chirp struct {
	Body string `json:"body"`
}

type Chirps struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) sendChirps(w http.ResponseWriter, req *http.Request) {

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

	decoder := json.NewDecoder(req.Body)
	params := chirp{}
	
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		http.Error(w, "Chirp too long", http.StatusNotFound)
		return
	}

	params.Body = filterChirp(params.Body)

	ChirpDb, err := cfg.DB.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: ValidatedID,
	})
	if err != nil {
    	http.Error(w, "Cannot Create Chirp", http.StatusNotFound)
        return
	}

	respBody := Chirps{
		ID:       ChirpDb.ID,
		CreatedAt: ChirpDb.CreatedAt,
		UpdatedAt: ChirpDb.UpdatedAt,
		Body:    ChirpDb.Body,
		UserID: ChirpDb.UserID,
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


func filterChirp(chirp string) string {
    bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

    words := strings.Split(chirp, " ")
    for i, word := range words {
        for _, banned := range bannedWords {
            // Compare case-insensitively
            if strings.EqualFold(word, banned) {
                words[i] = "****"
            }
        }
    }

    return strings.Join(words, " ")
}

