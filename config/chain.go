package config

type ContractConfig struct {
	Endpoint         string `json:"endpoint"`
	ContractAddr     string `json:"caddr"`
	GatewayAddr      string `json:"gaddr"`
	GatewaySecretKey string `json:"gatewaysk"`
}
