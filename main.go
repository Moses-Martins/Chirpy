package main

import (
	"os"  
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	_ "github.com/lib/pq"

	"database/sql"
    "github.com/joho/godotenv" 
	"github.com/Moses-Martins/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB *database.Queries
	Platform string
	JwtSecret string
}

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("SECRET")
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig {
		DB: dbQueries,
		fileserverHits: atomic.Int32{},
		Platform: platform,
		JwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.singleChirps)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUsers)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.sendChirps)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.UserUpdate)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirp)


	Server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err = Server.ListenAndServe()
	if err != nil {
		log.Fatal("There happens to be an error: ", err)
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}


