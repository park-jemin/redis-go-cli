package main

import (
	"strings"
	"testing"
)

// Test SET, GET, and DEL functionality for strings
func TestSetGetDel(t *testing.T) {
	store := NewStore()

	// Test SET
	if err := store.Set("key1", "value1"); err != nil {
		t.Fatalf("SET key1 failed: %v", err)
	}

	// Test GET
	if value, exists, err := store.Get("key1"); err != nil {
		t.Fatalf("GET key1 failed: %v", err)
	} else if !exists || value != "value1" {
		t.Errorf("Expected 'value1', got '%s'", value)
	}

	// Test DEL
	if deleted := store.Del("key1"); deleted != 1 {
		t.Fatalf("DEL key1 failed: expected 1 deleted, got %d", deleted)
	}

	// Test GET after DEL
	if _, exists, _ := store.Get("key1"); exists {
		t.Errorf("Expected key1 to be deleted")
	}
}

// Test overwriting an existing key
func TestSetOverwriteExistingKey(t *testing.T) {
	store := NewStore()

	// Set a key with a value
	if err := store.Set("key1", "value1"); err != nil {
		t.Fatalf("SET key1 failed: %v", err)
	}

	// Overwrite the key with a new value
	if err := store.Set("key1", "newValue"); err != nil {
		t.Fatalf("SET key1 (overwrite) failed: %v", err)
	}

	// Get the updated value
	value, exists, err := store.Get("key1")
	if err != nil || !exists || value != "newValue" {
		t.Errorf("Expected 'newValue', got '%s', exists: %v, error: %s", value, exists, err)
	}
}

// Test DEL with multiple keys, including non-existent ones
func TestDelMultipleKeys(t *testing.T) {
	store := NewStore()

	// Set multiple keys
	store.Set("key1", "value1")
	store.Set("key2", "value2")
	store.Set("key3", "value3")

	// Delete two existing keys and one non-existent key
	deleted := store.Del("key1", "key2", "nonExistentKey")

	// Ensure 2 keys were deleted
	if deleted != 2 {
		t.Errorf("Expected 2 keys deleted, got %d", deleted)
	}

	// Check the keys no longer exist
	if _, exists, _ := store.Get("key1"); exists {
		t.Errorf("Expected key1 to be deleted")
	}

	if _, exists, _ := store.Get("key2"); exists {
		t.Errorf("Expected key2 to be deleted")
	}

	// Check that key3 still exists
	if _, exists, _ := store.Get("key3"); !exists {
		t.Errorf("Expected key3 to still exist")
	}
}

// Test DEL then GET on deleted keys
func TestGetAfterDel(t *testing.T) {
	store := NewStore()

	if err := store.Set("key1", "value1"); err != nil {
		t.Fatalf("SET key1 failed: %v", err)
	}

	// Delete the key
	if deleted := store.Del("key1"); deleted != 1 {
		t.Fatalf("DEL key1 failed: expected 1 deletion, got %d", deleted)
	}

	// Try to GET the deleted key
	if _, exists, _ := store.Get("key1"); exists {
		t.Errorf("Expected key1 to be deleted")
	}
}

// Test SET same value multiple times
func TestSetSameValueMultipleTimes(t *testing.T) {
	store := NewStore()

	// Set key with the same value multiple times
	if err := store.Set("key1", "value1"); err != nil {
		t.Fatalf("SET key1 failed: %v", err)
	}

	// Set again with the same value
	if err := store.Set("key1", "value1"); err != nil {
		t.Fatalf("SET key1 (same value) failed: %v", err)
	}

	// Ensure the value remains the same
	value, exists, err := store.Get("key1")
	if err != nil || !exists || value != "value1" {
		t.Errorf("Expected 'value1', got '%s', exists: %v", value, exists)
	}
}

// Test LPush and LPop functionality
func TestLPushLPop(t *testing.T) {
	store := NewStore()

	// Test LPush
	count, err := store.LPush("list1", "one", "two", "three")
	if err != nil || count != 3 {
		t.Fatalf("LPUSH failed: %v, count: %d", err, count)
	}

	// Test LPop
	val, found, err := store.LPop("list1")
	if err != nil || !found || val != "three" {
		t.Fatalf("Expected 'three', got '%s', found: %v, err: %v", val, found, err)
	}

	// Continue popping
	val, _, _ = store.LPop("list1")
	if val != "two" {
		t.Errorf("Expected 'two', got '%s'", val)
	}

	val, _, _ = store.LPop("list1")
	if val != "one" {
		t.Errorf("Expected 'one', got '%s'", val)
	}

	// LPop on empty list
	_, found, _ = store.LPop("list1")
	if found {
		t.Errorf("Expected list to be empty")
	}
}

// Test LLen functionality
func TestLLen(t *testing.T) {
	store := NewStore()

	store.LPush("list1", "one", "two", "three")

	// Test LLen on non-empty list
	if length, err := store.LLen("list1"); err != nil || length != 3 {
		t.Fatalf("Expected length 3, got %d, error: %v", length, err)
	}

	// Test LLen on empty list
	store.LPush("list2", "four")
	store.LPop("list2")
	if length, err := store.LLen("list2"); err != nil || length != 0 {
		t.Fatalf("Expected length 0 for empty list, got %d, error: %v", length, err)
	}

	// Test LLen on non-existent key
	if length, err := store.LLen("nonExistent"); err != nil {
		t.Fatalf("Expected length 0 for nonexistent key, got %d, error: %v", length, err)
	}
}

// Test LRange functionality
func TestLRange(t *testing.T) {
	store := NewStore()

	store.LPush("list1", "one", "two", "three")

	// Test LRANGE within bounds
	result, err := store.LRange("list1", 0, 2)
	if err != nil {
		t.Fatalf("LRANGE failed: %v", err)
	}
	expected := []string{"three", "two", "one"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d elements, got %d", len(expected), len(result))
	}
	for i, val := range result {
		if val != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], val)
		}
	}

	// Test LRANGE with out-of-bounds indexes
	result, err = store.LRange("list1", -100, 100)
	if err != nil {
		t.Fatalf("LRANGE failed: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected all elements, got %d", len(result))
	}
}

// Test HSET and HGET functionality for hashes
func TestHSetHGet(t *testing.T) {
	store := NewStore()

	// Test HSET
	err := store.HSet("hash1", "field1", "value1")
	if err != nil {
		t.Fatalf("HSET failed: %v", err)
	}

	// Test HGET
	val, found, err := store.HGet("hash1", "field1")
	if err != nil || !found || val != "value1" {
		t.Fatalf("HGET failed: expected 'value1', got '%s'", val)
	}

	// HGET on non-existing field
	_, found, _ = store.HGet("hash1", "non_existent")
	if found {
		t.Errorf("Expected non-existing field to return false")
	}
}

// Test type-checking (WRONGTYPE error)
func TestWrongType(t *testing.T) {
	store := NewStore()

	// Set a key as string
	store.Set("string1", "value1")

	// Try to perform list operations on a string key
	_, err := store.LPush("string1", "invalid")
	if err == nil || strings.HasPrefix(err.Error(), "WRONGTYPE") {
		t.Errorf("Expected WRONGTYPE error, got %v", err)
	}

	// Set a key as a list
	store.LPush("list1", "value1")

	// Try to perform hash operations on a list key
	err = store.HSet("list1", "field1", "value1")
	if err == nil || strings.HasPrefix(err.Error(), "WRONGTYPE") {
		t.Errorf("Expected WRONGTYPE error, got %v", err)
	}

	// Set a key as a hash
	store.HSet("hash1", "field1", "value1")

	// Try to perform string operations on a hash key
	err = store.Set("hash1", "value1")

	if err == nil || strings.HasPrefix(err.Error(), "WRONGTYPE") {
		t.Errorf("Expected WRONGTYPE error, got %v", err)
	}
}
