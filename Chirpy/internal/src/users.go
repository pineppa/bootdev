package src

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JsonWebHook struct {
	Event string   `json:"event"`
	Data  JsonData `json:"data"`
}
type JsonData struct {
	UserID string `json:"user_id"`
}

func (apiCfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user JsonUser
	if err := ReadBodyToJson(r.Body, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Generates the hashed password for the user
	hashedPass, err := auth.HashPassword(user.HashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Creates a new user with the passed Email and Password information
	dbUser, err := apiCfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: user.Email,
		HashedPassword: sql.NullString{
			String: hashedPass,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Error creating the user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := apiCfg.SetToken(dbUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	refreshTokenString, err := apiCfg.SetRefreshToken(dbUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var jsonResp = JsonUser{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		TokenString:  tokenString,
		RefreshToken: refreshTokenString,
		IsChirpyRed:  false,
	}

	// Give a code 201 to inform that the User has been created
	log.Println("A new user has been created: ", dbUser.Email)
	ReturnJsonWithBody(w, jsonResp, http.StatusCreated, http.StatusInternalServerError)
}

func (apiCfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {

	var user JsonUser
	if err := ReadBodyToJson(r.Body, &user); err != nil {
		log.Println("Error parsing information: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	dbUser, err := apiCfg.DbQueries.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := auth.CheckPasswordHash(user.HashedPassword, dbUser.HashedPassword.String); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var jsonResp = JsonUser{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		TokenString:  dbUser.Token.String,
		RefreshToken: dbUser.RefreshToken.String,
		IsChirpyRed:  dbUser.IsChirpyRed,
	}

	// Give a code 201 to inform that the User has been created
	log.Println("The user has been logged in: ", dbUser.Email)
	ReturnJsonWithBody(w, jsonResp, http.StatusOK, http.StatusInternalServerError)
}

func (apiCfg *ApiConfig) HandlePutUsers(w http.ResponseWriter, r *http.Request) {
	// Command used to update the email or password information
	var jsonUser JsonUser
	if err := ReadBodyToJson(r.Body, &jsonUser); err != nil {
		log.Println("Failed to read the incoming content: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbUser, status, err := apiCfg.GetValidatedUser(r.Header)
	if status != http.StatusOK {
		log.Println("Failed to validate the user: ", err)
		w.WriteHeader(status)
		return
	}

	hashedPassword, err := auth.HashPassword(jsonUser.HashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiCfg.DbQueries.UpdateEmailPassword(context.Background(), database.UpdateEmailPasswordParams{
		Email:          jsonUser.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
		Token:          dbUser.Token,
	})

	ReturnJsonWithBody(w, jsonUser, http.StatusOK, http.StatusBadRequest)
}

func (apiCfg *ApiConfig) GetValidatedUser(header http.Header) (database.User, int, error) {
	str, err := auth.GetBearerToken(header)
	if err != nil {
		log.Printf("Failed retrieving the Bearer token: %v\n", err)
		return database.User{}, http.StatusUnauthorized, err
	}
	uu, err := auth.ValidateJWT(str, os.Getenv("SECRETS"))
	if err != nil {
		log.Printf("JWT does not match: %v\n", err)
		return database.User{}, http.StatusUnauthorized, err
	}

	dbUser, err := apiCfg.DbQueries.GetUserByToken(context.Background(), sql.NullString{String: str, Valid: true})
	if err != nil || uu != dbUser.ID {
		log.Printf("Failed to retreve the user: %v\n", err)
		return database.User{}, http.StatusUnauthorized, err
	}
	return dbUser, http.StatusOK, nil
}

func (apiCfg *ApiConfig) HandleRedChirp(w http.ResponseWriter, r *http.Request) {

	authKey, err := auth.GetAPIKey(r.Header)
	if err != nil || authKey != apiCfg.Polka {
		log.Println("Failed to match the API key: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var jsonWebHook JsonWebHook
	if err := ReadBodyToJson(r.Body, &jsonWebHook); err != nil {
		log.Println("Failed to read the incoming content: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check the actual request
	if jsonWebHook.Event != "user.upgraded" {
		log.Println("Wrong event content")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Get user UUID from string
	uu, err := uuid.Parse(jsonWebHook.Data.UserID)
	if err != nil {
		log.Println("Failed user to parse UUID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update status
	if err := apiCfg.DbQueries.UpdateRedStatus(context.Background(), database.UpdateRedStatusParams{
		IsChirpyRed: true,
		ID:          uu,
	}); err != nil {
		log.Println("Failed to find the user")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
