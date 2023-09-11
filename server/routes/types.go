package routes

import "math/big"

type IPayPayment struct {
	Nonce    *big.Int
	Balance  *big.Int
	SizeByte uint64
	FreeByte uint64
	Expire   uint64
}

func toInt64(s string) int64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Int64()
}

func toUint64(s string) uint64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Uint64()
}

func toBigInt(s string) *big.Int {
	b := new(big.Int)
	b.SetString(s, 10)
	return b
}
