package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/marwans200/cyberguard/internal/config"
)

// DB is the global database connection pool.
var DB *sql.DB

// Connect establishes the PostgreSQL connection pool and pings the server.
func Connect(cfg *config.Config) {
	var err error
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", cfg.DSN())
		if err == nil {
			DB.SetMaxOpenConns(25)
			DB.SetMaxIdleConns(5)
			DB.SetConnMaxLifetime(5 * time.Minute)
			if pingErr := DB.Ping(); pingErr == nil {
				log.Println("[db] connected to PostgreSQL")
				return
			}
		}
		log.Printf("[db] waiting for PostgreSQL... attempt %d/10", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("[db] could not connect to PostgreSQL: %v", err)
}
