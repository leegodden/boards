package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leegodden/boards/server/internal/config"
)

type DB struct {
	*pgxpool.Pool
}

func Connect(cfg config.DatabaseConfig) (*DB, error) {
	url := buildConnectionURL(cfg)
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")
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