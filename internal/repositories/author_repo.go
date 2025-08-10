package repositories

import (
	"errors"
	"fmt"
	"lms/internal/models"
	"lms/internal/views"
	"log"

	"gorm.io/gorm"
)

type AuthorRepo struct {
	DB *gorm.DB
}

func (ap *AuthorRepo) GetById(id uint) (*models.Author, error) {
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

func (ap *AuthorRepo) DeleteById(id uint) error {
	result := ap.DB.Delete(&models.Author{}, id)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return ErrInternal
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (ap *AuthorRepo) Insert(author *models.Author) error {
	if err := ap.DB.Create(author).Error; err != nil {
		if isInternalError(err) {
			return ErrInternal
		}
		return err
	}
	return nil
}

func (ap *AuthorRepo) Update(author *models.Author) error {
	result := ap.DB.Model(&models.Author{}).
		Where("id = ?", author.ID).
		Updates(author)
	if result.Error != nil {
		log.Println(result.Error.Error())
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (ar *AuthorRepo) Total() (int64, error) {
	var total int64
	err := ar.DB.Model(&models.Author{}).Count(&total).Error
	if err != nil {
		log.Println(err.Error())
		return 0, ErrInternal
	}
	return total, nil
}

func (ar *AuthorRepo) Search(data *models.SearchData) ([]models.Author, error) {
	var authors []models.Author
	query := ar.DB.Model(&models.Author{})
	s := "%" + data.Term + "%"

	query.Where("full_name ILIKE ?", s)

	query.Count(&data.Pagination.Total)
	if data.SortBy != "" {
		query.Order(fmt.Sprintf("%s %s", data.SortBy, data.Dir))
	}

	result := query.Offset(data.Pagination.Offset).Limit(data.Pagination.Limit).Find(&authors)

	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, ErrInternal
	}
	return authors, nil
}

func (ar *AuthorRepo) RecentAuthors(limit int) ([]models.Author, error) {
	var authors []models.Author
	err := ar.DB.Model(&models.Author{}).
		Order("created_at DESC").
		Limit(limit).
		Find(&authors).Error
	if err != nil {
		log.Println(err.Error())
		return nil, ErrInternal
	}
	return authors, nil
}

func (ar *AuthorRepo) PopularAuthors(limit int) ([]models.Author, error) {
	var authors []models.Author
	err := ar.DB.Table("popular_authors_view").
		Limit(limit).
		Scan(&authors).Error

	if err != nil {
		log.Println(err.Error())
		return nil, ErrInternal
	}
	return authors, nil
}

func (ar *AuthorRepo) ErrorToFromError(err error) views.Errors {
	errors := make(views.Errors)
	return errors
}
