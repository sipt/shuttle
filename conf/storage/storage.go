package storage

import (
	"context"
	"fmt"
)

type NewStorageFunc func(map[string]string) (IStorage, error)

var storageCreator = make(map[string]NewStorageFunc)

// Register: register {key: storage}
func Register(key string, f NewStorageFunc) {
	storageCreator[key] = f
}

// Get: get storage by key
func Get(key string, params map[string]string) (IStorage, error) {
	f, ok := storageCreator[key]
	if !ok {
		return nil, fmt.Errorf("storage not support: %s", key)
	}
	return f(params)
}

type IStorage interface {
	Name() string
	// Load: load config from disk? As JSON? YAML? TOML?
	Load() (data []byte, err error)
	// Save: save config to file? upload to server?
	Save(data []byte) (err error)
	// RegisterNotify
	RegisterNotify(ctx context.Context, notify func()) error
}
