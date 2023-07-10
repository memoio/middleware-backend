package api

import (
	"math/big"
	"time"
)

type ObjectInfo struct {
	Bucket      string
	Name        string
	Size        int64
	Cid         string
	ModTime     time.Time
	CType       string
	UserDefined map[string]string
}

type ObjectOptions struct {
	Size         int64
	Public       bool
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}

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

type BuyPackage struct {
	Pkgid     uint64
	Amount    int64
	Starttime uint64
	Chainid   string
}

type FileInfo struct {
	ID         int         `gorm:"primarykey"`
	ChainID    int         `gorm:"uniqueIndex:composite;column:chainid"`
	Address    string      `gorm:"uniqueIndex:composite"`
	SType      StorageType `gorm:"uniqueIndex:composite;column:stype"`
	Mid        string      `gorm:"uniqueIndex:composite"`
	Name       string      `gorm:"uniqueIndex:composite"`
	Size       int64
	ModTime    time.Time `gorm:"column:modtime"`
	Public     bool
	UserDefine string `gorm:"column:userdefine"`
}

func (FileInfo) TableName() string {
	return "fileinfo"
}
