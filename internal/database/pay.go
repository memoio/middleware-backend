package database

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type CheckPay struct {
	lw sync.Mutex
	ds store.KVStore

	contractAddr common.Address
	sellerAddr   common.Address

	pool map[common.Address]*PayCheck
}

func NewCheckPay(ds store.KVStore) *CheckPay {
	return &CheckPay{
		ds:           ds,
		contractAddr: contractAddr,
		sellerAddr:   sellerAddr,
		pool:         make(map[common.Address]*PayCheck),
	}
}

func (u *CheckPay) Check(ctx context.Context, info api.CheckInfo) error {
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
		p, err = u.loadPay(info.Buyer)
		if err != nil {
			return err
		}
	}

	p.Sign = info.Sign
	p.Duration = 1
	p.Nonce = info.Nonce.Uint64()
	p.Size = info.CheckSize.Uint64()
	p.UploadSize += info.FileSize.Uint64()

	u.pool[info.Buyer] = p
	p.Save(u.ds)

	return nil
}

func (u *CheckPay) create(buyer common.Address) (*PayCheck, error) {
	chk, err := generateCheck(buyer)
	if err != nil {
		return nil, err
	}

	p := &PayCheck{
		Check:      *chk,
		UploadSize: 0,
	}

	return p, p.Save(u.ds)
}

func (u *CheckPay) Size(buyer common.Address) uint64 {
	p, ok := u.pool[buyer]
	if !ok {
		var err error
		p, err = u.loadPay(buyer)
		if err != nil {
			return 0
		}
	}

	return p.UploadSize
}

func (u *CheckPay) loadPay(buyer common.Address) (*PayCheck, error) {
	key := store.NewKey(u.contractAddr.String(), buyer.String())
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

func (u *CheckPay) getCheck(ctx context.Context, buyer common.Address) api.CheckInfo {
	res := api.CheckInfo{}

	p, ok := u.pool[buyer]
	if !ok {
		var err error
		p, err = u.loadPay(buyer)
		if err != nil {
			return res
		}
	}

	return api.CheckInfo{
		CheckSize: big.NewInt(int64(p.Check.Size)),
		Sign:      p.Check.Sign,
		Nonce:     big.NewInt(int64(p.Check.Nonce)),
	}
}
