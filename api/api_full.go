package api

import (
	"context"
	"io"
	"math/big"
)

type IGateway interface {
	GetStoreType(context.Context) StorageType
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
	GetSapceCheckHash(ctx context.Context, checksize uint64, nonce *big.Int) Check
	GetTrafficCheckHash(ctx context.Context, checksize uint64, nonce *big.Int) Check

	BuySpace(ctx context.Context, buyer string, size uint64) (Transaction, error)
	BuyTraffic(ctx context.Context, buyer string, size uint64) (Transaction, error)
	ApproveTsHash(ctx context.Context, pt PayType, sender string, buyValue *big.Int) (Transaction, error)
	Allowance(ctx context.Context, pt PayType, buyer string) (*big.Int, error)
	CashTrafficCheck(context.Context, CheckInfo) (string, error)
	CashSpaceCheck(context.Context, CheckInfo) (string, error)
}

type IDataBase interface {
	ListObjects(context.Context, string, StorageType) ([]interface{}, error)
	GetObjectInfo(context.Context, string, string, StorageType) (interface{}, error)
	GetObjectInfoById(context.Context, int) (interface{}, error)
	PutObject(context.Context, FileInfo) error
	DeleteObject(context.Context, int) error

	AddUser(context.Context, USerInfo) error
	SelectUser(context.Context, string) (USerInfo, error)
	DeleteUser(context.Context, int) error
	ListUsers(context.Context) ([]USerInfo, error)
	GetUser(context.Context, int) (USerInfo, error)
}

type IDataStore interface {
	GetSpaceInfo(context.Context, string) (CheckInfo, error)
	GetTrafficInfo(context.Context, string) (CheckInfo, error)
	Upload(context.Context, CheckInfo) error
	Download(context.Context, CheckInfo) error
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

type IPublicKey interface {
	Commitment(d []byte) (G1, error)
	GenrateProof(rnd Fr, d []byte) (Proof, error)
	VerifyProof(rnd Fr, commit G1, pf Proof) error
}

type KVStore interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
}
