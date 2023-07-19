package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/logs"
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
	contractAddr     common.Address
	endpoint         string
	gatewayAddr      common.Address
	gatewaySecretKey string

	proxyAddr common.Address
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

	return &Contract{
		contractAddr:     instanceAddr,
		endpoint:         endPoint,
		gatewaySecretKey: cfc.GatewaySecretKey,
		chainID:          chainID,
		proxyAddr:        proxyAddr,
	}, nil
}

func (c *Contract) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	res := new(big.Int)
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return res, err
	}
	defer client.Close()

	instanceIns, err := inst.NewInstance(c.contractAddr, client)
	if err != nil {
		return res, err
	}

	tokenAddr, err := instanceIns.Instances(&bind.CallOpts{From: com.AdminAddr}, com.TypeToken)
	if err != nil {
		return res, err
	}
	tokenIns, err := token.NewToken(tokenAddr, client)
	if err != nil {
		return res, err
	}
	erc20Addr, err := tokenIns.GetTA(&bind.CallOpts{From: com.AdminAddr}, 0)
	if err != nil {
		return res, err
	}

	erc20Ins, err := erc.NewERC20(erc20Addr, client)
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

func (c *Contract) CheckContract() error {
	privateKey, err := crypto.HexToECDSA(c.gatewaySecretKey)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to decode gateway sk: %v", err)}
		logger.Error(lerr)
		return lerr
	}

	pk := privateKey.Public()
	pubKeyECDSA, ok := pk.(*ecdsa.PublicKey)

	if !ok {
		lerr := logs.ContractError{Message: "error casting public key to ECDSA"}
		logger.Error(lerr)
		return lerr
	}
	gatewayaddr := crypto.PubkeyToAddress(*pubKeyECDSA)
	if gatewayaddr != c.gatewayAddr {
		lerr := logs.ContractError{Message: fmt.Sprintf("gateway address and private key do not match %s", gatewayaddr)}
		logger.Error(lerr)
		return lerr
	}

	return nil
}
