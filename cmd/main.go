package main

import (
	"go-course-2024/internal/pkg/server"
	"go-course-2024/internal/pkg/storage"
	"log"
	"time"
)

type Response struct {
	Value string `json:"value"`
}

func main() {
	myStorage := storage.NewStorage()
	myStorage.StartCleanup(time.Minute) 
	s := server.NewServer(":8090", myStorage)
	myStorage.Set("asdf", "asdf", 10*time.Second) 

	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
