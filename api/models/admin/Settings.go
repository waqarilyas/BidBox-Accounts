package admin

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Settings struct {
	Timeframe        string  `json:"timeframe"`
	Maintainence     bool    `json:"maintainence"`
	Leverage         int     `json:"leverage"`
	Layers           int     `json:"layers"`
	ProfitPercentage float64 `json:"profit_percentage"`
}

type Conditions struct {
	Capital    int `json:"capital" example:"300"`
	Positions  int `json:"positions" example:"2"`
	StopLoss   int `json:"stop_loss" example:"80"`
	TakeProfit int `json:"take_profit" example:"2"`
	Leverage   int `json:"leverage" example:"80"`
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

func (s *Settings) UpdateMaintainance(db *gorm.DB, val bool) (*Settings, error) {
	db = db.Model(&Settings{}).
		UpdateColumns(
			map[string]interface{}{
				"maintainance": val,
			},
		)
	if db.Error != nil {
		return &Settings{}, db.Error
	}
	return s, nil

}

func (s *Settings) UpdateTimeframe(db *gorm.DB, val string) (*Settings, error) {
	db = db.Model(&Settings{}).
		UpdateColumns(
			map[string]interface{}{
				"timeframe": val,
			},
		)
	if db.Error != nil {
		return &Settings{}, db.Error
	}
	return s, nil

}

func (s *Settings) UpdateProfitPercentage(db *gorm.DB, val float64) (*Settings, error) {
	db = db.Model(&Settings{}).
		UpdateColumns(
			map[string]interface{}{
				"profit_percentage": val,
			},
		)
	if db.Error != nil {
		return &Settings{}, db.Error
	}
	return s, nil
}

func (s *Settings) UpdateLayers(db *gorm.DB, val int) (*Settings, error) {
	db = db.Model(&Settings{}).
		UpdateColumns(
			map[string]interface{}{
				"layers": val,
			},
		)
	if db.Error != nil {
		return &Settings{}, db.Error
	}
	return s, nil
}

func (s *Settings) UpdateLeverage(db *gorm.DB, val int) (*Settings, error) {
	db = db.Model(&Settings{}).
		UpdateColumns(
			map[string]interface{}{
				"leverage": val,
			},
		)
	if db.Error != nil {
		return &Settings{}, db.Error
	}
	return s, nil
}

func (s *Settings) GetTimeframe(db *gorm.DB) (string, error) {
	// db = db.Model(&Settings{})
	// if db.Error != nil {
	// 	return &Settings{}, db.Error
	// }
	// return , nil

	setting := &Settings{}
	err := db.Debug().Model(&Conditions{}).Find(&setting).Error
	if err != nil {
		return "", err
	}
	return setting.Timeframe, nil
}

func (s *Settings) GetSettings(db *gorm.DB) (*Settings, error) {
	setting := Settings{}
	err := db.Model(Settings{}).Take(&setting).Error
	if err != nil {
		return &Settings{}, err
	}
	return &setting, nil
}
