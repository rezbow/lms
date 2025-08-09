package repositories

import (
	"errors"
	"fmt"
	"lms/internal/models"
	"lms/internal/views"

	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

type BookRepo struct {
	DB *gorm.DB
}

func (bp *BookRepo) GetById(id uint) (*models.Book, error) {
	var book models.Book
	result := bp.DB.Preload("Loans").Preload("Author").Preload("Categories").First(&book, id)

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

func (bp *BookRepo) Total() int64 {
	var total int64
	bp.DB.Model(&models.Book{}).Count(&total)
	return total
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

func (bp *BookRepo) Search(
	data *models.SearchData,
) ([]models.Book, error) {
	var books []models.Book
	query := bp.DB.Model(&models.Book{})

	s := "%" + data.Term + "%"

	query.Where("books.title ILIKE ?", s).
		Or("CAST(books.author_id as TEXT) ILIKE ?", s).
		Or("books.isbn ILIKE ?", s).
		Or("books.publisher ILIKE ?", s)

	query.Count(&data.Pagination.Total)

	if data.SortBy != "" {
		query.Order(fmt.Sprintf("%s %s", data.SortBy, data.Dir))
	}

	result := query.
		Offset(data.Pagination.Offset).
		Limit(data.Pagination.Limit).
		Find(&books)

	if result.Error != nil {
		return nil, ErrInternal
	}

	data.Pagination.CalculateTotalPage()

	return books, nil
}
