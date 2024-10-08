package storage
import (
	"go.uber.org/zap"
)
type Storage struct {
	inner map[string]string
	logger *zap.Logger
}
func NewStorage() Storage {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("created new storage")
	return Storage{
	inner: make(map[string]string),
		logger: logger,
	}
}
func (r *Storage) Set(key, value string) {
	r.inner[key] = value
	r.logger.Info("set value in storage", zap.String("key", key), zap.String("value", value))
}
func (r Storage) Get(key string) *string {
	res, ok := r.inner[key]
	if !ok {
		r.logger.Warn("key not found in storage", zap.String("key", key))
		return nil
	}
	r.logger.Info("retrieved value from storage", zap.String("key", key), zap.String("value", res))
	return &res
}
func (r Storage) GetKind(key string) string {
	res, ok := r.inner[key]
	if !ok {
		r.logger.Warn("key not found when determining kind", zap.String("key", key))
		return ""
	}
	isNumber := true
	for _, char := range res {
		if char < '0' || char > '9' {
			isNumber = false
			break
		}
	}
	kind := "S" 
	if isNumber {
		kind = "D"
	}
	r.logger.Info("determined kind of value", zap.String("key", key), zap.String("kind", kind))
	return kind
}
