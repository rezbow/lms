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
	Validator *validator.Validate
}

func (lh *LoanHandler) Index(ctx *gin.Context) {
	pagination, err := readPagination(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	loans, err := lh.Repo.All(pagination)
	if err != nil {
		render(ctx, commonViews.ServerError(err.Error()), "server error")
		return
	}
	render(ctx, loanViews.LoanSearch(loans), "members")
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

	fmt.Println(updatedLoan)

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

	redirect(ctx, fmt.Sprintf("/loans/%d", loanId))
}

func (lh *LoanHandler) Search(ctx *gin.Context) {
	status := ctx.Query("status")

	bookId, err := readIntFromQuery(ctx.Query("book"))
	if err != nil {
		notfound(ctx)
		return
	}

	memberId, err := readIntFromQuery(ctx.Query("member"))
	if err != nil {
		notfound(ctx)
		return
	}
	pagination, err := readPagination(ctx)
	if err != nil {
		fmt.Println(pagination, err.Error())
		notfound(ctx)
		return
	}

	loans, err := lh.Repo.Filter(
		&models.LoanFilter{
			BookId:   uint(bookId),
			MemberId: uint(memberId),
			Status:   status,
		},
		pagination,
	)

	if err != nil {
		serverError(ctx)
		return
	}

	fmt.Println(loans)

	render(ctx, loanViews.LoanSearch(loans), "loans")
}
