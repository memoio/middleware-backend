package filedns

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	com "github.com/memoio/contractsv2/common"
	inst "github.com/memoio/contractsv2/go_contracts/instance"
	"github.com/memoio/did-solidity/go-contracts/proxy"
	"github.com/memoio/go-did/mfile"
	dtypes "github.com/memoio/go-did/types"
)

var (
	blockNumberKey = []byte("block number")
	blockNumber    *big.Int
)

type Dumper struct {
	client          *ethclient.Client
	contractABI     abi.ABI
	contractAddress common.Address
	// store           MapStore

	eventNameMap map[common.Hash]string
	indexedMap   map[common.Hash]abi.Arguments
}

func NewMfileDumper(chain string) (dumper *Dumper, err error) {
	dumper = &Dumper{
		// store:        store,
		eventNameMap: make(map[common.Hash]string),
		indexedMap:   make(map[common.Hash]abi.Arguments),
	}

	instanceAddr, endpoint := com.GetInsEndPointByChain(chain)

	dumper.client, err = ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		return dumper, err
	}

	// new instanceIns
	instanceIns, err := inst.NewInstance(instanceAddr, dumper.client)
	if err != nil {
		return dumper, err
	}

	dumper.contractAddress, err = instanceIns.Instances(&bind.CallOpts{}, com.TypeFileDid)
	if err != nil {
		return dumper, err
	}

	dumper.contractABI, err = abi.JSON(strings.NewReader(proxy.IFileDidABI))
	if err != nil {
		return dumper, err
	}

	for name, event := range dumper.contractABI.Events {
		dumper.eventNameMap[event.ID] = name

		var indexed abi.Arguments
		for _, arg := range dumper.contractABI.Events[name].Inputs {
			if arg.Indexed {
				indexed = append(indexed, arg)
			}
		}
		dumper.indexedMap[event.ID] = indexed
	}

	return dumper, nil
}

func (d *Dumper) DumpMfileDID() error {
	// var last *big.Int
	for {
		events, err := d.client.FilterLogs(context.TODO(), ethereum.FilterQuery{
			FromBlock: blockNumber,
			Addresses: []common.Address{d.contractAddress},
		})
		if err != nil {
			return err
		}

		log.Printf("Got %d logs", len(events))
		for _, event := range events {
			eventName, ok1 := d.eventNameMap[event.Topics[0]]
			if !ok1 {
				continue
			}
			switch eventName {
			case "RegisterMfileDid":
				err = d.HandleRegisterMfileDid(event)
			case "DeactivateMfileDid":
				err = d.HandleDeactivateMfileDid(event)
			case "ChangeFtype":
				err = d.HandleChangeFtype(event)
			case "ChangeController":
				err = d.HandleChangeController(event)
			case "ChangePrice":
				err = d.HandleChangePrice(event)
			case "ChangeKeywords":
				err = d.HandleChangeKeywords(event)
			case "BuyRead", "GrantRead":
				err = d.HandleAddRead(event)
			case "DeactivateRead":
				err = d.HandleDeactivateRead(event)
			default:
				continue
			}
			if err != nil {
				log.Println(err.Error())
				break
			}

			blockNumber = big.NewInt(int64(event.BlockNumber) + 1)
		}

		Searcher.Flush()
		DIDStore.SetLastBlockNumber(blockNumber)

		time.Sleep(2 * time.Minute)
	}
	// return nil
}

func (d *Dumper) unpack(log types.Log, out interface{}) error {
	eventName := d.eventNameMap[log.Topics[0]]
	indexed := d.indexedMap[log.Topics[0]]
	err := d.contractABI.UnpackIntoInterface(out, eventName, log.Data)
	if err != nil {
		return err
	}

	return abi.ParseTopics(out, indexed, log.Topics[1:])
}

type RegisterMfileDid struct {
	MfileDid string
}

func (d *Dumper) HandleRegisterMfileDid(log types.Log) error {
	var out RegisterMfileDid
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	resolver, err := mfile.NewMfileDIDResolver("dev")
	if err != nil {
		return err
	}
	document, err := resolver.Resolve("did:mfile:" + out.MfileDid)
	if err != nil {
		return err
	}
	document.Read = nil

	AddShareFile(document.ID.Identifier, document.Keywords)
	return DIDStore.Set(crypto.Keccak256Hash([]byte(out.MfileDid)), *document)
}

type DeactivateMfileDid struct {
	MfileDid   string
	Deactivate bool
}

func (d *Dumper) HandleDeactivateMfileDid(log types.Log) error {
	var out DeactivateMfileDid
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	if !out.Deactivate {
		return nil
	}
	return DIDStore.Delete(crypto.Keccak256Hash([]byte(out.MfileDid)))
}

type ChangeFtype struct {
	MfileDid common.Hash
	FType    uint8
}

func (d *Dumper) HandleChangeFtype(log types.Log) error {
	var out ChangeFtype
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	document, err := DIDStore.Get(out.MfileDid)
	if err != nil {
		return err
	}

	if out.FType == 0 {
		document.Type = "private"
	} else {
		document.Type = "public"
	}
	return DIDStore.Set(out.MfileDid, document)
}

type ChangeController struct {
	MfileDid   common.Hash
	Controller string
}

func (d *Dumper) HandleChangeController(log types.Log) error {
	var out ChangeController
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	document, err := DIDStore.Get(out.MfileDid)
	if err != nil {
		return err
	}

	document.Controller.Identifier = out.Controller
	return DIDStore.Set(out.MfileDid, document)
}

type ChangePrice struct {
	MfileDid common.Hash
	Price    *big.Int
}

func (d *Dumper) HandleChangePrice(log types.Log) error {
	var out ChangePrice
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	document, err := DIDStore.Get(out.MfileDid)
	if err != nil {
		return err
	}

	document.Price = out.Price.Int64()
	return DIDStore.Set(out.MfileDid, document)
}

type ChangeKeywords struct {
	MfileDid common.Hash
	Keywords []string
}

func (d *Dumper) HandleChangeKeywords(log types.Log) error {
	var out ChangeKeywords
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	document, err := DIDStore.Get(out.MfileDid)
	if err != nil {
		return err
	}

	document.Keywords = out.Keywords
	UpdateShareFile(document.ID.Identifier, document.Keywords)
	return DIDStore.Set(out.MfileDid, document)
}

type AddRead struct {
	MfileDid common.Hash
	MemoDid  string
}

func (d *Dumper) HandleAddRead(log types.Log) error {
	var out AddRead
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	document, err := DIDStore.Get(out.MfileDid)
	if err != nil {
		return err
	}

	did, _ := dtypes.ParseMemoDID("did:memo:" + out.MemoDid)
	document.Read = append(document.Read, *did)
	return DIDStore.Set(out.MfileDid, document)
}

type DeactivateRead struct {
	MfileDid common.Hash
	MemoDid  string
}

func (d *Dumper) HandleDeactivateRead(log types.Log) error {
	var out DeactivateRead
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	document, err := DIDStore.Get(out.MfileDid)
	if err != nil {
		return err
	}

	for index, read := range document.Read {
		if out.MemoDid == read.Identifier {
			document.Read = append(document.Read[:index], document.Read[index+1:]...)
			break
		}
	}
	return DIDStore.Set(out.MfileDid, document)
	// return nil
}
