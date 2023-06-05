package mefs

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/utils"
	mclient "github.com/memoio/go-mefs-v2/api/client"
	"github.com/memoio/go-mefs-v2/build"
	mcode "github.com/memoio/go-mefs-v2/lib/code"
	metag "github.com/memoio/go-mefs-v2/lib/etag"
	mtypes "github.com/memoio/go-mefs-v2/lib/types"
)

var _ gateway.IGateway = (*Mefs)(nil)

type Mefs struct {
	addr    string
	headers http.Header
}

func NewGateway() (gateway.IGateway, error) {
	repoDir := os.Getenv("MEFS_PATH")
	addr, headers, err := mclient.GetMemoClientInfo(repoDir)
	if err != nil {
		return nil, err
	}
	napi, closer, err := mclient.NewUserNode(context.Background(), addr, headers)
	if err != nil {
		return nil, err
	}
	defer closer()
	_, err = napi.ShowStorage(context.Background())
	if err != nil {
		return nil, err
	}

	return &Mefs{
		addr:    addr,
		headers: headers,
	}, nil
}

func (m *Mefs) MakeBucketWithLocation(ctx context.Context, bucket string) error {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return err
	}
	defer closer()
	opts := mcode.DefaultBucketOptions()

	_, err = napi.CreateBucket(ctx, bucket, opts)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mefs) GetBucketInfo(ctx context.Context, bucket string) (bi mtypes.BucketInfo, err error) {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return bi, err
	}
	defer closer()

	bi, err = napi.HeadBucket(ctx, bucket)
	if err != nil {
		return bi, err
	}
	return bi, nil
}

// func (m *Mefs) QueryPrice(ctx context.Context) (string, error) {
// 	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer closer()

// 	res, err := napi.ConfigGet(ctx, "order.price")
// 	if err != nil {
// 		return "", err
// 	}

// 	bs, err := json.MarshalIndent(res, "", "\t")
// 	if err != nil {
// 		return "", err
// 	}

// 	var out bytes.Buffer
// 	err = json.Indent(&out, bs, "", "\t")
// 	if err != nil {
// 		return "", err
// 	}

// 	return out.String(), nil
// }

func (m *Mefs) PutObject(ctx context.Context, bucket, object string, r io.Reader, opt gateway.ObjectOptions) (objInfo gateway.ObjectInfo, err error) {
	err = m.MakeBucketWithLocation(ctx, bucket)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			return objInfo, err
		}
	} else {
		log.Println("create bucket ", bucket)
		for !m.CheckBucket(ctx, bucket) {
			time.Sleep(10 * time.Second)
		}
	}

	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return objInfo, err
	}
	defer closer()

	poo := mtypes.CidUploadOption()
	for k, v := range opt.UserDefined {
		poo.UserDefined[k] = v
	}
	moi, err := napi.PutObject(ctx, bucket, object, r, poo)
	if err != nil {
		log.Println(err)
		return objInfo, err
	}

	etag, _ := metag.ToString(moi.ETag)

	return gateway.ObjectInfo{
		Bucket:      bucket,
		Name:        moi.Name,
		Size:        int64(moi.Size),
		Cid:         etag,
		ModTime:     time.Unix(moi.GetTime(), 0),
		UserDefined: moi.UserDefined,
	}, nil
}

func (m *Mefs) GetObject(ctx context.Context, objectName string, writer io.Writer, opt gateway.ObjectOptions) error {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return err
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", objectName)
	if err != nil {
		return err
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

func (m *Mefs) GetObjectInfo(ctx context.Context, cid string) (gateway.ObjectInfo, error) {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return gateway.ObjectInfo{}, err
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", cid)
	if err != nil {
		return gateway.ObjectInfo{}, err
	}
	ctype := utils.TypeByExtension(objInfo.Name)
	if objInfo.UserDefined["content-type"] != "" {
		ctype = objInfo.UserDefined["content-type"]
	}
	return gateway.ObjectInfo{
		Name:  objInfo.Name,
		Size:  int64(objInfo.Size),
		CType: ctype,
	}, nil
}

func (m *Mefs) ListObjects(ctx context.Context, bucket string) ([]gateway.ObjectInfo, error) {
	var loi []gateway.ObjectInfo
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return loi, err
	}
	defer closer()
	mloi, err := napi.ListObjects(ctx, bucket, mtypes.ListObjectsOptions{MaxKeys: 1000})
	if err != nil {
		return loi, err
	}

	for _, oi := range mloi.Objects {
		etag, _ := metag.ToString(oi.ETag)
		loi = append(loi, gateway.ObjectInfo{
			Bucket:      bucket,
			Name:        oi.GetName(),
			ModTime:     time.Unix(oi.GetTime(), 0).UTC(),
			Size:        int64(oi.Size),
			Cid:         etag,
			UserDefined: oi.UserDefined,
		})
	}

	return loi, nil
}

func (m *Mefs) DeleteObject(ctx context.Context, bucket, object string) error {
	napi, closer, err := mclient.NewUserNode(ctx, m.addr, m.headers)
	if err != nil {
		return err
	}
	defer closer()

	err = napi.DeleteObject(ctx, bucket, object)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mefs) CheckBucket(ctx context.Context, bucket string) bool {
	bi, err := m.GetBucketInfo(ctx, bucket)
	if err != nil {
		return false
	}

	return bi.Confirmed
}
