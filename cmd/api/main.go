package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"auth-service/internal/config"
	"auth-service/internal/server"
)

func main() {
	_ = godotenv.Load()
	config := config.Load()

	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("db ping failed:", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	handler := server.New(db)

	srv := &http.Server{
		Handler:      handler,
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("🚀 Server running on http://localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
