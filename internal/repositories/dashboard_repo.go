package repositories

import (
	"lms/internal/models"

	"gorm.io/gorm"
)

type DashboardRepo struct {
	DB *gorm.DB
}

func (dp *DashboardRepo) BookCount() (int64, error) {
	var total int64
	err := dp.DB.Model(&models.Book{}).Count(&total).Error
	if err != nil {
		return 0, ErrInternal
	}
	return total, nil
}

func (dp *DashboardRepo) MemberCount() (int64, error) {
	var total int64
	err := dp.DB.Model(&models.Member{}).Count(&total).Error
	if err != nil {
		return 0, ErrInternal
	}
	return total, nil

}

func (dp *DashboardRepo) ActiveLoanCount() (int64, error) {
	var total int64
	err := dp.DB.Model(&models.Loan{}).Where("status = 'borrowed'").Count(&total).Error
	if err != nil {
		return 0, ErrInternal
	}
	return total, nil
}

func (dr *DashboardRepo) ActiveMembersCount() (int64, error) {
	var total int64
	err := dr.DB.Table("active_members_view").Count(&total).Error
	if err != nil {
		return 0, ErrInternal
	}
	return total, nil
}

func (dp *DashboardRepo) OverdueLoanCount() (int64, error) {
	var total int64
	err := dp.DB.Model(&models.Loan{}).Where("status = 'overdue'").Count(&total).Error
	if err != nil {
		return 0, ErrInternal
	}
	return total, nil
}

func (dp *DashboardRepo) MostPopularBooks() ([]models.PopularBook, error) {
	var books []models.PopularBook
	err := dp.DB.Table("popular_books_view").Scan(&books).Error
	if err != nil {
		return nil, ErrInternal
	}
	return books, nil
}

func (dr *DashboardRepo) ActiveMembers() ([]models.ActiveMember, error) {
	var members []models.ActiveMember
	err := dr.DB.Table("active_members_view").Scan(&members).Error
	if err != nil {
		return nil, ErrInternal
	}
	return members, nil
}

func (dp *DashboardRepo) PopularCategories() ([]models.PopularCategory, error) {
	var categories []models.PopularCategory
	err := dp.DB.Table("popular_categories_view").Scan(&categories).Error
	if err != nil {
		return nil, ErrInternal
	}
	return categories, nil
}

func (dp *DashboardRepo) UpcomingLoans(limit int) ([]models.Loan, error) {
	var loans []models.Loan
	err := dp.DB.
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

/*
func (dp *DashboardRepo) Search(term string, pagination *utils.Pagination) ([]any, error) {
	result := models.DashboardSearchResult{}
	s := "%" + term + "%"
	tx := dp.DB.Model(&models.Book{}).
		Where("title ILIKE ?", term).
		Count(&pagination.Total)

}
*/
