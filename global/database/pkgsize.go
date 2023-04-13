package database

import (
	"fmt"
	"time"
)

type PkgInfo struct {
	ID        int
	Address   string
	SType     uint8
	Hashid    string
	Size      int64
	UTime     time.Time
	IsUpdated bool
}

func (p PkgInfo) Insert() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	a := AddressInfo{Address: p.Address}
	flag := a.isExist()
	if !flag {
		err = a.Insert()
		if err != nil {
			return DataBase{err.Error()}
		}
	}

	message := `
		INSERT INTO PkgInfo (address_id, stype, hash, size, time, is_updated)
		VALUES (
			(SELECT id FROM address WHERE address = ?),?,?,?,?,?
		  );
	`
	ai, err := QueryPkgSize(p.Address, p.SType)
	if err != nil {
		return err
	}

	err = ai.UpdateSize(p.Size)
	if err != nil {
		return err
	}

	return execSql(db, message, p.Address, p.SType, p.Hashid, p.Size, p.UTime, p.IsUpdated)
}
