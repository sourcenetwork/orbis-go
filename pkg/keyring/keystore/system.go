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
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

func init() {
	Register("os", initSystemKeystore)
}

var _ Keystore = (*systemKeystore)(nil)

// systemKeystore is a keyring that utilizies the
// built in key management system of the OS.
type systemKeystore struct {
	// service is the service name to use when using the system keyring
	service string
	// index uses the fileKeystore to create an index of existing keys
	// because the systemKeystore implementation doesn't provide any
	// List/Query functionality. TODO: Native List functions
	index *fileKeystore
}

// OpenSystemKeystore opens the system keyring managed by the OS.
// dir is the path to the index
// service is the system store prefix
func OpenSystemKeystore(dir string, service string) (*systemKeystore, error) {
	// the file keyring is just used as an index, and doesn't store
	// the actual private value, so we can use a FixedString password
	fk, err := OpenFileKeystore(dir, FixedStringPrompt("secret"))
	if err != nil {
		return nil, err
	}
	return &systemKeystore{
		service: service,
		index:   fk,
	}, nil
}

func initSystemKeystore(args ...any) (Keystore, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgs
	}

	dir, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("bad string arg: %w", ErrInvalidArgs)
	}
	dir = filepath.Join(dir, keyringOSDirName)

	service, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidArgs
	}

	return OpenSystemKeystore(dir, service)
}

func (s *systemKeystore) Set(name string, key []byte) error {
	enc := base64.StdEncoding.EncodeToString(key)
	err := s.index.Set(name, []byte(name))
	if err != nil {
		return err
	}
	err = keyring.Set(s.service, name, enc)
	if err != nil {
		// cleanup the index entry we just created
		err = s.index.Delete(name)
		if err != nil { // if this happens we're in a bad state (index will be corrupt)
			panic("failed to maintain os keyring index") // this shouldn't really happen tho.
		}
	}
	return nil
}

func (s *systemKeystore) Get(name string) ([]byte, error) {
	enc, err := keyring.Get(s.service, name)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(enc)))
	n, err := base64.StdEncoding.Decode(dst, []byte(enc))
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

func (s *systemKeystore) Delete(user string) error {
	err := keyring.Delete(s.service, user)
	if err != nil {
		return err
	}
	// if this fails the index will be in a bad state
	// but this is more forgivable.
	return s.index.Delete(user)
}

func (s *systemKeystore) List() ([]Info, error) {
	indexInfos, err := s.index.List()
	if err != nil {
		return nil, fmt.Errorf("couldn't get index list: %w", err)
	}

	var infos []Info
	for _, info := range indexInfos {
		key, err := s.Get(info.Name)
		if err != nil {
			return nil, fmt.Errorf("couldn't get key: %w", err)
		}
		infos = append(infos, Info{
			Name: info.Name,
			Data: key,
		})
	}
	return infos, nil
}
