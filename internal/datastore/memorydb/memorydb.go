package memorydb

import (
	"errors"
	"sync"
)

var KeyDoesNotExist = errors.New("key does not exist")

type MemoryDB struct {
	data  map[string][][]byte
	mutex sync.RWMutex
}

func New() *MemoryDB {
	return &MemoryDB{
		data:  make(map[string][][]byte),
		mutex: sync.RWMutex{},
	}
}

func (db *MemoryDB) Has(key string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	_, ok := db.data[key]
	return ok
}

func (db *MemoryDB) Get(key string) ([][]byte, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	v, ok := db.data[key]
	if !ok {
		return nil, KeyDoesNotExist
	}
	return v, nil
}

func (db *MemoryDB) Put(key string, value [][]byte) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.data[key] = value
	return nil
}

func (db *MemoryDB) Update(key string, updater func([][]byte) ([][]byte, error)) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	v, ok := db.data[key]
	if !ok {
		return KeyDoesNotExist
	}

	newValue, err := updater(v)
	if err != nil {
		return err
	}
	db.data[key] = newValue
	return nil
}

func (db *MemoryDB) List() ([]string, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (db *MemoryDB) Delete(key string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	delete(db.data, key)
	return nil
}
