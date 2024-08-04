// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package keyring

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
)

func init() {
	Register("file", initFileKeyring)
	Register("test", initTestKeyring)
}

const (
	keyringFileDirName = "keyring-file"
	keyringTestDirName = "keyring-test"
	keyringOSDirName   = "keyring-os"

	fileExtension = ".keyinfo"
)

var _ Keyring = (*fileKeyring)(nil)

var keyEncryptionAlgorithm = jwa.PBES2_HS512_A256KW

// fileKeyring is a keyring that stores keys in encrypted files.
type fileKeyring struct {
	// dir is the keystore root directory
	dir string
	// password is the user defined password used to generate encryption keys
	password []byte
	// prompt func is used to retrieve the user password
	prompt PromptFunc
}

// OpenFileKeyring opens the keyring in the given directory.
func OpenFileKeyring(dir string, prompt PromptFunc) (*fileKeyring, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &fileKeyring{
		dir:    dir,
		prompt: prompt,
	}, nil
}

func initFileKeyring(args ...any) (Keyring, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of args: %w", ErrInvalidArgs)
	}

	dir, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("bad string arg: %w", ErrInvalidArgs)
	}
	dir = filepath.Join(dir, keyringFileDirName)
	prompt := TerminalPrompt

	kr, err := OpenFileKeyring(dir, prompt)
	return kr, err
}

func initTestKeyring(args ...any) (Keyring, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	dir, ok := args[0].(string)
	if !ok {
		return nil, ErrInvalidArgs
	}
	dir = filepath.Join(dir, keyringTestDirName)
	prompt := FixedStringPrompt("secret")

	kr, err := OpenFileKeyring(dir, prompt)
	return kr, err
}

func (f *fileKeyring) Set(name string, key []byte) error {
	password, err := f.promptPassword()
	if err != nil {
		return err
	}
	cipher, err := jwe.Encrypt(key, jwe.WithKey(keyEncryptionAlgorithm, password))
	if err != nil {
		return err
	}
	return os.WriteFile(f.filepath(name), cipher, 0755)
}

func (f *fileKeyring) Get(name string) ([]byte, error) {
	cipher, err := os.ReadFile(f.filepath(name))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	password, err := f.promptPassword()
	if err != nil {
		return nil, err
	}
	return jwe.Decrypt(cipher, jwe.WithKey(keyEncryptionAlgorithm, password))
}

func (f *fileKeyring) Delete(name string) error {
	err := os.Remove(f.filepath(name))
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}

func (f *fileKeyring) List() ([]Info, error) {
	var infos []Info
	// walk the director and filter for our keyinfo files
	filepath.WalkDir(f.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == fileExtension {
			name := getFilename(d.Name())
			key, err := f.Get(name)
			if err != nil {
				return err
			}
			infos = append(infos, Info{
				Name: name,
				Key:  key,
			})
		}
		return nil
	})
	return infos, nil
}

// promptPassword returns the password from the user.
//
// If the password has been previously prompted it will be remembered.
func (f *fileKeyring) promptPassword() ([]byte, error) {
	if len(f.password) > 0 {
		return f.password, nil
	}
	password, err := f.prompt("Enter keystore password:")
	if err != nil {
		return nil, err
	}
	f.password = password
	return password, nil
}

func (f *fileKeyring) filepath(name string) string {
	return filepath.Join(f.dir, name+fileExtension)
}

func getFilename(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}
