// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/testcontext"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/uplink"
)

func TestLoadEncryptionKeyIntoEncryptionAccess(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		passphrase := make([]byte, rand.Intn(100)+1)
		_, err := rand.Read(passphrase)
		require.NoError(t, err)

		expectedKey, err := storj.NewKey(passphrase)
		require.NoError(t, err)
		ctx := testcontext.New(t)
		filename := ctx.File("encryption.key")
		err = ioutil.WriteFile(filename, expectedKey[:], os.FileMode(0400))
		require.NoError(t, err)
		defer ctx.Cleanup()

		access, err := loadEncryptionAccess(filename)
		require.NoError(t, err)
		require.Equal(t, *expectedKey, access.Key)
	})

	t.Run("error", func(t *testing.T) {
		ctx := testcontext.New(t)
		defer ctx.Cleanup()
		filename := ctx.File("encryption.key")

		_, err := loadEncryptionAccess(filename)
		require.Error(t, err)
	})
}

func TestSaveLoadEncryptionKey(t *testing.T) {
	var inputKey string
	{
		randKey := make([]byte, rand.Intn(storj.KeySize)*3+1)
		_, err := rand.Read(randKey)
		require.NoError(t, err)
		inputKey = string(randKey)
	}

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	filename := ctx.File("storj-test-cmd-uplink", "encryption.key")
	err := uplink.SaveEncryptionKey(inputKey, filename)
	require.NoError(t, err)

	access, err := loadEncryptionAccess(filename)
	require.NoError(t, err)

	if len(inputKey) > storj.KeySize {
		require.Equal(t, []byte(inputKey[:storj.KeySize]), access.Key[:])
	} else {
		require.Equal(t, []byte(inputKey), access.Key[:len(inputKey)])
	}
}
