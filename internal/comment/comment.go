package comment

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFetchingComment = errors.New("failed to fetch comment by id")
	ErrNotImplemented  = errors.New("not implemented")
)

// Comment - a representation of the comment
// structure for the service
type Comment struct {
	ID     string
	Slug   string
	Body   string
	Author string
}

type Store interface {
	GetComment(context.Context, string) (Comment, error)
	PostComment(context.Context, Comment) (Comment, error)
	DeleteComment(context.Context, string) error
	UpdateComment(context.Context, string, Comment) (Comment, error)
}

// Service - struct on which all logic will be built on top of
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	fmt.Println("retrieving comment")
	comment, err := s.Store.GetComment(ctx, id)
	if err != nil {
		fmt.Println(err)
		return Comment{}, ErrFetchingComment
	}

	return comment, nil
}

func (s *Service) PutComment(ctx context.Context, id string, updatedComment Comment) (Comment, error) {
	cmt, err := s.Store.UpdateComment(ctx, id, updatedComment)
	if err != nil {
		return Comment{}, err
	}

	return cmt, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	if err := s.Store.DeleteComment(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *Service) PostComment(ctx context.Context, comment Comment) (Comment, error) {
	cmt, err := s.Store.PostComment(ctx, comment)
	if err != nil {
		return Comment{}, fmt.Errorf("could not post comment: %w", err)
	}
	return cmt, nil
}
