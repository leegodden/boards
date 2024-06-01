package db

import (
	"fmt"

	"github.com/Wave-95/boards/server/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	fmt.Println("Successfully run db migrations.")
	return nil
}

func buildConnectionURL(dbConfig config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name)
}