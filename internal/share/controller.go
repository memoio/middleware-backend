package share

import (
	"time"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
)

type CreateShareRequest struct {
	MID         string          `josn:"mid"`
	Name        string          `json:"name"`
	Key         string          `json:"key"`
	SType       api.StorageType `json:"type"`
	ExpiredTime int64           `josn:"expire"`
}

func CreateShare(address string, chainID int, request CreateShareRequest) (string, error) {
	// 查看是否支持该存储模式
	_, ok := ApiMap["/"+request.SType.String()]
	if !ok {
		return "", logs.StorageNotSupport{}
	}

	// 查看文件是否存在，且属于该用户
	fileInfo, err := GetFileInfo(address, chainID, request.MID, request.SType, request.Name)
	if err != nil {
		return "", err
	}

	share := GetShareByUniqueIndex(address, chainID, request.MID, request.SType, request.Name)
	if share != nil && share.ShareID != "" && share.IsAvailable() {
		baseUrl := "https://ethdrive.net"
		config, err := config.ReadFile()
		if err == nil {
			baseUrl = config.EthDriveUrl
		}
		return baseUrl + "/s/" + share.ShareID, nil
	}

	newShare := ShareObjectInfo{
		Address:     address,
		ChainID:     chainID,
		MID:         request.MID,
		SType:       request.SType,
		FileName:    fileInfo.Name,
		Key:         request.Key,
		ExpiredTime: -1,
	}

	if request.ExpiredTime > 0 {
		newShare.ExpiredTime = time.Now().Unix() + request.ExpiredTime
	}

	id, err := newShare.CreateShare()
	if err != nil {
		return "", err
	}

	if err = database.DataBase.Model(&fileInfo).Update("shared", true).Error; err != nil {
		return "", logs.DataBaseError{Message: err.Error()}
	}

	baseUrl := "https://ethdrive.net"
	config, err := config.ReadFile()
	if err == nil {
		baseUrl = config.EthDriveUrl
	}
	return baseUrl + "/share/" + id, nil
}

type UpdateShareRequest struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func UpdateShare(share *ShareObjectInfo, request UpdateShareRequest) error {
	return share.UpdateShare(request.Attribute, request.Value)
}

type DeleteShareRequest struct {
	MID   string          `josn:"mid"`
	Name  string          `json:"name"`
	SType api.StorageType `json:"type"`
}

func DeleteShare(address string, chainID int, share *ShareObjectInfo) error {
	if share.Address != address || share.ChainID != chainID {
		return logs.NoPermission{Message: "can't delete"}
	}

	return share.DeleteShare()
}

func GetShare(address string, chainID int, share *ShareObjectInfo) (*ShareObjectInfo, error) {
	// if !CanRead(address, chainID, share) {

	// }

	return share, nil
}

func SaveShare(address string, chainID int, share *ShareObjectInfo) error {
	info, err := GetFileInfo(share.Address, share.ChainID, share.MID, share.SType, share.FileName)
	if err != nil {
		return err
	}

	info.Address = address
	info.ChainID = chainID

	_, err = database.Put(info)
	return err
}
