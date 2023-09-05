package ipfs

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	shapi "github.com/ipfs/go-ipfs-api"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
)

var _ api.IGateway = (*Ipfs)(nil)

func ChunkerSize(size string) shapi.AddOpts {
	return func(rb *shapi.RequestBuilder) error {
		rb.Option("chunker", size)
		return nil
	}
}

type Ipfs struct {
	st   api.StorageType
	host string
}

func NewGateway() (api.IGateway, error) {
	cf, err := config.ReadFile()
	if err != nil {
		return nil, err
	}

	return &Ipfs{
		host: cf.Storage.Ipfs.Host,
		st:   api.IPFS,
	}, nil
}
func (i *Ipfs) GetStoreType(ctx context.Context) api.StorageType {
	return i.st
}

func (i *Ipfs) PutObject(ctx context.Context, bucket, object string, r io.Reader, opts api.ObjectOptions) (objInfo api.ObjectInfo, err error) {
	sh := shapi.NewShell(i.host)
	cidvereion := shapi.CidVersion(1)
	chunkersize := ChunkerSize("size-253952")
	hash, err := sh.Add(r, cidvereion, chunkersize)
	if err != nil {
		return objInfo, err
	}

	return api.ObjectInfo{
		Bucket:      bucket,
		Name:        object,
		Size:        int64(opts.Size),
		Cid:         hash,
		ModTime:     time.Now(),
		UserDefined: opts.UserDefined,
		SType:       i.st,
	}, nil
}

func (i *Ipfs) GetObject(ctx context.Context, cid string, w io.Writer, opts api.ObjectOptions) error {
	sh := shapi.NewShell(i.host)
	r, err := sh.Cat(cid)
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(r)
	w.Write(data)
	return nil
}

func (i *Ipfs) GetObjectInfo(ctx context.Context, cid string) (api.ObjectInfo, error) {
	result := api.ObjectInfo{}
	sh := shapi.NewShell(i.host)
	objects, err := sh.List(cid)
	if err != nil {
		return result, err
	}
	ctype := utils.TypeByExtension(objects[0].Name)
	return api.ObjectInfo{
		Name:  objects[0].Name,
		Size:  int64(objects[0].Size),
		CType: ctype,
	}, nil
}

func (i *Ipfs) DeleteObject(ctx context.Context, address, mid string) error {
	return logs.StorageError{Message: "ipfs not support delete option"}
}
