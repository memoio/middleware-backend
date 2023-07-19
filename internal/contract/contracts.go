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
	"github.com/memoio/contractsv2/go_contracts/erc"
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
}

func NewContract(cfc config.ContractConfig) *Contract {
	return &Contract{
		contractAddr:     common.HexToAddress(cfc.ContractAddr),
		endpoint:         cfc.Endpoint,
		gatewayAddr:      common.HexToAddress(cfc.GatewayAddr),
		gatewaySecretKey: cfc.GatewaySecretKey,
	}
}
func NewContracts(cfc map[int]config.ContractConfig) map[int]*Contract {
	res := make(map[int]*Contract)

	for chainid, cfg := range cfc {
		res[chainid] = &Contract{
			contractAddr:     common.HexToAddress(cfg.ContractAddr),
			endpoint:         cfg.Endpoint,
			gatewayAddr:      common.HexToAddress(cfg.GatewayAddr),
			gatewaySecretKey: cfg.GatewaySecretKey,
		}
	}

	return res
}

func (c *Contract) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	res := new(big.Int)
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return res, err
	}
	defer client.Close()

	erc20Ins, err := erc.NewERC20(c.contractAddr, client)
	if err != nil {
		return res, err
	}

	bal, err := erc20Ins.BalanceOf(&bind.CallOpts{
		From: c.gatewayAddr,
	}, common.HexToAddress(addr))
	if err != nil {
		return res, err
	}
	return res.Set(bal), nil
}

// func (c *Contract) GetStoreAllSize() *big.Int {
// 	var out []interface{}
// 	err := c.CallContract(&out, "getStoreAllSize")
// 	if err != nil {
// 		logger.Error(err)
// 		return nil
// 	}

// 	available := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
// 	return available
// }

func (c *Contract) Call(name string, args ...interface{}) ([]interface{}, error) {
	var out []interface{}
	err := c.CallContract(&out, name, args...)
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
