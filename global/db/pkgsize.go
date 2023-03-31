package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PkgInfo struct {
	ID        int
	Address   string
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
	flag, err := p.isExist(db)
	if err != nil {
		return err
	}

	if flag {
		log.Println("address and hashid is exist")
		return ErrAlreadyExist
	}

	message := `
		INSERT INTO PkgInfo (address_id, hash, size, time, is_updated)
		SELECT id, ?, ?, ?, ? FROM address WHERE address=?
	`
	ai, err := QueryPkgSize(p.Address)
	if err != nil {
		return err
	}

	err = ai.UpdateSize(p.Size)
	if err != nil {
		return err
	}

	return execSql(db, message, p.Hashid, p.Size, p.UTime, p.IsUpdated, p.Address)
}

func (p PkgInfo) isExist(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow(`
	SELECT COUNT(*)
	FROM PkgInfo
	JOIN Address ON PkgInfo.address_id = Address.id
	WHERE Address.address = ? AND PkgInfo.hash = ?
	`, p.Address, p.Hashid).Scan(&count)

	if err != nil {
		log.Println(err)
		return false, err

	}
	return count > 0, nil
}
