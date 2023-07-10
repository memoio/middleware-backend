package database

import (
	"context"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"gorm.io/gorm"
)

var logger = logs.Logger("database")

type DataStore struct {
	*gorm.DB
}

func (d *DataStore) ListObjects(ctx context.Context, address string, st api.StorageType) ([]interface{}, error) {
	var fileInfos []FileInfo
	var ifileInfos []interface{}
	err := DataBase.Where("address = ? and stype = ?", address, st).Find(&fileInfos).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return nil, lerr
	}

	for _, v := range fileInfos {
		ifileInfos = append(ifileInfos, v)
	}

	return ifileInfos, nil
}

func (d *DataStore) GetObjectInfo(ctx context.Context, address, mid string, st api.StorageType) (interface{}, error) {
	var result FileInfo
	err := DataBase.Where("address = ? and mid = ? and stype = ?", address, mid, st).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}

func (d *DataStore) GetObjectInfoById(ctx context.Context, id int) (interface{}, error) {
	var result FileInfo
	err := DataBase.Where("id = ?", id).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}

type FileInfo struct {
	ID         int                 `gorm:"primarykey"`
	ChainID    int                 `gorm:"uniqueIndex:composite;column:chainid"`
	Address    string              `gorm:"uniqueIndex:composite"`
	SType      storage.StorageType `gorm:"uniqueIndex:composite;column:stype"`
	Mid        string              `gorm:"uniqueIndex:composite"`
	Name       string              `gorm:"uniqueIndex:composite"`
	Size       int64
	ModTime    time.Time `gorm:"column:modtime"`
	Public     bool
	UserDefine string `gorm:"column:userdefine"`
}

func (FileInfo) TableName() string {
	return "fileinfo"
}

func Put(fi FileInfo) (bool, error) {
	if err := DataBase.Create(&fi).Error; err != nil {
		return false, err
	}
	return true, nil
}

func Get(chain int, mid string, st storage.StorageType) (map[string]FileInfo, error) {
	var fileInfos []FileInfo
	var result = make(map[string]FileInfo)
	err := DataBase.Where("chainid = ? and mid = ? and stype = ?", chain, mid, st).Find(&fileInfos).Error
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfos {
		result[file.Address] = file
	}

	return result, err
}

func List(chain int, address string, st storage.StorageType) ([]FileInfo, error) {
	var fileInfos []FileInfo
	err := DataBase.Where("chainid = ? and address = ? and stype = ?", chain, address, st).Find(&fileInfos).Error
	if err != nil {
		return nil, err
	}

	return fileInfos, nil
}

func Delete(chain int, address, mid string, stype storage.StorageType) error {
	return DataBase.Delete(&FileInfo{}, "chainid = ? and address = ? and stype = ?", chain, address, stype).Error
}
