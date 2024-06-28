package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leegodden/boards/server/internal/config"
)

// DB is a wrapper around pgxpool.Pool. It embeds *pgxpool.Pool
// so that all of its method set is promoted to be accessible directly from DB instances.
type DB struct {
	*pgxpool.Pool
}

// Establishes a new db connection based on the provided configuration from config.go.
// Returns a pointer to the DB struct with the established connection or an error.
func Connect(cfg config.DatabaseConfig) (*DB, error) {

// Build the connection URL
	url := buildConnectionURL(cfg)
 
	// Establish a new connection to db
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")

// Return a pointer to a new DB instance so we can work directly with the same db 
// connection across different parts of our program
	return &DB{db}, nil
}

func buildConnectionURL(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}