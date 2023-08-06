package comment

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	errToReturn error
}

func (m mockStore) GetComment(_ context.Context, id string) (Comment, error) {
	if m.errToReturn != nil {
		return Comment{}, m.errToReturn
	}

	return Comment{
		Slug:   "some-slug",
		ID:     id,
		Author: "Some Author",
		Body:   "some body",
	}, nil
}

func (m mockStore) UpdateComment(_ context.Context, _ string, updatedComment Comment) (Comment, error) {
	if m.errToReturn != nil {
		return Comment{}, m.errToReturn
	}

	return updatedComment, nil
}

func (m mockStore) DeleteComment(_ context.Context, _ string) error {
	return m.errToReturn
}

func (m mockStore) PostComment(_ context.Context, comment Comment) (Comment, error) {
	if m.errToReturn != nil {
		return Comment{}, m.errToReturn
	}
	return comment, nil
}

func TestService_GetComment(t *testing.T) {
	t.Run("can get existing comment", func(t *testing.T) {
		svc := NewService(mockStore{})
		cmt, err := svc.GetComment(context.Background(), "1234")
		assert.NoError(t, err)
		assert.Equal(t, "1234", cmt.ID)
	})

	t.Run("cannot get non-existent comment", func(t *testing.T) {
		svc := NewService(mockStore{errToReturn: errors.New("comment not found")})
		_, err := svc.GetComment(context.Background(), "1234")
		assert.NotNil(t, err)
	})
}
