package main

import (
	"net/http"
	"encoding/json"
	"log"
)



func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {

	ChirpDb, err := cfg.DB.GetChirps(req.Context())
	if err != nil {
    	http.Error(w, "Cannot Retrieve Chirps", http.StatusNotFound)
        return
	}

	respBody := make([]Chirps, 0, len(ChirpDb))
    for _, dbChirp := range ChirpDb {
        chirpResp := Chirps{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        }
        respBody = append(respBody, chirpResp)
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