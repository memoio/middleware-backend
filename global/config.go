package global

import "github.com/ethereum/go-ethereum/common"

var (
	ContractAddr     = common.HexToAddress("0x2A0B376CC39eB2019e43207d00ee2c34878ca36D")
	Endpoint         = "https://chain.metamemo.one:8501"
	GatewayAddr      = common.HexToAddress("0x31e7829Ea2054fDF4BCB921eDD3a98a825242267")
	GatewaySecretKey = "8a87053d296a0f0b4600173773c8081b12917cef7419b2675943b0aa99429b62"

	PayTopic     = "0xc0e3b3bf3b856068b6537f07e399954cb5abc4fade906ee21432a8ded3c36ec8"
	StorageTopic = "0x01bed2f7b5b5c577e49071502e2c9985655c9ca5c1ab432156d99d199f6b1912"
)
