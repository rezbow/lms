package repositories

import (
	"errors"
	"lms/internal/models"

	"gorm.io/gorm"
)

type ActivityRepo struct {
	DB *gorm.DB
}

func (ar *ActivityRepo) Add(activityLog *models.ActivityLog) error {
	err := ar.DB.Create(&activityLog).Error
	if err != nil {
		return ErrInternal
	}
	return nil
}

func (ar *ActivityRepo) GetById(id uint) (*models.ActivityLog, error) {
	var activityLog models.ActivityLog
	err := ar.DB.Model(&models.ActivityLog{}).First(&activityLog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &activityLog, nil
}

func (ar *ActivityRepo) Recent(limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := ar.DB.Model(&models.ActivityLog{}).Limit(limit).Order("created_at desc").Find(&logs).Error
	if err != nil {
		return nil, ErrInternal
	}
	return logs, nil
}
