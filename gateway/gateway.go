package gateway

import (
	"context"
	"io"
	"time"

	"github.com/memoio/backend/config"
	logging "github.com/memoio/backend/global/log"
	"github.com/memoio/backend/utils"
	metag "github.com/memoio/go-mefs-v2/lib/utils/etag"
)

var logger = logging.Logger("gateway")

type Gateway struct {
	Mefs *Mefs
	Ipfs *Ipfs
}

func NewGateway(c *config.Config) *Gateway {
	g := &Gateway{}
	g.getMemofs()
	g.Ipfs = NewIpfsClient(c.Storage.Ipfs.Host)
	return g
}

func (g *Gateway) getMemofs() error {
	var err error
	g.Mefs, err = newMefs()
	if err != nil {
		return err
	}
	return nil
}

func (g Gateway) PutObject(ctx context.Context, bucket, object string, r io.Reader, storage StorageType, opt ObjectOptions) (ObjectInfo, error) {
	if storage == MEFS {
		logger.Debug("mefs put object")
		err := g.getMemofs()
		if err != nil {
			return ObjectInfo{}, err
		}
		moi, err := g.Mefs.PutObject(ctx, bucket, object, r, opt.UserDefined)
		if err != nil {
			return ObjectInfo{}, err
		}
		etag, _ := metag.ToString(moi.ETag)
		ctype := utils.TypeByExtension(object)
		if moi.UserDefined["content-type"] != "" {
			ctype = moi.UserDefined["content-type"]
		}

		oi := ObjectInfo{
			Address: bucket,
			Name:    moi.Name,
			ModTime: time.Unix(moi.GetTime(), 0),
			Size:    int64(moi.Size),
			Cid:     etag,
			CType:   ctype,
		}

		return oi, nil
	} else if storage == IPFS {
		logger.Debug("ipfs put object")
		cid, err := g.Ipfs.Putobject(r)
		if err != nil {
			return ObjectInfo{}, err
		}
		ctype := utils.TypeByExtension(object)
		if opt.UserDefined["content-type"] != "" {
			ctype = opt.UserDefined["content-type"]
		}
		oi := ObjectInfo{
			Address: bucket,
			Name:    object,
			ModTime: time.Now(),
			Cid:     cid,
			CType:   ctype,
		}
		return oi, nil
	}
	return ObjectInfo{}, StorageNotSupport{}
}

func (g Gateway) GetObject(ctx context.Context, cid string, storage StorageType, w io.Writer, opt ObjectOptions) error {
	if storage == MEFS {
		err := g.getMemofs()
		if err != nil {
			return err
		}
		err = g.Mefs.GetObject(ctx, cid, w)
		if err != nil {
			return err
		}
	} else if storage == IPFS {
		data, err := g.Ipfs.GetObject(cid)
		if err != nil {
			return err
		}
		w.Write(data)
	}
	return StorageNotSupport{}
}

func (g *Gateway) ListObjects(ctx context.Context, address string, storage StorageType) (ListObjectsInfo, error) {
	if storage == MEFS {
		err := g.getMemofs()
		if err != nil {
			return ListObjectsInfo{}, err
		}
		return g.Mefs.ListObjects(ctx, address)
	}
	if storage == IPFS {
		return g.Ipfs.ListObjects(ctx, address)
	}

	return ListObjectsInfo{}, StorageNotSupport{}
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

func (g *Gateway) GetBalanceInfo(ctx context.Context, address string, storage StorageType) (string, error) {
	if storage == MEFS {
		err := g.getMemofs()
		if err != nil {
			return "", err
		}
		return g.Mefs.GetBalanceInfo(ctx, address)
	} else if storage == IPFS {

	}
	return "", StorageNotSupport{}
}

func (g *Gateway) GetPrice(ctx context.Context, adrress, size, time string) (string, error) {
	return "", NotImplemented{}
}
