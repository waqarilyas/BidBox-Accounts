package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Exchanges struct {
	Id       uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"size:255;not null" json:"name"`
	ImageSrc string    `gorm:"not null;unique" json:"image_src"`
	Short    string    `gorm:"not null;unique" json:"short"`
	IsActive bool      `gorm:"" json:"is_active"`
}

func (e *Exchanges) FindAllExchanges(db *gorm.DB) (*[]Exchanges, error) {
	Exchange := []Exchanges{}
	err := db.Debug().Model(&Exchanges{}).Limit(100).Find(&Exchange).Error
	if err != nil {
		return &[]Exchanges{}, err
	}
	return &Exchange, nil
}

func (e *Exchanges) GetExchangeByShort(db *gorm.DB, short string) (*Exchanges, error) {
	err := db.Debug().Model(Exchanges{}).Where("short = ?", short).Take(&e).Error
	if err != nil {
		return &Exchanges{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Exchanges{}, errors.New("no exchange found with given short")
	}
	return e, nil

}
