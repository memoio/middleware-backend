package database

import (
	"math/big"
	"sync"

	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type SendStorage struct {
	lw   sync.Mutex
	ds   store.KVStore
	pool map[string]*StorageCheck
}

func NewSender(ds store.KVStore) *SendStorage {
	ss := &SendStorage{
		ds:   ds,
		pool: make(map[string]*StorageCheck),
	}

	return ss
}

func (s *SendStorage) AddStorage(chain int, address string, st storage.StorageType, size *big.Int, hashid string) error {
	if size.Sign() <= 0 {
		err := logs.DataBaseError{Message: "size should be larger than zero"}
		logger.Error(err)
		return err
	}

	pkey := string(getKey(storagePrefix, address, st, chain))

	s.lw.Lock()
	defer s.lw.Unlock()

	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadStorage(chain, address, st)
		if err != nil {
			return err
		}
		p = schk
		s.pool[pkey] = p
	}

	p.Add(hashid, size)

	data, err := p.Serialize()
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	key := getKey(storagePrefix, address, st, chain)
	s.ds.Put(key, data)

	return nil
}

func (s *SendStorage) DelStorage(chain int, address string, st storage.StorageType, size *big.Int, hashid string) error {
	if size.Sign() <= 0 {
		err := logs.DataBaseError{Message: "size should be larger than zero"}
		logger.Error(err)
		return err
	}

	pkey := string(getKey(storagePrefix, address, st, chain))

	s.lw.Lock()
	defer s.lw.Unlock()

	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadStorage(chain, address, st)
		if err != nil {
			return err
		}
		p = schk
		s.pool[pkey] = p
	}

	p.Del(hashid, size)

	data, err := p.Serialize()
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	key := store.NewKey(address, st)
	s.ds.Put(key, data)

	return nil
}

func (s *SendStorage) GetStorage(chain int, address string, st storage.StorageType) (*big.Int, error) {
	pkey := string(getKey(storagePrefix, address, st, chain))
	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadStorage(chain, address, st)
		if err != nil {
			return nil, err
		}

		s.pool[pkey] = schk
		return schk.Size(), nil
	}

	return p.Size(), nil
}

func (s *SendStorage) ResetStorage(chain int, address string, st storage.StorageType) error {
	schk, err := s.create(chain, address, st)
	if err != nil {
		logger.Error(err)
		return err
	}

	s.pool[address] = schk
	return nil
}

func (s *SendStorage) loadStorage(chain int, address string, st storage.StorageType) (*StorageCheck, error) {
	key := getKey(storagePrefix, address, st, chain)
	data, err := s.ds.Get(key)
	if err != nil {
		schk, err := s.create(chain, address, st)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		return schk, nil
	}

	schk := new(StorageCheck)
	err = schk.Deserialize(data)
	if err != nil {
		return nil, logs.DataBaseError{Message: err.Error()}
	}

	s.pool[address] = schk
	return schk, nil
}

func (s *SendStorage) create(chain int, address string, st storage.StorageType) (*StorageCheck, error) {
	sc := generateCheck(chain, address, st)

	return sc, sc.Save(s.ds)
}

func (s *SendStorage) GetAllStorage() []*StorageCheck {
	var res []*StorageCheck

	for _, sc := range s.pool {
		res = append(res, sc)
	}

	return res
}
