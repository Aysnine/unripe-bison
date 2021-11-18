package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/storage/redis"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectPG(connString string) *pgxpool.Pool {
	// Database connect timing
	start := time.Now()

	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = db.QueryRow(context.Background(), "select 'Database connected!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	// Database connect timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + greeting)

	return db
}

func ConnectRedis(connString string) *redis.Storage {
	// Database connect timing
	start := time.Now()

	store := redis.New(redis.Config{
		URL:   connString,
		Reset: false,
	})

	// Database connect timing
	stop := time.Now()
	fmt.Println(fmt.Sprintf("[duration=%v] ", stop.Sub(start).String()) + "redis item connected!")

	return store
}
