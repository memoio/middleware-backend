package database

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	contractAddr common.Address
	sellerAddr   common.Address
)

func generateCheck(buyer common.Address) (*Check, error) {
	c := &Check{
		ContractAddr: contractAddr,
		Buyer:        buyer,
	}

	return c, nil
}
