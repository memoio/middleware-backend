package global

import (
	"database/sql"
	"fmt"
	"log"
)

func (d *DataBase) AddAddress(address string) bool {
	_, err := d.execSql("INSERT INTO address(address) values(?)")
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (d *DataBase) QueryAddress() {
	d.querySql("SELECT * FROM address")
}

func (d *DataBase) execSql(message string, args ...any) (sql.Result, error) {
	stmt, err := d.DB.Prepare(message)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(args)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func (d *DataBase) querySql(message string) {
	rows, err := d.DB.Query(message)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var address string
		err = rows.Scan(&id, &address)
		if err != nil {
			return
		}
		fmt.Println(address)
	}
}
