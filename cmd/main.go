package main

import (
	"go-course-2024/internal/pkg/server"
	"go-course-2024/internal/pkg/storage"
	"log"
)

type Response struct {
	Value string `json:"value"`
}

func main() {
	myStorage := storage.NewStorage()
	s := server.NewServer(":8090", myStorage)
	myStorage.Set("asdf", "asdf")

	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
