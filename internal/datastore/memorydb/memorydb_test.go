package memorydb_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/zihaolam/ethereum-parser/internal/datastore/memorydb"
)

func TestMain(t *testing.T) {
	var db *memorydb.MemoryDB
	t.Run("New", func(t *testing.T) {
		db = memorydb.New()
		if db == nil {
			t.Fatal("expected a non-nil MemoryDB instance")
		}
	})

	t.Run("Has", func(t *testing.T) {
		key := "testKey"

		if db.Has(key) {
			t.Fatal("expected Has to return false for a non-existent key")
		}

		db.Put(key, [][]byte{[]byte("value")})

		if !db.Has(key) {
			t.Fatal("expected Has to return true for an existing key")
		}
	})

	t.Run("Get", func(t *testing.T) {
		key := "testKey2"
		value := [][]byte{[]byte("value")}

		_, err := db.Get(key)
		if !errors.Is(err, memorydb.KeyDoesNotExist) {
			t.Fatal("expected KeyDoesNotExist error for a non-existent key")
		}

		db.Put(key, value)

		gotValue, err := db.Get(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(gotValue[0], value[0]) {
			t.Fatalf("expected %v, got %v", value, gotValue)
		}
	})

	t.Run("Put", func(t *testing.T) {
		key := "testKey"
		value := [][]byte{[]byte("value")}

		err := db.Put(key, value)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		gotValue, err := db.Get(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(gotValue[0], value[0]) {
			t.Fatalf("expected %v, got %v", value, gotValue)
		}
	})

	t.Run("Update", func(t *testing.T) {
		key := "testKey"
		initialValue := [][]byte{[]byte("initial")}
		updatedValue := [][]byte{[]byte("updated")}

		err := db.Put(key, initialValue)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = db.Update(key, func(val [][]byte) ([][]byte, error) {
			return updatedValue, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		gotValue, err := db.Get(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Equal(gotValue[0], updatedValue[0]) {
			t.Fatalf("expected %v, got %v", updatedValue, gotValue)
		}

		// Test Update on non-existent key
		err = db.Update("nonExistentKey", func(val [][]byte) ([][]byte, error) {
			return nil, errors.New("error updating")
		})
		if !errors.Is(err, memorydb.KeyDoesNotExist) {
			t.Fatal("expected KeyDoesNotExist error for non-existent key")
		}
	})

	t.Run("List", func(t *testing.T) {
		newdb := memorydb.New()
		keys, err := newdb.List()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(keys) != 0 {
			t.Fatalf("expected no keys, got %d", len(keys))
		}

		newdb.Put("key1", [][]byte{[]byte("value1")})
		newdb.Put("key2", [][]byte{[]byte("value2")})

		keys, err = newdb.List()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(keys) != 2 {
			t.Fatalf("expected 2 keys, got %d", len(keys))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		key := "testKey"
		value := [][]byte{[]byte("value")}

		db.Put(key, value)

		err := db.Delete(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = db.Get(key)
		if !errors.Is(err, memorydb.KeyDoesNotExist) {
			t.Fatal("expected KeyDoesNotExist error after deleting key")
		}
	})
}
