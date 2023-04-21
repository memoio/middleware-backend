package gateway

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/contract"
	db "github.com/memoio/backend/global/database"
	"github.com/memoio/backend/internal/storage"
	"golang.org/x/crypto/sha3"
)

const (
	payTopic     = "0xc0e3b3bf3b856068b6537f07e399954cb5abc4fade906ee21432a8ded3c36ec8"
	storageTopic = "0x01bed2f7b5b5c577e49071502e2c9985655c9ca5c1ab432156d99d199f6b1912"

	payratio = 100
)

var (
	contractAddr     = common.HexToAddress("0x2A0B376CC39eB2019e43207d00ee2c34878ca36D")
	endpoint         = "https://chain.metamemo.one:8501"
	GatewayAddr      = common.HexToAddress("0x31e7829Ea2054fDF4BCB921eDD3a98a825242267")
	GatewaySecretKey = "8a87053d296a0f0b4600173773c8081b12917cef7419b2675943b0aa99429b62"
)

func (g Gateway) QueryPrice(ctx context.Context, address, size, time string) (string, error) {
	err := g.getMemofs()
	if err != nil {
		return "", err
	}
	log.Println(time)

	if time == "" {
		time = "365"
	}

	price, err := g.Mefs.QueryPrice(ctx)
	if err != nil {
		return "", err
	}

	pr := new(big.Int)
	pr.SetString(price, 10)

	ssize := new(big.Int)
	ssize.SetString(size, 10)

	stime := new(big.Int)
	stime.SetString(time, 10)

	if stime.Cmp(big.NewInt(365)) < 0 {
		return "", StorageError{Message: "at least storage 365 days"}
	}
	stime.Mul(stime, big.NewInt(86400))

	bi, err := g.Mefs.GetBucketInfo(ctx, address)
	if err != nil {
		log.Println("get balance info error")
		return "", StorageError{Message: "get balance info"}
	}

	segment := new(big.Int)
	segment.Mul(ssize, big.NewInt(int64(bi.DataCount+bi.ParityCount)))

	amount := new(big.Int)
	amount.Mul(pr, stime)
	amount.Mul(amount, segment)
	amount.Div(amount, big.NewInt(248000))
	amount.Div(amount, big.NewInt(int64(bi.DataCount)))

	return amount.String(), nil
}

func (g Gateway) Pay(ctx context.Context, to common.Address, cid string, amount, size *big.Int) bool {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return false
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, GatewayAddr)
	if err != nil {
		return false
	}
	log.Println("nonce: ", nonce)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("chainID: ", chainID)

	storeOrderPayFnSignature := []byte("storeOrderpay(address,string,uint256,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(storeOrderPayFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedHashLen := common.LeftPadBytes(big.NewInt(int64(len([]byte(cid)))).Bytes(), 32)
	paddedHashOffset := common.LeftPadBytes(big.NewInt(32*4).Bytes(), 32)
	paddedMd5 := common.LeftPadBytes([]byte(cid), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	PaddedSize := common.LeftPadBytes(size.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedHashOffset...)
	data = append(data, paddedAmount...)
	data = append(data, PaddedSize...)
	data = append(data, paddedHashLen...)
	data = append(data, paddedMd5...)

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(1000)
	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, data)

	privateKey, err := crypto.HexToECDSA(GatewaySecretKey)
	if err != nil {
		log.Println("get privateKey error: ", err)
		return false
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println("signedTx error: ", err)
		return false
	}

	return g.sendTransaction(ctx, signedTx, "pay")
}

func (g *Gateway) sendTransaction(ctx context.Context, signedTx *types.Transaction, ttype string) bool {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Println(err)
		return false
	}
	defer client.Close()

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println("waiting tx complete...")
	time.Sleep(30 * time.Second)

	receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
	if err != nil {
		log.Println(err)
		return false
	}
	if receipt.Status != 1 {
		log.Println("Status not right")
		log.Println(receipt.Logs)
		log.Println(receipt)
		return false
	}

	if len(receipt.Logs) == 0 {
		log.Println("no logs")
		return false
	}

	if len(receipt.Logs[0].Topics) == 0 {
		log.Println("no topics")
		return false
	}
	var topic string
	switch ttype {
	case "storage":
		topic = storageTopic
	case "pay":
		topic = payTopic
	}

	if receipt.Logs[0].Topics[0].String() != topic {
		log.Println("topic not right: ", receipt.Logs[0].Topics[0].String())
		return false
	}

	return true
}

func (g *Gateway) verify(ctx context.Context, storage storage.StorageType, address, date, cid string, size *big.Int) bool {
	flag := g.memverify(ctx, storage, address, cid, size)
	if !flag {
		return g.perverify(ctx, address, date, cid, size)
	}

	return true
}

func (g *Gateway) memverify(ctx context.Context, storage storage.StorageType, address, cid string, size *big.Int) bool {
	if !g.checkStorage(ctx, storage, address, size) {
		return false
	}
	return g.updateStorage(ctx, address, cid, size)
}

func (g *Gateway) perverify(ctx context.Context, address, date, cid string, size *big.Int) bool {
	price, err := g.QueryPrice(ctx, address, size.String(), date)
	if err != nil {
		log.Println("price error", err)
		return false
	}
	pri := new(big.Int)
	pri.SetString(price, 10)
	pri.Mul(pri, big.NewInt(payratio))

	balance := contract.BalanceOf(ctx, address)

	log.Println("Price", price)
	log.Println("Balance", balance)

	if balance.Cmp(pri) < 0 {
		log.Printf("allow: %d, price: %d, allowance not enough\n", balance, pri)
		return false
	}

	return g.Pay(ctx, contractAddr, cid, pri, size)
}

func (g *Gateway) checkStorage(ctx context.Context, storage storage.StorageType, address string, size *big.Int) bool {
	si, err := g.GetPkgSize(ctx, storage, address)
	if err != nil {
		return false
	}

	return si.Buysize+si.Free > size.Int64()+si.Used
}

func (g *Gateway) updateStorage(ctx context.Context, address, cid string, size *big.Int) bool {
	pi := db.PkgInfo{
		Address:   address,
		Hashid:    cid,
		Size:      size.Int64(),
		IsUpdated: false,
		UTime:     time.Now(),
	}

	err := pi.Insert()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
