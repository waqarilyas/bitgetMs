package models

import "github.com/jinzhu/gorm"

type CoinPair struct {
	Id   int    `json:"id"`
	Coin string `json:"coin"`
}

func (cp *CoinPair) GetAllCoins(db *gorm.DB) (*[]CoinPair, error) {
	Coins := []CoinPair{}
	err := db.Debug().Model(&CoinPair{}).Limit(100).Find(&Coins).Error
	if err != nil {
		return &[]CoinPair{}, err
	}
	return &Coins, nil
}

func (cp *CoinPair) SaveCoinPair(db *gorm.DB) (*CoinPair, error) {
	err := db.Debug().Create(cp).Error
	if err != nil {
		return &CoinPair{}, err
	}
	return cp, nil
}

func (cp *CoinPair) UpdateCoinPair(db *gorm.DB, coin string) (*CoinPair, error) {
	db = db.Debug().Model(&CoinPair{}).Where("coin = ?", coin).Take(&CoinPair{}).UpdateColumns(
		map[string]interface{}{
			"coin": cp.Coin,
		},
	)
	if db.Error != nil {
		return &CoinPair{}, db.Error
	}
	cp.Coin = coin
	return cp, nil
}

func (cp *CoinPair) DeleteCoinPair(db *gorm.DB, coin string) (int64, error) {
	db = db.Debug().Model(&CoinPair{}).Where("coin = ?", coin).Take(&CoinPair{}).Delete(&CoinPair{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
