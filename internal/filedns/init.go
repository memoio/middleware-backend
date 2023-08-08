package filedns

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/go-ego/riot/types"
	dtypes "github.com/memoio/go-did/types"
)

func InitFileDns(opts types.EngineOpts) {
	// 初始化Searcher
	InitSearchEngine(opts)

	// 初始化DIDStore
	var err error
	DIDStore, err = NewBadgerStore("./did", nil)
	if err != nil {
		panic(err.Error())
	}

	// 获取已经处理过的最后一个log的block number
	blockNumber, err = DIDStore.GetLastBlockNumber()
	if err != nil {
		panic(err.Error())
	}

	// 从DIDStore中读取MfileDID,并更新Searcher
	err = DIDStore.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if string(item.Key()) == string(blockNumberKey) {
				continue
			}

			err := item.Value(func(v []byte) error {
				var document dtypes.MfileDIDDocument
				err = json.Unmarshal(v, &document)
				if err != nil {
					return err
				}

				AddShareFile(document.ID.Identifier, document.Keywords)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
}
