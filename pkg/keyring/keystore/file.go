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
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
)

func init() {
	Register("file", initFileKeystore)
	Register("test", initTestKeystore)
}

const (
	keyringFileDirName = "keyring-file"
	keyringTestDirName = "keyring-test"
	keyringOSDirName   = "keyring-os"

	fileExtension = ".keyinfo"
)

var _ Keystore = (*fileKeystore)(nil)

var keyEncryptionAlgorithm = jwa.PBES2_HS512_A256KW

// fileKeystore is a keyring that stores keys in encrypted files.
type fileKeystore struct {
	// dir is the keystore root directory
	dir string
	// password is the user defined password used to generate encryption keys
	password []byte
	// prompt func is used to retrieve the user password
	prompt PromptFunc
	// basic mutex
	mu sync.Mutex
}

// OpenFileKeystore opens the keyring in the given directory.
func OpenFileKeystore(dir string, prompt PromptFunc) (*fileKeystore, error) {
	dir = os.ExpandEnv(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &fileKeystore{
		dir:    dir,
		prompt: prompt,
	}, nil
}

func initFileKeystore(args ...any) (Keystore, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of args: %w", ErrInvalidArgs)
	}

	dir, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("bad string arg: %w", ErrInvalidArgs)
	}
	dir = filepath.Join(dir, keyringFileDirName)
	prompt := TerminalPrompt

	kr, err := OpenFileKeystore(dir, prompt)
	return kr, err
}

func initTestKeystore(args ...any) (Keystore, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	dir, ok := args[0].(string)
	if !ok {
		return nil, ErrInvalidArgs
	}
	dir = filepath.Join(dir, keyringTestDirName)
	prompt := FixedStringPrompt("secret")

	kr, err := OpenFileKeystore(dir, prompt)
	return kr, err
}

func (f *fileKeystore) Set(name string, key []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	password, err := f.promptPassword()
	if err != nil {
		return err
	}

	return f.set(name, key, password)
}

func (f *fileKeystore) set(name string, key []byte, password []byte) error {
	cipher, err := jwe.Encrypt(key, jwe.WithKey(keyEncryptionAlgorithm, password))
	if err != nil {
		return err
	}
	return os.WriteFile(f.filepath(name), cipher, 0755)
}

func (f *fileKeystore) Get(name string) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	password, err := f.promptPassword()
	if err != nil {
		return nil, err
	}

	return f.get(name, password)
}

func (f *fileKeystore) get(name string, password []byte) ([]byte, error) {
	cipher, err := os.ReadFile(f.filepath(name))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	return jwe.Decrypt(cipher, jwe.WithKey(keyEncryptionAlgorithm, password))
}

func (f *fileKeystore) Delete(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// gate the action by the password
	_, err := f.promptPassword()
	if err != nil {
		return err
	}

	return f.delete(name)
}

func (f *fileKeystore) delete(name string) error {
	err := os.Remove(f.filepath(name))
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}

func (f *fileKeystore) List() ([]Info, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	password, err := f.promptPassword()
	if err != nil {
		return nil, err
	}

	var infos []Info
	// walk the director and filter for our keyinfo files
	filepath.WalkDir(f.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == fileExtension {
			name := getFilename(d.Name())
			key, err := f.get(name, password)
			if err != nil {
				return err
			}
			infos = append(infos, Info{
				Name: name,
				Data: key,
			})
		}
		return nil
	})
	return infos, nil
}

// promptPassword returns the password from the user.
//
// If the password has been previously prompted it will be remembered.
func (f *fileKeystore) promptPassword() ([]byte, error) {
	if len(f.password) > 0 {
		return f.password, nil
	}
	password, err := f.prompt("Enter keystore password")
	if err != nil {
		return nil, err
	}

	keyhashFileName := filepath.Join(f.dir, "keyhash")
	// confirm
	if _, err := os.Stat(keyhashFileName); errors.Is(err, os.ErrNotExist) {
		err = f.confirmPassword(password, 1)
		if err != nil {
			return nil, err
		}
	} else if err == nil {
		keyhashFileName := filepath.Join(f.dir, "keyhash")
		keyhashOrig, err := os.ReadFile(keyhashFileName)
		if err != nil {
			return nil, err
		}
		err = f.verifyPassword(password, keyhashOrig, 2)
		if err != nil {
			return nil, err
		}
	} else {
		// original os stat err
		return nil, err
	}

	f.password = password
	return password, nil
}

func (f *fileKeystore) confirmPassword(password []byte, attempts int) error {
	keyhashFileName := filepath.Join(f.dir, "keyhash")
	if attempts > 3 {
		return fmt.Errorf("too many attempts")
	}
	passwordConfirm, err := f.prompt("Confirm keystore password")
	if err != nil {
		return err
	}

	if string(passwordConfirm) != string(password) {
		attempts += 1
		return f.confirmPassword(password, attempts)
	}

	// create keyhash file
	keyhash := hash(password)
	return os.WriteFile(keyhashFileName, keyhash, 0600)
}

func (f *fileKeystore) verifyPassword(password []byte, pwhash []byte, attempts int) error {
	keyhash := hash(password)
	if !bytes.Equal(keyhash, pwhash) {
		if attempts > 3 {
			return fmt.Errorf("too many attempts")
		}
		password, err := f.prompt(fmt.Sprintf("Enter keystore password (attempt: %d/%d)", attempts, 3))
		if err != nil {
			return err
		}
		return f.verifyPassword(password, pwhash, attempts+1)
	}

	return nil
}

func hash(data []byte) []byte {
	hash := sha256.New()
	return hash.Sum(data)
}

func (f *fileKeystore) filepath(name string) string {
	path := filepath.Join(f.dir, name+fileExtension)
	return path
}

func getFilename(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}
