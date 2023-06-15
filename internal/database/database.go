package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

var logger = logs.Logger("database")

type FileInfo struct {
	Id         int
	ChainID    int
	Address    string
	SType      storage.StorageType
	Name       string
	Mid        string
	Size       int64
	ModTime    time.Time
	Public     bool
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
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
	}

	err = db.Ping()
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
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
		chainid INTEGER,
		address TEXT,
		stype INTEGER,
		name TEXT,
		mid TEXT,
		size INTEGER,
		modtime DATETIME,
		public BOOLEAN,
		userdefine TEXT,
		UNIQUE (chainid, address, stype, mid) ON CONFLICT IGNORE
	);
	`

	_, err = db.Exec(sqlMessage)
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
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
        INSERT INTO fileinfo (chainid, address, stype, name, mid, size, modtime, public, userdefine)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err = db.Exec(sqlStmt, fi.ChainID, fi.Address, fi.SType, fi.Name, fi.Mid, fi.Size, fi.ModTime, fi.Public, fi.UserDefine)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Get(chain int, mid string, st storage.StorageType) (map[string]FileInfo, error) {
	result := make(map[string]FileInfo)
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return result, err
	}
	defer db.Close()

	sqlStmt := `
	SELECT * FROM fileinfo
	WHERE chainid=? AND mid=? AND stype=?
	`

	var fi FileInfo
	var rows *sql.Rows

	rows, err = db.Query(sqlStmt, chain, mid, st)
	if err != nil {
		if err == sql.ErrNoRows {
			lerr := logs.DataBaseError{Message: fmt.Sprintf("no such record:chainid=%d mid=%s, stype=%d", chain, mid, st)}
			logger.Errorf(lerr.Message)
			return result, lerr
		}
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr.Message)
		return result, lerr
	}

	for rows.Next() {
		err := rows.Scan(&fi.Id, &fi.ChainID, &fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.ModTime, &fi.Public, &fi.UserDefine)
		if err != nil {
			lerr := logs.DataBaseError{Message: err.Error()}
			logger.Error(lerr.Message)
			return result, err
		}
		result[fi.Address] = fi
	}

	return result, nil
}

func List(chain int, address string, st storage.StorageType) ([]FileInfo, error) {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer db.Close()

	sqlStmt := `
        SELECT chainid, address, stype, name, mid, size, modtime, public, userdefine
        FROM fileinfo
        WHERE chainid=? AND address=? AND stype=?
    `
	rows, err := db.Query(sqlStmt, chain, address, st)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileList []FileInfo
	for rows.Next() {
		var fi FileInfo
		err := rows.Scan(&fi.ChainID, &fi.Address, &fi.SType, &fi.Name, &fi.Mid, &fi.Size, &fi.ModTime, &fi.Public, &fi.UserDefine)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, fi)
	}

	return fileList, nil
}

func Delete(chain int, address, mid string, stype storage.StorageType) error {
	db, err := OpenDataBase()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer db.Close()

	sqlStmt := `
	DELETE FROM fileinfo
	WHERE chainid=? AND address=? AND mid=? AND stype=?
`
	res, err := db.Exec(sqlStmt, chain, address, mid, stype)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		err = logs.DataBaseError{Message: "delete object failed"}
		return err
	}

	return nil
}
