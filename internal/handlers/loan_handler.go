package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	commonViews "lms/internal/views/common"
	loanViews "lms/internal/views/loans"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoanHandler struct {
	Repo      *repositories.LoanRepo
	LogRepo   *repositories.ActivityRepo
	Validator *validator.Validate
}

func (lh *LoanHandler) Index(ctx *gin.Context) {
	totalLoans, err := lh.Repo.Total()
	totalActiveLoan, err := lh.Repo.TotalWhereStatus("borrowed")
	totalReturnedLoan, err := lh.Repo.TotalWhereStatus("returned")
	totalOverdueLoans, err := lh.Repo.TotalOverdueLoans()
	recentLoans, err := lh.Repo.RecentLoans(5)
	overdueLoans, err := lh.Repo.OverdueLoans(5)
	upcomingLoans, err := lh.Repo.UpcomingLoans(5)
	if err != nil {
		serverError(ctx)
		return
	}
	data := models.LoanDashboard{
		TotalLoans:         totalLoans,
		TotalActiveLoans:   totalActiveLoan,
		TotalReturnedLoans: totalReturnedLoan,
		TotalOverdueLoans:  totalOverdueLoans,
		RecentLoans:        recentLoans,
		OverdueLoans:       overdueLoans,
		UpcomingLoan:       upcomingLoans,
	}

	render(ctx, loanViews.LoanDashboard(&data), "members")
}

func (lh *LoanHandler) GetById(ctx *gin.Context) {
	loanId, err := readID(ctx)
	if err != nil {
		render(ctx, commonViews.NotFound(), "404:((")
		return
	}
	loan, err := lh.Repo.GetById(loanId)
	if err != nil {
		if err == repositories.ErrNotFound {
			render(ctx, commonViews.NotFound(), "404 :((")
			return
		}
		render(ctx, commonViews.ServerError(err.Error()), "internal server error")
		return
	}
	render(ctx, loanViews.LoanInfo(loan), "loan")
}

func (lh *LoanHandler) DeleteById(ctx *gin.Context) {
	loanId, err := readID(ctx)
	if err != nil {
		render(ctx, commonViews.NotFound(), "404:((")
		return
	}
	if err := lh.Repo.DeleteById(loanId); err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
		} else if err == repositories.ErrInternal {
			serverError(ctx)
		} else {
			redirect(ctx, fmt.Sprintf("/loans/%d", loanId))
		}
		return
	}

	err = LogStaffActivity(
		lh.LogRepo,
		ctx,
		models.ActivityTypeDeleteLoan,
		loanId,
		models.EntityTypeLoan,
		"deleted loan",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, "/loans")
}

func (lh *LoanHandler) AddPage(ctx *gin.Context) {
	render(ctx, loanViews.LoanForm(nil, views.Errors{}, "/loans/add"), "add loan")
}

func (lh *LoanHandler) Add(ctx *gin.Context) {
	var userInput struct {
		BookId     uint      `form:"bookId" binding:"required" validate:"required,min=1"`
		MemberId   uint      `form:"memberId" binding:"required" validate:"required,min=1"`
		BorrowDate time.Time `form:"borrowDate" binding:"required" time_format:"2006-01-02" validate:"required"`
		DueDate    time.Time `form:"dueDate" binding:"required" time_format:"2006-01-02" validate:"required"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
		return
	}

	if err := lh.Validator.Struct(&userInput); err != nil {
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
		return
	}

	loan := models.Loan{
		BookId:     userInput.BookId,
		MemberId:   userInput.MemberId,
		BorrowDate: userInput.BorrowDate,
		DueDate:    userInput.DueDate,
	}

	if err := lh.Repo.Insert(&loan); err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
			return
		}
		render(ctx, loanViews.LoanForm(nil, lh.Repo.ConvertErrorToFormError(err), "/loans/add"), "add loan")
		return
	}

	err := LogStaffActivity(
		lh.LogRepo,
		ctx,
		models.ActivityTypeAddLoan,
		loan.ID,
		models.EntityTypeLoan,
		"added loan",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/loans/%d", loan.ID))
}

func (lh *LoanHandler) EditPage(ctx *gin.Context) {
	loanId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	loan, err := lh.Repo.GetById(loanId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	render(ctx, loanViews.LoanForm(loan, views.Errors{}, fmt.Sprintf("/loans/%d/edit", loan.ID)), "add loan")
}

func (lh *LoanHandler) Update(ctx *gin.Context) {

	loanId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	var userInput struct {
		BookId     uint      `form:"bookId" binding:"omitempty" validate:"omitempty,min=1"`
		MemberId   uint      `form:"memberId" binding:"omitempty" validate:"omitempty,min=1"`
		BorrowDate time.Time `form:"borrowDate" binding:"omitempty" time_format:"2006-01-02"`
		DueDate    time.Time `form:"dueDate" binding:"omitempty" time_format:"2006-01-02" `
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		formError(ctx, err)
		return
	}

	if err := lh.Validator.Struct(&userInput); err != nil {
		formError(ctx, err)
		return
	}

	updatedLoan := models.Loan{
		ID:         loanId,
		BookId:     userInput.BookId,
		MemberId:   userInput.MemberId,
		BorrowDate: userInput.BorrowDate,
		DueDate:    userInput.DueDate,
	}

	if err := lh.Repo.Update(&updatedLoan); err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
			return
		} else if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		formError(ctx, err)
		return
	}

	err = LogStaffActivity(
		lh.LogRepo,
		ctx,
		models.ActivityTypeUpdateLoan,
		loanId,
		models.EntityTypeLoan,
		"updated loan",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/loans/%d", loanId))

}

func (lh *LoanHandler) ReturnPage(ctx *gin.Context) {
	loanId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	loan, err := lh.Repo.GetById(loanId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	render(ctx, loanViews.LoanReturnForm(loan), "add loan")
}

func (lh *LoanHandler) ReturnLoan(ctx *gin.Context) {

	loanId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	var userInput struct {
		ReturnDate time.Time `form:"returnDate" binding:"required" time_format:"2006-01-2"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		formError(ctx, err)
		return
	}

	err = lh.Repo.Update(&models.Loan{ID: loanId, ReturnDate: &userInput.ReturnDate, Status: models.StatusReturned})

	if err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
			return
		} else if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		formError(ctx, err)
		return
	}

	err = LogStaffActivity(
		lh.LogRepo,
		ctx,
		models.ActivityTypeReturnLoan,
		loanId,
		models.EntityTypeLoan,
		"marked loan as returned",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/loans/%d", loanId))
}

func (lh *LoanHandler) Search(ctx *gin.Context) {
	searchData, err := readSearchData(ctx, "/loans/search")
	if err != nil {
		notfound(ctx)
		return
	}

	if !searchData.Valid(models.LoanSafeSortList) {
		notfound(ctx)
		return
	}

	loans, err := lh.Repo.Search(searchData)

	if err != nil {
		serverError(ctx)
		return
	}

	render(ctx, loanViews.LoanSearch(loans, searchData), "loans")
}
