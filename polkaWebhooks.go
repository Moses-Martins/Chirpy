package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"github.com/Moses-Martins/Chirpy/internal/auth"
)

type WebhookStruct struct {
    Event string `json:"event"`
  	Data struct {
    	UserID uuid.UUID `json:"user_id"`
  	} `json:"data"`
}

func (cfg *apiConfig) PolkaWebhooks(w http.ResponseWriter, req *http.Request) {
	
	token_string, err := auth.GetAPIKey(req.Header)
    if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("missing authorization header"))
		return
	}

	if token_string != cfg.PolkaKey {
		w.WriteHeader(401)
		return
	}


	decoder := json.NewDecoder(req.Body)
	params := WebhookStruct{}
	
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	
	_, err = cfg.DB.UpgradeUserToChirpyRed(req.Context(), params.Data.UserID) 
	if err != nil {
		log.Printf("Error upgrading user: %s", err)
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)





}