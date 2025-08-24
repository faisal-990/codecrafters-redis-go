package db

import (
	"errors"
	"time"
)

// Storage defines generic storage behaviour
type Storage interface {
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Rpush(key string, value []string) (int, error)
	Lrange(key, start, stop string) ([]string, error)
}

// entry represents a single key's value + expiration (for string keys)
type entry struct {
	value      string
	expiration time.Time // zero means no expiration
}

// DB is an in-memory map with expiration support for strings + separate lists map
type DB struct {
	data  map[string]entry    // string keys with optional expiry
	lists map[string][]string // list keys (no expiry for now)
}

var Instance *DB

// Init initializes the DB instance
func Init() {
	Instance = &DB{
		data:  make(map[string]entry),
		lists: make(map[string][]string),
	}
}

// Set stores a key with optional TTL (for string values)
func (d *DB) Set(key, value string, ttl time.Duration) error {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	d.data[key] = entry{
		value:      value,
		expiration: exp,
	}
	// If this key previously was a list, remove the list to avoid type confusion
	delete(d.lists, key)
	return nil
}

// Get retrieves a key, respecting expiration (for string values)
func (d *DB) Get(key string) (string, error) {
	e, ok := d.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	// lazy expiration check
	if !e.expiration.IsZero() && time.Now().After(e.expiration) {
		delete(d.data, key)
		return "", errors.New("key expired")
	}

	return e.value, nil
}

// -------------------- LIST COMMANDS --------------------

// Rpush appends a value to the tail (right) of the list at key.
// - Creates a new list if key doesn't exist (as list or string).
// - If key exists as a non-expired string, returns WRONGTYPE error.
// - If key exists as an expired string, deletes it and creates a new list.
// Returns the length of the list after the push.
func (d *DB) Rpush(key string, value []string) (int, error) {
	// If key currently holds a string, check expiration and enforce type
	if e, ok := d.data[key]; ok {
		// lazy expire
		if !e.expiration.IsZero() && time.Now().After(e.expiration) {
			// expired string -> remove and allow list creation
			delete(d.data, key)
		} else {
			// non-expired string -> WRONGTYPE
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	// Append to existing list (or create a new one)
	d.lists[key] = append(d.lists[key], value...)
	return len(d.lists[key]), nil
}

func (d *DB) Lrange(key string, start, stop int) ([]string, error) {
	list, ok := d.lists[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	length := len(list)

	// handle negative indices (like Redis)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	// clamp indices
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop || start >= length {
		return []string{}, nil
	}

	return list[start : stop+1], nil
}
