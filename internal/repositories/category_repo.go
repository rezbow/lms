package repositories

import (
	"errors"
	"lms/internal/models"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	DB *gorm.DB
}

func (cr *CategoryRepo) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	result := cr.DB.Preload("Books").Where("slug = ?", slug).First(&category)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &category, nil
}

func (cr *CategoryRepo) DeleteBySlug(slug string) error {
	result := cr.DB.Where("slug = ?", slug).Delete(&models.Category{})
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (cr *CategoryRepo) Insert(category *models.Category) error {
	result := cr.DB.Create(category)
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	}
	return nil
}

func (cr *CategoryRepo) Update(category *models.Category) error {
	result := cr.DB.Model(&models.Category{}).Where("id = ?", category.ID).Updates(category)
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

/*
func (cr *CategoryRepo) ConvertErrorsToFormErrors(err error) views.Errors {
	errors := make(views.Errors)
	pgErr := extractPQError(err)
	switch {
	case pgErr.Code == pgerrcode.ForeignKeyViolation:
		errors["authorId"] = "cannot find author"
	case pgErr.ConstraintName == "books_isbn_key":
		errors["isbn"] = "a books with this isbn already exists"
	case pgErr.Message == "INVALID_TOTAL_COPIES":
		errors["totalCopies"] = "invalid total copies"
	case pgErr.Message == "INVALID_AVAILABLE_COPIES":
		errors["availableCopies"] = "invalid available copies"
	default:
		errors["_"] = err.Error()
	}
	return errors
}
*/
