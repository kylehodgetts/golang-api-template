package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service CommentService
	Router  *mux.Router
	Server  *http.Server
}

func NewHandler(service CommentService) *Handler {
	h := &Handler{
		Service: service,
		Router:  mux.NewRouter(),
	}

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "i am alive")
	})

	h.Router.Use(JSONMiddleware, LoggingMiddleware, TimeoutMiddleware)
	h.Router.HandleFunc("/api/v1/comment", JWTAuth(h.PostComment)).Methods(http.MethodPost)
	h.Router.HandleFunc("/api/v1/comment/{id}", h.GetComment).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/v1/comment/{id}", JWTAuth(h.PutComment)).Methods(http.MethodPut)
	h.Router.HandleFunc("/api/v1/comment/{id}", JWTAuth(h.DeleteComment)).Methods(http.MethodDelete)
}

func (h *Handler) Serve() error {
	// Graceful handling of running and shutting down server
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// give server 15 seconds to finish handling requests before shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	h.Server.Shutdown(ctx)
	log.Println("shut down gracefully")
	return nil
}
