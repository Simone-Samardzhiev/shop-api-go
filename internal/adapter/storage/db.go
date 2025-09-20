package storage

import (
	"database/sql"
	"shop-api-go/internal/adapter/config"
)

// New opens a new connection to sql.DB.
func New(config *config.Database) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdleConnection)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	return db, nil
}
