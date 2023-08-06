package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"

	"github.com/kylehodgetts/go-rest-api-v2/internal/comment"
)

type CommentService interface {
	PostComment(context.Context, comment.Comment) (comment.Comment, error)
	GetComment(context.Context, string) (comment.Comment, error)
	PutComment(context.Context, string, comment.Comment) (comment.Comment, error)
	DeleteComment(context.Context, string) error
}

type PostCommentRequest struct {
	Slug   string `json:"slug" validate:"required"`
	Author string `json:"author" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func convertPostCommentRequestToComment(request PostCommentRequest) comment.Comment {
	return comment.Comment{
		Slug:   request.Slug,
		Author: request.Author,
		Body:   request.Body,
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var postCommentRequest PostCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&postCommentRequest); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(postCommentRequest); err != nil {
		http.Error(w, "not a valid comment", http.StatusBadRequest)
		return
	}

	postedComment, err := h.Service.PostComment(r.Context(), convertPostCommentRequestToComment(postCommentRequest))
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(postedComment); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentId, ok := vars["id"]
	if !ok || commentId == "" {
		log.Print("no id provided in request")
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	// TODO: differentiate between error and comment not found
	c, err := h.Service.GetComment(r.Context(), commentId)
	if err != nil {
		log.Print(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(c); err != nil {
		log.Print(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PutComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentId, ok := vars["id"]
	if !ok || commentId == "" {
		log.Print("no id provided in request")
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	var (
		c   comment.Comment
		err error
	)
	if err = json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err = h.Service.PutComment(r.Context(), commentId, c)
	if err != nil {
		log.Print(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(c); err != nil {
		log.Print(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentId, ok := vars["id"]
	if !ok || commentId == "" {
		log.Print("no id provided in request")
		http.Error(w, "no id provided", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteComment(r.Context(), commentId); err != nil {
		log.Print(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
