//go:build e2e

package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"

	"github.com/stretchr/testify/assert"

	"github.com/go-resty/resty/v2"
)

func createJWTToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("supersecuresigningkey"))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func TestPostComment(t *testing.T) {
	t.Run("can post comment", func(t *testing.T) {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "bearer "+createJWTToken()).
			SetBody(`{"slug": "/", "author": "Kyle", "body": "Hello World"}`).
			Post("http://localhost:8080/api/v1/comment")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
	})
	t.Run("cannot post comment without JWT", func(t *testing.T) {
		client := resty.New()
		resp, err := client.R().SetBody(`{"slug": "/", "author": "Kyle", "body": "Hello World"}`).Post("http://localhost:8080/api/v1/comment")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}
