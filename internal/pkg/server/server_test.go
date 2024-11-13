package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"go-course-2024/internal/pkg/server"
	"go-course-2024/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	myStorage := storage.NewStorage()
	s := server.NewServer(":8090", myStorage)

	errChan := make(chan error)

	go func() {
		errChan <- s.Start() 
	}()
	if err := <-errChan; err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	t.Run("TestRoot", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		s.NewAPI().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Server is running", response["message"])
	})

	t.Run("TestHealth", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		s.NewAPI().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestSetAndGetScalar", func(t *testing.T) {
		w := httptest.NewRecorder()
		entry := server.Entry{Value: "testValue"}
		body, _ := json.Marshal(entry)
		req, _ := http.NewRequest(http.MethodPut, "/scalar/set/testKey", bytes.NewBuffer(body))
		s.NewAPI().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/scalar/get/testKey", nil)
		s.NewAPI().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response server.Entry
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "testValue", response.Value)
	})
}