package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	Client *sqlx.DB
}

type DatabaseConnectionArgs struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(args DatabaseConnectionArgs) (*Database, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		args.Host,
		args.Port,
		args.Username,
		args.DBName,
		args.Password,
		args.SSLMode,
	)

	dbConnnection, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	return &Database{
		Client: dbConnnection,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.PingContext(ctx)
}
