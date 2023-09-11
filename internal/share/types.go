package share

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"github.com/segmentio/ksuid"
)

type ShareObjectInfo struct {
	ShareID     string              `json:"shareid" gorm:"primaryKey"`
	Address     string              `json:"address" gorm:"uniqueIndex:uni"`
	ChainID     int                 `json:"chainid" gorm:"uniqueIndex:uni"`
	MID         string              `json:"mid" gorm:"uniqueIndex:uni"`
	SType       storage.StorageType `json:"type" gorm:"uniqueIndex:uni"`
	FileName    string              `json:"filename"`
	ExpiredTime int64               `json:"expire"`
}

var MemoCache = new(sync.Map)

func InitShareTable() error {
	return database.GlobalDataBase.AutoMigrate(&ShareObjectInfo{})
}

func (s *ShareObjectInfo) CreateShare() (string, error) {
	uuid, err := ksuid.NewRandom()
	if err != nil {
		return "", err
	}
	s.ShareID = uuid.String()

	err = database.GlobalDataBase.Create(s).Error
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return "", logs.DataBaseError{Message: "Alread created the share"}
		}
		return "", logs.DataBaseError{Message: err.Error()}
	}

	return s.ShareID, nil
}

func GetShareByUniqueIndex(address string, chainid int, mid string, stype storage.StorageType) *ShareObjectInfo {
	var share ShareObjectInfo
	if err := database.GlobalDataBase.Where("address = ? and chain_id = ? and m_id = ? and s_type = ?", address, chainid, mid, stype).Find(&share).Error; err != nil {
		return nil
	}
	return &share
}

func GetShareByID(shareID string) *ShareObjectInfo {
	var share ShareObjectInfo
	if err := database.GlobalDataBase.Where("share_id = ?", shareID).First(&share).Error; err != nil {
		return nil
	}
	return &share
}

func (s *ShareObjectInfo) IsAvailable() bool {
	if s.ExpiredTime > 0 && time.Now().Unix() > s.ExpiredTime {
		// 考虑删除失效的分享
		s.DeleteShare()
		return false
	}

	_, err := GetFileInfo(s.Address, s.ChainID, s.MID, s.SType)
	if err != nil {
		// 文件已删除，考虑删除失效的分享
		s.DeleteShare()
		return false
	}

	return true
}

func (s *ShareObjectInfo) Source() (api.FileInfo, error) {
	return GetFileInfo(s.Address, s.ChainID, s.MID, s.SType)
}

func (s *ShareObjectInfo) CanDownload(address string, chainID int) error {
	// now anyone can download
	return nil
}

func (s *ShareObjectInfo) UpdateShare(attr string, value string) error {
	var err error
	switch attr {
	case "password":
		err = database.GlobalDataBase.Model(s).Update(attr, value).Error
		if err != nil {
			err = logs.DataBaseError{Message: err.Error()}
		}
	default:
		err = errors.New("unsupport attribute")
	}
	return err
}

func (s *ShareObjectInfo) DeleteShare() error {
	err := database.GlobalDataBase.Delete(s).Error
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}
	return nil
}

func GetFileInfo(address string, chainID int, mid string, stype storage.StorageType) (api.FileInfo, error) {
	fileInfos, err := database.Get(chainID, mid, stype)
	if err != nil {
		return api.FileInfo{}, logs.DataBaseError{Message: err.Error()}
	}
	for key, file := range fileInfos {
		if file.Public {
			return file, nil
		}
		if key == address {
			return file, nil
		}
	}

	return api.FileInfo{}, logs.NoPermission{Message: "can't access the file"}
}

func ListShares(address string, chainID int) ([]ShareObjectInfo, error) {
	var shares []ShareObjectInfo
	err := database.GlobalDataBase.Where("address = ? and chain_id = ?", address, chainID).Find(&shares).Error
	if err != nil {
		return nil, logs.DataBaseError{Message: err.Error()}
	}

	return shares, nil
}
