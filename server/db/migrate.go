package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/leegodden/boards/server/internal/config"
)

func Migrate(dbConfig config.DatabaseConfig) error {
	url := buildConnectionURL(dbConfig)
	m, err := migrate.New(
		"file:db/migrations",
		url)

	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migration changes")
			return nil
		}
		return err
	}
	fmt.Println("Successfully ran db migrations.")
	return nil
}