package database

import (
	"os"
	"time"

	"github.com/memoio/backend/api"
	"github.com/memoio/go-mefs-v2/lib/backend/kv"
	"github.com/memoio/go-mefs-v2/lib/backend/wrap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

// var logger = logs.Logger("share")

func init() {
	logger.Debug("Initializing database")

	db, err := gorm.Open(sqlite.Open("backend.db"), &gorm.Config{})
	if err != nil {
		logger.Panicf("Failed to connect to database: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panicf("Failed to get sql database: %s", err.Error())
	}

	// 设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// 设置超时时间
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	err = sqlDB.Ping()
	if err != nil {
		logger.Panicf("Failed to ping database: %s", err.Error())
	}

	DataBase = db

	DataBase.AutoMigrate(&api.FileInfo{})
}

func NewDataStore(st string) (*DataStore, error) {
	res := &DataStore{}

	opt := kv.DefaultOptions
	bpath := "./datastore/" + st
	err := os.MkdirAll(bpath, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return res, err
	}
	ds, err := kv.NewBadgerStore(bpath, &opt)
	if err != nil {
		logger.Error(err)
		return res, err
	}

	dss := wrap.NewKVStore("upload", ds)

	up := NewCheckPay(dss)

	dss = wrap.NewKVStore("download", ds)
	down := NewCheckPay(dss)
	return &DataStore{DataBase, up, down}, nil
}
