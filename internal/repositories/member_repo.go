package repositories

import (
	"errors"
	"lms/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type MemberRepo struct {
	DB *gorm.DB
}

func (mr *MemberRepo) GetById(id string) (*models.Member, error) {
	var member models.Member
	result := mr.DB.First(&member, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &member, nil
}

var (
	ErrMemberInvalidPhoneNumber   = errors.New("Invalid phone number")
	ErrMemberDuplicatePhoneNumber = errors.New("Duplicate phone number")
	ErrMemberDuplicateEmail       = errors.New("Duplicate email")
	ErrMemberInvalidStatus        = errors.New("Invalid Status")
	ErrMemberInvalidEmail         = errors.New("Invalid Email")
)

func (mr *MemberRepo) constraintViolationError(err *pgconn.PgError) error {
	switch err.ConstraintName {
	case "members_phone_validation_check":
		return ErrMemberInvalidPhoneNumber
	case "members_phone_key":
		return ErrMemberDuplicatePhoneNumber
	case "members_email_key":
		return ErrMemberDuplicateEmail
	case "members_status_check":
		return ErrMemberInvalidStatus
	case "members_email_validation_check":
		return ErrMemberInvalidEmail
	default:
		return errors.New(err.Message)
	}
}

func (mr *MemberRepo) Insert(member *models.Member) error {
	if err := mr.DB.Create(member).Error; err != nil {
		pgerr := extractPQError(err)
		if pgerrcode.IsIntegrityConstraintViolation(pgerr.Code) {
			return mr.constraintViolationError(pgerr)
		}
		return ErrInternal
	}
	return nil
}

func (mr *MemberRepo) DeleteById(id string) error {
	result := mr.DB.Delete(&models.Member{}, id)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (mr *MemberRepo) All(page int, pageSize int) ([]models.Member, error) {
	var members []models.Member
	offset := (page - 1) * pageSize
	if err := mr.DB.Limit(pageSize).Offset(offset).Find(&members).Error; err != nil {
		return nil, ErrInternal
	}
	return members, nil
}

func (mr *MemberRepo) Update(member *models.Member) error {
	if err := mr.DB.Model(&models.Member{}).Where("id = ?", member.ID).Updates(member).Error; err != nil {
		pgErr := extractPQError(err)
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return mr.constraintViolationError(pgErr)
		}
		return ErrInternal
	}
	return nil
}

func (mr *MemberRepo) Filter(
	name string,
	phone string,
	email string,
	limit int,
	page int,
	total *int64,
) ([]models.Member, error) {
	var members []models.Member
	query := mr.DB.Model(&models.Member{})

	if name != "" {
		query.Where("members.name ILIKE ?", "%"+name+"%")
	}

	if phone != "" {
		query.Where("members.phone ILIKE ?", "%"+phone+"%")
	}

	if email != "" {
		query.Where("members.phone ILIKE ?", "%"+email+"%")
	}

	query.Count(total)
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&members).Error; err != nil {
		return nil, ErrInternal
	}
	return members, nil
}

func (mr *MemberRepo) Total() int64 {
	var total int64
	mr.DB.Model(&models.Book{}).Count(&total)
	return total
}
