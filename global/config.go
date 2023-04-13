package global

import "github.com/ethereum/go-ethereum/common"

var (
	Endpoint = "https://chain.metamemo.one:8501"

	ContractAddr   = common.HexToAddress("0x2A0B376CC39eB2019e43207d00ee2c34878ca36D")
	ContractAddrV2 = common.HexToAddress("0x6E4D5070af7A545fBa9506D28AD42Cc6D9d38A7a")

	GatewayAddr      = common.HexToAddress("0x31e7829Ea2054fDF4BCB921eDD3a98a825242267")
	GatewaySecretKey = "8a87053d296a0f0b4600173773c8081b12917cef7419b2675943b0aa99429b62"

	PayTopic     = "0xc0e3b3bf3b856068b6537f07e399954cb5abc4fade906ee21432a8ded3c36ec8"
	StorageTopic = "0x63fbca6586cb6d6fcf9fe8ab7daf3ffaf7fdad8f5d2ab29109fe71599b10d800"
	BuyTopic     = "0x9393f0a0a85953b7957a62d1ced4afd964332dad208249e1db83ce254babfccc"
)
