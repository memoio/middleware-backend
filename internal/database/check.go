package database

import (
	"math/big"
	"sync"
	"time"

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
		logger.Info("fileinfo", fi)
		return true, nil
	default:
		logger.Error("fail to write fileinfo to database")
		return false, logs.DataBaseError{Message: "fail to write fileinfo to database"}
	}
}

func (w *WriteCheck) Read() error {
	for _, flist := range w.pool {
		for {
			logger.Info("read fi ", flist.fi)
			select {
			case fi, ok := <-flist.fi:
				if !ok {
					break
				}

				res, err := Put(fi)
				if err != nil {
					logger.Errorf("failed to write file info to database: %v", err)
					return err
				}
				if !res {
					return logs.DataBaseError{Message: "write to database error"}
				}
				flist.Size.Add(flist.Size, big.NewInt(fi.Size))

			case <-time.After(1 * time.Second):
				break
			}
		}
	}
	return nil
}
