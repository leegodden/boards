package main

import (
	"log"

	"github.com/leegodden/boards/server/db"
	"github.com/leegodden/boards/server/internal/config"
)

func main() {
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer db.Close()
}