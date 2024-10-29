package main

import (
	"GO-COURSE-2024/internal/pkg/server"
	"GO-COURSE-2024/internal/pkg/storage"
)

type Response struct {
	Value string `json:"value"`
}

func main() {
	myStorage := storage.NewStorage()
	s := server.New(":8090", &myStorage)
	myStorage.Set("asdf", "asdf")
	s.Start()
}
