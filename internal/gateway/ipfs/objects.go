package ipfs

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	shapi "github.com/ipfs/go-ipfs-api"
	"github.com/memoio/backend/config"
	db "github.com/memoio/backend/global/database"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/backend/utils"
)

var _ gateway.IGateway = (*Ipfs)(nil)

func ChunkerSize(size string) shapi.AddOpts {
	return func(rb *shapi.RequestBuilder) error {
		rb.Option("chunker", size)
		return nil
	}
}

type Ipfs struct {
	host string
}

func NewGateway() (gateway.IGateway, error) {
	cf, err := config.ReadFile()
	if err != nil {
		return nil, err
	}

	return &Ipfs{
		host: cf.Storage.Ipfs.Host,
	}, nil
}

func (i *Ipfs) PutObject(ctx context.Context, address, object string, r io.Reader, opt gateway.ObjectOptions) (objInfo gateway.ObjectInfo, err error) {
	sh := shapi.NewShell(i.host)
	cidvereion := shapi.CidVersion(1)
	chunkersize := ChunkerSize("size-253952")
	hash, err := sh.Add(r, cidvereion, chunkersize)
	if err != nil {
		return objInfo, err
	}
	oi := db.ObjectInfo{Address: address, Name: object, Size: opt.Size, Cid: hash}
	err = oi.Insert()
	if err != nil {
		return objInfo, err
	}

	return gateway.ObjectInfo{
		Address: address,
		Name:    object,
		Size:    int64(opt.Size),
		Cid:     hash,
		ModTime: time.Now(),
	}, nil
}

func (i *Ipfs) GetObject(ctx context.Context, cid string, w io.Writer, opts gateway.ObjectOptions) error {
	sh := shapi.NewShell(i.host)
	r, err := sh.Cat(cid)
	if err != nil {
		return err
	}
	data, _ := ioutil.ReadAll(r)
	w.Write(data)
	return nil
}

func (i *Ipfs) GetObjectInfo(ctx context.Context, cid string) (gateway.ObjectInfo, error) {
	sh := shapi.NewShell(i.host)
	objects, err := sh.List(cid)
	if err != nil {
		return gateway.ObjectInfo{}, err
	}
	ctype := utils.TypeByExtension(objects[0].Name)
	return gateway.ObjectInfo{
		Name:  objects[0].Name,
		Size:  int64(objects[0].Size),
		CType: ctype,
	}, nil
}

func (i *Ipfs) ListObjects(ctx context.Context, address string) ([]gateway.ObjectInfo, error) {
	ob, err := db.ListObjects(address)
	if err != nil {
		return []gateway.ObjectInfo{}, err
	}

	var objects []gateway.ObjectInfo

	for _, oj := range ob {
		objects = append(objects, toObjectInfo(oj))
	}
	return objects, nil
}

func (m *Ipfs) GetPkgSize(ctx context.Context, address string) (storage.StorageInfo, error) {
	ai, err := db.QueryPkgSize(address, uint8(storage.IPFS))
	if err != nil {
		if err == db.ErrNotExist {
			si, err := contract.GetPkgSize(storage.IPFS, address)
			if err != nil {
				return si, err
			}
			log.Println("si", si)
			ai = db.Storage{
				Address:    address,
				SType:      uint8(storage.IPFS),
				Buysize:    si.Buysize,
				Free:       si.Free,
				Used:       si.Used,
				Files:      si.Files,
				UpdateTime: time.Now(),
			}

			err = ai.Insert()
			if err != nil {
				return storage.StorageInfo{}, err
			}
			return si, nil
		}
		return storage.StorageInfo{}, err
	}

	return storage.StorageInfo{Storage: storage.IPFS.String(), Buysize: ai.Buysize, Used: ai.Used, Free: ai.Free, Files: ai.Files}, nil
}

func (m *Ipfs) UpdateStorage(ctx context.Context, address, cid string, size *big.Int) bool {
	pi := db.PkgInfo{
		Address:   address,
		SType:     uint8(storage.IPFS),
		Hashid:    cid,
		Size:      size.Int64(),
		IsUpdated: false,
		UTime:     time.Now(),
	}

	err := pi.Insert()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func toObjectInfo(o db.ObjectInfo) gateway.ObjectInfo {
	return gateway.ObjectInfo{
		Address: o.Address,
		Name:    o.Name,
		Cid:     o.Cid,
		Size:    o.Size,
	}
}
