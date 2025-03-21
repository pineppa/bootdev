package src

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func ReadBodyToJson[T any](body io.ReadCloser, jsonStruct *T) error {
	// Reads the body content
	content, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to load the request body: %w", err)
	}

	// Makes the body conent readable in the expected JSON Structure
	json.Unmarshal(content, &jsonStruct)
	return nil
}

func ReturnJsonWithBody[T any](w http.ResponseWriter, jsonStruct T, succesState int, failureState int) {
	writeBody, err := json.Marshal(jsonStruct)
	if err != nil {
		w.WriteHeader(failureState)
		return
	}
	w.WriteHeader(succesState)
	if succesState != http.StatusNoContent {
		w.Write(writeBody)
	}
}

func (apiCfg *ApiConfig) SetToken(dbUser database.User) (string, error) {
	// Creates an activation JWT token valid for an hour
	tokenString, err := auth.MakeJWT(dbUser.ID, apiCfg.Secrets, time.Hour)
	if err != nil {
		fmt.Printf("Error retreiving a JWT token: %v\n", err)
		return "", err
	}

	apiCfg.DbQueries.CreateToken(context.Background(), database.CreateTokenParams{
		Token:  tokenString,
		UserID: dbUser.ID,
	})

	// Assigns the JWT token to the newly created user
	if err := apiCfg.DbQueries.SetTokenByID(context.Background(), database.SetTokenByIDParams{
		ID:    dbUser.ID,
		Token: sql.NullString{String: tokenString, Valid: true},
	}); err != nil {
		fmt.Printf("Error retreiving a JWT token: %v\n", err)
		return "", err
	}
	return tokenString, nil
}

func (apiCfg *ApiConfig) SetRefreshToken(dbUser database.User) (string, error) {
	// Generate the refreshing token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Printf("Error retreiving a refreshing token: %v\n", err)
		return "", err
	}

	apiCfg.DbQueries.CreateToken(context.Background(), database.CreateTokenParams{
		Token:  refreshTokenString,
		UserID: dbUser.ID,
	})

	// Assigns the refresh token to the newly created user
	if err := apiCfg.DbQueries.SetRefreshTokenByID(context.Background(), database.SetRefreshTokenByIDParams{
		ID:           dbUser.ID,
		RefreshToken: sql.NullString{String: refreshTokenString, Valid: true},
	}); err != nil {
		fmt.Printf("Error retreiving a refresh token: %v\n", err)
		return "", err
	}
	return refreshTokenString, nil
}
