package wallet

import (
	"context"
	"crypto/ecdsa"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("wallet")

type LocalWallet struct {
	lw       sync.Mutex
	accounts map[common.Address]*ecdsa.PrivateKey
	keystore api.Keystore
}

func New(ks api.Keystore) *LocalWallet {
	lw := &LocalWallet{
		keystore: ks,
		accounts: make(map[common.Address]*ecdsa.PrivateKey),
	}

	return lw
}

func (w *LocalWallet) Find(addr common.Address) (*ecdsa.PrivateKey, error) {
	w.lw.Lock()
	defer w.lw.Unlock()

	pi, ok := w.accounts[addr]
	if ok {
		return pi, nil
	}

	ki, err := w.keystore.Get(addr.String())
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
	}

	pi, err = crypto.ToECDSA(ki)
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
	}

	w.accounts[addr] = pi

	return pi, nil
}

func (w *LocalWallet) WalletNew(ctx context.Context) (common.Address, error) {
	var addr common.Address
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return addr, lerr
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	addr, err = PrivateToAddr(privateKey)
	if err != nil {
		return addr, err
	}

	err = w.keystore.Put(addr.String(), privateKeyBytes)
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return addr, lerr
	}

	w.lw.Lock()
	w.accounts[addr] = privateKey
	w.lw.Unlock()

	return addr, nil
}

func (w *LocalWallet) WalletList(ctx context.Context) ([]common.Address, error) {
	as, err := w.keystore.List()
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
	}

	out := make([]common.Address, 0, len(as))

	for _, s := range as {
		addr := common.HexToAddress(s)
		out = append(out, addr)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].String() < out[j].String()
	})

	return out, nil
}

func (w *LocalWallet) WalletHas(ctx context.Context, addr common.Address) (bool, error) {
	_, err := w.keystore.Get(addr.String())
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return false, lerr
	}
	return true, nil
}

func (w *LocalWallet) WalletDelete(ctx context.Context, addr common.Address) error {
	err := w.keystore.Delete(addr.String())
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	w.lw.Lock()
	delete(w.accounts, addr)
	w.lw.Unlock()

	return nil
}

func (w *LocalWallet) WalletExport(ctx context.Context, addr common.Address) ([]byte, error) {
	ki, err := w.keystore.Get(addr.String())
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
	}

	return ki, nil
}

func (w *LocalWallet) WalletImport(ctx context.Context, ki []byte) (common.Address, error) {
	var addr common.Address

	pi, err := crypto.ToECDSA(ki)
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return addr, lerr
	}
	addr, err = PrivateToAddr(pi)
	if err != nil {
		return addr, err
	}
	err = w.keystore.Put(addr.String(), ki)
	if err != nil {
		lerr := logs.WalletError{Message: err.Error()}
		logger.Error(lerr)
		return addr, lerr
	}

	w.lw.Lock()
	w.accounts[addr] = pi
	w.lw.Unlock()

	return addr, nil
}
