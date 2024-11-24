package main

import (
    _ "github.com/lib/pq"
    "go-course-2024/internal/pkg/server"
    "go-course-2024/internal/pkg/storage"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    "database/sql"
    "fmt"
)

func main() {


    db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
	defer db.Close()


    myStorage := storage.NewStorage()
    const storageFilePath = "/root/storage_data/storage.json"

    if err := myStorage.LoadFromFile(storageFilePath); err != nil {
        log.Printf("Failed to load storage: %v", err)
    } else {
        log.Println("Storage loaded successfully")
    }

    myStorage.StartCleanup(time.Minute)

    s := server.NewServer(":8090", myStorage)

    go func() {
        if err := s.Start(); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop
    log.Println("Shutting down gracefully...")
    if err := myStorage.SaveToFile(storageFilePath); err != nil {
        log.Printf("Failed to save storage: %v", err)
    } else {
        log.Println("Storage saved successfully")
    }
}


func setupDatabase() (*sql.DB, error) {
	dsn := os.Getenv("POSTGRES") 
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS core (
		version bigserial PRIMARY KEY,
		timestamp bigint NOT NULL,
		payload JSONB NOT NULL
	)`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}
