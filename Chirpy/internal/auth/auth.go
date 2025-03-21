package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPass), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	nowTime := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(nowTime),
		ExpiresAt: jwt.NewNumericDate(nowTime.Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Parsed the wrong signing method\n")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("Failed to parse with Claims\n")
		return uuid.UUID{}, err
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			log.Printf("Failed to extract claims\n")
			return uuid.UUID{}, errors.New("invalid user ID format in token")
		}
		return userID, nil
	}
	log.Printf("Invalid claims\n")
	return uuid.UUID{}, errors.New("invalid token claims")
}

func GetBearerToken(headers http.Header) (string, error) {
	return GetAuthorizationToken(headers, "Bearer ")
}

func GetAPIKey(headers http.Header) (string, error) {
	return GetAuthorizationToken(headers, "ApiKey ")
}

func GetAuthorizationToken(headers http.Header, prefix string) (string, error) {
	// Retrieve the authorization Header
	authoriz := headers.Get("Authorization")
	if authoriz == "" {
		return "", errors.New("missing authorization header")
	}

	// Removes the prefix
	if !strings.HasPrefix(authoriz, prefix) {
		return "", errors.New("invalid Authorization header format")
	}

	// Extract the token string
	tokenString := strings.TrimSpace(strings.TrimPrefix(authoriz, prefix))
	if tokenString == "" {
		return "", errors.New("empty token string")
	}
	return tokenString, nil
}
