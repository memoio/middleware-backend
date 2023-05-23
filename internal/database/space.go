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

func (s *SendStorage) AddStorage(address string, st storage.StorageType, size *big.Int, hashid string) error {
	if size.Sign() <= 0 {
		err := logs.DataBaseError{Message: "size should be larger than zero"}
		logger.Error(err)
		return err
	}

	pkey := address + st.String()

	s.lw.Lock()
	defer s.lw.Unlock()

	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadStorage(address, st)
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

	key := store.NewKey(address, st)
	s.ds.Put(key, data)

	return nil
}

func (s *SendStorage) DelStorage(address string, st storage.StorageType, size *big.Int, hashid string) error {
	if size.Sign() <= 0 {
		err := logs.DataBaseError{Message: "size should be larger than zero"}
		logger.Error(err)
		return err
	}

	pkey := address + st.String()

	s.lw.Lock()
	defer s.lw.Unlock()

	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadStorage(address, st)
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

func (s *SendStorage) GetStorage(address string, st storage.StorageType) (*big.Int, error) {
	pkey := address + st.String()
	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadStorage(address, st)
		if err != nil {
			return nil, err
		}

		s.pool[pkey] = schk
		return schk.Size(), nil
	}

	return p.Size(), nil
}

func (s *SendStorage) ResetStorage(address string, st storage.StorageType) error {
	schk, err := s.create(address, st)
	if err != nil {
		logger.Error(err)
		return err
	}

	s.pool[address] = schk
	return nil
}

func (s *SendStorage) loadStorage(address string, st storage.StorageType) (*StorageCheck, error) {
	key := store.NewKey(address, st)
	data, err := s.ds.Get(key)
	if err != nil {
		schk, err := s.create(address, st)
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

func (s *SendStorage) create(address string, st storage.StorageType) (*StorageCheck, error) {
	sc := generateCheck(address, st)

	return sc, sc.Save(s.ds)
}

func (s *SendStorage) GetAllStorage() []*StorageCheck {
	var res []*StorageCheck

	for _, sc := range s.pool {
		res = append(res, sc)
	}

	return res
}
