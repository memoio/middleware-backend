package controller

import (
	"context"
	"encoding/json"
	"io"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
)

func (c *Controller) SetStore(store api.IGateway) {
	c.store = store
}

func (c *Controller) GetStorage(ctx context.Context) api.StorageType {
	return c.store.GetStoreType(ctx)
}

func (c *Controller) PutObject(ctx context.Context, address, object string, r io.Reader, opts ObjectOptions) (PutObjectResult, error) {
	result := PutObjectResult{}
	ci, err := c.canWrite(ctx, address, opts.Sign, uint64(opts.Size))
	if err != nil {
		return result, err
	}

	if opts.Area != "" {
		err = c.changeStore(ctx, opts.Area)
		if err != nil {
			return result, err
		}
	}

	oi, err := c.store.PutObject(ctx, address, object, r, api.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	userdefine, err := json.Marshal(oi.UserDefined)
	if err != nil {
		return result, err
	}

	fi := api.FileInfo{
		Address:    address,
		Name:       object,
		Mid:        oi.Cid,
		SType:      oi.SType,
		Size:       oi.Size,
		ModTime:    oi.ModTime,
		UserID:     oi.USerID,
		UserDefine: string(userdefine),
	}

	err = c.storeFileInfo(ctx, fi, ci)
	if err != nil {
		c.store.DeleteObject(ctx, address, oi.Name)
		return result, err
	}

	result.Mid = oi.Cid

	return result, nil
}

func (c *Controller) GetObject(ctx context.Context, address, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}

	ob, err := c.getObjectInfo(ctx, address, mid)
	if err != nil {
		return result, err
	}

	ci, err := c.canRead(ctx, address, opts.Sign, uint64(ob.Size))
	if err != nil {
		return result, err
	}

	err = c.store.GetObject(ctx, mid, w, api.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	result.Name = ob.Name
	result.CType = utils.TypeByExtension(ob.Name)
	result.Size = ob.Size

	err = c.datastore.Download(ctx, ci)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Controller) ListObjects(ctx context.Context, address string) (ListObjectsResult, error) {
	result := ListObjectsResult{}

	st := c.store.GetStoreType(ctx)
	loi, err := c.database.ListObjects(ctx, address, st)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Storage = st.String()

	for _, ioi := range loi {
		// userdefine := make(map[string]string)
		oi := ioi.(api.FileInfo)
		// err = json.Unmarshal([]byte(oi.UserDefine), &userdefine)
		// if err != nil {
		// 	lerr := logs.ControllerError{Message: fmt.Sprint("unmarshal userdefine error, ", err)}
		// 	logger.Error(lerr)
		// 	return result, lerr
		// }

		result.Objects = append(result.Objects, ObjectInfoResult{
			ID:      oi.ID,
			Name:    oi.Name,
			Size:    oi.Size,
			Mid:     oi.Mid,
			ModTime: oi.ModTime,
			Public:  oi.Public,
			// UserDefined: userdefine,
		})
	}

	return result, nil
}

func (c *Controller) DeleteObject(ctx context.Context, address string, id int) error {
	oi, err := c.getObjectInfoById(ctx, id)
	if err != nil {
		return err
	}

	if address != oi.Address {
		lerr := logs.ControllerError{Message: "address not right"}
		logger.Error(lerr)
		return lerr
	}

	err = c.store.DeleteObject(ctx, address, oi.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return c.database.DeleteObject(ctx, id)
		}
		return err
	}

	return c.database.DeleteObject(ctx, id)
}

func (c *Controller) GetSpaceCheckHash(ctx context.Context, address string, size uint64) (api.Check, error) {
	ci, err := c.getSpaceCheckInfo(ctx, address, size)
	if err != nil {
		return api.Check{}, err
	}

	return c.contract.GetSapceCheckHash(ctx, ci.FileSize.Uint64(), ci.Nonce), nil
}

func (c *Controller) GetTrafficCheckHash(ctx context.Context, address string, size uint64) (api.Check, error) {
	ci, err := c.getTrafficCheckInfo(ctx, address, size)
	if err != nil {
		return api.Check{}, err
	}

	return c.contract.GetTrafficCheckHash(ctx, ci.FileSize.Uint64(), ci.Nonce), nil
}

func (c *Controller) GetSpacePrice(ctx context.Context) (uint64, error) {
	out, err := c.contract.Call(ctx, "proxy", "spacePrice")
	if err != nil {
		return 0, err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

func (c *Controller) GetTrafficPrice(ctx context.Context) (uint64, error) {
	out, err := c.contract.Call(ctx, "proxy", "trafficPrice")
	if err != nil {
		return 0, err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

func (c *Controller) BuySpace(ctx context.Context, address string, size uint64) (api.Transaction, error) {
	return c.contract.BuySpace(ctx, address, size)
}

func (c *Controller) BuyTraffic(ctx context.Context, address string, size uint64) (api.Transaction, error) {
	return c.contract.BuyTraffic(ctx, address, size)
}

func (c *Controller) Approve(ctx context.Context, pt api.PayType, buyer string, value *big.Int) (api.Transaction, error) {
	return c.contract.ApproveTsHash(ctx, pt, buyer, value)
}

func (c *Controller) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	balance, err := c.contract.BalanceOf(ctx, address)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Controller) SpacePayInfo(ctx context.Context, address string) (IPayPayment, error) {
	out, err := c.contract.Call(ctx, "proxy", "spacePayInfo", common.HexToAddress(address))
	if err != nil {
		return IPayPayment{}, err
	}

	out0 := *abi.ConvertType(out[0], new(IPayPayment)).(*IPayPayment)

	return out0, err
}

func (c *Controller) TrafficPayInfo(ctx context.Context, address string) (IPayPayment, error) {
	out, err := c.contract.Call(ctx, "proxy", "trafficPayInfo", common.HexToAddress(address))
	if err != nil {
		return IPayPayment{}, err
	}

	out0 := *abi.ConvertType(out[0], new(IPayPayment)).(*IPayPayment)

	return out0, err
}

func (c *Controller) CashSpace(ctx context.Context, buyer string) (string, error) {
	check, err := c.datastore.GetSpaceInfo(ctx, buyer)
	if err != nil {
		return "", err
	}
	return c.contract.CashSpaceCheck(ctx, check)
}

func (c *Controller) CashTraffic(ctx context.Context, buyer string) (string, error) {
	check, err := c.datastore.GetTrafficInfo(ctx, buyer)
	if err != nil {
		return "", err
	}
	return c.contract.CashTrafficCheck(ctx, check)
}

func (c *Controller) Allowance(ctx context.Context, pt api.PayType, address string) (*big.Int, error) {
	return c.contract.Allowance(ctx, pt, address)
}

func (c *Controller) CheckReceipt(ctx context.Context, receipt string) error {
	return c.contract.CheckTrsaction(ctx, receipt)
}
