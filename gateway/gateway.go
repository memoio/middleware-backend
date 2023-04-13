package gateway

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/memoio/backend/config"
	"github.com/memoio/backend/contract"
	"github.com/memoio/backend/gateway/ipfs"
	"github.com/memoio/backend/gateway/mefs"
	"github.com/memoio/backend/gateway/types"
	"github.com/memoio/backend/global"
	db "github.com/memoio/backend/global/database"
	logging "github.com/memoio/backend/global/log"
)

type ObjectInfo = types.ObjectInfo

var logger = logging.Logger("gateway")

type Gateway struct {
	Mefs *mefs.Mefs
	Ipfs *ipfs.Ipfs
}

func NewGateway(c *config.Config) *Gateway {
	g := &Gateway{}
	g.getMemofs()
	g.Ipfs = ipfs.New(c.Storage.Ipfs.Host)
	return g
}

func (g Gateway) PutObject(ctx context.Context, address, object string, r io.Reader, storage StorageType, opts ObjectOptions) (ObjectInfo, error) {
	switch storage {
	case MEFS:
		return g.MefsPutObject(ctx, address, object, r, opts)
	case IPFS:
		return g.IpfsPutObject(ctx, address, object, r, opts)
	}
	return ObjectInfo{}, StorageNotSupport{}
}

func (g Gateway) GetObject(ctx context.Context, cid string, storage StorageType, w io.Writer, opt ObjectOptions) error {
	switch storage {
	case MEFS:
		return g.MefsGetObject(ctx, cid, w, opt)
	case IPFS:
		return g.IpfsGetObject(ctx, cid, w, opt)
	}
	return StorageNotSupport{}
}

func (g *Gateway) ListObjects(ctx context.Context, address string, storage StorageType) ([]ObjectInfo, error) {
	if storage == MEFS {
		err := g.getMemofs()
		if err != nil {
			return []ObjectInfo{}, err
		}
		return g.Mefs.ListObjects(ctx, address)
	}
	if storage == IPFS {
		return g.Ipfs.ListObjects(ctx, address)
	}

	return []ObjectInfo{}, StorageNotSupport{}
}

func (g *Gateway) GetObjectInfo(ctx context.Context, storage StorageType, cid string) (ObjectInfo, error) {
	if storage == MEFS {
		err := g.getMemofs()
		if err != nil {
			return ObjectInfo{}, err
		}
		return g.Mefs.GetObjectInfo(ctx, cid)
	} else if storage == IPFS {
		return g.Ipfs.GetObjectInfo(ctx, cid)
	}

	return ObjectInfo{}, StorageNotSupport{}
}

// func (g *Gateway) GetBalanceInfo(ctx context.Context, address string) (string, error) {
// 	err := g.getMemofs()
// 	if err != nil {
// 		return "", err
// 	}
// 	return g.Mefs.GetBalanceInfo(ctx, address)
// }

func (g *Gateway) GetPrice(ctx context.Context, address, size, time string) (string, error) {
	return "", NotImplemented{}
}

func (g *Gateway) GetPkgSize(ctx context.Context, storage StorageType, address string) (global.StorageInfo, error) {
	ai, err := db.QueryPkgSize(address, uint8(storage))
	if err != nil {
		if err == db.ErrNotExist {
			si, err := contract.GetPkgSize(uint8(storage), address)
			if err != nil {
				return si, err
			}
			log.Println("si", si)
			ai = db.Storage{
				Address:    address,
				Buysize:    si.Buysize,
				Free:       si.Free,
				Used:       si.Used,
				Files:      si.Files,
				UpdateTime: time.Now(),
			}

			err = ai.Insert()
			if err != nil {
				return global.StorageInfo{}, err
			}
			return si, nil
		}
		return global.StorageInfo{}, err
	}

	return global.StorageInfo{Buysize: ai.Buysize, Used: ai.Used, Free: ai.Free, Files: ai.Files}, nil
}

func (g *Gateway) TestPutobject(ctx context.Context, address, hashid string, size int64) error {
	if !g.checkStorage(ctx, MEFS, address, big.NewInt(size)) {
		return StorageError{Message: "storage not enough"}
	}

	pi := db.PkgInfo{
		Address:   address,
		Hashid:    hashid,
		Size:      size,
		IsUpdated: false,
		UTime:     time.Now(),
	}

	return pi.Insert()
}

func (g *Gateway) TestUpdatePkg(ctx context.Context, storage StorageType, address, hashid string, size int64) (global.StorageInfo, error) {
	if !contract.StoreOrderPkg(address, hashid, uint8(storage), big.NewInt(size)) {
		return global.StorageInfo{}, fmt.Errorf("update error")
	}

	si, err := contract.GetPkgSize(uint8(storage), address)
	if err != nil {
		return global.StorageInfo{}, err
	}

	return si, nil
}

func (g *Gateway) TestPay(ctx context.Context, address, hash string, amount, size int64) bool {
	return contract.StoreOrderPay(address, hash, big.NewInt(amount), big.NewInt(size))
}
