package repositories

import (
	"errors"
	"lms/internal/models"
	"lms/internal/views"

	"gorm.io/gorm"
)

type MemberRepo struct {
	DB *gorm.DB
}

func (mr *MemberRepo) GetById(id uint) (*models.Member, error) {
	var member models.Member
	result := mr.DB.Model(&models.Member{}).Preload("Loans").Find(&member, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &member, nil
}

func (mr *MemberRepo) Insert(member *models.Member) error {
	result := mr.DB.Create(member)
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	}
	return nil
}

func (mr *MemberRepo) DeleteById(id string) error {
	result := mr.DB.Delete(&models.Member{}, id)
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

func (mr *MemberRepo) All(page int, pageSize int) ([]models.Member, error) {
	var members []models.Member
	offset := (page - 1) * pageSize
	if err := mr.DB.Limit(pageSize).Offset(offset).Find(&members).Error; err != nil {
		return nil, ErrInternal
	}
	return members, nil
}

func (mr *MemberRepo) Update(member *models.Member) error {
	result := mr.DB.Model(&models.Member{}).Where("id = ?", member.ID).Updates(member)
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

func (mr *MemberRepo) HasActiveLoans(memberId uint) (bool, error) {
	var total int64
	err := mr.DB.Model(&models.Loan{}).Where("member_id = ?", memberId).Count(&total).Error
	if err != nil {
		return false, ErrInternal
	}
	return total > 0, nil
}

func (mr *MemberRepo) ConvertErrorsToFormErrors(err error) views.Errors {
	errors := make(views.Errors)
	pgErr := extractPQError(err)
	switch pgErr.ConstraintName {
	case "members_phone_validation_check":
		errors["phoneNumber"] = "invalid phone number"
	case "members_phone_key":
		errors["phoneNumber"] = "a member with this phone number already exists"
	case "members_email_key":
		errors["email"] = "a member with this email already exists"
	case "members_email_validation_check":
		errors["email"] = "invalid email"
	case "members_status_check":
		errors["status"] = "invalid status"
	case "members_national_id_key":
		errors["nationalId"] = "a memeber with this national id alreadyexists"
	default:
		errors["_"] = err.Error()
	}
	return errors
}
