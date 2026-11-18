package node

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func LoadOrCreateIdentity(path string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(path)

	if err == nil {
		return crypto.UnmarshalPrivateKey(data)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	key, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	data, err = crypto.MarshalPrivateKey(key)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return nil, err
	}

	return key, nil
}
