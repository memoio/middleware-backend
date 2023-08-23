package database

import (
	"context"
	"math/rand"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func (d *DataStore) AddUser(ctx context.Context, ui api.USerInfo) error {
	if err := DataBase.Create(&ui).Error; err != nil {
		return err
	}
	return nil
}

func (d *DataStore) SelectUser(ctx context.Context, area string) (api.USerInfo, error) {
	var result api.USerInfo
	var userInfos []api.USerInfo
	err := DataBase.Where("area = ?", area).Find(&userInfos).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return result, lerr
	}

	index := rand.Intn(len(userInfos))

	return userInfos[index], nil
}

func (d *DataStore) DeleteUser(ctx context.Context, id int) error {
	return DataBase.Delete(&api.USerInfo{}, "id = ?", id).Error
}

func (d *DataStore) ListUsers(ctx context.Context) ([]api.USerInfo, error) {
	var userInfos []api.USerInfo
	err := DataBase.Find(&userInfos).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return userInfos, lerr
	}

	return userInfos, nil
}

func (d *DataStore) GetUser(ctx context.Context, id int) (api.USerInfo, error) {
	var result api.USerInfo
	err := DataBase.Where("id = ?", id).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}
