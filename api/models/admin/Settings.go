package admin

import "github.com/jinzhu/gorm"

type Conditions struct {
	Capital    int
	Positions  int
	StopLoss   int
	TakeProfit int
	Leverage   int
}

func (c *Conditions) GetConditions(db *gorm.DB, limit int, offset int) (*[]Conditions, error) {
	Condition := []Conditions{}
	err := db.Debug().Model(&Conditions{}).Limit(limit).Offset(offset).Find(&Condition).Error
	if err != nil {
		return &[]Conditions{}, err
	}
	return &Condition, nil
}
