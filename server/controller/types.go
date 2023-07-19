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

type packageInfo struct {
	Time    uint64
	Kind    uint8
	Buysize *big.Int
	Amount  *big.Int
	State   uint8
}

type packageInfos struct {
	Pkgid int
	packageInfo
}

type userBuyPackage struct {
	Starttime uint64
	Endtime   uint64
	Kind      uint8
	Buysize   *big.Int
	Amount    *big.Int
	State     uint8
}

type flowSize struct {
	Used *big.Int
	Free *big.Int
}

type IPayPayment struct {
	Nonce    *big.Int
	Balance  *big.Int
	SizeByte uint64
	FreeByte uint64
	Expire   uint64
}
