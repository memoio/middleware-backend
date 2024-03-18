package da

import (
	"bytes"
	"context"
	"encoding/hex"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/database"
	"gorm.io/gorm"
)

type DAFileInfo struct {
	// gorm.Model
	Commit     bls12381.G1Affine
	Mid        string
	Size       int64
	Expiration int64
}

type DAFileInfoStore struct {
	gorm.Model
	Commit     string `gorm:"index;column:commit"`
	Mid        string `gorm:"index;column:mid"`
	Size       int64
	Expiration int64
}

func InitDAFileInfoTable() error {
	return database.GlobalDataBase.AutoMigrate(&DAFileInfoStore{})
}

func (f *DAFileInfo) CreateDAFileInfo() error {
	commitByte48 := f.Commit.Bytes()
	var info = &DAFileInfoStore{
		Commit:     hex.EncodeToString(commitByte48[:]),
		Mid:        f.Mid,
		Size:       f.Size,
		Expiration: f.Expiration,
	}
	return database.GlobalDataBase.Create(info).Error
}

func GetDAFileLength() (int64, error) {
	var length int64
	err := database.GlobalDataBase.Model(&DAFileInfoStore{}).Count(&length).Error
	return length, err
}

func GetRangeDAFileInfo(start uint, end uint) ([]DAFileInfo, error) {
	var files []DAFileInfoStore
	var result []DAFileInfo
	err := database.GlobalDataBase.Model(&DAFileInfoStore{}).Where("id >= ? and id <= ?", start, end).Find(&files).Error
	if err != nil {
		return nil, err
	}

	result = make([]DAFileInfo, len(files))
	for index, file := range files {
		var commit bls12381.G1Affine
		commitByte, err := hex.DecodeString(file.Commit)
		if err != nil {
			return nil, err
		}
		_, err = commit.SetBytes(commitByte)
		if err != nil {
			return nil, err
		}

		result[index] = DAFileInfo{
			Commit:     commit,
			Mid:        file.Mid,
			Size:       file.Size,
			Expiration: file.Expiration,
		}
	}
	return result, nil
}

func GetFileInfoByCommit(commit bls12381.G1Affine) (DAFileInfo, error) {
	var file DAFileInfoStore
	commitByte48 := commit.Bytes()
	err := database.GlobalDataBase.Model(&DAFileInfoStore{}).Where("\"commit\" = ?", hex.EncodeToString(commitByte48[:])).First(&file).Error

	return DAFileInfo{
		Commit:     commit,
		Mid:        file.Mid,
		Size:       file.Size,
		Expiration: file.Expiration,
	}, err
}

func GetFileByCommit(commit bls12381.G1Affine) ([]byte, error) {
	var file DAFileInfo
	var buf bytes.Buffer
	commitByte48 := commit.Bytes()
	err := database.GlobalDataBase.Model(&DAFileInfoStore{}).Where("\"commit\" = ?", hex.EncodeToString(commitByte48[:])).First(&file).Error
	if err != nil {
		return nil, err
	}

	err = daStore.GetObject(context.TODO(), file.Mid, &buf, api.ObjectOptions{})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
