package repositories

import (
	"errors"
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
	result := lp.DB.First(&loan, id)

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

func (lp *LoanRepo) All(page, pageSize int) ([]models.Loan, error) {
	var loans []models.Loan
	offset := (page - 1) * pageSize
	result := lp.DB.Limit(pageSize).Offset(offset).Find(&loans)
	if result.Error != nil {
		return nil, ErrInternal
	}
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

func (lr *LoanRepo) Filter(bookId, memberId int, status string, pagination *utils.Pagination) ([]models.Loan, error) {
	var loans []models.Loan
	query := lr.DB.Model(&models.Loan{})

	if bookId >= 0 {
		query.Where("book_id = ? ", bookId)
	}

	if memberId >= 0 {
		query.Where("member_id = ? ", memberId)
	}

	if status != "" {
		query.Where("status = ? ", status)
	}

	err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&loans).Error

	if err != nil {
		return nil, ErrInternal
	}

	return loans, nil

}
