package database

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/go-mefs-v2/lib/types/store"
)

type UploadPay struct {
	lw sync.Mutex
	ds store.KVStore

	contractAddr common.Address
	sellerAddr   common.Address

	pool map[common.Address]*PayCheck
}

func NewUploaderPay(ds store.KVStore) *UploadPay {
	return &UploadPay{
		ds:           ds,
		contractAddr: contractAddr,
		sellerAddr:   sellerAddr,
		pool:         make(map[common.Address]*PayCheck),
	}
}

func (u *UploadPay) Upload(ctx context.Context, buyer common.Address, sign []byte, nonce, checksize, size *big.Int) error {
	if size.Sign() <= 0 {
		lerr := logs.DataBaseError{Message: "size should be lager than zero"}
		logger.Error(lerr)
		return lerr
	}

	u.lw.Lock()
	defer u.lw.Unlock()

	p, ok := u.pool[buyer]
	if !ok {
		var err error
		p, err = u.loadPay(buyer)
		if err != nil {
			return err
		}
	}

	p.Sign = sign
	p.Duration = 1
	p.Nonce = nonce.Uint64()
	p.Size = checksize.Uint64()
	p.UploadSize += size.Uint64()

	u.pool[buyer] = p
	p.Save(u.ds)

	return nil
}

func (u *UploadPay) create(buyer common.Address) (*PayCheck, error) {
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

func (u *UploadPay) Size(buyer common.Address) uint64 {
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

func (u *UploadPay) loadPay(buyer common.Address) (*PayCheck, error) {
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
