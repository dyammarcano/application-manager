package cache

import (
	"github.com/dgraph-io/badger/v3"
	"sync"
)

type (
	processItem func(item *badger.Item) error

	V3Cache struct {
		db *badger.DB
		wg sync.WaitGroup
	}
)

// NewCache creates a new Badger database
func NewCache(path string) (*V3Cache, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	return &V3Cache{
		db: db,
		wg: sync.WaitGroup{},
	}, nil
}

// Close the Badger database
func (c *V3Cache) Close() error {
	return c.db.Close()
}

// Get a value from the Badger database
func (c *V3Cache) Get(key string) (string, error) {
	var value string
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
		return err
	})
	return value, err
}

// Set a value in the Badger database
func (c *V3Cache) Set(key string, value string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

// Delete a value from the Badger database
func (c *V3Cache) Delete(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (c *V3Cache) iterateDB(process processItem) error {
	return c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if err := process(item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *V3Cache) DeleteAll() error {
	return c.db.DropAll()
}

func (c *V3Cache) GetKeys() ([]string, error) {
	keys := make([]string, 0)
	err := c.iterateDB(func(item *badger.Item) error {
		key := string(item.Key())
		keys = append(keys, key)
		return nil
	})
	return keys, err
}

func (c *V3Cache) Length() (int, error) {
	var length int
	err := c.iterateDB(func(item *badger.Item) error {
		length++
		return nil
	})
	return length, err
}

func (c *V3Cache) Size() int {
	return len(c.db.Tables())
}

func (c *V3Cache) GetAll() (map[string]string, error) {
	values := make(map[string]string)
	err := c.iterateDB(func(item *badger.Item) error {
		var value string
		err := item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
		key := string(item.Key())
		values[key] = value

		return err
	})
	return values, err
}
