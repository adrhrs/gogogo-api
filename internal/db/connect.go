package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() (*pgxpool.Pool, error) {
	isInternal := true
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		isInternal = false
		dsn = "postgresql://gogogo_db_user:iM2hd6qpwVq071f6PLeJ8ehNqfSBSBIw@dpg-d12f9ejuibrs73f5kvbg-a.oregon-postgres.render.com/gogogo_db"
	}

	log.Println("dsn", dsn, isInternal)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping db: %w", err)
	}

	return pool, nil
}
