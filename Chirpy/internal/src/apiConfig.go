package src

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	Secrets        string
	Polka          string
}

type JsonUser struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"password"`
	TokenString    string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

type JsonToken struct {
	Token string `json:"token"`
}

func (apiCfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *ApiConfig) ReqHitsMetrics(w http.ResponseWriter, r *http.Request) {

	visitedHtml, err := template.ParseFiles("metrics.html")
	if err != nil {
		return
	}
	visitCount := struct {
		VisitCount int32
	}{
		VisitCount: apiCfg.fileserverHits.Load(),
	}

	w.Header().Set("Content-Type", "text/html")
	visitedHtml.Execute(w, visitCount)
}

func (apiCfg *ApiConfig) HandleResetDatabase(w http.ResponseWriter, r *http.Request) {
	if apiCfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
	}
	apiCfg.fileserverHits.Store(0)
	apiCfg.DbQueries.DeleteAllUsers(r.Context())
	apiCfg.DbQueries.DeleteAllChirps(r.Context())
	apiCfg.DbQueries.DeleteAllTokens(r.Context())
}

func (apiCfg *ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error retreiving the Bearer string token: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refrToken, err := apiCfg.DbQueries.GetTokenByToken(context.Background(), tokenString)

	if err != nil {
		log.Printf("Error refreshing token: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if refrToken.ExpiredAt.Valid {
		log.Println("The token has already expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	JWTTokenString, err := auth.MakeJWT(refrToken.UserID, apiCfg.Secrets, time.Hour)
	if err != nil {
		log.Printf("Error retreiving a JWT token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := apiCfg.DbQueries.SetTokenByID(context.Background(), database.SetTokenByIDParams{
		ID:    refrToken.UserID,
		Token: sql.NullString{String: JWTTokenString, Valid: true},
	}); err != nil {
		log.Printf("Error retreiving a JWT token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ReturnJsonWithBody(w, JsonToken{
		Token: JWTTokenString,
	}, http.StatusOK, http.StatusBadRequest)
}

func (apiCfg *ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("Error retrieving the Bearer Token: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := apiCfg.DbQueries.RevokeToken(r.Context(), tokenString); err != nil {
		log.Println("Error revoking the Token: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent)
}
