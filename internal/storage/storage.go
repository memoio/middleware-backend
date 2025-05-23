package storage

import (
	"math/big"
)

type StorageInfo struct {
	Storage string
	Used    int64
	Buysize int64
	Free    int64
	Files   int
}

type StorageType uint8

const (
	MEFS StorageType = iota
	IPFS
	QINIU
)

func (s StorageType) String() string {
	switch s {
	case MEFS:
		return "mefs"
	case IPFS:
		return "ipfs"
	case QINIU:
		return "qiniu"
	default:
		return "unknow storage"
	}
}

func StringToStorageType(s string) StorageType {
	storage := new(big.Int)
	storage.SetString(s, 10)
	return StorageType(storage.Uint64())
}

func Uint8ToStorageType(s uint8) StorageType {
	return StorageType(s)
}
