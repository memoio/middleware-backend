package market

import (
	"math/big"

	"github.com/go-ego/riot/types"
)

func InitDriveMarket(opts types.EngineOpts) {
	// 初始化Searcher
	InitSearchEngine(opts)

	err := InitNFTTable()
	if err != nil {
		panic(err.Error())
	}

	err = InitBlockNumberTable()
	if err != nil {
		panic(err.Error())
	}

	blockNumber = big.NewInt(GetLastBlockNumber())

	var page = 1
	const pageSize = 100
	for {
		nfts, err := ListNFT(page, pageSize, "", true)
		if err != nil {
			panic(err.Error())
		}

		// add doc to searcher
		for _, nft := range nfts {
			AddNFT(nft.TokenID, nft.Description, nft.Name)
		}

		if len(nfts) != pageSize {
			break
		}
		page++
	}
	FlushSearcher()
}
