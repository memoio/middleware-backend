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

type Config struct {
	Storage     StorageConfig  `json:"storage"`
	Contract    ContractConfig `json:"contract"`
	SecurityKey string         `json:"securityKey"`
	Domain      string         `json:"domain"`
	LensAPIUrl  string         `json:"lensAPIUrl"`
}

type StorageConfig struct {
	Mefs        MefsConfig       `json:"mefs"`
	Ipfs        IpfsConfig       `json:"ipfs"`
	Prices      map[string]int64 `json:"prices"`
	TrafficCost int64            `json:"traffic_cost"`
}

type ContractConfig struct {
	Endpoint         string `json:"endpoint"`
	ContractAddr     string `json:"caddr"`
	GatewayAddr      string `json:"gaddr"`
	GatewaySecretKey string `json:"gatewaysk"`
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

func newDefaultContractConfig() ContractConfig {
	return ContractConfig{
		Endpoint:     "https://chain.metamemo.one:8501",
		ContractAddr: "0xA78b166947487d93EA0e87e68132FC4609B00fA1",
		GatewayAddr:  "0x31e7829Ea2054fDF4BCB921eDD3a98a825242267",
	}
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
		Contract:    newDefaultContractConfig(),
		SecurityKey: newDefaultSecurityKeyConfig(),
		Domain:      newDefaultDomainConfig(),
		LensAPIUrl:  "https://api.lens.dev",
	}
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
	cfg := NewDefaultConfig()
	data, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		fmt.Println("Config load Failed, ", err)
		return
	}
	// Check if the file exists
	_, err = os.Stat(CONFIGPATH)
	if !os.IsNotExist(err) {
		return
	}

	err = ioutil.WriteFile(CONFIGPATH, data, 0644)
	if err != nil {
		fmt.Println("Config load Failed, ", err)
		return
	}
	fmt.Println("config init success!")
}
