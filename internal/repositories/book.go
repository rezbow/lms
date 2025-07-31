package repositories

import (
	"errors"
	"lms/internal/models"

	"gorm.io/gorm"
)

type BookRepo struct {
	DB *gorm.DB
}

var (
	ErrBookHasActiveLoan = errors.New("this book has active loans and cannot be deleted")
)

func (bp *BookRepo) GetById(id string) (*models.Book, error) {
	var book models.Book
	result := bp.DB.Joins("left join authors on books.author_id = authors.id").
		Preload("Author").
		First(&book, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &book, nil
}

func (bp *BookRepo) DeleteById(id string) error {
	result := bp.DB.Delete(&models.Book{}, id)
	if result.Error != nil {
		if result.Error == gorm.ErrForeignKeyViolated {
			return ErrBookHasActiveLoan
		}
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (bp *BookRepo) Insert(book *models.Book) error {
	if err := bp.DB.Create(book).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
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
	limit int,
	page int,
	total *int64,
) ([]models.Book, error) {
	var books []models.Book
	query := bp.DB.Model(&models.Book{}).
		Joins("inner join authors on books.author_id = authors.id").
		Preload("Author")

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
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
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
	if err := bp.DB.Limit(pageSize).Offset(offset).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (bp *BookRepo) Update(book *models.Book) error {
	if err := bp.DB.Model(&models.Book{}).Where("id = ?", book.ID).Updates(book).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return ErrAuthorIdNotFound
		}
		return err
	}
	return nil
}
