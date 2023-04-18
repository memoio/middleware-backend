package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Storage struct {
	ID         int
	Address    string
	SType      uint8
	Buysize    int64
	Free       int64
	Used       int64
	Files      int
	UpdateTime time.Time
}

func (s Storage) Insert() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	a := AddressInfo{Address: s.Address}
	flag := a.isExist()
	if !flag {
		err = a.Insert()
		if err != nil {
			return DataBase{err.Error()}
		}
	}
	message := `
	INSERT INTO storage (address_id, stype, pay_space, free_space, used_space, file_count, update_time)
	VALUES (
	  (SELECT id FROM address WHERE address = ?),?,?,?,?,?,?
	);
	`

	return execSql(db, message, s.Address, s.SType, s.Buysize, s.Free, s.Used, s.Files, s.UpdateTime)
}

func (s Storage) Update() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	message := "UPDATE storage SET pay_space=?, free_space=?, used_space=?, file_count=?, update_time=CURRENT_TIMESTAMP WHERE address_id=(SELECT id FROM address WHERE address=?) AND stype=?"

	return execSql(db, message, s.Buysize, s.Free, s.Used, s.Files, s.Address, s.SType)
}

func (a Storage) AddSize(fsize int64) error {
	a.Used += fsize
	a.Files += 1
	a.UpdateTime = time.Now()

	if a.Used > a.Buysize+a.Free {
		return DataBase{"storage not enough"}
	}

	return a.Update()
}

func QueryPkgSize(address string, storage uint8) (Storage, error) {
	db, err := OpenDataBase()
	if err != nil {
		return Storage{}, DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	row := db.QueryRow("SELECT pay_space, free_space, used_space, file_count FROM storage WHERE address_id IN (SELECT id FROM address WHERE address = ?) AND stype = ?", address, storage)

	var result Storage
	err = row.Scan(&result.Buysize, &result.Free, &result.Used, &result.Files)
	if err != nil {
		if err == sql.ErrNoRows {
			return Storage{}, ErrNotExist
		}
		return Storage{}, err
	}

	result.Address = address
	return result, nil
}
