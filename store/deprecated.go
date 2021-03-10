// Package store contains deprecated aliases for subpackages
package store // import "github.com/throttled/throttled/v2/store"

import (
	"github.com/throttled/throttled/v2/store/memstore"
)

// NewMemStore initializes a new memory-based store.
//
// Deprecated: Use github.com/throttled/throttled/v2/store/memstore instead.
func NewMemStore(maxKeys int) *memstore.MemStore {
	st, err := memstore.New(maxKeys)
	if err != nil {
		// As of this writing, `lru.New` can only return an error if you pass
		// maxKeys <= 0 so this should never occur.
		panic(err)
	}
	return st
}

