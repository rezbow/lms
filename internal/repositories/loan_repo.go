package repositories

import (
	"errors"
	"lms/internal/models"
	"strings"

	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

type LoanRepo struct {
	DB *gorm.DB
}

var (
	ErrLoanInvalidStatus = errors.New("status should be borrowed or returned")
)

func (lp *LoanRepo) GetById(id int) (*models.Loan, error) {
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

func (lp *LoanRepo) DeleteById(id int) error {
	result := lp.DB.Delete(&models.Loan{}, id)
	if result.Error != nil {
		return ErrInternal
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (lp *LoanRepo) Insert(loan *models.Loan) error {
	result := lp.DB.Create(loan)
	if result.Error != nil {
		return lp.handleInsertError(result.Error)
	}
	return nil
}

func (lp *LoanRepo) All(page, pageSize int) ([]models.Loan, error) {
	var loans []models.Loan
	offset := (page - 1) * pageSize
	result := lp.DB.Preload("Book").Preload("Member").Limit(pageSize).Offset(offset).Find(&loans)
	if result.Error != nil {
		return nil, ErrInternal
	}
	return loans, nil
}

func (lr *LoanRepo) Update(loan *models.Loan) error {
	result := lr.DB.Model(&models.Loan{}).Where("id = ?", loan.ID).Updates(loan)
	if result.Error != nil {
		return lr.handleInsertError(result.Error)
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (lr *LoanRepo) handleInsertError(err error) error {
	pgError := extractPQError(err)
	if pgError.Code == pgerrcode.ForeignKeyViolation {
		if idx := strings.Index(pgError.ConstraintName, "book"); idx >= 0 {
			return ErrBookIdNotFound
		} else {
			return ErrMemberIdNotFound
		}
	} else if pgError.Code == pgerrcode.CheckViolation {
		return ErrLoanInvalidStatus
	}
	return ErrInternal
}
