package database

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

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
	Id         int
	Address    string
	SType      storage.StorageType
	Name       string
	Mid        string
	Size       int64
	ModTime    time.Time
	UserDefine string
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
		modtime DATETIME,
		userdefine TEXT,
		UNIQUE (address, stype, name, mid) ON CONFLICT IGNORE
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

	logger.Info("put message: ", fi)
	sqlStmt := `
        INSERT INTO fileinfo (address, stype, name, mid, size, modtime, userdefine)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err = db.Exec(sqlStmt, fi.Address, fi.SType, fi.Name, fi.Mid, fi.Size, fi.ModTime, fi.UserDefine)
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
	err = db.QueryRow(sqlStmt, address, mid, st).Scan(&fi.Id, &fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.ModTime, &fi.UserDefine)
	if err != nil {
		if err == sql.ErrNoRows {
			lerr := logs.DataBaseError{Message: fmt.Sprintf("no such record: mid=%s, stype=%d", mid, st)}
			logger.Errorf(lerr.Message)
			return FileInfo{}, lerr
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
        SELECT address, stype, name, mid, size, modtime, userdefine
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
		err := rows.Scan(&fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.ModTime, &fi.UserDefine)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, fi)
	}

	return fileList, nil
}

func Delete(address, name string, stype storage.StorageType) (bool, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return false, err
	}
	defer db.Close()

	sqlStmt := `
	DELETE FROM fileinfo
	WHERE address=? AND name=? AND stype=?
`
	res, err := db.Exec(sqlStmt, address, name, stype)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}
