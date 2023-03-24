package contract

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	getPkgSizeAbi = `[{"constant":true,"inputs":[{"name":"to","type":"address"}],"name":"getPkgSize","outputs":[{"name":"used","type":"uint256"},{"name":"available","type":"uint256"},{"name":"total","type":"uint256"},{"name":"expires","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]`
)

func createAbi(cabi string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(cabi))
	if err != nil {
		fmt.Println(err)
	}
	return parsed
}
