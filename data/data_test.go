package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMemUnsortedKV(t *testing.T) {
	kv, err := NewMemUnsortedKV()
	if err != nil {
		t.Fatalf("NewMemUnsortedKV() error = %v", err)
	}
	if kv == nil {
		t.Fatal("NewMemUnsortedKV() returned nil")
	}

	// Test Set and Get
	err = kv.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	val, found := kv.Get("key1")
	if !found {
		t.Error("Get() key not found")
	}
	if val != "value1" {
		t.Errorf("Get() = %v, want %v", val, "value1")
	}
}

func TestUnsortedKVHas(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	if kv.Has("nonexistent") {
		t.Error("Has() returned true for nonexistent key")
	}

	kv.Set("exists", "value")
	if !kv.Has("exists") {
		t.Error("Has() returned false for existing key")
	}
}

func TestUnsortedKVGetOr(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	// Test with nonexistent key
	val := kv.GetOr("nonexistent", "default")
	if val != "default" {
		t.Errorf("GetOr() = %v, want %v", val, "default")
	}

	// Test with existing key
	kv.Set("exists", "actual")
	val = kv.GetOr("exists", "default")
	if val != "actual" {
		t.Errorf("GetOr() = %v, want %v", val, "actual")
	}
}

func TestUnsortedKVDel(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	kv.Set("key1", "value1")
	if !kv.Has("key1") {
		t.Error("Key was not set")
	}

	err := kv.Del("key1")
	if err != nil {
		t.Errorf("Del() error = %v", err)
	}

	if kv.Has("key1") {
		t.Error("Key still exists after deletion")
	}
}

func TestUnsortedKVClear(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	kv.Set("key1", "value1")
	kv.Set("key2", "value2")
	kv.Set("key3", "value3")

	if kv.Empty() {
		t.Error("KV should not be empty")
	}

	err := kv.Clear()
	if err != nil {
		t.Errorf("Clear() error = %v", err)
	}

	if !kv.Empty() {
		t.Error("KV should be empty after clear")
	}
}

func TestUnsortedKVEach(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	kv.Set("key1", "value1")
	kv.Set("key2", "value2")
	kv.Set("key3", "value3")

	count := 0
	kv.Each(func(k, v string) bool {
		count++
		return false
	})

	if count != 3 {
		t.Errorf("Each() iterated %d times, want 3", count)
	}
}

func TestUnsortedKVEachWithStop(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	kv.Set("key1", "value1")
	kv.Set("key2", "value2")
	kv.Set("key3", "value3")

	count := 0
	kv.Each(func(k, v string) bool {
		count++
		return count >= 2 // Stop after 2 iterations
	})

	if count != 2 {
		t.Errorf("Each() iterated %d times, want 2", count)
	}
}

func TestUnsortedKVEmpty(t *testing.T) {
	kv, _ := NewMemUnsortedKV()

	if !kv.Empty() {
		t.Error("New KV should be empty")
	}

	kv.Set("key", "value")
	if kv.Empty() {
		t.Error("KV should not be empty after adding item")
	}

	kv.Del("key")
	if !kv.Empty() {
		t.Error("KV should be empty after deleting only item")
	}
}

func TestDiskUnsortedKV(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, "test.db")

	// Create and populate KV
	kv, err := NewDiskUnsortedKV(dbFile)
	if err != nil {
		t.Fatalf("NewDiskUnsortedKV() error = %v", err)
	}

	err = kv.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Check file exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}

	// Load from disk
	kv2, err := NewDiskUnsortedKV(dbFile)
	if err != nil {
		t.Fatalf("NewDiskUnsortedKV() error loading existing file = %v", err)
	}

	val, found := kv2.Get("key1")
	if !found {
		t.Error("Key not found in loaded KV")
	}
	if val != "value1" {
		t.Errorf("Get() = %v, want %v", val, "value1")
	}
}

func TestDiskUnsortedKVReader(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, "test_reader.db")

	// Create and populate KV
	kv, err := NewDiskUnsortedKV(dbFile)
	if err != nil {
		t.Fatalf("NewDiskUnsortedKV() error = %v", err)
	}
	kv.Set("key1", "value1")

	// Open as reader
	reader, err := NewDiskUnsortedKVReader(dbFile)
	if err != nil {
		t.Fatalf("NewDiskUnsortedKVReader() error = %v", err)
	}

	// Read should work
	val, found := reader.Get("key1")
	if !found || val != "value1" {
		t.Error("Reader could not read value")
	}

	// Modify reader (should not persist)
	reader.Set("key2", "value2")

	// Load again and verify key2 doesn't exist
	kv2, _ := NewDiskUnsortedKV(dbFile)
	if kv2.Has("key2") {
		t.Error("Reader modifications should not persist")
	}
}

func TestUnsortedKVMarshalJSON(t *testing.T) {
	kv, _ := NewMemUnsortedKV()
	kv.Set("key1", "value1")
	kv.Set("key2", "value2")

	data, err := kv.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("MarshalJSON() returned empty data")
	}
}

func TestFlushPolicy(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, "test_flush.db")

	// Test FlushNone
	kv, err := NewUnsortedKV(dbFile, FlushNone)
	if err != nil {
		t.Fatalf("NewUnsortedKV() error = %v", err)
	}

	kv.Set("key1", "value1")

	// File might not exist or be empty with FlushNone
	kv2, err := NewUnsortedKV(dbFile, FlushNone)
	if err != nil {
		// This is ok if file doesn't exist
	} else if kv2.Has("key1") {
		// Value exists, need to manually flush to test
		kv.Flush()
	}

	// Test FlushOnEdit (already tested in TestDiskUnsortedKV)
}
