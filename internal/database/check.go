package database

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type Check struct {
	Value  *big.Int
	Hashid string
}

type StorageCheck struct {
	Address common.Address
	SType   storage.StorageType
	Size    *big.Int
	lw      sync.Mutex
	Ch      []Check
}

func (s *StorageCheck) Serialize() ([]byte, error) {
	return cbor.Marshal(s)
}

func (s *StorageCheck) Deserialize(b []byte) error {
	return cbor.Unmarshal(b, s)
}

func generateCheck(address common.Address, st storage.StorageType) *StorageCheck {
	return &StorageCheck{
		Address: address,
		SType:   st,
		Size:    big.NewInt(0),
		Ch:      make([]Check, 0),
	}
}

func (s *StorageCheck) Put(ch Check) {
	s.lw.Lock()
	defer s.lw.Unlock()
	s.Ch = append(s.Ch, ch)
	s.Size.Add(s.Size, ch.Value)
}

func (s *StorageCheck) Get() Check {
	s.lw.Lock()
	defer s.lw.Unlock()
	if s.Len() == 0 {
		return Check{}
	}
	ch := s.Ch[0]
	s.Ch = s.Ch[1:]
	s.Size.Sub(s.Size, ch.Value)
	return ch
}

func (s *StorageCheck) Len() int {
	return len(s.Ch)
}

func (s *StorageCheck) Save(ds store.KVStore) error {
	key := store.NewKey(s.Address.String(), s.SType.String())

	data, err := s.Serialize()
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	err = ds.Put(key, data)
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	return nil
}

type DeleteCheck struct {
	Address     common.Address
	StorageType storage.StorageType
	Ch          []Check
}
