package main

import (
	"go-course-2024/internal/pkg/server"
	"go-course-2024/internal/pkg/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	myStorage := storage.NewStorage()
	const storageFilePath = "/root/storage_data/storage.json"

	if err := myStorage.SaveToFile(storageFilePath); err != nil {
		log.Printf("Failed to save storage: %v", err)
	}
	log.Println("Storage saved. Exiting.")

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


}
