package repositories

import (
	"errors"
	"fmt"
	"lms/internal/models"
	"lms/internal/views"
	"log"
	"time"

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
	result := sr.DB.Where("role != 'admin'").
		Where("id = ?", id).
		Delete(&models.Staff{})
	if result.Error != nil {
		return ErrInternal
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (sr *StaffRepo) All(pagination *models.Pagination) ([]models.Staff, error) {
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

func (sr *StaffRepo) GetByRole(role string) ([]models.Staff, error) {
	var staff []models.Staff
	err := sr.DB.Model(&models.Staff{}).Where("role = ?", role).Find(&staff).Error
	if err != nil {
		log.Println(err.Error())
		return nil, ErrInternal
	}
	return staff, nil
}

func (sr *StaffRepo) RecordLogin(id uint) error {
	result := sr.DB.Model(&models.Staff{}).
		Where("id = ?", id).
		Update("last_login", time.Now())
	if result.Error != nil {
		return ErrInternal
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (sr *StaffRepo) Total() (int64, error) {
	var total int64
	err := sr.DB.Model(&models.Staff{}).Count(&total).Error
	if err != nil {
		log.Println(err.Error())
		return 0, ErrInternal
	}
	return total, nil
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

func (sr *StaffRepo) Search(
	data *models.SearchData,
) ([]models.Staff, error) {
	var staff []models.Staff
	query := sr.DB.Model(&models.Staff{})

	s := "%" + data.Term + "%"

	query.Where("full_name ILIKE ?", s).
		Or("CAST(id as TEXT) ILIKE ?", s).
		Or("username ILIKE ?", s).
		Or("phone_number ILIKE ?", s).
		Or("email ILIKE ?", s).
		Or("role ILIKE ?", s).
		Or("status ILIKE ?", s)

	query.Count(&data.Pagination.Total)

	if data.SortBy != "" {
		query.Order(fmt.Sprintf("%s %s", data.SortBy, data.Dir))
	}

	result := query.
		Offset(data.Pagination.Offset).
		Limit(data.Pagination.Limit).
		Find(&staff)

	if result.Error != nil {
		return nil, ErrInternal
	}

	data.Pagination.CalculateTotalPage()

	return staff, nil
}
