package repositories

import (
	"errors"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func (repo *Repo) All(target any, page, pageSize int) error {
	offset := (page - 1) * pageSize
	if err := repo.DB.Limit(pageSize).Offset(offset).Find(target).Error; err != nil {
		return ErrInternal
	}
	return nil
}

func (repo *Repo) GetById(target any, id int) error {
	result := repo.DB.First(&target, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}
	return nil
}

func (repo *Repo) DeleteById(target any, id int) error {
	result := repo.DB.Delete(target, id)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
