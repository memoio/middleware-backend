package database

import (
	"context"

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
	var fileInfos []api.FileInfo
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
	var result api.FileInfo
	err := DataBase.Where("address = ? and mid = ? and stype = ?", address, mid, st).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}

func (d *DataStore) GetObjectInfoById(ctx context.Context, id int) (interface{}, error) {
	var result api.FileInfo
	err := DataBase.Where("id = ?", id).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}

func (d *DataStore) DeleteObject(ctx context.Context, id int) error {
	return DataBase.Delete(&api.FileInfo{}, "id = ?", id).Error
}

func (d *DataStore) PutObject(ctx context.Context, fi api.FileInfo) error {
	if err := DataBase.Create(&fi).Error; err != nil {
		return err
	}
	return nil
}

// func GetFileByUniqueIndex(address string, chainid int, mid string, stype storage.StorageType) *FileInfo {
// 	var file FileInfo
// 	if err := DataBase.Where("address = ? and chain_id = ? and mid = ? and s_type = ?", address, chainid, mid, stype).Find(&file).Error; err != nil {
// 		return nil
// 	}
// 	return &file
// }

func Put(fi api.FileInfo) (bool, error) {
	if err := DataBase.Create(&fi).Error; err != nil {
		return false, err
	}
	return true, nil
}

func Get(chain int, mid string, st storage.StorageType) (map[string]api.FileInfo, error) {
	var fileInfos []api.FileInfo
	var result = make(map[string]api.FileInfo)
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

func List(chain int, address string, st storage.StorageType) ([]api.FileInfo, error) {
	var fileInfos []api.FileInfo
	err := DataBase.Where("chainid = ? and address = ? and stype = ?", chain, address, st).Find(&fileInfos).Error
	if err != nil {
		return nil, err
	}

	return fileInfos, nil
}

func Delete(chain int, address, mid string, stype storage.StorageType) error {
	return DataBase.Delete(&api.FileInfo{}, "chainid = ? and address = ? and stype = ?", chain, address, stype).Error
}
