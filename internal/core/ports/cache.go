package ports

import "time"

type CacheRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, value interface{}) error
	SetString(key, value string, expiration time.Duration) error
	GetString(key string) (string, error)
	Delete(key string) error
	GetClient() interface{}
	Exists(key string) (bool, error)
	Ping() error
	Close() error
}
