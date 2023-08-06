package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/kylehodgetts/go-rest-api-v2/internal/comment"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(commentRow CommentRow) comment.Comment {
	return comment.Comment{
		ID:     commentRow.ID,
		Slug:   commentRow.Slug.String,
		Body:   commentRow.Body.String,
		Author: commentRow.Author.String,
	}
}

func (d *Database) GetComment(ctx context.Context, uuid string) (comment.Comment, error) {
	getCommentQuery := `SELECT * FROM comments WHERE ID = $1`

	commentRow := CommentRow{}
	if err := d.Client.GetContext(ctx, &commentRow, getCommentQuery, uuid); err != nil {
		return comment.Comment{}, err
	}

	return convertCommentRowToComment(commentRow), nil
}

func (d *Database) PostComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	commentID, err := uuid.NewV4()
	if err != nil {
		return comment.Comment{}, fmt.Errorf("unable to generate new uuid for comment: %w", err)
	}

	cmt.ID = commentID.String()
	postRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO comments (id, slug, author, body) VALUES (:id, :slug, :author, :body)`,
		postRow,
	)

	if err != nil {
		return comment.Comment{}, fmt.Errorf("could not insert new comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return cmt, nil
}

func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("unable to delete comment with id %s: %w", id, err)
	}
	return nil
}

func (d *Database) UpdateComment(ctx context.Context, id string, cmt comment.Comment) (comment.Comment, error) {
	updateRow := CommentRow{
		ID:     id,
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(ctx, `UPDATE comments SET author = :author, slug=:slug, body=:body WHERE id = :id`, updateRow)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("could not update comment with id %s: %w", id, err)
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertCommentRowToComment(updateRow), nil
}
