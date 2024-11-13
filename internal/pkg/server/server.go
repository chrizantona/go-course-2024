package server

import (
	"go.uber.org/zap"
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
		ctx.JSON(http.StatusOK, gin.H{"status": "Healthy"})
	})

	engine.GET("/scalar/get/:key", r.handlerGetScalar)
	engine.PUT("/scalar/set/:key", r.handlerSetScalar)

	return engine
}

func (r *Server) handlerGetScalar(ctx *gin.Context) {
    key := ctx.Param("key")
    valuePtr, err := r.storage.Get(key)
    if err != nil {
        r.storage.Logger().Warn("Error retrieving key", zap.String("key", key), zap.Error(err))
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Key not found or expired"})
        return
    }
    if valuePtr == nil {
        r.storage.Logger().Warn("Key not found or expired", zap.String("key", key))
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Key not found or expired"})
        return
    }
    value, ok := (*valuePtr).(string)
    if !ok {
        r.storage.Logger().Error("Value type mismatch", zap.String("key", key))
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Value type mismatch"})
        return
    }
    r.storage.Logger().Info("Key retrieved", zap.String("key", key), zap.Any("value", value))
    ctx.JSON(http.StatusOK, gin.H{"value": value})
}



func (r *Server) handlerSetScalar(ctx *gin.Context) {
    r.storage.Logger().Info("Start handlerSetScalar")
    key := ctx.Param("key")
    r.storage.Logger().Info("Parsed key", zap.String("key", key))

    var entry Entry
    if err := ctx.ShouldBindJSON(&entry); err != nil {
        r.storage.Logger().Error("Invalid JSON in request", zap.Error(err))
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid JSON",
            "details": err.Error(), 
        })
        return
    }


    ttl := time.Duration(entry.TTL) * time.Second
    r.storage.Logger().Info("Parsed TTL", zap.Duration("ttl", ttl))


    if err := r.storage.Set(key, entry.Value, ttl); err != nil {
        r.storage.Logger().Error("Failed to set value in storage", zap.String("key", key), zap.Error(err))
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to set value in storage",
            "details": err.Error(),
        })
        return
    }

    r.storage.Logger().Info("Set value successfully", zap.String("key", key), zap.Any("value", entry.Value))
    ctx.JSON(http.StatusOK, gin.H{"status": "Value set successfully"})
}






func (r *Server) Start() error {
	return r.NewAPI().Run(r.host)
}
