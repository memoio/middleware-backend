package database

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type Check struct {
	ContractAddr common.Address
	Buyer        common.Address
	Nonce        uint64
	Size         uint64
	Duration     uint64
	Sign         []byte
}

type PayCheck struct {
	Check
	UploadSize uint64
}

func (p *PayCheck) Serialize() ([]byte, error) {
	return cbor.Marshal(p)
}

func (p *PayCheck) Deserialize(b []byte) error {
	return cbor.Unmarshal(b, p)
}

func (p *PayCheck) Save(ds store.KVStore) error {
	key := store.NewKey(p.ContractAddr.String(), p.Buyer.String())
	data, err := p.Serialize()
	if err != nil {
		return err
	}
	return ds.Put(key, data)
}
