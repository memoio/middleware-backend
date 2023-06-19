package share

import (
	"sync"
	"time"

	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/storage"
	"github.com/segmentio/ksuid"
	"golang.org/x/xerrors"
)

type Identity struct {
	Address string `json:"address"`
	ChainID int    `json:"chainid"`
}

type ShareObjectInfo struct {
	ShareID         string              `json:"shareid"`
	UserID          Identity            `json:"userid"`
	MID             string              `json:"mid"`
	SType           storage.StorageType `json:"type"`
	FileName        string              `json:"filename"`
	Password        string              `json:"password"`
	ExpiredTime     int64               `json:"expire"`
	Downloads       int                 `josn:"downloads"`
	RemainDownloads int                 `json:"remainDownloads"`
}

var shareObjectMap = make(map[string]*ShareObjectInfo)
var shareMap = make(map[Identity][]string)
var mutex sync.RWMutex

var MemoCache = new(sync.Map)

func (s *ShareObjectInfo) CreateShare() (string, error) {
	uuid, err := ksuid.NewRandom()
	if err != nil {
		return "", err
	}
	s.ShareID = uuid.String()

	mutex.Lock()
	defer mutex.Unlock()

	shareObjectMap[s.ShareID] = s
	shareMap[s.UserID] = append(shareMap[s.UserID], s.ShareID)

	return s.ShareID, nil
}

func GetShareByID(shareID string) *ShareObjectInfo {
	share := shareObjectMap[shareID]
	return share
}

func (s *ShareObjectInfo) IsAvailable() bool {
	if s.RemainDownloads == 0 || (s.ExpiredTime > 0 && time.Now().Unix() > s.ExpiredTime) {
		// 考虑删除失效的分享
		s.DeleteShare()
		return false
	}

	_, err := GetFileInfo(s.UserID.Address, s.UserID.ChainID, s.MID, s.SType)
	if err != nil {
		// 文件已删除，考虑删除失效的分享
		s.DeleteShare()
		return false
	}

	return true
}

func (s *ShareObjectInfo) Source() (database.FileInfo, error) {
	return GetFileInfo(s.UserID.Address, s.UserID.ChainID, s.MID, s.SType)
}

func (s *ShareObjectInfo) CanDownload(address string, chainID int) error {
	// now anyone can download
	return nil
}

func (s *ShareObjectInfo) DownloadBy(address string, chainID int) error {
	if !s.WasDownloadedBy(address, chainID) {
		s.Downloads++
		if s.RemainDownloads > 0 {
			s.RemainDownloads--
		}

		MemoCache.Store("download"+address+s.ShareID, struct{}{})
	}
	return nil
}

func (s *ShareObjectInfo) WasDownloadedBy(address string, chainID int) bool {
	_, ok := MemoCache.Load("download" + address + s.ShareID)
	return ok
}

func (s *ShareObjectInfo) UpdateShare(attr string, value string) error {
	switch attr {
	case "password":
		s.Password = value
	default:
		return xerrors.Errorf("unsupport attribute")
	}
	return nil
}

func (s *ShareObjectInfo) DeleteShare() error {
	mutex.Lock()
	defer mutex.Unlock()
	shareIDs := shareMap[s.UserID]
	for i, shareID := range shareIDs {
		if shareID == s.ShareID {
			shareIDs = append(shareIDs[:i], shareIDs[i+1:]...)
			break
		}
	}
	shareMap[s.UserID] = shareIDs
	delete(shareObjectMap, s.ShareID)
	return nil
}

func GetFileInfo(address string, chainID int, mid string, stype storage.StorageType) (database.FileInfo, error) {
	fileInfos, err := database.Get(chainID, mid, stype)
	if err != nil {
		return database.FileInfo{}, xerrors.Errorf("Can't find the file")
	}
	for key, file := range fileInfos {
		if file.Public {
			return file, nil
		}
		if key == address {
			return file, nil
		}
	}

	return database.FileInfo{}, xerrors.Errorf("Can't access the file")
}

func ListShares(address string, chainID int) ([]ShareObjectInfo, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	shareIDs, ok := shareMap[Identity{address, chainID}]
	if !ok {
		return nil, nil
	}

	var shares []ShareObjectInfo
	for _, id := range shareIDs {
		shares = append(shares, *shareObjectMap[id])
	}
	return shares, nil
}
