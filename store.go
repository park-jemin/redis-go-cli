package main

import (
	"container/list"
	"errors"
)

// DataType represents the type of data stored in the store
type DataType uint8

const (
	String DataType = iota + 1
	List
	Hash
)

// Store represents the in-memory store for key-value pairs, lists, and hashes
type Store struct {
	types  map[string]DataType          // For tracking the type of each key
	data   map[string]string            // For string keys/values
	lists  map[string]*list.List        // For lists
	hashes map[string]map[string]string // For hashes
}

// NewStore initializes a new Store
func NewStore() *Store {
	return &Store{
		types:  make(map[string]DataType),
		data:   make(map[string]string),
		lists:  make(map[string]*list.List),
		hashes: make(map[string]map[string]string),
	}
}

// TypeCheck verifies if a key is of the expected data type
func (s *Store) TypeCheck(key string, expected DataType) error {
	if dataType, found := s.types[key]; found && dataType != expected {
		return errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return nil
}

// Set upserts a key + (string) value pair to the store, returning the previous value
func (s *Store) Set(key, value string) error {
	if err := s.TypeCheck(key, String); err != nil {
		return err
	}
	s.data[key] = value
	s.types[key] = String
	return nil
}

// Get retrieves a (string) value by key from the store
func (s *Store) Get(key string) (string, bool, error) {
	if err := s.TypeCheck(key, String); err != nil {
		return "", false, err
	}
	value, exists := s.data[key]
	return value, exists, nil
}

// Del deletes a key from the store based on its type
func (s *Store) Del(keys ...string) uint {
	deleted := uint(0)

	for _, key := range keys {
		// Check if the key exists in the types map
		if dataType, found := s.types[key]; found {
			// Remove from the appropriate map based on the data type
			switch dataType {
			case String:
				delete(s.data, key)
			case List:
				delete(s.lists, key)
			case Hash:
				delete(s.hashes, key)
			}
			// Remove the key's entry from the types map
			delete(s.types, key)
			deleted++
		}
	}

	return deleted
}

