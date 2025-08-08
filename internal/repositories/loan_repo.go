package repositories

import (
	"errors"
	"fmt"
	"lms/internal/models"
	"lms/internal/utils"
	"lms/internal/views"

	"gorm.io/gorm"
)

type LoanRepo struct {
	DB *gorm.DB
}

func (lp *LoanRepo) GetById(id uint) (*models.Loan, error) {
	var loan models.Loan
	result := lp.DB.Preload("Book").Preload("Member").First(&loan, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &loan, nil
}

func (lp *LoanRepo) DeleteById(id uint) error {
	result := lp.DB.Delete(&models.Loan{}, id)
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

func (lp *LoanRepo) Insert(loan *models.Loan) error {
	result := lp.DB.Create(loan)
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	}
	return nil
}

func (lp *LoanRepo) All(pagination *utils.Pagination) ([]models.Loan, error) {
	var loans []models.Loan
	query := lp.DB.Model(&models.Loan{}).Count(&pagination.Total)
	result := query.Preload("Book").Preload("Member").Limit(pagination.Limit).Offset(pagination.Offset).Find(&loans)
	if result.Error != nil {
		return nil, ErrInternal
	}

	pagination.CalculateTotalPage()

	return loans, nil
}

func (lr *LoanRepo) Update(loan *models.Loan) error {
	result := lr.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(loan)
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

func (lr *LoanRepo) ConvertErrorToFormError(err error) views.Errors {
	errors := make(views.Errors)
	pgError := extractPQError(err)
	switch {
	case pgError.ConstraintName == "loans_book_id_fkey":
		errors["bookId"] = "book not found"
	case pgError.ConstraintName == "loans_member_id_fkey":
		errors["memberId"] = "member not found"
	case pgError.ConstraintName == "unique_active_loan":
		errors["memberId"] = "this member already borrowed this book"
	case pgError.Message == "OVER_BORROWING":
		errors["bookId"] = "no available copy to borrow"
	default:
		errors["_"] = err.Error()

	}
	return errors
}

func (lr *LoanRepo) Filter(filter *models.LoanFilter, pagination *utils.Pagination) ([]models.Loan, error) {
	var loans []models.Loan
	query := lr.DB.Model(&models.Loan{})

	if filter.BookId > 0 {
		query.Where("book_id = ? ", filter.BookId)
	}

	if filter.MemberId > 0 {
		query.Where("member_id = ? ", filter.MemberId)
	}

	if filter.Status != "" {
		query.Where("status = ? ", filter.Status)
	}

	err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&loans).Error

	if err != nil {
		return nil, ErrInternal
	}

	return loans, nil

}

func (lr *LoanRepo) Search(term string, pagination *utils.Pagination) ([]models.Loan, error) {
	var loans []models.Loan
	query := lr.DB.Model(&models.Loan{}).Preload("Member").Preload("Book")
	s := "%" + term + "%"

	query.Where("CAST(id as TEXT) ILIKE ?", s)

	query.Count(&pagination.Total)
	err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&loans).Error

	if err != nil {
		return nil, ErrInternal
	}

	pagination.CalculateTotalPage()
	fmt.Println(loans, pagination)

	return loans, nil

}
