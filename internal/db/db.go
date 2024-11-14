package db

import (
	"context"
	"database/sql"
	"time"
)

func NewDB (addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {

	db, err := sql.Open("postgres", addr)
	if err !=nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, err

}