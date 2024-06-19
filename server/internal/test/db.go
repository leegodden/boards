package test

import (
	"testing"

	"github.com/leegodden/boards/server/db"
	"github.com/leegodden/boards/server/internal/config"
)

func DB(t *testing.T) *db.DB {
	
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		t.Errorf("Issue loading config:%v", err)
		t.FailNow()
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		t.Errorf("Issue connecting db:%v", err)
		t.FailNow()
	}
	return db
}