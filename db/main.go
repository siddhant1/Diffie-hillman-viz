package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

func NewPostgresClient() *sql.DB {
	connStr := "user=postgres dbname=diffie sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
