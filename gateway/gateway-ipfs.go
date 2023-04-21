package gateway

import (
	"context"
	"io"
	"math/big"
	"time"

	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/backend/utils"
)

func (g Gateway) IpfsPutObject(ctx context.Context, address, object string, r io.Reader, opts ObjectOptions) (ObjectInfo, error) {
	logger.Debug("ipfs put object")
	size := big.NewInt(opts.Size)
	if !g.checkStorage(ctx, storage.IPFS, address, size) {
		return ObjectInfo{}, StorageError{Storage: storage.IPFS.String(), Message: "storage not enough"}
	}
	cid, err := g.Ipfs.Putobject(address, object, size.Int64(), r)
	if err != nil {
		return ObjectInfo{}, err
	}
	ctype := utils.TypeByExtension(object)
	if opts.UserDefined["content-type"] != "" {
		ctype = opts.UserDefined["content-type"]
	}

	if !g.updateStorage(ctx, address, cid, size) {
		return ObjectInfo{}, StorageError{Storage: storage.IPFS.String(), Message: "storage update error"}
	}
	oi := ObjectInfo{
		Address: address,
		Name:    object,
		ModTime: time.Now(),
		Cid:     cid,
		CType:   ctype,
	}
	return oi, nil
}

func (g Gateway) IpfsGetObject(ctx context.Context, cid string, w io.Writer, opt ObjectOptions) error {
	data, err := g.Ipfs.GetObject(cid)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}
