package gateway

import (
	"context"
	"io"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/global"
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
		Available: available.String(),
		Free:      free.String(),
		Used:      used.String(),
		Files:     files.String(),
	}
	log.Println(si)
	return si, nil
}
