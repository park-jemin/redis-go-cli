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

// LPush inserts values at the head of a list
func (s *Store) LPush(key string, values ...string) (int, error) {
	if err := s.TypeCheck(key, List); err != nil {
		return 0, err
	}
	if _, exists := s.lists[key]; !exists {
		s.lists[key] = list.New()
		s.types[key] = List
	}
	for _, value := range values {
		s.lists[key].PushFront(value)
	}
	return len(values), nil
}

// LPop removes and returns the first element of a list
func (s *Store) LPop(key string) (string, bool, error) {
	if err := s.TypeCheck(key, List); err != nil {
		return "", false, err
	}
	if list, found := s.lists[key]; found && list.Len() > 0 {
		value := list.Remove(list.Front()).(string)
		return value, true, nil
	}
	return "", false, nil
}

func (s *Store) LLen(key string) (int, error) {
	if err := s.TypeCheck(key, List); err != nil {
		return 0, err
	}
	if list, found := s.lists[key]; found {
		return list.Len(), nil
	}
	return 0, nil
}

// LRange returns a subrange of elements from the list at the given key
func (s *Store) LRange(key string, start, stop int) ([]string, error) {
	// Type check to ensure the key is a list
	if err := s.TypeCheck(key, List); err != nil {
		return nil, err
	}

	list, found := s.lists[key]
	if !found || list.Len() == 0 {
		return []string{}, nil
	}

	length := list.Len()

	// Handle negative indexes
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	// If start is greater than the length of the list, return an empty array
	if start >= length {
		return []string{}, nil
	}

	// Bound the stop index to the end of the list
	if stop >= length {
		stop = length - 1
	}

	// If start is greater than stop, return an empty array
	if start > stop {
		return []string{}, nil
	}

	// Collect elements within the specified range
	result := []string{}
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		if i >= start && i <= stop {
			result = append(result, e.Value.(string))
		}
		if i > stop {
			break
		}
		i++
	}

	return result, nil
}

// HSet sets a field in a hash
func (s *Store) HSet(key, field, value string) error {
	if err := s.TypeCheck(key, Hash); err != nil {
		return err
	}
	if _, exists := s.hashes[key]; !exists {
		s.hashes[key] = make(map[string]string)
		s.types[key] = Hash
	}
	s.hashes[key][field] = value
	return nil
}

// HGet retrieves a field value from a hash
func (s *Store) HGet(key, field string) (string, bool, error) {
	if err := s.TypeCheck(key, Hash); err != nil {
		return "", false, err
	}
	if hash, found := s.hashes[key]; found {
		value, found := hash[field]
		return value, found, nil
	}
	return "", false, nil
}
