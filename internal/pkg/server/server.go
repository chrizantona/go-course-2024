package server

import (
	"encoding/json"
	"go-course-2024/internal/pkg/storage"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
)

type Response struct {
	Value string `json:"value"`
}

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value string `json:"value"`
	TTL   int64  `json:"ttl"` 
}


func NewServer(host string, st *storage.Storage) *Server {
	s := &Server{
		host:    host,
		storage: st,
	}

	return s
}

func (r *Server) NewAPI() *gin.Engine {
	engine := gin.New()
	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server is running"})
	})

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	engine.GET("/scalar/get/:key", r.handlerGetScalar)
	engine.PUT("/scalar/set/:key", r.handlerSetScalar)

	return engine
}

func (r *Server) handlerGetScalar(ctx *gin.Context) {
	key := ctx.Param("key")
	valuePtr, err := r.storage.Get(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if valuePtr == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	value, ok := (*valuePtr).(string)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) 
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: value})
}



func (r *Server) handlerSetScalar(ctx *gin.Context) {
	key := ctx.Param("key")
	var entry Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&entry); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ttl := time.Duration(entry.TTL) * time.Second
	if err := r.storage.Set(key, entry.Value, ttl); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}



func (r *Server) Start() error {
	return r.NewAPI().Run(r.host)
}
