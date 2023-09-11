package database

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"gorm.io/gorm"
)

var logger = logs.Logger("database")

var _ api.IDataBase = (*DataBase)(nil)

type DataBase struct {
	*gorm.DB
}

func (d *DataBase) ListObjects(ctx context.Context, address string, st api.StorageType) ([]interface{}, error) {
	var fileInfos []api.FileInfo
	var ifileInfos []interface{}
	err := d.Where("address = ? and stype = ?", address, st).Find(&fileInfos).Error
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

func (d *DataBase) GetObjectInfo(ctx context.Context, address, mid string, st api.StorageType) (interface{}, error) {
	var result api.FileInfo
	err := d.Where("address = ? and mid = ? and stype = ?", address, mid, st).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}

func (d *DataBase) GetObjectInfoById(ctx context.Context, id int) (interface{}, error) {
	var result api.FileInfo
	err := d.Where("id = ?", id).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}

func (d *DataBase) DeleteObject(ctx context.Context, id int) error {
	return d.Delete(&api.FileInfo{}, "id = ?", id).Error
}

func (d *DataBase) PutObject(ctx context.Context, fi api.FileInfo) error {
	if err := d.Create(&fi).Error; err != nil {
		return err
	}
	return nil
}
