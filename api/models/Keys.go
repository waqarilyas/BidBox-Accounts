package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// swagger:model Key
type Key struct {
	// ID of Key
	// in: uuid
	Keyid uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"key_id"`
	// Exchange of key
	// in: string
	Service string `gorm:"size:255;not null" json:"service"`
	// Api Key
	// in: string
	ApiKey string `gorm:"not null;unique" json:"api_key"`
	// Secret Key
	// in: string
	SecretKey string `gorm:"not null;unique" json:"secret"`
	// Passphrase
	// in: string
	Passphrase string `gorm:"" json:"passphrase"`
	// User Email
	// in: string
	UserEmail string `json:"user_email"`
	// Strategy
	// in: string
	Strategy string `json:"strategy"`
	// Mode
	// in: string
	Mode string `json:"mode"`
	// Auto Compound
	// in: string
	Compound bool `json:"compound"`
	// Trade Amount
	// in: int
	TradeAmount int `json:"trade_amount"`
	// Enable Key
	// in: bool
	Start bool `json:"start"`
}

func (u *Key) FindAllKeys(db *gorm.DB) (*[]Key, error) {
	Keys := []Key{}
	err := db.Debug().Model(&Key{}).Limit(100).Find(&Keys).Error
	if err != nil {
		return &[]Key{}, err
	}
	return &Keys, nil
}

func (u *Key) FindKeyById(db *gorm.DB, kid uuid.UUID) (*Key, error) {
	err := db.Debug().Model(Key{}).Where("keyid = ?", kid).Take(&u).Error
	if err != nil {
		return &Key{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Key{}, errors.New("Key not found")
	}
	return u, nil
}

func (u *Key) FindNoOfUsers(db *gorm.DB) (int, error) {
	var count int

	query := db.Model(Key{}).
		Group("user_email").
		Count(&count)
	if query.Error != nil {
		return 0, query.Error
	}

	return count, nil
}

func (u *Key) FindKeysByEmail(db *gorm.DB, email string, service string) (*Key, error) {
	Keys := Key{}
	err := db.Debug().Model(Key{}).Where("user_email = ? AND service = ?", email, service).Find(&Keys).Error
	if err != nil {
		return &Key{}, err
	}
	return &Keys, nil
}

func (u *Key) FindKeysByUserId(db *gorm.DB, uid uuid.UUID) (*[]Key, error) {
	Keys := []Key{}
	err := db.Debug().Model(Key{}).Where("uid = ?", uid).Find(&Keys).Error
	if err != nil {
		return &[]Key{}, err
	}
	return &Keys, nil
}

func (u *Key) FindKeysByUserEmail(db *gorm.DB, email string) (*[]Key, error) {
	Keys := []Key{}
	err := db.Debug().Model(Key{}).Where("user_email = ?", email).Find(&Keys).Error
	if err != nil {
		return &[]Key{}, err
	}
	return &Keys, nil
}

func (u *Key) FindKeyByUserIdAndShort(db *gorm.DB, uid uuid.UUID, service string) (*Key, error) {
	Keys := Key{}
	err := db.Debug().Model(Key{}).Where("uid = ? AND service = ? ", uid, service).Find(&Keys).Error
	if err != nil {
		return &Key{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Key{}, errors.New("no connected keys found")
	}
	return &Keys, nil
}

func (u *Key) FindKeyByUserEmailAndShort(db *gorm.DB, email string, service string) (*Key, error) {
	Keys := Key{}
	err := db.Debug().Model(Key{}).Where("user_email = ? AND service = ? ", email, service).Find(&Keys).Error
	if err != nil {
		return &Key{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Key{}, errors.New("no connected keys found")
	}
	return &Keys, nil
}

var strategy = []string{"cycle", "single", "stop_make", "stop_long", "stop_short"}
var modes = []string{"conservative", "aggressive"}

func (k *Key) Validate(prev *Key) error {
	st := false
	md := false
	if k.Strategy == "" {
		k.Strategy = prev.Strategy
	}
	if k.Mode == "" {
		k.Mode = prev.Mode
	}
	if k.Compound == prev.Compound {
		k.Compound = prev.Compound
	}
	if k.Start == prev.Start {
		k.Start = prev.Start
	}
	for _, v := range modes {
		if v == k.Mode {
			md = true
			break
		}
	}
	for _, v := range strategy {
		if v == k.Strategy {
			st = true
			break
		}
	}
	if st && md {
		return nil
	}
	return errors.New("strategy or mode is incorrect")
}

func (k *Key) UpdateKeySettings(db *gorm.DB, email string, service string) (*Key, error) {
	updated_key := Key{}
	db = db.Debug().Model(&Key{}).Where("user_email = ? AND service = ?", email, service).Take(&updated_key).
		UpdateColumns(
			map[string]interface{}{
				"strategy": k.Strategy,
				"mode":     k.Mode,
				"compound": k.Compound,
				"start":    k.Start,
			},
		)
	if db.Error != nil {
		return &Key{}, db.Error
	}
	return &updated_key, nil
}

func (k *Key) GetSettings(db *gorm.DB, email string) (*[]Key, error) {
	keys := []Key{}
	err := db.Debug().Model(&Key{}).Where("user_email = ?", email).Find(&keys).Error
	if err != nil {
		return &[]Key{}, err
	}
	return &keys, nil
}

func (k *Key) DeleteKey(db *gorm.DB, email string, service string) (int64, error) {
	db = db.Debug().Model(&Key{}).Where("user_email = ? AND service = ?", email, service).Take(&Key{}).Delete(&Key{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
