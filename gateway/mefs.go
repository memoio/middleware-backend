package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/utils"
	mclient "github.com/memoio/go-mefs-v2/api/client"
	"github.com/memoio/go-mefs-v2/build"
	mcode "github.com/memoio/go-mefs-v2/lib/code"
	mtypes "github.com/memoio/go-mefs-v2/lib/types"
	metag "github.com/memoio/go-mefs-v2/lib/utils/etag"
	"golang.org/x/crypto/sha3"
)

type Mefs struct {
	addr    string
	headers http.Header
}

func newMefs() (*Mefs, error) {
	repoDir := os.Getenv("MEFS_PATH")
	addr, headers, err := mclient.GetMemoClientInfo(repoDir)
	if err != nil {
		return nil, funcError(MEFS, newfunc, err)
	}
	napi, closer, err := mclient.NewUserNode(context.Background(), addr, headers)
	if err != nil {
		return nil, funcError(MEFS, newfunc, err)
	}
	defer closer()
	_, err = napi.ShowStorage(context.Background())
	if err != nil {
		return nil, funcError(MEFS, newfunc, err)
	}

	return &Mefs{
		addr:    addr,
		headers: headers,
	}, nil
}

func (m *Mefs) MakeBucketWithLocation(ctx context.Context, bucket string) error {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return funcError(MEFS, makefunc, err)
	}
	defer closer()
	opts := mcode.DefaultBucketOptions()

	_, err = napi.CreateBucket(ctx, bucket, opts)
	if err != nil {
		return funcError(MEFS, makefunc, err)
	}
	return nil
}

func (m *Mefs) GetBucketInfo(ctx context.Context, bucket string) (bi mtypes.BucketInfo, err error) {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return bi, StorageError{Message: err.Error()}
	}
	defer closer()

	bi, err = napi.HeadBucket(ctx, bucket)
	if err != nil {
		return bi, StorageError{Message: err.Error()}
	}
	return bi, nil
}

func (m *Mefs) QueryPrice(ctx context.Context) (string, error) {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return "", funcError(MEFS, pricefunc, err)
	}
	defer closer()

	res, err := napi.ConfigGet(ctx, "order.price")
	if err != nil {
		return "", funcError(MEFS, pricefunc, err)
	}

	bs, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return "", funcError(MEFS, pricefunc, err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, bs, "", "\t")
	if err != nil {
		return "", funcError(MEFS, pricefunc, err)
	}

	return out.String(), nil
}

func (m *Mefs) PutObject(ctx context.Context, bucket, object string, r io.Reader, UserDefined map[string]string) (objInfo mtypes.ObjectInfo, err error) {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return objInfo, funcError(MEFS, putfunc, err)
	}
	defer closer()

	poo := mtypes.CidUploadOption()
	for k, v := range UserDefined {
		poo.UserDefined[k] = v
	}
	moi, err := napi.PutObject(ctx, bucket, object, r, poo)
	if err != nil {
		log.Println(err)
		return objInfo, funcError(MEFS, putfunc, err)
	}
	return moi, nil
}

func (m *Mefs) GetObject(ctx context.Context, objectName string, writer io.Writer) error {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return funcError(MEFS, getfunc, err)
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", objectName)
	if err != nil {
		return funcError(MEFS, getfunc, err)
	}

	length := int64(objInfo.Size)

	stepLen := int64(build.DefaultSegSize * 16)
	stepAccMax := 16

	start := int64(0)
	end := length
	stepacc := 1
	for start < end {
		if stepacc > stepAccMax {
			stepacc = stepAccMax
		}

		readLen := stepLen*int64(stepacc) - (start % stepLen)
		if end-start < readLen {
			readLen = end - start
		}

		doo := mtypes.DownloadObjectOptions{
			Start:  start,
			Length: readLen,
		}

		data, err := napi.GetObject(ctx, "", objectName, doo)
		if err != nil {
			//log.Println("received length err is:", start, readLen, stepLen, err)
			break
		}
		writer.Write(data)
		start += int64(readLen)
		stepacc *= 2
	}

	return nil
}

func (m *Mefs) GetObjectInfo(ctx context.Context, cid string) (ObjectInfo, error) {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return ObjectInfo{}, funcError(MEFS, getfunc, err)
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", cid)
	if err != nil {
		return ObjectInfo{}, funcError(MEFS, getfunc, err)
	}
	ctype := utils.TypeByExtension(objInfo.Name)
	if objInfo.UserDefined["content-type"] != "" {
		ctype = objInfo.UserDefined["content-type"]
	}
	return ObjectInfo{
		Name:  objInfo.Name,
		Size:  int64(objInfo.Size),
		CType: ctype,
	}, nil
}

func (m *Mefs) ListObjects(ctx context.Context, address string) (ListObjectsInfo, error) {
	loi := ListObjectsInfo{}
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return loi, funcError(MEFS, listfunc, err)
	}
	defer closer()
	mloi, err := napi.ListObjects(ctx, address, mtypes.ListObjectsOptions{MaxKeys: 1000})
	if err != nil {
		return loi, funcError(MEFS, listfunc, err)
	}

	for _, oi := range mloi.Objects {
		etag, _ := metag.ToString(oi.ETag)
		loi.Objects = append(loi.Objects, ObjectInfo{
			Address:     address,
			Name:        oi.GetName(),
			ModTime:     time.Unix(oi.GetTime(), 0).UTC(),
			Size:        int64(oi.Size),
			Cid:         etag,
			UserDefined: oi.UserDefined,
		})
	}

	return loi, nil
}

func (m *Mefs) GetBalanceInfo(ctx context.Context, address string) (string, error) {
	err := m.MakeBucketWithLocation(ctx, address)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return "", funcError(MEFS, makefunc, err)
		}
	} else {
		log.Println("create bucket ", address)
		time.Sleep(20 * time.Second)
	}

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		log.Println("connect to eth error", err)
		return "", EthError{Message: err.Error()}
	}

	defer client.Close()

	addr := common.HexToAddress(address)
	balanceOfFnSignature := []byte("balanceOf(address)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(balanceOfFnSignature)
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
		return "", EthError{Message: err.Error()}
	}
	bal := new(big.Int)
	bal.SetBytes(result)
	log.Printf("address: %s,balance: %s\n", addr, bal)

	return bal.String(), nil
}

func (m *Mefs) DeleteObject(ctx context.Context, address, object string) error {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return funcError(MEFS, deletefunc, err)
	}
	defer closer()
	err = napi.DeleteObject(ctx, address, object)
	if err != nil {
		return funcError(MEFS, deletefunc, err)
	}
	return nil
}
