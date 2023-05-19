package database

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type SendStorage struct {
	lw   sync.Mutex
	ds   store.KVStore
	pool map[common.Address]*StorageCheck
}

func NewSender(ds store.KVStore) *SendStorage {
	ss := &SendStorage{
		ds:   ds,
		pool: make(map[common.Address]*StorageCheck),
	}

	return ss
}

func (s *SendStorage) AddStorage(address common.Address, st storage.StorageType, size *big.Int, hashid string) error {
	if size.Sign() <= 0 {
		return logs.DataBaseError{Message: "deposit value should be larger than zero"}
	}

	s.lw.Lock()
	defer s.lw.Unlock()

	p, ok := s.pool[address]
	if !ok {
		key := store.NewKey(address.String(), st)
		data, err := s.ds.Get(key)
		if err == nil {
			schk := new(StorageCheck)
			err = schk.Deserialize(data)
			if err != nil {
				return logs.DataBaseError{Message: err.Error()}
			}
			p = schk
			s.pool[address] = p
		} else {
			schk, err := s.create(address, st)
			if err != nil {
				return logs.DataBaseError{Message: err.Error()}
			}
			p = schk
			s.pool[address] = p
		}

	}

	ssize := new(big.Int).Set(p.Size)
	ssize.Add(ssize, size)

	p.Size.Set(ssize)
	p.Put(Check{
		Value:  size,
		Hashid: hashid,
	})

	data, err := p.Serialize()
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	key := store.NewKey(address.String(), st)
	s.ds.Put(key, data)

	return nil
}

func (s *SendStorage) GetStorage(address common.Address, st storage.StorageType) (*big.Int, error) {
	p, ok := s.pool[address]
	if !ok {
		key := store.NewKey(address.String(), st)
		data, err := s.ds.Get(key)
		if err != nil {
			schk, err := s.create(address, st)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			s.pool[address] = schk
		}

		schk := new(StorageCheck)
		err = schk.Deserialize(data)
		if err != nil {
			return nil, logs.DataBaseError{Message: err.Error()}
		}

		s.pool[address] = schk
		return p.Size, nil
	}
	return p.Size, nil
}

func (s *SendStorage) create(address common.Address, st storage.StorageType) (*StorageCheck, error) {
	sc := generateCheck(address, st)

	return sc, sc.Save(s.ds)
}
