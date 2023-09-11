package database

import (
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/storage"
	"gorm.io/gorm"
)

func Put(fi api.FileInfo) (bool, error) {
	if err := GlobalDataBase.Create(&fi).Error; err != nil {
		return false, err
	}
	return true, nil
}

func Get(chain int, mid string, st storage.StorageType) (map[string]api.FileInfo, error) {
	var fileInfos []api.FileInfo
	var result = make(map[string]api.FileInfo)
	err := GlobalDataBase.Where("chainid = ? and mid = ? and stype = ?", chain, mid, st).Find(&fileInfos).Error
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
	err := GlobalDataBase.Where("chainid = ? and address = ? and stype = ?", chain, address, st).Find(&fileInfos).Error
	if err != nil {
		return nil, err
	}

	return fileInfos, nil
}

func Delete(chain int, address, mid string, stype storage.StorageType) error {
	return GlobalDataBase.Delete(&api.FileInfo{}, "chainid = ? and address = ? and stype = ?", chain, address, stype).Error
}
