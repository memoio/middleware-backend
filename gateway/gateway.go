package gateway

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/contract"
	"github.com/memoio/backend/global"
	"github.com/memoio/backend/global/db"
	logging "github.com/memoio/backend/global/log"
	"golang.org/x/crypto/sha3"
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

func (g *Gateway) GetBalanceInfo(ctx context.Context, address string) (string, error) {
	err := g.getMemofs()
	if err != nil {
		return "", err
	}
	return g.Mefs.GetBalanceInfo(ctx, address)
}

func (g *Gateway) GetPrice(ctx context.Context, address, size, time string) (string, error) {
	return "", NotImplemented{}
}

func (g *Gateway) GetStorageInfo(ctx context.Context, address string) (global.StorageInfo, error) {
	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		log.Println("connect to eth error", err)
		return global.StorageInfo{}, err
	}
	defer client.Close()

	addr := common.HexToAddress(address)
	pkgSizeFnSignature := []byte("getPkgSize(address)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pkgSizeFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(addr.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return global.StorageInfo{}, err
	}

	if len(result) != 128 {
		return global.StorageInfo{}, StorageError{}
	}

	available := new(big.Int)
	available.SetBytes(result[0:32])
	free := new(big.Int)
	free.SetBytes(result[32:64])
	used := new(big.Int)
	used.SetBytes(result[64:96])
	files := new(big.Int)
	files.SetBytes(result[96:])

	si := global.StorageInfo{
		Available: available.Int64(),
		Free:      free.Int64(),
		Used:      used.Int64(),
		Files:     int(files.Int64()),
	}
	log.Println(si)
	return si, nil
}

func (g *Gateway) GetPkgSize(ctx context.Context, address string) (global.StorageInfo, error) {
	ai, err := db.QueryPkgSize(address)
	if err != nil {
		log.Println(err)
		if err == db.ErrNotExist {
			si, err := contract.GetPkgSize(address)
			if err != nil {
				return si, err
			}

			ai = db.AddressInfo{
				Address:    address,
				Available:  si.Available,
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
	}

	return global.StorageInfo{Available: ai.Available, Used: ai.Used, Free: ai.Free, Files: ai.Files}, nil
}

func (g *Gateway) TestPutobject(ctx context.Context, address, hashid string, size int64) error {
	if !g.checkStorage(ctx, address, big.NewInt(size)) {
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

func (g *Gateway) TestUpdatePkg(ctx context.Context, address, hashid string, size int64) (global.StorageInfo, error) {
	if !contract.StoreOrderPkg(address, hashid, big.NewInt(size)) {
		return global.StorageInfo{}, fmt.Errorf("update error")
	}

	si, err := contract.GetPkgSize(address)
	if err != nil {
		return global.StorageInfo{}, err
	}

	return si, nil
}

func (g *Gateway) TestPay(ctx context.Context, address, hash string, amount, size int64) bool {
	return contract.StoreOrderPay(address, hash, big.NewInt(amount), big.NewInt(size))
}

func (g *Gateway) TestGetPkgInfo(ctx context.Context) error {
	return contract.StoreGetPkgInfos()
}
