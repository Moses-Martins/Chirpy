package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/google/uuid"
)



func (cfg *apiConfig) singleChirps(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid chirpID", http.StatusBadRequest)
		return
	}

	ChirpDb, err := cfg.DB.GetChirps(req.Context())
	if err != nil {
    	http.Error(w, "Cannot Retrieve Chirps", http.StatusNotFound)
        return
	}

	dbToStruct := make([]Chirps, 0, len(ChirpDb))
    for _, dbChirp := range ChirpDb {
        chirpResp := Chirps{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        }
        dbToStruct = append(dbToStruct, chirpResp)
    }


	respBody, exist := findChirpByID(dbToStruct, id)

	if exist != true {
		w.WriteHeader(404)
		return
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


func findChirpByID(chirps []Chirps, id uuid.UUID) (Chirps, bool) {
    for _, chirp := range chirps {
        if chirp.ID == id {
            return chirp, true
        }
    }
    return Chirps{}, false
}
