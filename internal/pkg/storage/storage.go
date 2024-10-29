package storage

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"encoding/json"
)
	type Kind string
	const(
	KindInt Kind = "D"
	KindString Kind = "S"
	KindUnknown Kind = ""
	)
	type Storage struct {
		inner map[string]string
		listStorage map[string][]string
		logger *zap.Logger
	}
	func NewStorage() Storage {
		logger, _ := zap.NewProduction(zap.IncreaseLevel(zapcore.DPanicLevel))
		defer logger.Sync()
		logger.Info("created new storage")
		return Storage{
		listStorage: make(map[string][]string),
		inner: make(map[string]string),
			logger: logger,
		}
	}

	func (s *Storage) LPUSH(key string, elements ...string) int {
		s.listStorage[key] = append(elements, s.listStorage[key]...)
		s.logger.Info("LPUSH executed", zap.String("key", key), zap.Strings("elements", elements))
		return len(s.listStorage[key])
	}

	func (s *Storage) RPUSH(key string, elements ...string) int {
		s.listStorage[key] = append(s.listStorage[key], elements...)
		s.logger.Info("RPUSH executed", zap.String("key", key), zap.Strings("elements", elements))
		return len(s.listStorage[key])
	}
	func contains(slice []string, item string)bool {
		for _, element := range slice{
			if element == item{
				return true
			}
		}
		return false
	}

	func (s* Storage) RADDTOSET(key string, elements ...string) int{
		if _, exists := s.listStorage[key]; !exists {
			s.listStorage[key] = []string{}
		}
		for _, 	element := range elements{
			if !contains(s.listStorage[key],element){
				s.listStorage[key] = append(s.listStorage[key],element)
			}
		}
		s.logger.Info("RADDTOSET executed",zap.String("key",key),zap.Strings("elements", elements))
		return len(s.listStorage[key])
	}


	func (s *Storage) LPOP(key string, count ...int) []string {
		n := count[0]
		if n > len(s.listStorage[key]) {
			n = len(s.listStorage[key])
		}
		res := s.listStorage[key][:n]
		s.listStorage[key] = s.listStorage[key][n:]
		
		s.logger.Info("LPOP executed", zap.String("key", key), zap.Int("count", n))
		
		return res
	}


	func (s *Storage) RPOP(key string, count ...int) []string {
		n := count[0]
		if n > len(s.listStorage[key]) {
			n = len(s.listStorage[key])
		}
		res := s.listStorage[key][len(s.listStorage[key])-n:]
		s.listStorage[key] = s.listStorage[key][:len(s.listStorage[key])-n]

		s.logger.Info("RPOP executed", zap.String("key", key), zap.Int("count", n))
		
		return res
	}

	func (s *Storage) LSET(key string, index int, element string) error {
		if index < 0 || index >= len(s.listStorage[key]) {
			return fmt.Errorf("index out of range")
		}
		s.listStorage[key][index] = element
		s.logger.Info("LSET executed", zap.String("key", key), zap.Int("index", index), zap.String("element", element))
		return nil
	}
	func (s *Storage) LGET(key string, index int) (string, error) {
		if list, exists := s.listStorage[key]; !exists {
			return "", fmt.Errorf("key %s not found", key)
		} else if index < 0 || index >= len(list) {
			return "", fmt.Errorf("index out of range")
		}
		value := s.listStorage[key][index]
		s.logger.Info("LGET executed", zap.String("key", key), zap.Int("index", index), zap.String("value", value))
		return value, nil
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
	