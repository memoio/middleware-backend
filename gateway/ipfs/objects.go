package ipfs

import (
	"context"
	"io"
	"io/ioutil"

	shapi "github.com/ipfs/go-ipfs-api"
	"github.com/memoio/backend/gateway/types"
	db "github.com/memoio/backend/global/database"
	"github.com/memoio/backend/utils"
)

type ObjectInfo = types.ObjectInfo

func ChunkerSize(size string) shapi.AddOpts {
	return func(rb *shapi.RequestBuilder) error {
		rb.Option("chunker", size)
		return nil
	}
}

type Ipfs struct {
	host string
}

func New(host string) *Ipfs {
	return &Ipfs{
		host: host,
	}
}

func (i *Ipfs) Putobject(address, name string, size int64, r io.Reader) (string, error) {
	l := types.New("Putobject")
	sh := shapi.NewShell(i.host)
	cidvereion := shapi.CidVersion(1)
	chunkersize := ChunkerSize("size-253952")
	hash, err := sh.Add(r, cidvereion, chunkersize)
	if err != nil {
		return "", l.DealError(err)
	}
	oi := db.ObjectInfo{Address: address, Name: name, Size: size, Cid: hash}
	err = oi.Insert()
	if err != nil {
		return "", l.DealError(err)
	}

	return hash, nil
}

func (i *Ipfs) GetObject(cid string) ([]byte, error) {
	l := types.New("GetObject")
	sh := shapi.NewShell(i.host)
	r, err := sh.Cat(cid)
	if err != nil {
		return nil, l.DealError(err)
	}
	data, _ := ioutil.ReadAll(r)
	return data, nil
}

func (i *Ipfs) GetObjectInfo(ctx context.Context, cid string) (ObjectInfo, error) {
	l := types.New("GetObjectInfo")
	sh := shapi.NewShell(i.host)
	objects, err := sh.List(cid)
	if err != nil {
		return ObjectInfo{}, l.DealError(err)
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
		objects = append(objects, toObjectInfo(oj))
	}
	return objects, nil
}

func toObjectInfo(o db.ObjectInfo) ObjectInfo {
	return ObjectInfo{
		Address: o.Address,
		Name:    o.Name,
		Cid:     o.Cid,
		Size:    o.Size,
	}
}
