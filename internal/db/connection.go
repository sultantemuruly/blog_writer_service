package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, dsn string ) (*pgx.Conn, error) {
    conn, err := pgx.Connect(ctx, dsn)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }

	log.Println("Connected to the database successfully")
	return conn, nil
}