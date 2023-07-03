package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Conditions struct {
	Capital   int
	Positions int
	Leverage  int
	StopLoss  int
}

func (u *Conditions) FindAllConditions(db *gorm.DB) (*[]Conditions, error) {
	conditions := []Conditions{}
	err := db.Model(Conditions{}).Limit(100).Take(conditions).Error
	if err != nil {
		return &[]Conditions{}, err
	}
	return &conditions, nil
}

func (u *Conditions) FindCondition(db *gorm.DB, capital int) (*Conditions, error) {
	cond := Conditions{}
	err := db.Model(Conditions{}).Where("capital = ?", capital).Take(&cond).Error
	if err != nil {
		return &Conditions{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Conditions{}, errors.New("Key not found")
	}
	return &cond, nil
}
