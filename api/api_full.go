package api

import (
	"context"
	"io"
	"math/big"
)

type IGateway interface {
	PutObject(context.Context, string, string, io.Reader, ObjectOptions) (ObjectInfo, error)
	GetObject(context.Context, string, io.Writer, ObjectOptions) error
	DeleteObject(context.Context, string, string) error
	// ListObjects(context.Context, string) ([]ObjectInfo, error)
	// GetObjectInfo(context.Context, string) (ObjectInfo, error)
}

type IContract interface {
	Call(ctx context.Context, name, method string, args ...interface{}) ([]interface{}, error)
	Send(ctx context.Context, sender, name, method string, args ...interface{}) (string, error)
	BalanceOf(context.Context, string) (*big.Int, error)
	CheckTrsaction(context.Context, string) error
	GetStorePayHash(ctx context.Context, checksize uint64, nonce *big.Int) string
	GetReadPayHash(ctx context.Context, checksize uint64, nonce *big.Int) string

	BuySpace(ctx context.Context, buyer string, size uint64) (string, error)
	BuyTraffic(ctx context.Context, buyer string, size uint64) (string, error)
	Approve(ctx context.Context, pt PayType, sender string, buyValue *big.Int) (string, error)
	Allowance(ctx context.Context, pt PayType, buyer string) (*big.Int, error)
	CashTrafficCheck(ctx context.Context, sender string, nonce *big.Int, sizeByte uint64, sign []byte) (string, error)
	CashSpaceCheck(ctx context.Context, sender string, nonce *big.Int, sizeByte uint64, durationDay uint64, sign []byte) (string, error)
}

type IDataBase interface {
	ListObjects(context.Context, string, StorageType) ([]interface{}, error)
	GetObjectInfo(context.Context, string, string, StorageType) (interface{}, error)
	GetObjectInfoById(context.Context, int) (interface{}, error)
	PutObject(context.Context, FileInfo) error
	DeleteObject(context.Context, int) error

	GetUpSize(context.Context, string) (uint64, error)
	GetDownSize(context.Context, string) (uint64, error)
	Upload(context.Context, CheckInfo) error
	Download(context.Context, CheckInfo) error
	SpaceCheck(ctx context.Context, buyer string) CheckInfo
	TrafficCheck(ctx context.Context, buyer string) CheckInfo

	AddUser(context.Context, USerInfo) error
	SelectUser(context.Context, string) (USerInfo, error)
	DeleteUser(context.Context, int) error
	ListUsers(context.Context) ([]USerInfo, error)
	GetUser(context.Context, int) (USerInfo, error)
}

type IConfig interface {
	GetStore() interface{}
	GetContract() interface{}
}

type Keystore interface {
	Get(string) ([]byte, error)
	Put(string, []byte) error
	List() ([]string, error)
	Delete(string) error
}
