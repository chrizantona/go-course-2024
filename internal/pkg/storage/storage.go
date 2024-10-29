package storage

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"strconv"
)

type Kind string

const (
	KindInt     Kind = "D"
	KindString  Kind = "S"
	KindUnknown Kind = ""
)

type Storage struct {
	listStorage map[string][]string
	inner       map[string]string
	logger      *zap.Logger
}

func NewStorage() *Storage {
	logger, _ := zap.NewProduction(zap.IncreaseLevel(zapcore.DPanicLevel))
	defer logger.Sync()
	logger.Info("created new storage")
	return &Storage{
		listStorage: make(map[string][]string),
		inner:       make(map[string]string),
		logger:      logger,
	}
}

func (s *Storage) LPUSH(key string, elements ...string) error {
	if _, exists := s.inner[key]; exists {
		return fmt.Errorf("key %s already exists in inner storage", key)
	}
	s.listStorage[key] = append(elements, s.listStorage[key]...)
	s.logger.Info("LPUSH executed", zap.String("key", key), zap.Strings("elements", elements))
	return nil
}

func (s *Storage) RPUSH(key string, elements ...string) error {
	if _, exists := s.inner[key]; exists {
		return fmt.Errorf("key %s already exists in inner storage", key)
	}
	s.listStorage[key] = append(s.listStorage[key], elements...)
	s.logger.Info("RPUSH executed", zap.String("key", key), zap.Strings("elements", elements))
	return nil
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func (s *Storage) RADDTOSET(key string, elements ...string) error {
	if _, exists := s.inner[key]; exists {
		return fmt.Errorf("key %s already exists in inner storage", key)
	}
	for _, element := range elements {
		if !contains(s.listStorage[key], element) {
			s.listStorage[key] = append(s.listStorage[key], element)
		}
	}
	s.logger.Info("RADDTOSET executed", zap.String("key", key), zap.Strings("elements", elements))
	return nil
}

func (s *Storage) LPOP(key string, count ...int) ([]string, error) {
	n := count[0]
	if n > len(s.listStorage[key]) {
		n = len(s.listStorage[key])
	}
	res := s.listStorage[key][:n]
	s.listStorage[key] = s.listStorage[key][n:]

	s.logger.Info("LPOP executed", zap.String("key", key), zap.Int("count", n))

	return res, nil
}

func (s *Storage) RPOP(key string, count ...int) ([]string, error) {
	n := count[0]
	if n > len(s.listStorage[key]) {
		n = len(s.listStorage[key])
	}
	res := s.listStorage[key][len(s.listStorage[key])-n:]
	s.listStorage[key] = s.listStorage[key][:len(s.listStorage[key])-n]

	s.logger.Info("RPOP executed", zap.String("key", key), zap.Int("count", n))

	return res, nil
}

func (s *Storage) LSET(key string, index int, element string) error {
	if index < 0 || index >= len(s.listStorage[key]) {
		return fmt.Errorf("index out of range")
	}
	s.listStorage[key][index] = element
	s.logger.Info("LSET executed", zap.String("key", key), zap.Int("index", index), zap.String("element", element))
	return nil
}

func (s *Storage) LGET(key string, index uint) (string, error) {
	if index >= uint(len(s.listStorage[key])) {
		return "", fmt.Errorf("index %d out of bounds for key %s", index, key)
	}

	list, exists := s.listStorage[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	return list[index], nil
}

func (r *Storage) Set(key, value string) error {
	if _, exists := r.listStorage[key]; exists {
		return fmt.Errorf("key %s already exists in listStorage", key)
	}
	r.inner[key] = value
	r.logger.Info("set value in storage", zap.String("key", key), zap.String("value", value))
	return nil
}

func (r *Storage) Get(key string) (*string, error) {
	res, ok := r.inner[key]
	if !ok {
		r.logger.Warn("key not found in storage", zap.String("key", key))
		return nil, fmt.Errorf("key %s not found", key)
	}
	r.logger.Info("retrieved value from storage", zap.String("key", key), zap.String("value", res))
	return &res, nil
}

func (r *Storage) GetKind(key string) (Kind, error) {
	res, ok := r.inner[key]
	if !ok {
		r.logger.Warn("key not found when determining kind", zap.String("key", key))
		return KindUnknown, fmt.Errorf("key %s not found", key)
	}

	if _, err := strconv.Atoi(res); err == nil {
		kind := KindInt
		r.logger.Info("determined kind of value", zap.String("key", key), zap.String("kind", string(kind)))
		return kind, nil
	}

	kind := KindString
	r.logger.Info("determined kind of value", zap.String("key", key), zap.String("kind", string(kind)))
	return kind, nil
}

func (s *Storage) SaveToFile(filename string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("error marshalling storage: %v", err)
	}
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	s.logger.Info("Storage saved to file", zap.String("filename", filename))
	return nil
}

func (s *Storage) LoadFromFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	if err := json.Unmarshal(file, s); err != nil {
		return fmt.Errorf("error unmarshalling storage: %v", err)
	}
	s.logger.Info("Storage loaded from file", zap.String("filename", filename))
	return nil
}