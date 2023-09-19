package database

import (
	"context"
	"math/rand"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/logs"
)

func (d *DataBase) AddUser(ctx context.Context, ui api.USerInfo) error {
	if err := d.Create(&ui).Error; err != nil {
		return err
	}
	return nil
}

func (d *DataBase) SelectUser(ctx context.Context, area string) (api.USerInfo, error) {
	var result api.USerInfo
	var userInfos []api.USerInfo
	err := d.Where("area = ?", area).Find(&userInfos).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return result, lerr
	}
	var index = 0
	if len(userInfos) != 0 {
		index = rand.Intn(len(userInfos))
	} else {
		userInfos = append(userInfos, api.USerInfo{
			Api:   config.Cfg.Storage.Mefs.Api,
			Token: config.Cfg.Storage.Mefs.Token,
		})
	}

	return userInfos[index], nil
}

func (d *DataBase) DeleteUser(ctx context.Context, id int) error {
	return d.Delete(&api.USerInfo{}, "id = ?", id).Error
}

func (d *DataBase) ListUsers(ctx context.Context, area string) ([]api.USerInfo, error) {
	var userInfos []api.USerInfo
	query := d.Model(&api.USerInfo{})
	if area != "" {
		query = query.Where("area = ?", area)
	}
	err := query.Find(&userInfos).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(lerr)
		return userInfos, lerr
	}

	return userInfos, nil
}

func (d *DataBase) GetUser(ctx context.Context, id int) (api.USerInfo, error) {
	var result api.USerInfo
	err := d.Where("id = ?", id).Find(&result).Error
	if err != nil {
		lerr := logs.DataBaseError{Message: err.Error()}
		logger.Error(err)
		return result, lerr
	}

	return result, err
}
