package share

import (
	"errors"
	"time"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
	"github.com/segmentio/ksuid"
)

type ShareObjectInfo struct {
	ShareID     string          `json:"shareid" gorm:"primaryKey"`
	Address     string          `json:"address" gorm:"uniqueIndex:uni"`
	ChainID     int             `json:"chainid" gorm:"uniqueIndex:uni"`
	MID         string          `json:"mid" gorm:"uniqueIndex:uni"`
	SType       api.StorageType `json:"type" gorm:"uniqueIndex:uni"`
	FileName    string          `json:"filename" gorm:"uniqueIndex:uni"`
	Key         string          `json:"key"`
	ExpiredTime int64           `json:"expire"`
}

func InitShareTable() error {
	return database.DataBase.AutoMigrate(&ShareObjectInfo{})
}

func (s *ShareObjectInfo) CreateShare() (string, error) {
	uuid, err := ksuid.NewRandom()
	if err != nil {
		return "", err
	}
	s.ShareID = uuid.String()

	err = database.DataBase.Create(s).Error
	if err != nil {
		return "", logs.DataBaseError{Message: err.Error()}
	}

	return s.ShareID, nil
}

func GetShareByUniqueIndex(address string, chainid int, mid string, stype api.StorageType, name string) *ShareObjectInfo {
	var share ShareObjectInfo
	if name != "" {
		if err := database.DataBase.Where("address = ? and chain_id = ? and m_id = ? and s_type = ? and file_name = ?", address, chainid, mid, stype, name).Find(&share).Error; err != nil {
			return nil
		}
	} else {
		if err := database.DataBase.Where("address = ? and chain_id = ? and m_id = ? and s_type = ?", address, chainid, mid, stype).Find(&share).Error; err != nil {
			return nil
		}
	}
	return &share
}

func GetShareByID(shareID string) *ShareObjectInfo {
	var share ShareObjectInfo
	if err := database.DataBase.Where("share_id = ?", shareID).First(&share).Error; err != nil {
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

	_, err := GetFileInfo(s.Address, s.ChainID, s.MID, s.SType, s.FileName)
	if err != nil {
		// 文件已删除，考虑删除失效的分享
		s.DeleteShare()
		return false
	}

	return true
}

func (s *ShareObjectInfo) Source() (database.FileInfo, error) {
	return GetFileInfo(s.Address, s.ChainID, s.MID, s.SType, s.FileName)
}

func (s *ShareObjectInfo) CanDownload(address string, chainID int) error {
	// now anyone can download
	return nil
}

func (s *ShareObjectInfo) UpdateShare(attr string, value string) error {
	var err error
	switch attr {
	case "password":
		err = database.DataBase.Model(s).Update(attr, value).Error
		if err != nil {
			err = logs.DataBaseError{Message: err.Error()}
		}
	default:
		err = errors.New("unsupport attribute")
	}
	return err
}

func (s *ShareObjectInfo) DeleteShare() error {
	err := database.DataBase.Delete(s).Error
	if err != nil {
		return logs.DataBaseError{Message: err.Error()}
	}

	fileInfo, err := GetFileInfo(s.Address, s.ChainID, s.MID, s.SType, s.FileName)
	if err == nil {
		if err = database.DataBase.Model(&fileInfo).Update("shared", false).Error; err != nil {
			return logs.DataBaseError{Message: err.Error()}
		}
	}

	return nil
}

func GetFileInfo(address string, chainID int, mid string, stype api.StorageType, name string) (database.FileInfo, error) {
	var fileinfos []database.FileInfo
	var err error
	if name != "" {
		err = database.DataBase.Where("chainid = ? and mid = ? and stype = ? and name = ?", chainID, mid, stype, name).Find(&fileinfos).Error
	} else {
		err = database.DataBase.Where("chainid = ? and mid = ? and stype = ?", chainID, mid, stype).Find(&fileinfos).Error
	}
	if err != nil {
		return database.FileInfo{}, logs.DataBaseError{Message: err.Error()}
	}

	for _, file := range fileinfos {
		if file.Public {
			return file, nil
		}
		if file.Address == address {
			return file, nil
		}
	}

	return database.FileInfo{}, logs.NoPermission{Message: "can't access the file"}
}

func ListShares(address string, chainID int) ([]ShareObjectInfo, error) {
	var shares []ShareObjectInfo
	err := database.DataBase.Where("address = ? and chain_id = ?", address, chainID).Find(&shares).Error
	if err != nil {
		return nil, logs.DataBaseError{Message: err.Error()}
	}

	return shares, nil
}
