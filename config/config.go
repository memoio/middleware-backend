package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Storage StorageConfig `json:"storage"`
}

type StorageConfig struct {
	Mefs MefsConfig `json:"mefs"`
	Ipfs IpfsConfig `json:"ipfs"`
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
		Mefs: newDefaultMefsConfig(),
		Ipfs: newDefaultIpfsConfig(),
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		Storage: newDefaultStorageConfig(),
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

func ReadFile(file string) (*Config, error) {
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
