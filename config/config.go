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
	Storage     StorageConfig          `json:"storage"`
	Contracts   map[int]ContractConfig `json:"contracts"`
	Contract    ContractConfig         `json:"contract"`
	SecurityKey string                 `json:"securityKey"`
	Domain      string                 `json:"domain"`
	LensAPIUrl  string                 `json:"lensAPIUrl"`
	EthDriveUrl string                 `json:"ethDriveUrl"`
}

type StorageConfig struct {
	Mefs        MefsConfig       `json:"mefs"`
	Ipfs        IpfsConfig       `json:"ipfs"`
	Prices      map[string]int64 `json:"prices"`
	TrafficCost int64            `json:"traffic_cost"`
}

type MefsConfig struct {
	Api          string `json:"api"`
	Token        string `json:"token"`
	ContractAddr string `json:"contractaddr"`
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
		Api: "/ip4/127.0.0.1/tcp/5001",
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
		Chain: "dev",
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
		Storage:     newDefaultStorageConfig(),
		Contracts:   newDefaultContractsConfig(),
		Contract:    newDefaultContractConfig(),
		SecurityKey: newDefaultSecurityKeyConfig(),
		Domain:      newDefaultDomainConfig(),
		LensAPIUrl:  "https://api.lens.dev",
		EthDriveUrl: "https://ethdrive.net",
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
