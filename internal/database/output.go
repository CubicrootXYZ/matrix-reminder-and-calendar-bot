package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) GetOutputByID(outputID uint) (*Output, error) {
	var output Output
	err := service.db.First(&output, "id = ?", outputID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &output, err
}

func (service *service) UpdateOutput(output *Output) (*Output, error) {
	err := service.db.Save(output).Error

	return output, err
}
