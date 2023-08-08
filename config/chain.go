package config

type ContractConfig struct {
	Chain        string `json:"chain"`
	ContractAddr string `json:"caddr"`
	SellerAddr   string `json:"saddr"`
	// GatewaySecretKey string `json:"gatewaysk"`
}
