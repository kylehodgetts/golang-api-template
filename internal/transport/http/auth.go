package http

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
)

func JWTAuth(
	original func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Auth performed based on header sent in request
		authHeader, ok := r.Header["Authorization"]
		if !ok || authHeader == nil {
			log.Errorf("Authorization header not present")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		// Bearer token-string
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			log.Errorf("bearer token not present in Authorization header: %v", authHeaderParts)
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		// validate incoming token now we know that it's present
		if !validateToken(authHeaderParts[1]) {
			log.Errorf("token %s is invalid", authHeaderParts[1])
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		original(w, r)
	}
}

func validateToken(token string) bool {
	// this is the signing key that all tokens must be signed with
	// to pass validation
	mySigningKey := []byte("supersecuresigningkey")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}

		return mySigningKey, nil
	})

	return err == nil && parsedToken.Valid
}
