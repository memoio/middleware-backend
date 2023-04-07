package gateway

import (
	"context"
	"io"
	"io/ioutil"

	shapi "github.com/ipfs/go-ipfs-api"
	"github.com/memoio/backend/global/db"
	"github.com/memoio/backend/utils"
)

func ChunkerSize(size string) shapi.AddOpts {
	return func(rb *shapi.RequestBuilder) error {
		rb.Option("chunker", size)
		return nil
	}
}

type Ipfs struct {
	host string
}

func NewIpfsClient(host string) *Ipfs {
	return &Ipfs{
		host: host,
	}
}

func (i *Ipfs) Putobject(address, name string, size int64, r io.Reader) (string, error) {
	sh := shapi.NewShell(i.host)
	cidvereion := shapi.CidVersion(1)
	chunkersize := ChunkerSize("size-253952")
	hash, err := sh.Add(r, cidvereion, chunkersize)
	if err != nil {
		return "", funcError(IPFS, putfunc, err)
	}
	oi := db.ObjectInfo{Address: address, Name: name, Size: size, Cid: hash}
	err = oi.Insert()
	if err != nil {
		return "", funcError(IPFS, putfunc, err)
	}

	return hash, nil
}

func (i *Ipfs) GetObject(cid string) ([]byte, error) {
	sh := shapi.NewShell(i.host)
	r, err := sh.Cat(cid)
	if err != nil {
		return nil, funcError(IPFS, getfunc, err)
	}
	data, _ := ioutil.ReadAll(r)
	return data, nil
}

func (i *Ipfs) GetObjectInfo(ctx context.Context, cid string) (ObjectInfo, error) {
	sh := shapi.NewShell(i.host)
	objects, err := sh.List(cid)
	if err != nil {
		return ObjectInfo{}, funcError(IPFS, getinfofunc, err)
	}
	ctype := utils.TypeByExtension(objects[0].Name)
	return ObjectInfo{
		Name:  objects[0].Name,
		Size:  int64(objects[0].Size),
		CType: ctype,
	}, nil
}

func (i *Ipfs) ListObjects(ctx context.Context, address string) ([]ObjectInfo, error) {
	ob, err := db.ListObjects(address)
	if err != nil {
		return []ObjectInfo{}, err
	}

	var objects []ObjectInfo

	for _, oj := range ob {
		objects = append(objects, ObjectInfo(oj))
	}
	return objects, nil
}
