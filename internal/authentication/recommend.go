package auth

import "github.com/memoio/backend/internal/database"

type Recommend struct {
	Address     string `json:"address" gorm:"primaryKey"`
	Recommender string `json:"recommender"`
	Source      string `json:"source"`
}

func InitRecommendTable() error {
	return database.DataBase.AutoMigrate(&Recommend{})
}

func GetRecommend(address string) (*Recommend, error) {
	var recommend Recommend
	if err := database.DataBase.Model(&Recommend{}).Where("address = ?", address).First(&recommend).Error; err != nil {
		return nil, err
	}
	return &recommend, nil
}

func ListRecommend(recommender string) ([]Recommend, error) {
	var recommends []Recommend
	if recommender == "" {
		err := database.DataBase.Model(&Recommend{}).Find(&recommends).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := database.DataBase.Model(&Recommend{}).Where("recommender = ?", recommender).Find(&recommends).Error
		if err != nil {
			return nil, err
		}
	}
	return recommends, nil
}

func (r *Recommend) CreateRecommend() error {
	return database.DataBase.Create(r).Error
}
