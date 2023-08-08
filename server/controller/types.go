package controller

import (
	"math/big"
	"time"

	"github.com/memoio/backend/api"
)

type ObjectOptions api.ObjectOptions

type PutObjectResult struct {
	Mid string
}

type GetObjectResult struct {
	Name  string
	Size  int64
	CType string
}

type ListObjectsResult struct {
	Address string
	Storage string
	Objects []ObjectInfoResult
}

type ObjectInfoResult struct {
	ID      int
	Name    string
	Size    int64
	Mid     string
	Public  bool
	ModTime time.Time
	// UserDefined map[string]string
}

type IPayPayment struct {
	Nonce    *big.Int
	Balance  *big.Int
	SizeByte uint64
	FreeByte uint64
	Expire   uint64
}
