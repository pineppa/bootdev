package src

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type JsonBody struct {
	Body        string    `json:"body"`
	CleanedBody string    `json:"cleaned_body"`
	Valid       bool      `json:"valid"`
	UserID      uuid.UUID `json:"user_id"`
}

type JsonChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func fixString(str string) string {
	strJoin := strings.Fields(str)
	for i, elem := range strJoin {
		strElem := strings.ToLower(elem)
		if strings.Compare(strElem, "kerfuffle") == 0 || strings.Compare(strElem, "sharbert") == 0 || strings.Compare(strElem, "fornax") == 0 {
			strJoin[i] = "****"
		}
	}
	return strings.Join(strJoin, " ")
}

func (apiCfg *ApiConfig) HandlePostChirps(w http.ResponseWriter, r *http.Request) {

	var jsonContent JsonBody
	if err := ReadBodyToJson(r.Body, &jsonContent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	dbUser, status, err := apiCfg.GetValidatedUser(r.Header)
	if status != http.StatusOK {
		log.Println("Failed to validate the user: ", err)
	}

	if len(jsonContent.Body) > 140 {
		log.Println("Chirp content too long")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	jsonContent.CleanedBody = fixString(jsonContent.Body)

	chirp, err := apiCfg.DbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   jsonContent.Body,
		UserID: dbUser.ID,
	})
	if err != nil {
		log.Printf("Failed to create the chirp: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var jsonChirp = JsonChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	ReturnJsonWithBody(w, jsonChirp, http.StatusCreated, http.StatusBadRequest)
}

func (apiCfg *ApiConfig) HandleSingleChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDVal := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDVal)
	if err != nil {
		log.Println("Failed to parse Chirp")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := apiCfg.DbQueries.GetChirpById(r.Context(), chirpID)
	if err != nil {
		log.Println("Failed to retrieve Chirp")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonChirp := JsonChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	ReturnJsonWithBody(w, jsonChirp, http.StatusOK, http.StatusBadRequest)
}

func (apiCfg *ApiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {

	str, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed retrieving the Bearer token: %v\n", err)
		if strings.HasPrefix(err.Error(), "missing") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	uu, err := auth.ValidateJWT(str, os.Getenv("SECRETS"))
	if err != nil {
		log.Printf("JWT does not match: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpIDVal := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDVal)
	if err != nil {
		log.Println("Failed to parse Chirp")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := apiCfg.DbQueries.DeleteChirp(context.Background(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: uu,
	}); err != nil {
		log.Println("Chirp not found: ", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (apiCfg *ApiConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {

	authorId := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")
	sortParam = strings.ToLower(sortParam)

	if sortParam != "asc" && sortParam != "desc" {
		log.Println("sort parameter wrongly passed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirps, err := apiCfg.getChirps(authorId)
	if err != nil {
		log.Printf("Failed to retrieve all chirps: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortParam == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	var arrChirps = make([]JsonChirp, len(chirps))
	for i, ch := range chirps {
		arrChirps[i] = JsonChirp{
			ID:        ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body:      ch.Body,
			UserID:    ch.UserID,
		}
	}

	ReturnJsonWithBody(w, arrChirps, http.StatusOK, http.StatusBadRequest)
}

func (apiCfg *ApiConfig) getChirps(authorId string) ([]database.Chirp, error) {

	if authorId == "" {
		return apiCfg.DbQueries.GetAllGetChirps(context.Background())
	}
	uu, err := uuid.Parse(authorId)
	if err != nil {
		log.Println("Failed to parse the author idea: ", err)
	}
	return apiCfg.DbQueries.GetAllGetChirpsByUser(context.Background(), uu)
}
