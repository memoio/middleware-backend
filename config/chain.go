package config

type ContractConfig struct {
	Chain string `json:"chain"`
	// ContractAddr     string `json:"caddr"`
	// GatewayAddr      string `json:"gaddr"`
	GatewaySecretKey string `json:"gatewaysk"`
}
