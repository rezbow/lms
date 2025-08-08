package handlers

import (
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type DashboardHanlder struct {
	Repo    *repositories.DashboardRepo
	LogRepo *repositories.ActivityRepo
}

func (dh *DashboardHanlder) Dashboard(ctx *gin.Context) {
	session := sessions.Default(ctx)

	totalBooks, err := dh.Repo.BookCount()
	totalMembers, err := dh.Repo.MemberCount()
	totatlActiveLoans, err := dh.Repo.ActiveLoanCount()
	totalOverdueLoans, err := dh.Repo.OverdueLoanCount()
	popularBooks, err := dh.Repo.MostPopularBooks()
	activeMembers, err := dh.Repo.ActiveMembers()
	popularCategories, err := dh.Repo.PopularCategories()
	upcomingLoans, err := dh.Repo.UpcomingLoans()
	recentActivities, err := dh.LogRepo.Recent(10)

	if err != nil {
		serverError(ctx)
		return
	}

	data := models.Dashboard{
		TotalBooks:        totalBooks,
		TotalMembers:      totalMembers,
		TotalActiveLoans:  totatlActiveLoans,
		TotalOverdueLoans: totalOverdueLoans,
		PopularBooks:      popularBooks,
		ActiveMembers:     activeMembers,
		PopularCategory:   popularCategories,
		UpcomingLoans:     upcomingLoans,
		RecentActivities:  recentActivities,
	}

	render(ctx, views.Dashboard(&data, session), "dashboard")
}

/*
func (dh *DashboardHanlder) Search(ctx *gin.Context) {
	term := ctx.Query("q")
	pagination, err := readPagination(ctx, "/search?q="+term)
	if err != nil {
		notfound(ctx)
		return
	}
}
*/
