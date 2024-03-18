package config

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type DataAccessConfig struct {
	UserSecurityKey      string
	SubmitterSecurityKey string
}

func newDefaultDataAccessConfig() DataAccessConfig {
	sk1, _ := crypto.GenerateKey()
	sk2, _ := crypto.GenerateKey()

	return DataAccessConfig{
		UserSecurityKey:      hex.EncodeToString(crypto.FromECDSA(sk1)),
		SubmitterSecurityKey: hex.EncodeToString(crypto.FromECDSA(sk2)),
	}
}
