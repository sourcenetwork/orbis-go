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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileKeystoreDirect(t *testing.T) {
	prompt := FixedStringPrompt("secret")

	kr, err := OpenFileKeystore(t.TempDir(), prompt)
	require.NoError(t, err)

	// seed the file keyring to resolve the password
	err = kr.Set("peer_key", []byte("abc"))
	require.NoError(t, err)

	// password should be remembered
	assert.Equal(t, []byte("secret"), kr.password)

	// clean the state so the testKeystore function
	// uses an empty keyring
	err = kr.Delete("peer_key")
	require.NoError(t, err)

	testKeystore(t, kr)
}

func TestBackendFileKeystoreOpen(t *testing.T) {
	kr, err := Open("file", t.TempDir())
	require.NoError(t, err)

	// get the file keyring struct from the interface
	// so we can stub the prompt func
	krf, ok := kr.(*fileKeystore)
	require.True(t, ok)
	krf.prompt = FixedStringPrompt("secret")

	testKeystore(t, kr)
}

func TestBackendTestKeystoreOpen(t *testing.T) {
	kr, err := Open("test", t.TempDir())
	require.NoError(t, err)

	testKeystore(t, kr)
}

func testKeystore(t *testing.T, kr Keystore) {
	err := kr.Set("peer_key", []byte("abc"))
	require.NoError(t, err)

	err = kr.Set("node_key", []byte("123"))
	require.NoError(t, err)

	peerKey, err := kr.Get("peer_key")
	require.NoError(t, err)
	assert.Equal(t, []byte("abc"), peerKey)

	nodeKey, err := kr.Get("node_key")
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), nodeKey)

	err = kr.Delete("node_key")
	require.NoError(t, err)

	_, err = kr.Get("node_key")
	assert.ErrorIs(t, err, ErrNotFound)

	// add another entry so there is more than 1 key
	// when calling List
	err = kr.Set("random_key", []byte("xyz"))
	require.NoError(t, err)

	infos, err := kr.List()
	require.NoError(t, err)
	require.Equal(t, []Info{
		{
			Name: "peer_key",
			Data: []byte("abc"),
		},
		{
			Name: "random_key",
			Data: []byte("xyz"),
		},
	}, infos)
}
