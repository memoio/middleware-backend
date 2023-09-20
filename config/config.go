package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

const CONFIGPATH = "./config.json"

var Cfg *Config

type Config struct {
	SwagHost    string                 `json:"swaghost"`
	Storage     StorageConfig          `json:"storage"`
	Contracts   map[int]ContractConfig `json:"contracts"`
	Contract    ContractConfig         `json:"contract"`
	SecurityKey string                 `json:"securityKey"`
	Domain      string                 `json:"domain"`
	EthDriveUrl string                 `json:"ethDriveUrl"`
	// LensAPIUrl  string                 `json:"lensAPIUrl"`
}

type StorageConfig struct {
	Mefs        MefsConfig       `json:"mefs"`
	Ipfs        IpfsConfig       `json:"ipfs"`
	Prices      map[string]int64 `json:"prices"`
	TrafficCost int64            `json:"traffic_cost"`
}

type MefsConfig struct {
	Api   string `json:"api"`
	Token string `json:"token"`
}

type IpfsConfig struct {
	Host string `json:"host"`
}

func newDefaultIpfsConfig() IpfsConfig {
	return IpfsConfig{
		Host: "127.0.0.1:5002",
	}
}

func newDefaultMefsConfig() MefsConfig {
	return MefsConfig{
		Api:   "/ip4/192.168.1.46/tcp/26812",
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.OzYt9dIYMEuUoaEQeXah2wkPUZ1O6Yya7mwuuIAP89s",
	}
}

func newDefaultStorageConfig() StorageConfig {
	return StorageConfig{
		Mefs:        newDefaultMefsConfig(),
		Ipfs:        newDefaultIpfsConfig(),
		Prices:      map[string]int64{"mefs": 25000, "ipfs": 25000},
		TrafficCost: 25000,
	}
}

func newDefaultContractsConfig() map[int]ContractConfig {
	cfg := map[int]ContractConfig{
		985: {
			Chain: "dev",
		},
	}
	return cfg
}

func newDefaultContractConfig() ContractConfig {
	cfg := ContractConfig{
		Chain:      "dev",
		SellerAddr: "0xdFF2A42524df7574361A90aac9141DE3f4D8eA02",
	}
	return cfg
}

func newDefaultSecurityKeyConfig() string {
	return hex.EncodeToString(crypto.Keccak256([]byte(time.Now().String())))
}

func newDefaultDomainConfig() string {
	return "memo.io"
}

func NewDefaultConfig() *Config {
	return &Config{
		SwagHost:    "localhost:8090",
		Storage:     newDefaultStorageConfig(),
		Contracts:   newDefaultContractsConfig(),
		Contract:    newDefaultContractConfig(),
		SecurityKey: newDefaultSecurityKeyConfig(),
		Domain:      newDefaultDomainConfig(),
		EthDriveUrl: "https://ethdrive.net",
		// LensAPIUrl:  "https://api.lens.dev",
	}
}

func (cfg *Config) GetStore() interface{} {
	return cfg.Storage
}

func (cfg *Config) GetContract() interface{} {
	return cfg.Contract
}

func (cfg *Config) WriteFile(file string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	configString, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(f, string(configString))
	return err
}

func ReadFile() (*Config, error) {
	file := CONFIGPATH
	cfg := NewDefaultConfig()
	if file != "" {
		rawConfig, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if len(rawConfig) == 0 {
			return cfg, nil
		}

		err = json.Unmarshal(rawConfig, &cfg)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func init() {
	Cfg = NewDefaultConfig()
	data, err := json.MarshalIndent(Cfg, "", "	")
	if err != nil {
		fmt.Println("Config load Failed, ", err)
		return
	}
	// Check if the file exists
	_, err = os.Stat(CONFIGPATH)
	if !os.IsNotExist(err) {
		Cfg, err = ReadFile()
		if err != nil {
			fmt.Println("Load config error, ", err)
		}
		return
	}

	err = ioutil.WriteFile(CONFIGPATH, data, 0644)
	if err != nil {
		fmt.Println("Config load Failed, ", err)
		return
	}
	fmt.Println("config init success!")
}
