package mefs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/memoio/backend/gateway/types"
	"github.com/memoio/backend/utils"
	mclient "github.com/memoio/go-mefs-v2/api/client"
	"github.com/memoio/go-mefs-v2/build"
	mcode "github.com/memoio/go-mefs-v2/lib/code"
	metag "github.com/memoio/go-mefs-v2/lib/etag"
	mtypes "github.com/memoio/go-mefs-v2/lib/types"
)

type ObjectInfo = types.ObjectInfo

type Mefs struct {
	addr    string
	headers http.Header
}

func New() (*Mefs, error) {
	l := types.New("New")

	repoDir := os.Getenv("MEFS_PATH")
	addr, headers, err := mclient.GetMemoClientInfo(repoDir)
	if err != nil {
		return nil, l.DealError(err)
	}
	napi, closer, err := mclient.NewUserNode(context.Background(), addr, headers)
	if err != nil {
		return nil, l.DealError(err)
	}
	defer closer()
	_, err = napi.ShowStorage(context.Background())
	if err != nil {
		return nil, l.DealError(err)
	}

	return &Mefs{
		addr:    addr,
		headers: headers,
	}, nil
}

func (m *Mefs) MakeBucketWithLocation(ctx context.Context, bucket string) error {
	l := types.New("MakeBucketWithLocation")
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return l.DealError(err)
	}
	defer closer()
	opts := mcode.DefaultBucketOptions()

	_, err = napi.CreateBucket(ctx, bucket, opts)
	if err != nil {
		return l.DealError(err)
	}
	return nil
}

func (m *Mefs) GetBucketInfo(ctx context.Context, bucket string) (bi mtypes.BucketInfo, err error) {
	l := types.New("GetBucketInfo")
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return bi, l.DealError(err)
	}
	defer closer()

	bi, err = napi.HeadBucket(ctx, bucket)
	if err != nil {
		return bi, l.DealError(err)
	}
	return bi, nil
}

func (m *Mefs) QueryPrice(ctx context.Context) (string, error) {
	l := types.New("QueryPrice")
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return "", l.DealError(err)
	}
	defer closer()

	res, err := napi.ConfigGet(ctx, "order.price")
	if err != nil {
		return "", l.DealError(err)
	}

	bs, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return "", l.DealError(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, bs, "", "\t")
	if err != nil {
		return "", l.DealError(err)
	}

	return out.String(), nil
}

func (m *Mefs) PutObject(ctx context.Context, address, object string, r io.Reader, UserDefined map[string]string) (objInfo mtypes.ObjectInfo, err error) {
	l := types.New("PutObject")
	err = m.MakeBucketWithLocation(ctx, address)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return objInfo, l.DealError(err)
		}
	} else {
		log.Println("create bucket ", address)
		for !m.CheckBucket(ctx, address) {
			time.Sleep(10 * time.Second)
		}
	}

	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return objInfo, l.DealError(err)
	}
	defer closer()

	poo := mtypes.CidUploadOption()
	for k, v := range UserDefined {
		poo.UserDefined[k] = v
	}
	moi, err := napi.PutObject(ctx, address, object, r, poo)
	if err != nil {
		log.Println(err)
		return objInfo, l.DealError(err)
	}
	return moi, nil
}

func (m *Mefs) GetObject(ctx context.Context, objectName string, writer io.Writer) error {
	l := types.New("PutObject")
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return l.DealError(err)
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", objectName)
	if err != nil {
		return l.DealError(err)
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
	l := types.New("PutObject")
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return ObjectInfo{}, l.DealError(err)
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", cid)
	if err != nil {
		return ObjectInfo{}, l.DealError(err)
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

func (m *Mefs) ListObjects(ctx context.Context, address string) ([]ObjectInfo, error) {
	l := types.New("PutObject")
	var loi []ObjectInfo
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return loi, l.DealError(err)
	}
	defer closer()
	mloi, err := napi.ListObjects(ctx, address, mtypes.ListObjectsOptions{MaxKeys: 1000})
	if err != nil {
		return loi, l.DealError(err)
	}

	for _, oi := range mloi.Objects {
		etag, _ := metag.ToString(oi.ETag)
		loi = append(loi, ObjectInfo{
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

// func (m *Mefs) GetBalanceInfo(ctx context.Context, address string) (string, error) {
// 	bal := contract.BalanceOf(common.HexToAddress(address))
// 	log.Printf("address: %s,balance: %s\n", address, bal)

// 	return bal.String(), nil
// }

func (m *Mefs) DeleteObject(ctx context.Context, address, object string) error {
	l := types.New("PutObject")
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return l.DealError(err)
	}
	defer closer()
	err = napi.DeleteObject(ctx, address, object)
	if err != nil {
		return l.DealError(err)
	}
	return nil
}

func (m *Mefs) CheckBucket(ctx context.Context, address string) bool {
	bi, err := m.GetBucketInfo(ctx, address)
	if err != nil {
		return false
	}

	return bi.Confirmed
}
