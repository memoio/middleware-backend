package share

import (
	"encoding/base64"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/storage"
	"golang.org/x/xerrors"
)

type Identity struct {
	Address string
	ChainID int
}

type ShareObjectInfo struct {
	ShareID         string
	UserID          Identity
	MID             string
	SType           storage.StorageType
	FileName        string
	Password        string
	ExpiredTime     int64
	Downloads       int
	RemainDownloads int
}

var shareObjectMap = make(map[string]*ShareObjectInfo)
var shareMap = make(map[Identity][]string)
var mutex sync.RWMutex

var MemoCache = new(sync.Map)

func (s *ShareObjectInfo) CreateShare() (string, error) {
	data := crypto.Keccak256([]byte(s.UserID.Address), []byte(strconv.Itoa(s.UserID.ChainID)), []byte(s.MID))
	s.ShareID = base64.StdEncoding.EncodeToString(data)

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
	if s.RemainDownloads == 0 || time.Now().Unix() > s.ExpiredTime {
		// 考虑删除失效的分享
		return false
	}

	_, err := database.Get(s.UserID.ChainID, s.UserID.Address, s.MID, s.SType)
	if err != nil {
		// 文件已删除，考虑删除失效的分享
		return false
	}

	return true
}

func (s *ShareObjectInfo) Source() (database.FileInfo, error) {
	return database.Get(s.UserID.ChainID, s.UserID.Address, s.MID, s.SType)
}

// func (s *ShareObjectInfo) Preview() error {
// 	return nil
// }

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
