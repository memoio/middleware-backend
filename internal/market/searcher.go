package market

import (
	"strconv"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var Searcher = riot.Engine{}

var DefaultSearcherOpts = types.EngineOpts{
	IDOnly: true,
	Using:  1,
	IndexerOpts: &types.IndexerOpts{
		IndexType: types.DocIdsIndex,
	},
	GseDict:       "./data/dictionary.txt",
	StopTokenFile: "./data/stop_tokens.txt",
}

func InitSearchEngine(opts types.EngineOpts) {
	Searcher.Init(opts)
}

func AddNFT(tokenID int64, description string, name string) {
	// TODO: 使用其他分词器
	var text = description + ";" + name
	Searcher.IndexDoc(strconv.FormatInt(tokenID+1, 10), types.DocData{Content: text})
}

func UpdateNFT(tokenID int64, description string, name string) {
	// TODO: 使用其他分词器
	var text = description + ";" + name
	Searcher.IndexDoc(strconv.FormatInt(tokenID+1, 10), types.DocData{Content: text}, true)
}

func DeleteNFT(tokenID int64) {
	Searcher.RemoveDoc(strconv.FormatInt(tokenID+1, 10), true)
}

func FlushSearcher() {
	Searcher.Flush()
}

func Search(text string, page int, pageSize int) types.ScoredIDs {
	res := Searcher.SearchID(types.SearchReq{
		Text: text,
		RankOpts: &types.RankOpts{
			OutputOffset: (page - 1) * pageSize,
			MaxOutputs:   pageSize,
		},
	})

	return res.Docs
}
