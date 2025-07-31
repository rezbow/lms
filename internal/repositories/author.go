package repositories

import (
	"errors"
	"lms/internal/models"

	"gorm.io/gorm"
)

type AuthorRepo struct {
	DB *gorm.DB
}

func (ap *AuthorRepo) GetById(id string) (*models.Author, error) {
	var author models.Author
	result := ap.DB.First(&author, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &author, nil
}

func (ap *AuthorRepo) DeleteById(id string) error {
	result := ap.DB.Delete(&models.Author{}, id)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (ap *AuthorRepo) Insert(author *models.Author) error {
	if err := ap.DB.Create(author).Error; err != nil {
		return err
	}
	return nil
}
