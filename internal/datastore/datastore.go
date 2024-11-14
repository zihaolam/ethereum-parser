package datastore

type DataStore interface {
	// Returns all keys in the datastore.
	List() ([]string, error)

	// Checks if a key exists in the data store
	Has(key string) bool

	// Retrieves the value of the given key.
	// Returns error if key does not exist
	Get(key string) ([][]byte, error)

	// Stores a given value into the datastore.
	// Values are stored as slices of slices of bytes so that multiple values can be stored under the same key.
	Put(key string, value [][]byte) error

	// Updates the value of the given key using the provided update function.
	Update(key string, updater func([][]byte) ([][]byte, error)) error

	// Deletes key value pair from the datastore. Safely returns even if key does not exist.
	Delete(key string) error
}
