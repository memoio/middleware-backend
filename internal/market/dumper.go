package market

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// blockNumberKey = []byte("block number")
	blockNumber *big.Int = big.NewInt(0)
	// NFTABI      string
)

type Dumper struct {
	client          *ethclient.Client
	contractABI     abi.ABI
	contractAddress common.Address

	eventNameMap map[common.Hash]string
	indexedMap   map[common.Hash]abi.Arguments
}

func NewDriveNFTDumper() (dumper *Dumper, err error) {
	dumper = &Dumper{
		contractAddress: common.HexToAddress("0xd895f9cb9fcBb6fC9fE2c1B1041E506B3247Bf09"),
		eventNameMap:    make(map[common.Hash]string),
		indexedMap:      make(map[common.Hash]abi.Arguments),
	}

	dumper.client, err = ethclient.DialContext(context.TODO(), "https://chain.metamemo.one:8501")
	if err != nil {
		return dumper, err
	}

	file, err := os.Open("./data/DriveNFT.json")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	dumper.contractABI, err = abi.JSON(file)
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

func (d *Dumper) DumperDriveNFT() error {
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
			case "Transfer":
				err = d.HandleTransferNFT(event)
			default:
				continue
			}
			if err != nil {
				log.Println(err.Error())
				break
			}

			blockNumber = big.NewInt(int64(event.BlockNumber) + 1)
		}

		if len(events) > 0 {
			SetLastBlockNumber(blockNumber.Int64())
			FlushSearcher()
		}

		time.Sleep(5 * time.Minute)
	}
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

type Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type NFTData struct {
	DataUrl     string `json:"dataurl"`
	Mid         string `json:"mid"`
	Description string `json:"desc"`
	Name        string `json:"name"`
	// Time        string `json:"time"`
	// Class       string `json:"class"`
	// Price       string `json:"price"`
	// Currency    string `json:"currency"`
	// Size        int    `json:"size"`
	// FileType    string `json:"filetype"`
	// FileName    string `json:"filename"`
	// ExternUrl   string `json:"exturl"`
}

func (d *Dumper) HandleTransferNFT(log types.Log) error {
	var out Transfer
	err := d.unpack(log, &out)
	if err != nil {
		return err
	}

	if out.From.Big().Cmp(big.NewInt(0)) == 0 {
		results, err := d.CallContract("tokenURI", out.TokenId)
		if err != nil {
			return err
		}

		var info NFTData
		nftInfo := results[0].(string)
		err = json.Unmarshal([]byte(nftInfo), &info)
		if err != nil {
			return err
		}

		AddNFT(out.TokenId.Int64(), info.Description, info.Name)

		return MintNFT(out.TokenId.Int64(), out.To.Hex(), info)
	}

	if out.To.Big().Cmp(big.NewInt(0)) == 0 {
		DeleteNFT(out.TokenId.Int64())

		return BurnNFT(out.TokenId.Int64())
	}

	return TransferNFT(out.TokenId.Int64(), out.To.Hex())
}

func (d *Dumper) CallContract(method string, params ...interface{}) ([]interface{}, error) {
	data, err := d.contractABI.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	callMSG := ethereum.CallMsg{
		To:   &d.contractAddress,
		Data: data,
	}

	result, err := d.client.CallContract(context.TODO(), callMSG, nil)
	if err != nil {
		return nil, err
	}

	res, err := d.contractABI.Unpack(method, result)
	if err != nil {
		return nil, err
	}

	return res, nil
}
