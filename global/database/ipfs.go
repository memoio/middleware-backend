package database

import (
	"fmt"
)

type ObjectInfo struct {
	Address string
	Name    string
	Size    int64
	Cid     string
}

func (o ObjectInfo) Insert() error {
	db, err := OpenDataBase()
	if err != nil {
		return DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	message := "INSERT INTO ipfsobjects (address_id, name, size, hash) VALUES ((SELECT id FROM address WHERE address = ?), '', ?, ?)"

	return execSql(db, message, o.Address, o.Size, o.Cid)
}

func ListObjects(address string) ([]ObjectInfo, error) {
	db, err := OpenDataBase()
	if err != nil {
		return nil, DataBase{fmt.Sprintf("connect database error %s", err)}
	}
	defer db.Close()

	message := "SELECT name, size, hash FROM ipfsobjects WHERE address_id = (SELECT id FROM address WHERE address = ?)"

	rows, err := db.Query(message, address)
	if err != nil {
		return nil, DataBase{fmt.Sprintf("list ipfs objects error %s", err)}
	}
	defer rows.Close()

	var result []ObjectInfo
	for rows.Next() {
		var object ObjectInfo
		err = rows.Scan(&object.Name, &object.Size, &object.Cid)

		if err != nil {
			return nil, DataBase{fmt.Sprintf("list ipfs objects error %s", err)}
		}

		result = append(result, object)
	}

	if err = rows.Err(); err != nil {
		return nil, DataBase{fmt.Sprintf("list ipfs objects error %s", err)}
	}

	return result, nil
}
