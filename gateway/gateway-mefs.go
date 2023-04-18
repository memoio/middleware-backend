package gateway

import (
	"context"
	"errors"
	"io"
	"math/big"
	"time"

	"github.com/memoio/backend/contract"
	"github.com/memoio/backend/gateway/mefs"
	"github.com/memoio/backend/utils"
	metag "github.com/memoio/go-mefs-v2/lib/etag"
)

func (g *Gateway) getMemofs() error {
	var err error
	g.Mefs, err = mefs.New()
	if err != nil {
		return err
	}
	return nil
}

func (g Gateway) MefsPutObject(ctx context.Context, address, object string, r io.Reader, opts ObjectOptions) (ObjectInfo, error) {
	logger.Debug("mefs put object")
	err := g.getMemofs()
	if err != nil {
		return ObjectInfo{}, err
	}
	date := opts.UserDefined["X-Amz-Meta-Date"]
	if date == "" {
		date = "365"
	}

	moi, err := g.Mefs.PutObject(ctx, address, object, r, opts.UserDefined)
	if err != nil {
		return ObjectInfo{}, err
	}

	etag, _ := metag.ToString(moi.ETag)
	size := big.NewInt(int64(moi.Size))

	flag := g.verify(ctx, MEFS, address, date, etag, size)
	if !flag {
		g.Mefs.DeleteObject(ctx, address, object)
		return ObjectInfo{}, err
	}

	ctype := utils.TypeByExtension(object)

	if moi.UserDefined["content-type"] != "" {
		ctype = moi.UserDefined["content-type"]
	}

	oi := ObjectInfo{
		Address: address,
		Name:    moi.Name,
		ModTime: time.Unix(moi.GetTime(), 0),
		Size:    int64(moi.Size),
		Cid:     etag,
		CType:   ctype,
	}

	return oi, nil
}

func (g Gateway) MefsGetObject(ctx context.Context, cid string, w io.Writer, opt ObjectOptions) error {
	err := g.getMemofs()
	if err != nil {
		return err
	}
	err = g.Mefs.GetObject(ctx, cid, w)
	if err != nil {
		return err
	}
	return nil
}

func (g Gateway) MefsDeleteObject(ctx context.Context, address, mid string) error {
	err := g.getMemofs()
	if err != nil {
		return err
	}
	oi, err := g.Mefs.GetObjectInfo(ctx, mid)
	if err != nil {
		return err
	}

	if oi.Address != address {
		return nil
	}

	err = g.Mefs.DeleteObject(ctx, address, oi.Name)
	if err != nil {
		return err
	}

	r := contract.StoreOrderPkgExpiration(address, mid, uint8(MEFS), big.NewInt(oi.Size))
	if !r {
		return errors.New("contract error")
	}
	return nil
}
