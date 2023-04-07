package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type AddressInfo struct {
	ID         int
	Address    string
	Available  int64
	Free       int64
	Used       int64
	Files      int
	UpdateTime time.Time
}

func isAddressExist(db *sql.DB, address string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM address WHERE address = ?", address).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a AddressInfo) Insert() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	exist, err := isAddressExist(db, a.Address)
	if err != nil {
		log.Println(err)
		return err
	}

	if exist {
		log.Println("address already exist")
		return ErrAlreadyExist
	}

	message := `
	INSERT INTO address (address, pay_space, free_space, used_space, file_count, update_time)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	return execSql(db, message, a.Address, a.Available, a.Free, a.Used, a.Files, a.UpdateTime)
}

func (a AddressInfo) Update() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	exist, err := isAddressExist(db, a.Address)
	if err != nil {
		log.Println(err)
		return DataBase{err.Error()}
	}

	if !exist {
		log.Println("address not exist")
		return ErrNotExist
	}

	message := `
	UPDATE address
	SET pay_space = ?, free_space = ?, used_space = ?, file_count = ?, update_time = ?
	WHERE address = ?
	`

	return execSql(db, message, a.Available, a.Free, a.Used, a.Files, a.UpdateTime, a.Address)
}

func (a AddressInfo) UpdateSize(fsize int64) error {
	a.Used += fsize
	a.Files += 1
	a.UpdateTime = time.Now()

	if a.Used > a.Available+a.Free {
		return DataBase{"storage not enough"}
	}

	return a.Update()
}

func QueryPkgSize(address string) (AddressInfo, error) {
	db, err := OpenDataBase()
	if err != nil {
		return AddressInfo{}, DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, address, pay_space, free_space, used_space, file_count, update_time FROM address WHERE address = ?", address)

	var result AddressInfo
	err = row.Scan(&result.ID, &result.Address, &result.Available, &result.Free, &result.Used, &result.Files, &result.UpdateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return AddressInfo{}, ErrNotExist
		}
		return AddressInfo{}, err
	}

	return result, nil
}
