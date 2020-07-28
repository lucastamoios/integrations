package storage

import (
	"sync"
	"time"
)

type HashStorage interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Del(key string)
	Expire(key string, seconds time.Duration)
}

type HashValue struct {
	Value string
	ExpiresAt *time.Time
}

type MapStorage struct {
	storage sync.Map
}

func NewHashStorage() *MapStorage{
	return &MapStorage{sync.Map{}}
}

func (ms *MapStorage) Set(key, value string) {
	v := HashValue{value, nil}
	ms.storage.Store(key, v)
}

func (ms *MapStorage) Get(key string) (string, bool) {
	ms.checkExpiration(key)
	value, ok := ms.storage.Load(key)
	if !ok {
		return "", false
	}
	return value.(HashValue).Value, ok
}

func (ms *MapStorage) Del(key string) {
	ms.storage.Delete(key)
}

func (ms *MapStorage) Expire(key string, duration time.Duration) {
	if value, ok := ms.storage.Load(key); ok {
		v := value.(HashValue).Value
		e := time.Now().Add(duration)
		ms.storage.Store(key, HashValue{v, &e})
	}
}

func (ms *MapStorage) checkExpiration (key string) {
	if value, ok := ms.storage.Load(key); ok && value.(HashValue).ExpiresAt.Before(time.Now()){
		ms.storage.Delete(key)
	}
}