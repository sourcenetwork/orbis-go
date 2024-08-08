// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package keystore

import (
	"sync"
)

var (
	mu       sync.Mutex
	keyrings map[string]InitFunc = map[string]InitFunc{}
)

type Info struct {
	Name string
	Data []byte
}

/*

keyring := keyring.Open(...)
keyring.Set("alice", key)

*/

// Keystore provides a simple set/get interface for a keyring service.
type Keystore interface {
	// Set stores the given key in the keystore under the given name.
	//
	// If a key with the given name already exists it will be overriden.
	Set(name string, key []byte) error
	// Get returns the key with the given name from the keystore.
	//
	// If a key with the given name does not exist `ErrNotFound` is returned.
	Get(name string) ([]byte, error)
	// Delete removes the key with the given name from the keystore.
	//
	// If a key with that name does not exist `ErrNotFound` is returned.
	Delete(name string) error
	// List all keys in the Keystore, only public information
	List() ([]Info, error)
}

type InitFunc func(args ...any) (Keystore, error)

func Register(backend string, initFunc InitFunc) error {
	mu.Lock()
	defer mu.Unlock()

	if backend == "" {
		return ErrInvalidBackend
	}
	if initFunc == nil {
		return ErrInvalidInitFunc
	}
	_, exist := keyrings[backend]
	if exist {
		return ErrBackendExists
	}

	keyrings[backend] = initFunc
	return nil
}

func Open(backend string, args ...any) (Keystore, error) {
	mu.Lock()
	defer mu.Unlock()

	if backend == "" {
		return nil, ErrInvalidBackend
	}
	initFn, exists := keyrings[backend]
	if !exists {
		return nil, ErrBackendNotFound
	}

	return initFn(args...)
}
