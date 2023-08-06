//go:build integration

package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kylehodgetts/go-rest-api-v2/internal/comment"
)

func TestCommentDatabase(t *testing.T) {
	connectionArgs := DatabaseConnectionArgs{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_DB"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}

	t.Run("test create and get comment", func(t *testing.T) {
		fmt.Println("testing the creation of comments")
		db, err := NewDatabase(connectionArgs)
		assert.NoError(t, err)
		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "slug",
			Author: "author",
			Body:   "Body",
		})
		assert.NoError(t, err)

		newCmt, err := db.GetComment(context.Background(), cmt.ID)
		assert.NoError(t, err)
		assert.Equal(t, cmt, newCmt)
	})

	t.Run("test delete comment", func(t *testing.T) {
		db, err := NewDatabase(connectionArgs)
		assert.NoError(t, err)
		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "new-slug",
			Author: "new author",
			Body:   "body",
		})
		assert.NoError(t, err)

		err = db.DeleteComment(context.Background(), cmt.ID)
		assert.NoError(t, err)

		_, err = db.GetComment(context.Background(), cmt.ID)
		assert.Error(t, err)
	})

	t.Run("test update comment", func(t *testing.T) {
		db, err := NewDatabase(connectionArgs)
		assert.NoError(t, err)
		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "new-slug",
			Author: "new author",
			Body:   "body",
		})
		assert.NoError(t, err)

		updatedCmt, err := db.UpdateComment(context.Background(), cmt.ID, comment.Comment{
			Slug:   "new-slug",
			Author: "new author",
			Body:   "new body",
		})
		assert.NoError(t, err)

		queryUpdatedComment, err := db.GetComment(context.Background(), updatedCmt.ID)
		assert.NoError(t, err)
		assert.Equal(t, updatedCmt.Body, queryUpdatedComment.Body)
	})
}
