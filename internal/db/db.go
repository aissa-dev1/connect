package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func Connect() *pgxpool.Pool {
	once.Do(func() {
		connectPool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONNECT_URL"))

		if err != nil {
			log.Fatalf("Unable to connect to database %v\n", err)
		}

		pool = connectPool
		fmt.Println("Database connected successfully")
	})

	return pool
}

func Pool() *pgxpool.Pool {
	if pool == nil {
		log.Fatalf("Database pool not initialized. Call Connect() first.")
	}

	return pool
}
