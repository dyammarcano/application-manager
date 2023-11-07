package cache

import (
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

//func TestNewCache(t *testing.T) {
//	dir, err := os.MkdirTemp("", "test")
//	if err != nil {
//		t.Errorf("Error creating temp dir: %v", err)
//	}
//
//	cache, err := NewCache(dir)
//	if err != nil {
//		t.Errorf("Error creating cache: %v", err)
//	}
//
//	defer func() {
//		if err := cache.Close(); err != nil {
//			t.Errorf("Error closing cache")
//		}
//
//		<-time.After(1 * time.Second)
//
//		if err := os.RemoveAll(dir); err != nil {
//			t.Errorf("Error removing temp dir: %v", err)
//		}
//	}()
//
//	if err := cache.Set("test", "test"); err != nil {
//		t.Errorf("Error setting cache: %v", err)
//	}
//
//	if value, err := cache.Get("test"); err != nil {
//		t.Errorf("Error getting cache: %v", err)
//	} else {
//		if value != "test" {
//			t.Errorf("Value is different")
//		}
//	}
//
//	if err := cache.Delete("test"); err != nil {
//		t.Errorf("Error deleting cache: %v", err)
//	}
//
//	value := 1000000
//
//	for i := 0; i < value; i++ {
//		go func() {
//			if err := cache.Set(ulid.Make().String(), "test"); err != nil {
//				t.Errorf("Error setting cache: %v", err)
//			}
//		}()
//	}
//
//	<-time.After(1 * time.Second)
//
//	if records, err := cache.GetAll(); err != nil {
//		t.Errorf("Error getting cache: %v", err)
//	} else {
//		if len(records) != value {
//			t.Errorf("Value is different")
//		}
//	}
//
//	length, _ := cache.Length()
//	t.Logf("V3Cache length: %d", length)
//}

func TestNewCache1_000_000(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	assert.Nil(t, err)

	cache, err := NewCache(dir)
	assert.Nil(t, err)

	defer func() {
		err := cache.Close()
		assert.Nil(t, err)
	}()

	value := 1000000

	wg := sync.WaitGroup{}

	for i := 0; i < value; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)

			err := cache.Set(ulid.Make().String(), "test")
			assert.Nil(t, err)
		}()
	}

	wg.Wait()

	records, err := cache.GetAll()
	assert.Nil(t, err)

	assert.Equal(t, value, len(records))
}
