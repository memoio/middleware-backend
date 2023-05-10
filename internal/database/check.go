package database

import (
	"math/big"
	"sync"

	"github.com/memoio/backend/internal/logs"
)

type WriteCheck struct {
	lw   sync.Mutex
	pool map[string]*FileInfoList
}

func NewWriteCheck() *WriteCheck {
	wc := WriteCheck{
		pool: map[string]*FileInfoList{},
	}
	return &wc
}

func (w *WriteCheck) Write(fi FileInfo) (bool, error) {
	address := fi.Address
	w.lw.Lock()
	defer w.lw.Unlock()

	var ch chan FileInfo
	p, ok := w.pool[fi.Address]
	if !ok {
		flist := &FileInfoList{
			Size: new(big.Int),
			fi:   make(chan FileInfo, 1000),
		}
		w.pool[address] = flist
		ch = flist.fi
	} else {
		ch = p.fi
	}

	select {
	case ch <- fi:

		return true, nil
	default:
		logger.Error("fail to write fileinfo to database")
		return false, logs.DataBaseError{Message: "fail to write fileinfo to database"}
	}
}

func (w *WriteCheck) Read() error {
	w.lw.Lock()
	defer w.lw.Unlock()
	for _, flist := range w.pool {
		for fi := range flist.fi {
			_, err := Put(fi)
			if err != nil {
				logger.Errorf("failed to write file info to database: %v", err)
			}
			flist.Size.Add(flist.Size, big.NewInt(fi.Size))
		}
	}
	return nil
}
