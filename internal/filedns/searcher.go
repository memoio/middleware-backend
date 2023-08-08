package filedns

import (
	"strings"

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

func AddShareFile(mid string, keywords []string) {
	// TODO: get filename
	// TODO: 使用其他分词器
	var text = strings.Join(keywords, ";")
	// log.Printf("Add Doc( ID: %s, Content: %s)\n", mid, text)
	Searcher.IndexDoc(mid, types.DocData{Content: text})
}

func UpdateShareFile(mid string, keywords []string) {
	// TODO: get filename
	// TODO: 使用其他分词器
	var text = strings.Join(keywords, ";")
	// log.Printf("Update Doc( ID: %s, Content: %s )\n", mid, text)
	Searcher.IndexDoc(mid, types.DocData{Content: text}, true)
}

func Search(text string, page int, pageSize int) types.ScoredIDs {
	res := Searcher.SearchID(types.SearchReq{
		Text: text,
		RankOpts: &types.RankOpts{
			OutputOffset: page * pageSize,
			MaxOutputs:   pageSize,
		},
	})

	return res.Docs
}
