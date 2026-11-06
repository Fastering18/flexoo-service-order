package config

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB   *sql.DB
	once sync.Once
)

func ConnectDatabase() {
	once.Do(func() {
		_ = godotenv.Load()

		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = os.Getenv("POSTGRES_URL")
		}
		if dsn == "" {
			log.Println("DATABASE_URL not set, running with fallback store")
			return
		}

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Database connection error: %v\n", err)
			return
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Printf("Database ping error: %v\n", err)
			return
		}

		DB = db
		runSchemaFile(db)
	})
}

func runSchemaFile(db *sql.DB) {
	paths := []string{
		"src/database/schema.sql",
		"database/schema.sql",
		"../src/database/schema.sql",
		"/var/task/src/database/schema.sql",
	}
	for _, p := range paths {
		if content, err := os.ReadFile(p); err == nil {
			if _, err := db.Exec(string(content)); err != nil {
				log.Printf("Schema init warning: %v\n", err)
			}
			return
		}
	}
}

