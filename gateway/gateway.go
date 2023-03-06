package gateway

import (
	"context"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/config"
	logging "github.com/memoio/backend/global/log"
	"github.com/memoio/backend/utils"
	metag "github.com/memoio/go-mefs-v2/lib/utils/etag"
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

func (g *Gateway) getMemofs() error {
	var err error
	g.Mefs, err = newMefs()
	if err != nil {
		return err
	}
	return nil
}

func (g Gateway) PutObject(ctx context.Context, address, object string, r io.Reader, storage StorageType, opts ObjectOptions) (ObjectInfo, error) {
	if storage == MEFS {
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

		flag := g.verify(ctx, address, date, etag, size)
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
	} else if storage == IPFS {
		logger.Debug("ipfs put object")
		size := big.NewInt(opts.Size)
		if !g.checkStorage(ctx, address, size) {
			return ObjectInfo{}, StorageError{Storage: storage.String(), Message: "storage not enough"}
		}
		cid, err := g.Ipfs.Putobject(r)
		if err != nil {
			return ObjectInfo{}, err
		}
		ctype := utils.TypeByExtension(object)
		if opts.UserDefined["content-type"] != "" {
			ctype = opts.UserDefined["content-type"]
		}

		if !g.updateStorage(ctx, address, cid, size) {
			return ObjectInfo{}, StorageError{Storage: storage.String(), Message: "storage update error"}
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
		return nil
	} else if storage == IPFS {
		data, err := g.Ipfs.GetObject(cid)
		if err != nil {
			return err
		}
		w.Write(data)
		return nil
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

func (g *Gateway) GetBalanceInfo(ctx context.Context, address string ) (string, error) {
	err := g.getMemofs()
	if err != nil {
		return "", err
	}
	return g.Mefs.GetBalanceInfo(ctx, address)
}

func (g *Gateway) GetPrice(ctx context.Context, address, size, time string) (string, error) {
	return "", NotImplemented{}
}

func (g *Gateway) GetStorageInfo(ctx context.Context, address string) (StorageInfo, error) {
	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		log.Println("connect to eth error", err)
		return StorageInfo{}, err
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
		return StorageInfo{}, err
	}

	if len(result) != 128 {
		return StorageInfo{}, StorageError{}
	}

	available := new(big.Int)
	available.SetBytes(result[0:32])
	free := new(big.Int)
	free.SetBytes(result[32:64])
	used := new(big.Int)
	used.SetBytes(result[64:96])
	files := new(big.Int)
	files.SetBytes(result[96:])

	si := StorageInfo{
		Available: available.String(),
		Free:      free.String(),
		Used:      used.String(),
		Files:     files.String(),
	}
	log.Println(si)
	return si, nil
}
