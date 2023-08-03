package contract

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/contractsv2/go_contracts/erc"
	"github.com/memoio/middleware-contracts/go-contracts/control"
	"github.com/memoio/middleware-contracts/go-contracts/proxy"
)

var (
	ksp = filepath.Join("./", "keystore")
)

var logger = logs.Logger("contract")

func createAbi(cabi string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(cabi))
	if err != nil {
		logger.Error(err)
	}
	return parsed
}

func getContractABI(name string) abi.ABI {
	switch name {
	case "control":
		return createAbi(control.ControlABI)
	case "proxy":
		return createAbi(proxy.ProxyABI)
	case "erc20":
		return createAbi(erc.ERC20ABI)
	}
	return abi.ABI{}
}

func (c *Contract) CallContract(ctx context.Context, results *[]interface{}, name, method string, args ...interface{}) error {
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return err
	}
	defer client.Close()

	logger.Infof("CallContract %s %s %s %s", name, method, args, c.proxyAddr)
	if results == nil {
		results = new([]interface{})
	}

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(method, args...)
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{
		To:   &c.proxyAddr,
		Data: encodeData,
	}

	result, err := client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return err
	}

	if len(*results) == 0 {
		res, err := contractABI.Unpack(method, result)
		*results = res
		return err
	}
	res := *results
	return contractABI.UnpackIntoInterface(res[0], method, result)
}

