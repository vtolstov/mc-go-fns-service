package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func DatabaseConnect(cfg *DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", cfg.Name)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(int(cfg.ConnMax))
	db.SetMaxIdleConns(int(cfg.ConnMax / 2))
	db.SetConnMaxLifetime(time.Duration(cfg.ConnLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	return db, nil
}
