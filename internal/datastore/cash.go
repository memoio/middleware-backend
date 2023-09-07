package datastore

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

type CashCheck struct {
	lw sync.Mutex

	ds   api.KVStore
	pool map[common.Address]*PayCheck
}

func NewCheckPay(ds api.KVStore) *CashCheck {
	return &CashCheck{
		ds:   ds,
		pool: make(map[common.Address]*PayCheck),
	}
}

func (u *CashCheck) check(ctx context.Context, ct CheckType, info api.CheckInfo) error {
	if info.FileSize.Sign() <= 0 {
		lerr := logs.DataBaseError{Message: "size should be lager than zero"}
		logger.Error(lerr)
		return lerr
	}

	u.lw.Lock()
	defer u.lw.Unlock()

	p, ok := u.pool[info.Buyer]
	if !ok {
		var err error
		p, err = u.loadPay(ctx, info.Buyer)
		if err != nil {
			return err
		}
	}
	if ct == SPACE {
		p.space.Sign = info.Sign
		p.space.Duration = 1
		p.space.Nonce = info.Nonce.Uint64()
		p.space.Size += info.FileSize.Uint64()
	} else {
		p.traffic.Sign = info.Sign
		p.traffic.Duration = 1
		p.traffic.Nonce = info.Nonce.Uint64()
		p.traffic.Size += info.FileSize.Uint64()
	}

	u.pool[info.Buyer] = p
	p.Save(u.ds)

	return nil
}

func (u *CashCheck) create(buyer common.Address) (*PayCheck, error) {
	p := &PayCheck{
		space:        &Check{Size: 0},
		traffic:      &Check{Size: 0},
		Buyer:        buyer,
		ContractAddr: contractAddr,
	}

	return p, p.Save(u.ds)
}

func (u *CashCheck) loadPay(ctx context.Context, buyer common.Address) (*PayCheck, error) {
	key := newKey(buyer.String())
	data, err := u.ds.Get(key)
	if err != nil {
		pchk, err := u.create(buyer)
		if err != nil {
			return nil, err
		}
		u.pool[buyer] = pchk
		return pchk, nil
	}

	pchk := new(PayCheck)
	err = pchk.Deserialize(data)
	if err != nil {
		return nil, err
	}
	u.pool[buyer] = pchk

	return pchk, nil
}

func (u *CashCheck) getCheck(ctx context.Context, ct CheckType, buyer common.Address) (api.CheckInfo, error) {
	res := api.CheckInfo{}

	p, ok := u.pool[buyer]
	if !ok {
		var err error
		p, err = u.loadPay(ctx, buyer)
		if err != nil {
			return res, err
		}
	}
	if ct == SPACE {
		return api.CheckInfo{
			FileSize: big.NewInt(int64(p.space.Size)),
			Sign:     p.space.Sign,
			Nonce:    big.NewInt(int64(p.space.Nonce)),
		}, nil
	} else {
		return api.CheckInfo{
			FileSize: big.NewInt(int64(p.traffic.Size)),
			Sign:     p.traffic.Sign,
			Nonce:    big.NewInt(int64(p.traffic.Nonce)),
		}, nil
	}

}
