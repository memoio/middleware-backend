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

const (
	storagePrefix = "storage"
	payPrefix     = "pay"
)

type PayCheck struct {
	ChainID int
	Address common.Address
	SType   storage.StorageType
	Size    *big.Int
	Value   *big.Int
	hash    []string
	lw      sync.Mutex
}

func (p *PayCheck) Serialize() ([]byte, error) {
	return cbor.Marshal(p)
}

func (p *PayCheck) Deserialize(b []byte) error {
	return cbor.Unmarshal(b, p)
}

func newPayCheck(chain int, address string, st storage.StorageType) *PayCheck {
	return &PayCheck{
		ChainID: chain,
		Address: common.HexToAddress(address),
		SType:   st,
		Size:    big.NewInt(0),
		Value:   big.NewInt(0),
	}
}

func (s *PayCheck) Save(ds store.KVStore) error {
	key := getKey(payPrefix, s.Address.Hex(), s.SType, s.ChainID)

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

func (s *PayCheck) Add(hash string, size, value *big.Int) error {
	s.lw.Lock()
	defer s.lw.Unlock()

	s.Value.Add(s.Value, value)
	s.Size.Add(s.Size, size)
	s.hash = append(s.hash, hash)
	return nil
}

func (s *PayCheck) Hash() string {
	var res string
	for _, hash := range s.hash {
		res += hash
	}

	return res
}

type StorageCheck struct {
	ChainID int
	Address common.Address
	SType   storage.StorageType
	AddSize *big.Int
	addhash []string
	DelSize *big.Int
	delhash []string
	lw      sync.Mutex
}

func (s *StorageCheck) Serialize() ([]byte, error) {
	return cbor.Marshal(s)
}

func (s *StorageCheck) Deserialize(b []byte) error {
	return cbor.Unmarshal(b, s)
}

// func (s *SendStorage) GetHash()
func generateCheck(chain int, address string, st storage.StorageType) *StorageCheck {
	return &StorageCheck{
		ChainID: chain,
		Address: common.HexToAddress(address),
		SType:   st,
		AddSize: big.NewInt(0),
		DelSize: big.NewInt(0),
	}
}

func (s *StorageCheck) Save(ds store.KVStore) error {
	key := getKey(storagePrefix, s.Address.Hex(), s.SType, s.ChainID)

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

func (s *StorageCheck) Add(hash string, size *big.Int) error {
	s.lw.Lock()
	defer s.lw.Unlock()

	s.AddSize.Add(s.AddSize, size)
	s.addhash = append(s.addhash, hash)
	return nil
}

func (s *StorageCheck) Del(hash string, size *big.Int) error {
	s.lw.Lock()
	defer s.lw.Unlock()

	s.DelSize.Add(s.DelSize, size)
	s.delhash = append(s.delhash, hash)
	return nil
}

func (s *StorageCheck) Size() *big.Int {
	result := new(big.Int).Set(s.AddSize)
	return result.Sub(result, s.DelSize)
}

func (s *StorageCheck) AddHash() string {
	var res string
	for _, hash := range s.addhash {
		res += hash
	}

	return res
}

func (s *StorageCheck) DelHash() string {
	var res string
	for _, hash := range s.delhash {
		res += hash
	}

	return res
}

func getKey(prefix, address string, st storage.StorageType, chain int) []byte {
	return store.NewKey(prefix, address, st.String(), chain)
}
