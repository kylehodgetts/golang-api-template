package main

import (
	"fmt"
	"os"

	"github.com/kylehodgetts/go-rest-api-v2/internal/comment"

	"github.com/kylehodgetts/go-rest-api-v2/internal/db"

	transportHttp "github.com/kylehodgetts/go-rest-api-v2/internal/transport/http"
)

// Run - Reponsible for instantiation and startup of application
func Run() error {
	fmt.Println("application start")
	database, err := db.NewDatabase(db.DatabaseConnectionArgs{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_DB"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("SSL_MODE"),
	})
	if err != nil {
		fmt.Println("Failed to connect to database")
		return err
	}

	if err := database.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	commentService := comment.NewService(database)
	httpHandler := transportHttp.NewHandler(commentService)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("GO REST API")
}
