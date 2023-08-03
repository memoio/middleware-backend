package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
	com "github.com/memoio/contractsv2/common"
	"github.com/memoio/contractsv2/go_contracts/erc"
	inst "github.com/memoio/contractsv2/go_contracts/instance"
	"github.com/memoio/contractsv2/go_contracts/token"
)

type PackageInfo struct {
	Time    uint64
	Kind    uint8
	Buysize *big.Int
	Amount  *big.Int
	State   uint8
}

type UserBuyPackage struct {
	Starttime uint64
	Endtime   uint64
	Kind      uint8
	Buysize   *big.Int
	Amount    *big.Int
	State     uint8
}

type FlowSize struct {
	Used *big.Int
	Free *big.Int
}

type Contract struct {
	contractAddr common.Address
	endpoint     string
	seller       common.Address

	erc20     common.Address
	tokenAddr common.Address
	proxyAddr common.Address
	storeAddr common.Address
	readAddr  common.Address
	chainID   *big.Int
}

func NewContract(cfc config.ContractConfig) (*Contract, error) {
	instanceAddr, endPoint := com.GetInsEndPointByChain(cfc.Chain)

	client, err := ethclient.DialContext(context.TODO(), endPoint)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		chainID = big.NewInt(666)
	}

	instanceIns, err := inst.NewInstance(instanceAddr, client)
	if err != nil {
		return nil, err
	}

	proxyAddr, err := instanceIns.Instances(&bind.CallOpts{From: instanceAddr}, com.TypeMiddlewareProxy)
	if err != nil {
		return nil, err
	}
	readPayAddr, err := instanceIns.Instances(&bind.CallOpts{From: instanceAddr}, com.TypeReadPay)
	if err != nil {
		return nil, err
	}
	storePayAddr, err := instanceIns.Instances(&bind.CallOpts{From: instanceAddr}, com.TypeStorePay)
	if err != nil {
		return nil, err
	}

	tokenAddr, err := instanceIns.Instances(&bind.CallOpts{From: instanceAddr}, com.TypeToken)
	if err != nil {
		return nil, err
	}
	tokenIns, err := token.NewToken(tokenAddr, client)
	if err != nil {
		return nil, err
	}
	erc20Addr, err := tokenIns.GetTA(&bind.CallOpts{From: com.AdminAddr}, 0)
	if err != nil {
		return nil, err
	}

	seller, err := utils.GetSeller(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Contract{
		contractAddr: instanceAddr,
		endpoint:     endPoint,
		seller:       common.HexToAddress(seller),

		erc20:     erc20Addr,
		tokenAddr: tokenAddr,
		chainID:   chainID,
		proxyAddr: proxyAddr,
		storeAddr: storePayAddr,
		readAddr:  readPayAddr,
	}, nil
}

func (c *Contract) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	res := new(big.Int)
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return res, err
	}
	defer client.Close()

	erc20Ins, err := erc.NewERC20(c.erc20, client)
	if err != nil {
		return res, err
	}

	bal, err := erc20Ins.BalanceOf(&bind.CallOpts{
		From: c.proxyAddr,
	}, common.HexToAddress(addr))
	if err != nil {
		return res, err
	}
	return res.Set(bal), nil
}

func (c *Contract) Call(ctx context.Context, name, method string, args ...interface{}) ([]interface{}, error) {
	var out []interface{}
	err := c.CallContract(ctx, &out, name, method, args...)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return out, lerr
	}

	return out, nil
}

func (c *Contract) GetStorePayHash(ctx context.Context, checksize uint64, nonce *big.Int) string {
	hash := com.GetCashCheckHash(c.storeAddr, c.seller, checksize, nonce)
	return hexutil.Encode(hash)
}

func (c *Contract) GetReadPayHash(ctx context.Context, checksize uint64, nonce *big.Int) string {
	hash := com.GetCashCheckHash(c.readAddr, c.seller, checksize, nonce)
	return hexutil.Encode(hash)
}

func (c *Contract) GetStoreAddr(ctx context.Context) string {
	return c.storeAddr.String()
}

func (c *Contract) GetReadAddr(ctx context.Context) string {
	return c.readAddr.String()
}
