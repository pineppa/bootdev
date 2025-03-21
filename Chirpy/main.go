package main

import (
	"chirpy/internal/database"
	"chirpy/internal/src"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		fmt.Printf("Failed to load the environment: %v\n", err)
		return
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to load the database: %v\n", err)
		return
	}

	var smux = http.NewServeMux()
	var apiCfg = src.ApiConfig{}

	apiCfg.DbQueries = database.New(db)
	apiCfg.Platform = os.Getenv("PLATFORM")
	apiCfg.Secrets = os.Getenv("SECRETS")
	apiCfg.Polka = os.Getenv("POLKA_KEY")

	InitHandlers(smux, &apiCfg)

	s := InitServer(smux)
	s.ListenAndServe()
}

func InitServer(smux *http.ServeMux) http.Server {
	return http.Server{
		Handler: smux,
		Addr:    ":8080",
	}
}

func InitHandlers(smux *http.ServeMux, apiCfg *src.ApiConfig) {
	handler := http.FileServer(http.Dir("."))
	handler = http.StripPrefix("/app", apiCfg.MiddlewareMetricsInc(handler))

	smux.Handle("/app/", handler)
	smux.HandleFunc("GET /api/healthz", MyHandler)
	smux.HandleFunc("GET /admin/metrics", apiCfg.ReqHitsMetrics)
	smux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	smux.HandleFunc("POST /admin/reset", apiCfg.HandleResetDatabase)
	smux.HandleFunc("POST /api/chirps", apiCfg.HandlePostChirps)
	smux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	smux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleSingleChirp)
	smux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandleDeleteChirp)
	smux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	smux.HandleFunc("PUT /api/users", apiCfg.HandlePutUsers)
	smux.HandleFunc("POST /api/refresh", apiCfg.HandleRefresh)
	smux.HandleFunc("POST /api/revoke", apiCfg.HandleRevoke)
	smux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandleRedChirp)
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
