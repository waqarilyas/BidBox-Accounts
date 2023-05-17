package admin

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Conditions struct {
	Capital    int `json:"capital"`
	Positions  int `json:"positions"`
	StopLoss   int `json:"stop_loss"`
	TakeProfit int `json:"take_profit"`
	Leverage   int `json:"leverage"`
}

func (c *Conditions) Validate(prev *Conditions) {
	if c.Leverage == 0 {
		c.Leverage = prev.Leverage
	}
	if c.Positions == 0 {
		c.Positions = prev.Positions
	}
	if c.StopLoss == 0 {
		c.StopLoss = prev.StopLoss
	}
	if c.TakeProfit == 0 {
		c.TakeProfit = prev.TakeProfit
	}
}

func (c *Conditions) GetConditions(db *gorm.DB, limit int, offset int) (*[]Conditions, int, error) {
	Condition := []Conditions{}
	err := db.Debug().Model(&Conditions{}).Order("capital").Limit(limit).Offset(offset).Find(&Condition).Error
	if err != nil {
		return &[]Conditions{}, 0, err
	}
	var count int
	db.Model(&Conditions{}).Count(&count)

	return &Condition, count, nil
}

func (c *Conditions) FindConditionById(db *gorm.DB, capital int) (*Conditions, error) {
	cond := &Conditions{}
	err := db.Debug().Model(Conditions{}).Where("capital = ?", capital).Take(&cond).Error
	if err != nil {
		return &Conditions{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Conditions{}, errors.New("condition not found")
	}
	return cond, nil
}

func (c *Conditions) UpdateConditions(db *gorm.DB, capital int) (*Conditions, error) {
	db = db.Debug().Model(&Conditions{}).Where("capital = ?", capital).Take(&Conditions{}).
		UpdateColumns(
			map[string]interface{}{
				"take_profit": c.TakeProfit,
				"stop_loss":   c.StopLoss,
				"leverage":    c.Leverage,
				"positions":   c.Positions,
			},
		)
	if db.Error != nil {
		return &Conditions{}, db.Error
	}
	return c, nil
}
