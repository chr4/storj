// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"runtime"

	"github.com/zeebo/errs"
	"golang.org/x/crypto/argon2"

	"storj.io/storj/pkg/storj"
)

func sha256hmac(key, data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// DeriveRootKey derives a root key for some path using the salt for the bucket and
// a password from the user. See the password key derivation design doc.
func DeriveRootKey(password, salt []byte, path storj.Path) (*storj.Key, error) {
	mixedSalt, err := sha256hmac(password, salt)
	if err != nil {
		return nil, err
	}

	pathSalt, err := sha256hmac(mixedSalt, []byte(path))
	if err != nil {
		return nil, err
	}

	// use a time of 1, 64MB of ram, and all of the cores.
	keyData := argon2.IDKey(password, pathSalt, 1, 64<<10, uint8(runtime.GOMAXPROCS(-1)), 32)
	if len(keyData) != len(storj.Key{}) {
		return nil, errs.New("invalid output from argon2id")
	}

	var key storj.Key
	copy(key[:], keyData)
	return &key, nil
}

// DeriveDefaultPassword combines a salt from the project with a user password to create
// a default password used for DeriveRootKey in the case that a single secret should be
// used to derive all of the bucket level secrets. See the password key derivation design
// doc.
func DeriveDefaultPassword(password, salt []byte) ([]byte, error) {
	mixedSalt, err := sha256hmac(password, salt)
	if err != nil {
		return nil, err
	}

	// use a time of 1, 64MB of ram, and all of the cores.
	return argon2.IDKey(password, mixedSalt, 1, 64<<10, uint8(runtime.GOMAXPROCS(-1)), 32), nil
}
