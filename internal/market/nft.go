package market

import (
	"fmt"

	"github.com/memoio/backend/internal/database"
)

type NFT struct {
	TokenID int64  `json:"tokenid" gorm:"primaryKey"`
	Owner   string `json:"owner"`
	Sales   int64  `json:"sales"`
	NFTData
}

func InitNFTTable() error {
	return database.DataBase.AutoMigrate(&NFT{})
}

func MintNFT(tokenID int64, owner string, data NFTData) error {
	nft := NFT{
		TokenID: tokenID + 1,
		Owner:   owner,
		NFTData: data,
		Sales:   0,
	}
	return database.DataBase.Create(&nft).Error
}

func BurnNFT(tokenID int64) error {
	return database.DataBase.Delete(&NFT{}, tokenID+1).Error
}

func TransferNFT(tokenID int64, to string) error {
	var nft NFT
	err := database.DataBase.Where("token_id = ?", tokenID+1).First(&nft).Error
	if err != nil {
		return err
	}

	modifiedNFT := NFT{
		TokenID: nft.TokenID,
		Owner:   to,
		NFTData: nft.NFTData,
		Sales:   nft.Sales + 1,
	}

	return database.DataBase.Model(&nft).Updates(modifiedNFT).Error
}

func ListNFT(page int, pageSize int, order string, ascend bool) ([]NFT, error) {
	var orderOptions string
	switch order {
	case "sales":
		orderOptions = "sales"
	// case "views":
	// 	orderOptions = "views"
	case "date":
		orderOptions = "time"
	default:
		orderOptions = "token_id"
	}

	if !ascend {
		orderOptions += " desc"
	}

	fmt.Println(ascend, orderOptions)

	var nfts []NFT
	err := database.DataBase.Order(orderOptions).Offset((page - 1) * pageSize).Limit(pageSize).Find(&nfts).Error
	if err != nil {
		return nil, err
	}

	for index := range nfts {
		nfts[index].TokenID -= 1
	}

	return nfts, nil
}

// 只需保留最新未处理过的block number。
// 重启时，无需再次处理处理过的event
type BlockNumber struct {
	ID          int64 `gorm:"primaryKey"`
	BlockNumber int64
}

const blockNumberID = 1

func InitBlockNumberTable() error {
	err := database.DataBase.AutoMigrate(&BlockNumber{})
	if err != nil {
		return err
	}

	var number int64
	database.DataBase.Model(&BlockNumber{}).Count(&number)
	if number == 0 {
		return database.DataBase.Create(&BlockNumber{ID: blockNumberID, BlockNumber: 0}).Error
	}
	return nil
}

func SetLastBlockNumber(blockNumber int64) error {
	var blockNumberData = BlockNumber{
		ID:          blockNumberID,
		BlockNumber: blockNumber,
	}
	return database.DataBase.Model(&BlockNumber{ID: blockNumberID}).Updates(&blockNumberData).Error
}

func GetLastBlockNumber() int64 {
	var blockNumberData BlockNumber
	err := database.DataBase.Model(&BlockNumber{}).First(&blockNumberData).Error
	if err != nil {
		return 0
	}

	return blockNumberData.BlockNumber
}
