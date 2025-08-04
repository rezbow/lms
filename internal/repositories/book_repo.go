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
	filter *models.BookFilter,
	pagination *utils.Pagination,
	total *int64,
) ([]models.Book, error) {
	var books []models.Book
	query := bp.DB.Model(&models.Book{})

	if filter.Title != "" {
		query.Where("books.title ILIKE ?", "%"+filter.Title+"%")
	}

	if filter.AuthorId > 0 {
		query.Where("books.author_id = ?", filter.AuthorId)
	}

	if filter.ISBN != "" {
		query.Where("books.isbn = ?", filter.ISBN)
	}

	if filter.Publisher != "" {
		query.Where("books.publisher ILIKE ?", "%"+filter.Publisher+"%")
	}

	if filter.Language != "" {
		query.Where("books.language ILIKE ? ", "%"+filter.Language+"%")
	}

	if filter.Translator != "" {
		query.Where("books.translator ILIKE ? ", "%"+filter.Translator+"%")
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

func (bp *BookRepo) All(pagination *utils.Pagination) ([]models.Book, error) {
	var books []models.Book
	if err := bp.DB.Preload("Loans").Limit(pagination.Limit).Offset(pagination.Offset).Find(&books).Error; err != nil {
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
