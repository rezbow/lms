package repositories

import (
	"errors"
	"lms/internal/models"
	"lms/internal/utils"
	"lms/internal/views"

	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

type BookRepo struct {
	DB *gorm.DB
}

func (bp *BookRepo) GetById(id uint) (*models.Book, error) {
	var book models.Book
	result := bp.DB.Preload("Loans").First(&book, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &book, nil
}

func (bp *BookRepo) DeleteById(id uint) error {
	result := bp.DB.Delete(&models.Book{}, id)
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

func (bp *BookRepo) Insert(book *models.Book) error {
	result := bp.DB.Create(book)
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	}
	return nil
}

func (bp *BookRepo) Filter(
	title string,
	author string,
	isbn string,
	pagination *utils.Pagination,
	total *int64,
) ([]models.Book, error) {
	var books []models.Book
	query := bp.DB.Model(&models.Book{}).Preload("Loans")

	if title != "" {
		query.Where("books.title_fa ILIKE ?", "%"+title+"%").
			Or("books.title_en ILIKE ?", "%"+title+"%")
	}

	if author != "" {
		query.Where("authors.name_fa ILIKE ?", "%"+author+"%").
			Or("authors.name_en ILIKE ?", "%"+author+"%")
	}

	if isbn != "" {
		query.Where("isbn = ?", isbn)
	}

	query.Count(total)

	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (bp *BookRepo) Total() int64 {
	var total int64
	bp.DB.Model(&models.Book{}).Count(&total)
	return total
}

func (bp *BookRepo) All(page int, pageSize int) ([]models.Book, error) {
	var books []models.Book
	offset := (page - 1) * pageSize
	if err := bp.DB.Preload("Loan").Limit(pageSize).Offset(offset).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (bp *BookRepo) Update(book *models.Book) error {

	result := bp.DB.Model(&models.Book{}).Where("id = ?", book.ID).Updates(book)
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

func (bp *BookRepo) ConvertErrorsToFormErrors(err error) views.Errors {
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
