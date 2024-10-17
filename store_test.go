package main

import "testing"

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
