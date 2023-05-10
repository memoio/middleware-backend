package database

import (
	"database/sql"
	"fmt"
	"math/big"

	_ "github.com/mattn/go-sqlite3"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

var logger = logs.Logger("database")

type FileInfoList struct {
	Size *big.Int
	fi   chan FileInfo
}

type FileInfo struct {
	Id      int
	Address string
	SType   storage.StorageType
	Name    string
	Mid     string
	Size    int64
	OnChain bool
}

func init() {
	err := createTable()
	if err != nil {
		logger.Error("init db failed")
	}
}

func OpenDataBase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./backend.db")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return db, nil
}

func createTable() error {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer db.Close()

	sqlMessage := `
	CREATE TABLE IF NOT EXISTS fileinfo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT,
		stype INTEGER,
		name TEXT,
		mid TEXT,
		size INTEGER,
		onchain BOOLEAN,
		UNIQUE (address, stype, mid) ON CONFLICT IGNORE
	);
	`

	_, err = db.Exec(sqlMessage)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func Put(fi FileInfo) (bool, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return false, err
	}
	defer db.Close()

	sqlStmt := `
        INSERT INTO fileinfo (address, stype, name, mid, size, onchain)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err = db.Exec(sqlStmt, fi.Address, fi.SType, fi.Name, fi.Mid, fi.Size, fi.OnChain)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Get(address, mid string, st storage.StorageType) (FileInfo, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return FileInfo{}, err
	}
	defer db.Close()

	sqlStmt := `
	SELECT * FROM fileinfo
	WHERE address=? AND mid=? AND stype=?
`
	var fi FileInfo
	err = db.QueryRow(sqlStmt, address, mid, st).Scan(&fi.Id, &fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.OnChain)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("no such object: address=%s, mid=%s, stype=%d", address, mid, st)
			return FileInfo{}, logs.DataBaseError{Message: fmt.Sprintf("no such record: address=%s, mid=%s, stype=%d", address, mid, st)}
		}
		return FileInfo{}, err
	}
	return fi, nil
}

func List(address string, st storage.StorageType) ([]FileInfo, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer db.Close()

	sqlStmt := `
        SELECT address, stype, name, mid, size, onchain
        FROM fileinfo
        WHERE address=? AND stype=?
    `
	rows, err := db.Query(sqlStmt, address, st)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileList []FileInfo
	for rows.Next() {
		var fi FileInfo
		err := rows.Scan(&fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.OnChain)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, fi)
	}

	return fileList, nil
}

func Delete(address, mid string, stype storage.StorageType) (bool, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return false, err
	}
	defer db.Close()

	sqlStmt := `
	DELETE FROM fileinfo
	WHERE address=? AND mid=? AND stype=?
`
	res, err := db.Exec(sqlStmt, address, mid, stype)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func GetNotOnChain() ([]FileInfo, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer db.Close()

	sqlStmt := `
	SELECT * FROM fileinfo
	WHERE onchain=0
`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []FileInfo
	for rows.Next() {
		var fi FileInfo
		err = rows.Scan(&fi.Id, &fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.OnChain)
		if err != nil {
			return nil, err
		}
		result = append(result, fi)
	}
	return result, nil
}
