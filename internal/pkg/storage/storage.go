package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Kind string

const (
	KindInt     Kind = "D"
	KindString  Kind = "S"
	KindUnknown Kind = ""
	KindDict    Kind = "M"
)

type value struct {
	v         interface{}
	expiresAt int64
}

type Storage struct {
	listStorage map[string][]string
	inner       map[string]interface{}
	expiration  map[string]int64
	logger      *zap.Logger
	mu          sync.RWMutex 
}


func (s *Storage) Logger() *zap.Logger {
    return s.logger
}

func NewStorage() *Storage {
    logger, _ := zap.NewProduction(zap.IncreaseLevel(zapcore.DPanicLevel))
    return &Storage{
        listStorage: make(map[string][]string),
        inner:       make(map[string]interface{}),  
        expiration:  make(map[string]int64),       
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

func (r *Storage) Set(key string, value interface{}, ttl time.Duration) error {
    if r.inner == nil {
        r.logger.Error("inner map is nil, initializing")
        r.inner = make(map[string]interface{})
    }
    if r.expiration == nil {
        r.logger.Error("expiration map is nil, initializing")
        r.expiration = make(map[string]int64)
    }

    r.logger.Info("Setting key in storage", zap.String("key", key), zap.Any("value", value))

    if _, exists := r.listStorage[key]; exists {
        return fmt.Errorf("key %s already exists in listStorage", key)
    }

    r.inner[key] = value
    r.expiration[key] = time.Now().Add(ttl).UnixMilli()
    r.logger.Info("Set value successfully", zap.String("key", key), zap.Any("value", value))

    return nil
}




func (r *Storage) Get(key string) (*interface{}, error) {
	res, ok := r.inner[key]
	if !ok {
		r.logger.Warn("key not found in storage", zap.String("key", key))
		return nil, fmt.Errorf("key %s not found", key)
	}

	if time.Now().UnixMilli() >= r.expiration[key] {
		delete(r.inner, key)     
		delete(r.expiration, key)   
		r.logger.Info("key expired and deleted", zap.String("key", key))
		return nil, fmt.Errorf("key %s has expired", key)
	}

	r.logger.Info("retrieved value from storage", zap.String("key", key), zap.Any("value", res))
	return &res, nil
}


func (s *Storage) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			s.cleanExpiredKeys()
		}
	}()
}

func (s *Storage) cleanExpiredKeys() {
	s.logger.Info("cleaning expired keys")
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, expirationTime := range s.expiration {
		if time.Now().UnixMilli() >= expirationTime {
			delete(s.inner, key)  
			delete(s.expiration, key)   
			s.logger.Info("deleted expired key", zap.String("key", key))
		}
	}
}



func (r *Storage) GetKind(key string) (Kind, error) {
	res, ok := r.inner[key]
	if !ok {
		r.logger.Warn("key not found when determining kind", zap.String("key", key))
		return KindUnknown, fmt.Errorf("key %s not found", key)
	}

	strValue, ok := res.(string)
	if ok {
		if _, err := strconv.Atoi(strValue); err == nil {
			kind := KindInt
			r.logger.Info("determined kind of value", zap.String("key", key), zap.String("kind", string(kind)))
			return kind, nil
		}

		kind := KindString
		r.logger.Info("determined kind of value", zap.String("key", key), zap.String("kind", string(kind)))
		return kind, nil
	}

	if _, ok := res.(map[string]interface{}); ok {
		kind := KindDict
		r.logger.Info("determined kind of value", zap.String("key", key), zap.String("kind", string(kind)))
		return kind, nil
	}

	r.logger.Warn("unknown kind for key", zap.String("key", key))
	return KindUnknown, fmt.Errorf("unknown kind for key %s", key)
}

func (s *Storage) SaveToFile(filename string) error {
    s.logger.Info("Attempting to save storage to file", zap.String("filename", filename))


    if len(s.inner) == 0 {
        s.logger.Warn("No data to save. Storage is empty")
        return nil
    }


    saveData := struct {
        Inner      map[string]interface{} `json:"inner"`
        Expiration map[string]int64       `json:"expiration"`
    }{
        Inner:      s.inner,
        Expiration: s.expiration,
    }


    data, err := json.Marshal(saveData)
    if err != nil {
        s.logger.Error("Error marshalling storage", zap.Error(err))
        return fmt.Errorf("error marshalling storage: %v", err)
    }


    if err := ioutil.WriteFile(filename, data, 0666); err != nil {
        s.logger.Error("Error writing to file", zap.String("filename", filename), zap.Error(err))
        return fmt.Errorf("error writing to file: %v", err)
    }

    s.logger.Info("Data successfully saved to file", zap.String("filename", filename))
    return nil
}


func (s *Storage) LoadFromFile(filename string) error {
    s.logger.Info("Attempting to load storage from file", zap.String("filename", filename))


    file, err := ioutil.ReadFile(filename)
    if err != nil {
        s.logger.Error("Error reading file", zap.String("filename", filename), zap.Error(err))
        return fmt.Errorf("error reading file: %v", err)
    }

    loadData := struct {
        Inner      map[string]interface{} `json:"inner"`
        Expiration map[string]int64       `json:"expiration"`
    }{}


    if err := json.Unmarshal(file, &loadData); err != nil {
        s.logger.Error("Error unmarshalling storage", zap.Error(err))
        return fmt.Errorf("error unmarshalling storage: %v", err)
    }

    s.inner = loadData.Inner
    s.expiration = loadData.Expiration

    s.logger.Info("Data successfully loaded from file", zap.String("filename", filename))
    return nil
}
