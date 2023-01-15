package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) GetInputByID(inputID uint) (*Input, error) {
	var input Input
	err := service.db.First(&input, "id = ?", inputID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &input, err
}
