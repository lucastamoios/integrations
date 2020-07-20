package storage

import (
	"log"
	"sync"
)

type HashStorage interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Del(key string)
	Pop(key string) (string, bool)  // Same as Get, but also removes from storage
	Log(s string)
}

type MapStorage struct {
	storage sync.Map
}

func NewHashStorage() *MapStorage{
	return &MapStorage{sync.Map{}}
}

func (ms *MapStorage) Set(key, value string) {
	ms.storage.Store(key, value)
}

func (ms *MapStorage) Get(key string) (string, bool) {
	value, ok := ms.storage.Load(key)
	if !ok {
		return "", false
	}
	return value.(string), ok
}

func (ms *MapStorage) Pop(key string) (string, bool) {
	value, ok := ms.storage.Load(key)
	if !ok {
		return "", false
	}
	ms.storage.Delete(key)
	return value.(string), ok
}

func (ms *MapStorage) Del(key string) {
	ms.storage.Delete(key)
}


func (ms *MapStorage) Log(s string) {
	log.Printf("%s: %v\n", s, ms.storage)
}