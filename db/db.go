package db

import (
	"errors"
	"time"
)

// Storage defines generic storage behaviour
type Storage interface {
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
}

// entry represents a single key's value + expiration
type entry struct {
	value      string
	expiration time.Time // zero means no expiration
}

// DB is an in-memory map with mutex + expiration support
type DB struct {
	data map[string]entry
}

var Instance *DB

// Init initializes the DB instance
func Init() {
	Instance = &DB{
		data: make(map[string]entry),
	}
}

// Set stores a key with optional TTL
func (d *DB) Set(key, value string, ttl time.Duration) error {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	d.data[key] = entry{
		value:      value,
		expiration: exp,
	}
	return nil
}

// Get retrieves a key, respecting expiration
func (d *DB) Get(key string) (string, error) {
	e, ok := d.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	// check if expired
	// if no time is set than zero value of time is present in golang by default
	if !e.expiration.IsZero() && time.Now().After(e.expiration) {
		// expired → delete key

		delete(d.data, key)

		return "", errors.New("key expired")
	}

	return e.value, nil
}
