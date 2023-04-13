package database

import (
	"fmt"
)

type AddressInfo struct {
	ID      int
	Address string
}

func (a AddressInfo) Insert() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	message := `
	INSERT INTO address (address)
	VALUES (?)
	`

	return execSql(db, message, a.Address)
}

func (a AddressInfo) isExist() bool {
	db, err := OpenDataBase()
	if err != nil {
		return false
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM address WHERE address = ?"
	row := db.QueryRow(query, a.Address)
	var count int
	if err := row.Scan(&count); err != nil {
		return false
	}
	return count == 1
}
