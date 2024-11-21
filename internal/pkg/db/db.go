package db

import (
	"amartha/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(cfg config.Database) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Address, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)

	db, err := sqlx.Connect(cfg.Type, connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}