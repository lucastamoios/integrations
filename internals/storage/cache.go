package storage

import "sync"

type HashStorage interface {
	Set(key, value string)
	Get(key string) (string, bool)
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
	return value.(string), ok
}
