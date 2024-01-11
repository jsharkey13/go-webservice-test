package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initDb() {
	// Tried using database/sql alone, but it does not support the jsonb[] type we use in places!
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		os.Exit(1) // FIXME:fatal
	} else {
		db = dbPool
	}
}
