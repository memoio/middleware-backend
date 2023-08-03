package utils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/wallet"
)

func GetSeller(ctx context.Context) (string, error) {
	ks, err := wallet.NewKeyRepo(api.KeystorePath)
	if err != nil {
		return "", err
	}
	wl := wallet.New(ks)
	seller := config.Cfg.Contract.SellerAddr

	if seller == "" {
		return "", fmt.Errorf("seller is empty")
	}
	res, err := wl.WalletHas(ctx, common.HexToAddress(seller))
	if err != nil || !res {
		return "", fmt.Errorf("sellet not in wallet %s", err)
	}
	return seller, nil
}
