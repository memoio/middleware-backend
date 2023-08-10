package database

import (
	"math/big"
	"sync"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type SendPay struct {
	lw   sync.Mutex
	ds   store.KVStore
	pool map[string]*PayCheck
}

func NewSenderPay(ds store.KVStore) *SendPay {
	return &SendPay{
		ds:   ds,
		pool: make(map[string]*PayCheck),
	}
}

func (s *SendPay) AddPay(chain int, address string, st api.StorageType, size, value *big.Int, hashid string) error {
	if size.Sign() <= 0 || value.Sign() <= 0 {
		err := logs.DataBaseError{Message: "size or amount should be larger than zero"}
		logger.Error(err)
		return err
	}

	pkey := string(getKey(payPrefix, address, st, chain))

	s.lw.Lock()
	defer s.lw.Unlock()

	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadPay(chain, address, st)
		if err != nil {
			return err
		}
		p = schk
		s.pool[pkey] = p
	}

	p.Add(hashid, size, value)

	data, err := p.Serialize()
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	key := getKey(payPrefix, address, st, chain)
	s.ds.Put(key, data)

	return nil
}

func (s *SendPay) loadPay(chain int, address string, st api.StorageType) (*PayCheck, error) {
	key := getKey(payPrefix, address, st, chain)
	data, err := s.ds.Get(key)
	if err != nil {
		schk, err := s.create(chain, address, st)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		return schk, nil
	}

	schk := new(PayCheck)
	err = schk.Deserialize(data)
	if err != nil {
		return nil, logs.DataBaseError{Message: err.Error()}
	}

	s.pool[address] = schk
	return schk, nil
}

func (s *SendPay) create(chain int, address string, st api.StorageType) (*PayCheck, error) {
	sc := newPayCheck(chain, address, st)

	return sc, sc.Save(s.ds)
}

func (s *SendPay) GetAllStorage() []*PayCheck {
	var res []*PayCheck

	for _, sc := range s.pool {
		res = append(res, sc)
	}

	return res
}

func (s *SendPay) ResetPay(chain int, address string, st api.StorageType) error {
	pchk, err := s.create(chain, address, st)
	if err != nil {
		logger.Error(err)
		return err
	}

	s.pool[address] = pchk
	return nil
}

func (s *SendPay) Size(chain int, address string, st api.StorageType) (*big.Int, error) {
	pkey := string(getKey(payPrefix, address, st, chain))
	p, ok := s.pool[pkey]
	if !ok {
		schk, err := s.loadPay(chain, address, st)
		if err != nil {
			return nil, err
		}

		s.pool[pkey] = schk
		return schk.Value, nil
	}

	return p.Value, nil
}
