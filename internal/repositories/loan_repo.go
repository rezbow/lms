package repositories

import (
	"errors"
	"fmt"
	"lms/internal/models"
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

func (lp *LoanRepo) TotalOverdueLoans() (int64, error) {
	var total int64
	err := lp.DB.Model(&models.Loan{}).
		Where("due_date < CURRENT_DATE").
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (lp *LoanRepo) OverdueLoans(limit int) ([]models.Loan, error) {
	var loans []models.Loan
	err := lp.DB.
		Model(&models.Loan{}).
		Preload("Member").
		Preload("Book").
		Where("return_date is NULL").
		Where("status = 'borrowed'").
		Where("due_date < CURRENT_TIMESTAMP").
		Order("due_date DESC").
		Limit(limit).
		Find(&loans).Error
	if err != nil {
		return nil, ErrInternal
	}
	return loans, nil
}

func (lp *LoanRepo) UpcomingLoans(limit int) ([]models.Loan, error) {
	var loans []models.Loan
	err := lp.DB.
		Model(&models.Loan{}).
		Preload("Member").
		Preload("Book").
		Where("return_date is NULL").
		Where("status = 'borrowed'").
		Where("due_date >= CURRENT_TIMESTAMP").
		Where("due_date <= CURRENT_DATE + INTERVAL '5 days' ").
		Order("due_date DESC").
		Limit(limit).
		Find(&loans).Error
	if err != nil {
		return nil, ErrInternal
	}
	return loans, nil
}

func (lp *LoanRepo) RecentLoans(limit int) ([]models.Loan, error) {
	var loans []models.Loan
	err := lp.DB.Model(&models.Loan{}).
		Preload("Member").
		Preload("Book").
		Order("borrow_date DESC").
		Limit(limit).
		Find(&loans).Error
	if err != nil {
		return nil, ErrInternal
	}
	return loans, nil
}

func (lp *LoanRepo) TotalWhereStatus(is string) (int64, error) {
	var total int64
	err := lp.DB.Model(&models.Loan{}).Where("status = ?", is).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (lp *LoanRepo) Total() (int64, error) {
	var total int64
	err := lp.DB.Model(&models.Loan{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (lp *LoanRepo) All(data *models.SearchData) ([]models.Loan, error) {
	var loans []models.Loan
	query := lp.DB.Model(&models.Loan{}).Count(&data.Pagination.Total)
	result := query.Preload("Book").Preload("Member").Limit(data.Pagination.Limit).Offset(data.Pagination.Offset).Find(&loans)
	if result.Error != nil {
		return nil, ErrInternal
	}

	data.Pagination.CalculateTotalPage()

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

func (lr *LoanRepo) Search(data *models.SearchData) ([]models.Loan, error) {
	var loans []models.Loan
	s := "%" + data.Term + "%"
	query := lr.DB.Model(&models.Loan{}).Preload("Member").Preload("Book")

	query.Where("CAST(id as TEXT) ILIKE ?", s)

	query.Count(&data.Pagination.Total)
	if data.SortBy != "" {
		query.Order(fmt.Sprintf("%s %s", data.SortBy, data.Dir))
	}
	err := query.Limit(data.Pagination.Limit).Offset(data.Pagination.Offset).Find(&loans).Error

	if err != nil {
		return nil, ErrInternal
	}

	data.Pagination.CalculateTotalPage()

	return loans, nil

}
