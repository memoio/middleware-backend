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
	Send(ctx context.Context, name, method string, args ...interface{}) (string, error)
	BalanceOf(context.Context, string) (*big.Int, error)
	StoreBuyPkg(context.Context, string, BuyPackage) (string, error)
	CheckTrsaction(context.Context, string) error
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
