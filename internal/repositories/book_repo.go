package repositories

import (
	"errors"
	"lms/internal/models"
	"lms/internal/utils"

	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

type BookRepo struct {
	DB *gorm.DB
}

var (
	ErrBookHasActiveLoan = errors.New("this book has active loans and cannot be deleted")
)

func (bp *BookRepo) GetById(id int) (*models.Book, error) {
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

func (bp *BookRepo) DeleteById(id int) error {
	result := bp.DB.Delete(&models.Book{}, id)
	if result.Error != nil {
		pgErr := extractPQError(result.Error)
		if pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrBookHasActiveLoan
		}
		return ErrInternal
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (bp *BookRepo) Insert(book *models.Book) error {
	if err := bp.DB.Create(book).Error; err != nil {
		pgErr := extractPQError(err)
		if pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrAuthorIdNotFound
		}
		return err
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
	if err := bp.DB.Model(&models.Book{}).Where("id = ?", book.ID).Updates(book).Error; err != nil {
		pgErr := extractPQError(err)
		if pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrAuthorIdNotFound
		}
		return err
	}
	return nil
}
