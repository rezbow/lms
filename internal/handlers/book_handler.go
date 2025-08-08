package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	bookViews "lms/internal/views/books"
	loanViews "lms/internal/views/loans"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookHandler struct {
	BookRepo     *repositories.BookRepo
	CategoryRepo *repositories.CategoryRepo
	LogRepo      *repositories.ActivityRepo
	Validator    *validator.Validate
}

func (bh *BookHandler) Get(ctx *gin.Context) {

	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, bookViews.Book(book), book.Title)
}

func (bh *BookHandler) Delete(ctx *gin.Context) {

	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	err = bh.BookRepo.DeleteById(bookId)
	if err != nil {
		switch err {
		case repositories.ErrNotFound:
			notfound(ctx)
			return
		case repositories.ErrInternal:
			serverError(ctx)
			return
		default:
			redirect(ctx, fmt.Sprintf("/books/%d", bookId))
			return

		}
	}

	err = LogStaffActivity(
		bh.LogRepo,
		ctx,
		models.ActivityTypeDeleteBook,
		bookId,
		models.EntityTypeBook,
		"deleted Book",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, "/books")
}

func (bh *BookHandler) AddLoanPage(ctx *gin.Context) {
	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	render(ctx, loanViews.LoanForm(&models.Loan{BookId: book.ID}, views.Errors{}, "/loans/add"), "add loan")

}

func (bh *BookHandler) AddPage(ctx *gin.Context) {
	categories, err := bh.CategoryRepo.All()
	if err != nil {
		serverError(ctx)
		return
	}
	render(ctx, bookViews.BookForm(nil, views.Errors{}, "/books/add", categories), "add book")
}

func (bh *BookHandler) Add(ctx *gin.Context) {
	var bookForm struct {
		Title       string `form:"title" binding:"required" validate:"required,min=1,max=100"`
		ISBN        string `form:"isbn" binding:"required" validate:"required,max=20"`
		TotalCopies int    `form:"totalCopies" binding:"required" validate:"required,min=1"`
		AuthorId    uint   `form:"authorId" binding:"required" validate:"required,min=1"`
		// optional
		Publisher  *string `form:"publisher" binding:"omitempty" validate:"omitempty,max=30"`
		Language   *string `form:"language" binding:"omitempty" validate:"omitempty,max=30"`
		Summary    *string `form:"summary" binding:"omitempty" validate:"omitempty,max=1000"`
		Translator *string `form:"translator" binding:"omitempty" validate:"omitempty,max=50"`
		Categories []uint  `form:"categories" binding:"required"`
	}

	if err := ctx.ShouldBind(&bookForm); err != nil {
		render(ctx, bookViews.BookForm(nil, parseValidationErrors(err), "/books/add", nil), "add book")
		return
	}

	if err := bh.Validator.Struct(bookForm); err != nil {
		render(ctx, bookViews.BookForm(nil, parseValidationErrors(err), "/books/add", nil), "add book")
		return
	}

	book := models.Book{
		Title:           bookForm.Title,
		ISBN:            bookForm.ISBN,
		TotalCopies:     bookForm.TotalCopies,
		AvailableCopies: bookForm.TotalCopies,
		AuthorId:        bookForm.AuthorId,
	}

	if bookForm.Publisher != nil {
		book.Publisher = *bookForm.Publisher
	}

	if bookForm.Summary != nil {
		book.Summary = *bookForm.Summary
	}

	if bookForm.Translator != nil {
		book.Translator = *bookForm.Translator
	}

	if bookForm.Language != nil {
		book.Language = *bookForm.Language
	}

	if bookForm.Categories != nil {
		for _, id := range bookForm.Categories {
			book.Categories = append(book.Categories, &models.Category{ID: id})
		}
	}

	if err := bh.BookRepo.Insert(&book); err != nil {
		errors := bh.BookRepo.ConvertErrorsToFormErrors(err)
		if errors["_"] != "" {
			serverError(ctx)
			return
		}
		render(ctx, bookViews.BookForm(nil, errors, "/books/add", nil), book.Title)
		return
	}

	err := LogStaffActivity(
		bh.LogRepo,
		ctx,
		models.ActivityTypeAddBook,
		book.ID,
		models.EntityTypeBook,
		"added Book",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/books/%d", book.ID))
}

func (bh *BookHandler) Search(ctx *gin.Context) {

	term := ctx.Query("q")

	pagination, err := readPagination(ctx, "/books/search?q="+term)
	if err != nil {
		notfound(ctx)
		return
	}

	books, err := bh.BookRepo.Search(
		term,
		pagination,
	)

	if err != nil {
		serverError(ctx)
		return
	}

	data := views.SearchData{
		Term:       term,
		BaseUrl:    "/books/search",
		Pagination: pagination,
		Sort:       "",
		Direction:  "",
	}

	render(ctx, bookViews.BookSearch(books, &data), "search results")
}

func (bh *BookHandler) Index(ctx *gin.Context) {
	pagination, err := readPagination(ctx, "/books/?")

	books, err := bh.BookRepo.All(pagination)
	if err != nil {
		serverError(ctx)
		return
	}

	data := views.SearchData{
		Term:       "",
		BaseUrl:    "/books/search",
		Pagination: pagination,
		Sort:       "",
		Direction:  "",
	}

	render(ctx, bookViews.BookSearch(books, &data), "books")

}

func (bh *BookHandler) EditPage(ctx *gin.Context) {
	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, bookViews.BookForm(book, views.Errors{}, fmt.Sprintf("/books/%d/edit", book.ID), nil), book.Title)
}

func (bh *BookHandler) Update(ctx *gin.Context) {

	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	var bookUpdateForm struct {
		Title       *string `form:"title" binding:"omitempty" validate:"omitempty,min=1,max=100"`
		ISBN        *string `form:"isbn" binding:"omitempty" validate:"omitempty,min=1"`
		TotalCopies *int    `form:"totalCopies" binding:"omitempty" validate:"omitempty,min=1"`
		AuthorId    *uint   `form:"authorId" binding:"omitempty" validate:"omitempty,min=1"`
		// optional
		Publisher  *string `form:"publisher" binding:"omitempty" validate:"omitempty,max=30"`
		Language   *string `form:"language" binding:"omitempty" validate:"omitempty,max=30"`
		Summary    *string `form:"summary" binding:"omitempty" validate:"omitempty,max=1000"`
		Translator *string `form:"translator" binding:"omitempty" validate:"omitempty,max=50"`
	}

	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	endpoint := fmt.Sprintf("/books/%d/edit", bookId)

	if err := ctx.ShouldBind(&bookUpdateForm); err != nil {
		render(ctx, bookViews.BookForm(book, parseValidationErrors(err), endpoint, nil), book.Title)
		return
	}

	if err := bh.Validator.Struct(&bookUpdateForm); err != nil {
		render(ctx, bookViews.BookForm(book, parseValidationErrors(err), endpoint, nil), book.Title)
		return
	}

	if bookUpdateForm.Title != nil {
		book.Title = *bookUpdateForm.Title
	}

	if bookUpdateForm.ISBN != nil {
		book.ISBN = *bookUpdateForm.ISBN
	}

	if bookUpdateForm.TotalCopies != nil {
		book.TotalCopies = *bookUpdateForm.TotalCopies
	}

	if bookUpdateForm.TotalCopies != nil {
		book.AuthorId = *bookUpdateForm.AuthorId
	}

	if bookUpdateForm.Publisher != nil {
		book.Publisher = *bookUpdateForm.Publisher
	}

	if bookUpdateForm.Summary != nil {
		book.Summary = *bookUpdateForm.Summary
	}

	if bookUpdateForm.Translator != nil {
		book.Translator = *bookUpdateForm.Translator
	}

	if bookUpdateForm.Language != nil {
		book.Language = *bookUpdateForm.Language
	}

	if err := bh.BookRepo.Update(book); err != nil {
		errors := bh.BookRepo.ConvertErrorsToFormErrors(err)
		if errors["_"] != "" {
			serverError(ctx)
			return
		}
		render(ctx, bookViews.BookForm(book, errors, endpoint, nil), book.Title)
		return
	}

	err = LogStaffActivity(
		bh.LogRepo,
		ctx,
		models.ActivityTypeUpdateBook,
		book.ID,
		models.EntityTypeBook,
		"updated book",
	)

	if err != nil {
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/books/%d", bookId))
}
