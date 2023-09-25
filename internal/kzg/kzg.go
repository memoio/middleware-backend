package kzg

import (
	"os"

	mkzg "github.com/memoio/middleware/kzg"
)

var kzg_path = "./data/kzg.txt"

func NewKzg() (*mkzg.PublicKey, error) {
	fi, err := os.ReadFile(kzg_path)
	if err != nil {
		return nil, err
	}

	pk := &mkzg.PublicKey{}
	err = pk.Unmarshal(fi)
	if err != nil {
		return nil, err
	}

	return pk, nil
}
