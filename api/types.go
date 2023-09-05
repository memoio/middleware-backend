package api

import (
	"math/big"
	"time"

	bls12377 "github.com/consensys/gnark-crypto/ecc/bls12-377"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr/kzg"
	"github.com/ethereum/go-ethereum/common"
)

type ObjectInfo struct {
	SType       StorageType
	USerID      int
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
	Sign         string
	Area         string
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}

type SignMessage struct {
	// StorePayAddr common.Address
	Size uint64
	// Nonce        *big.Int
	Sign string
}

type CheckInfo struct {
	Buyer    common.Address
	FileSize *big.Int
	Nonce    *big.Int
	Sign     []byte
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
	UserID     int    `gorm:"column:userid"`
}

func (FileInfo) TableName() string {
	return "fileinfo"
}

type USerInfo struct {
	ID    int    `gorm:"primarykey"`
	Area  string `gorm:"uniqueIndex:composite;column:area"`
	Api   string `gorm:"uniqueIndex:composite;column:api"`
	Token string `gorm:"uniqueIndex:composite;column:token"`
}

func (USerInfo) TableName() string {
	return "userinfo"
}

type PayType uint8

const (
	StorePay PayType = iota
	ReadPay
)

func StringToPayType(s string) PayType {
	switch s {
	case "space":
		return ReadPay
	case "traffic":
		return StorePay
	}
	return StorePay
}

type G1 = bls12377.G1Affine
type G2 = bls12377.G2Affine
type GT = bls12377.GT
type Fr = fr.Element

type Proof = kzg.OpeningProof
