package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Positions struct {
	Id              int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	Symbol          string    `json:"symbol"`
	Leverage        string    `json:"leverage"`
	OpenPrice       string    `json:"open_price"`
	LiqPrice        string    `json:"liq_price"`
	TakeProfit      string    `json:"take_profit"`
	MarkPrice       string    `json:"mark_price"`
	StopLoss        string    `json:"stop_loss"`
	UnrealizedPl    string    `json:"unrealized_pl"`
	Side            string    `json:"side"`
	Size            string    `json:"size"`
	Margin          string    `json:"margin"`
	UserEmail       string    `gorm:"not null" json:"user_email"`
	Status          string    `gorm:"default:'opened'" json:"status"`
	Exchange        string    `json:"exchange"`
	LastUpdatePrice string    `json:"last_update_price"`
}

func (position *Positions) CreateNewPosition(db *gorm.DB) (*Positions, error) {
	err := db.Create(&position).Error
	if err != nil {
		return &Positions{}, err
	}
	return position, nil
}

func (u *Positions) GetOpenPositionsByExchange(db *gorm.DB, exchange string) (*[]Positions, error) {
	positions := []Positions{}
	err := db.Model(&Positions{}).Where("exchange = ? AND status = ?", exchange, "opened").Find(&positions).Error
	if err != nil {
		return &[]Positions{}, err
	}
	return &positions, nil
}

func (position *Positions) UpdateOrCreatePosition(db *gorm.DB) (*Positions, error) {
	db.Table("positions")
	existingPosition := &Positions{}
	err := db.Where("symbol = ? AND user_email = ? AND side = ? AND exchange = ?", position.Symbol, position.UserEmail, position.Side, position.Exchange).First(existingPosition).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return position.CreateNewPosition(db)
	}

	existingPosition.OpenPrice = position.OpenPrice
	existingPosition.LiqPrice = position.LiqPrice
	existingPosition.TakeProfit = position.TakeProfit
	existingPosition.StopLoss = position.StopLoss
	existingPosition.UnrealizedPl = position.UnrealizedPl
	existingPosition.MarkPrice = position.MarkPrice
	existingPosition.Side = position.Side
	existingPosition.Size = position.Size
	existingPosition.Margin = position.Margin
	existingPosition.Status = position.Status
	existingPosition.Exchange = position.Exchange
	err = db.Save(existingPosition).Error
	if err != nil {
		return nil, err
	}

	return existingPosition, nil
}
