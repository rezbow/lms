package repositories

import (
	"errors"
	"lms/internal/models"
	"lms/internal/utils"
	"lms/internal/views"

	"gorm.io/gorm"
)

type StaffRepo struct {
	DB *gorm.DB
}

func (sr *StaffRepo) GetById(id uint) (*models.Staff, error) {
	var staff models.Staff
	result := sr.DB.First(&staff, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &staff, nil
}

func (sr *StaffRepo) GetByEmail(email string) (*models.Staff, error) {
	var staff models.Staff
	result := sr.DB.Where("email = ?", email).Find(&staff)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return &staff, nil
}

func (sr *StaffRepo) Insert(staff *models.Staff) error {
	result := sr.DB.Create(staff)
	if result.Error != nil {
		if isInternalError(result.Error) {
			return ErrInternal
		}
		return result.Error
	}
	return nil
}

func (sr *StaffRepo) DeleteById(id uint) error {
	result := sr.DB.Delete(&models.Staff{}, id)
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

func (sr *StaffRepo) All(pagination *utils.Pagination) ([]models.Staff, error) {
	var staff []models.Staff
	if err := sr.DB.Limit(pagination.Limit).Offset(pagination.Offset).Find(&staff).Error; err != nil {
		return nil, ErrInternal
	}
	return staff, nil
}

func (sr *StaffRepo) Update(staff *models.Staff) error {
	result := sr.DB.Model(&models.Staff{}).Where("id = ?", staff.ID).Updates(staff)
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

func (sr *StaffRepo) Filter(
	filter *models.StaffFilter,
	pagination *utils.Pagination,
	total *int64,
) ([]models.Staff, error) {
	var staff []models.Staff
	query := sr.DB.Model(&models.Staff{})

	if filter.FullName != "" {
		query.Where("full_name ILIKE ?", "%"+filter.FullName+"%")
	}

	if filter.Username != "" {
		query.Where("username ILIKE ?", "%"+filter.Username+"%")
	}

	if filter.Role != "" {
		query.Where("role = ?", filter.Role)
	}

	if filter.Status != "" {
		query.Where("status = ?", filter.Status)
	}

	query.Count(total)

	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&staff).Error; err != nil {
		return nil, ErrInternal
	}
	return staff, nil
}

func (sr *StaffRepo) Total() int64 {
	var total int64
	sr.DB.Model(&models.Staff{}).Count(&total)
	return total
}

func (sr *StaffRepo) ConvertErrorsToFormErrors(err error) views.Errors {
	errors := make(views.Errors)
	pgErr := extractPQError(err)
	switch pgErr.ConstraintName {
	case "staff_phone_number_key":
		errors["phoneNumber"] = "a staff with this phone number already exists"
	case "staff_email_key":
		errors["email"] = "a staff with this email already exists"
	case "staff_role_check":
		errors["role"] = "invalid role"
	case "staff_username_unique":
		errors["username"] = "username already exists"
	case "phone_number_check":
		errors["phoneNumber"] = "invalid phone number"
	case "email_check":
		errors["email"] = "invalid email"
	default:
		errors["_"] = err.Error()
	}
	return errors
}
