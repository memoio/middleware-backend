package gateway

import (
	"context"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

type ObjectInfo struct {
	Address     string
	Name        string
	Size        int64
	Cid         string
	ModTime     time.Time
	CType       string
	UserDefined map[string]string
}

type ObjectOptions struct {
	Size         int64
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}

type IGateway interface {
	// objects
	PutObject(context.Context, string, string, io.Reader, ObjectOptions) (ObjectInfo, error)
	GetObject(context.Context, string, io.Writer, ObjectOptions) error
	ListObjects(context.Context, string) ([]ObjectInfo, error)
	GetObjectInfo(context.Context, string) (ObjectInfo, error)

	// contract
	PayForSize(context.Context, string, *big.Int, *big.Int) bool
	GetPkgSize(context.Context, string) (storage.StorageInfo, error)
	UpdateStorage(context.Context, string, string, *big.Int) bool
	// database
}

func PutObject(ctx context.Context, api IGateway, address, object string, r io.Reader, opts ObjectOptions) (ObjectInfo, error) {
	date := opts.UserDefined["X-Amz-Meta-Date"]
	if date == "" {
		date = "365"
	}
	days := new(big.Int)
	days.SetString(date, 10)
	size := big.NewInt(opts.Size)
	res := CheckStorage(ctx, api, address, size)
	if !res {
		if !api.PayForSize(ctx, address, days, size) {
			log.Println("Error: storage not enough and pay for size faild")
			return ObjectInfo{}, logs.GatewayError{Message: "storage not enough and pay for size faild"}
		}
	}
	log.Println("Check Storage Passed ", size)
	ob, err := api.PutObject(ctx, address, object, r, opts)
	if err != nil {
		return ObjectInfo{}, err
	}

	result := api.UpdateStorage(ctx, address, ob.Cid, size)
	if !result {
		log.Printf("Error: update storage failed, address: %s, cid %s, size %d\n", address, ob.Cid, size)
	}

	return ob, nil
}

func CheckStorage(ctx context.Context, api IGateway, address string, size *big.Int) bool {
	si, err := api.GetPkgSize(ctx, address)
	if err != nil {
		return false
	}

	return si.Buysize+si.Free > size.Int64()+si.Used
}
