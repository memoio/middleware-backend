package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/memoio/backend/api"
)

type keyRepo struct {
	path string
}

type Key struct {
	Address     string
	SecretValue []byte
}

// todo: create a file for verifing password before reopen
func NewKeyRepo(path string) (api.Keystore, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	// create a file for test password

	return &keyRepo{
		path,
	}, nil
}

func (k keyRepo) Put(name string, ki []byte) error {
	key := &Key{
		Address:     name,
		SecretValue: ki,
	}

	keyjson, err := encryptKey(key)
	if err != nil {
		return err
	}

	path := joinPath(k.path, name)
	_, err = os.Stat(path)
	if err == nil {
		// exist
		return nil
	}

	return writeKeyFile(path, keyjson)
}

func (k *keyRepo) Get(name string) ([]byte, error) {
	var res []byte
	path := joinPath(k.path, name)

	keyjson, err := ioutil.ReadFile(path)
	if err != nil {
		return res, err
	}

	key, err := decryptKey(keyjson)
	if err != nil {
		return res, err
	}

	if strings.Compare(key.Address, name) != 0 {
		return res, fmt.Errorf("key content mismatch: have peer %x, want %x", key.Address, name)
	}

	err = json.Unmarshal(key.SecretValue, &res)
	if err != nil {
		return res, fmt.Errorf("decoding key '%s': %w", name, err)
	}

	return res, nil
}

func (k *keyRepo) List() ([]string, error) {
	dir, err := os.Open(k.path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(files))
	for _, f := range files {
		keys = append(keys, string(f.Name()))
	}
	return keys, nil
}

func (k *keyRepo) Delete(name string) error {
	_, err := k.Get(name)
	if err != nil {
		return err
	}

	keyPath := joinPath(k.path, name)

	err = os.Remove(keyPath)
	if err != nil {
		return err
	}

	return nil
}
func joinPath(dir string, filename string) (path string) {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(dir, filename)
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	err := os.MkdirAll(filepath.Dir(file), dirPerm)
	if err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func writeKeyFile(file string, content []byte) error {
	name, err := writeTemporaryKeyFile(file, content)
	if err != nil {
		return err
	}
	return os.Rename(name, file)
}
