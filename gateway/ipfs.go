package gateway

import (
	"context"
	"io"
	"io/ioutil"
	"math/big"
	"os"

	shapi "github.com/ipfs/go-ipfs-api"
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

func (i *Ipfs) Putobject(r io.Reader) (string, error) {
	sh := shapi.NewShell(i.host)
	cidvereion := shapi.CidVersion(1)
	chunkersize := ChunkerSize("size-253952")
	hash, err := sh.Add(r, cidvereion, chunkersize)
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

func (i *Ipfs) ListObjects(ctx context.Context, address string) (ListObjectsInfo, error) {
	return ListObjectsInfo{}, NotImplemented{}
}

func (i *Ipfs) saveFileInfo(ctx context.Context, address, object, cid string, opts ObjectOptions) error {
	err := creatFile("address/" + address)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(address, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	size := big.NewInt(opts.Size)
	info := object + " " + size.String() + " " + cid + "\n"
	_, err = f.Write([]byte(info))
	return err
}

func creatFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	err = os.MkdirAll(path, os.ModePerm)
	return err
}
