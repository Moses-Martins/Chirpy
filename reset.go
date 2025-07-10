package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) ResetUsers(w http.ResponseWriter, req *http.Request) {
	if cfg.Platform != "dev" {
        http.Error(w, "Forbidden: not allowed in this environment", http.StatusForbidden)
        return
    }

	err := cfg.DB.DeleteAllUsers(req.Context())
    if err != nil {
        log.Printf("Error deleting users: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted"))
}