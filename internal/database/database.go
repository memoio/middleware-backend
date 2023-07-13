package database

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"gorm.io/gorm"
)

var logger = logs.Logger("database")

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
	Shared     bool
	UserDefine string `gorm:"column:userdefine"`
}

func (FileInfo) TableName() string {
	return "fileinfo"
}

// func GetFileByUniqueIndex(address string, chainid int, mid string, stype storage.StorageType) *FileInfo {
// 	var file FileInfo
// 	if err := DataBase.Where("address = ? and chain_id = ? and mid = ? and s_type = ?", address, chainid, mid, stype).Find(&file).Error; err != nil {
// 		return nil
// 	}
// 	return &file
// }

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

	if len(fileInfos) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, file := range fileInfos {
		result[file.Address] = file
	}

	return result, err
}

func GetPublic(chain int, mid string, st storage.StorageType) (map[string]FileInfo, error) {
	var fileInfos []FileInfo
	var result = make(map[string]FileInfo)
	err := DataBase.Where("chainid = ? and mid = ? and stype = ? and public = true", chain, mid, st).Find(&fileInfos).Error
	if err != nil {
		return nil, err
	}

	if len(fileInfos) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, file := range fileInfos {
		result[file.Address] = file
	}

	return result, err
}

func GetById(id int) (FileInfo, error) {
	var result = FileInfo{}
	err := DataBase.Where("id = ?", id).Find(&result).Error
	if err != nil {
		return result, err
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
