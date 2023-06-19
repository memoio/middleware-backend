package share

import (
	"time"

	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/storage"
	"golang.org/x/xerrors"
)

type CreateShareRequest struct {
	MID             string              `josn:"mid"`
	SType           storage.StorageType `json:"type"`
	Password        string              `json:"password"`
	ExpiredTime     int64               `josn:"expire"`
	RemainDownloads int                 `json:"downloads"`
}

func CreateShare(address string, chainID int, request CreateShareRequest) (string, error) {
	// 查看文件是否存在，且属于该用户
	fileInfo, err := GetFileInfo(address, chainID, request.MID, request.SType)
	if err != nil {
		return "", err
	}

	newShare := ShareObjectInfo{
		UserID:          Identity{address, chainID},
		MID:             request.MID,
		SType:           request.SType,
		FileName:        fileInfo.Name,
		Password:        request.Password,
		ExpiredTime:     -1,
		RemainDownloads: -1,
	}

	if request.RemainDownloads > 0 {
		newShare.RemainDownloads = request.RemainDownloads
	}

	if request.ExpiredTime > 0 {
		newShare.ExpiredTime = time.Now().Unix() + request.ExpiredTime
	}

	id, err := newShare.CreateShare()
	if err != nil {
		return "", xerrors.Errorf("Failed to create share link")
	}

	baseUrl := "https://ethdrive.net"
	config, err := config.ReadFile()
	if err == nil {
		baseUrl = config.EthDriveUrl
	}
	return baseUrl + "/s/" + id, nil
}

type UpdateShareRequest struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func UpdateShare(share *ShareObjectInfo, request UpdateShareRequest) error {
	return share.UpdateShare(request.Attribute, request.Value)
}

func DeleteShare(address string, chainID int, share *ShareObjectInfo) error {
	if share.UserID.Address != address || share.UserID.ChainID != chainID {
		return xerrors.Errorf("It's not your share link, can't delete")
	}

	return share.DeleteShare()
}

func GetShare(address string, chainID int, share *ShareObjectInfo, password string) (*ShareObjectInfo, error) {
	var unlocked = true
	if share.Password != "" {
		_, unlocked = MemoCache.Load("unlock" + address + share.ShareID)
		// 当前用户未输入相应的密码解锁
		if !unlocked {
			if share.Password == password {
				unlocked = true
				MemoCache.Store("unlock"+address+share.ShareID, struct{}{})
			}
		}
	}

	if !unlocked {
		return nil, xerrors.Errorf("Please enter the correct password")
	}

	return share, nil
}

func SaveShare(address string, chainID int, share *ShareObjectInfo) error {
	info, err := GetFileInfo(share.UserID.Address, share.UserID.ChainID, share.MID, share.SType)
	if err != nil {
		return err
	}

	info.Address = address
	info.ChainID = chainID

	_, err = database.Put(info)
	return err
}
