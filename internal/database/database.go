package database

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

var logger = logs.Logger("database")

type FileInfo struct {
	ID         int `gorm:"primarykey"`
	ChainID    int `gorm:"column:chainid"`
	Address    string
	SType      storage.StorageType `gorm:"column:stype"`
	Mid        string
	Name       string
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
