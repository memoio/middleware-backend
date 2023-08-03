package database

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/config"
)

var (
	contractAddr common.Address
	sellerAddr   common.Address
)

func init() {
	sellerAddr = common.HexToAddress(config.Cfg.Contract.SellerAddr)
}

func generateCheck(buyer common.Address) (*Check, error) {
	c := &Check{
		ContractAddr: contractAddr,
		Buyer:        buyer,
	}

	return c, nil
}
