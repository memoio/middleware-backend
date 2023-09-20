package datastore

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
)

var (
	contractAddr common.Address
	sellerAddr   common.Address
)

func init() {
	sellerAddr = common.HexToAddress(config.Cfg.Contract.SellerAddr)
}

type Check struct {
	Nonce    uint64
	Size     uint64
	Duration uint64
	Sign     []byte
}

type PayCheck struct {
	ContractAddr common.Address
	Buyer        common.Address
	space        Check
	traffic      Check
}

func (p *PayCheck) Serialize() ([]byte, error) {
	return cbor.Marshal(p)
}

func (p *PayCheck) Deserialize(b []byte) error {
	return cbor.Unmarshal(b, p)
}

// save paycheck into ds
func (p *PayCheck) Save(ds api.KVStore) error {
	key := newKey(p.Buyer.String())
	data, err := p.Serialize()
	if err != nil {
		return err
	}
	return ds.Put(key, data)
}
