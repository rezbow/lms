package internal

import (
	"lms/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	bookHandler *handlers.BookHandler,
	authorHandler *handlers.AuthorHandler,
	memberHandler *handlers.MemberHandler,
	loanHandler *handlers.LoanHandler,
) *gin.Engine {
	r := gin.Default()

	r.HTMLRender = &TemplRender{}

	r.GET("/", handlers.Get)

	// books resource
	books := r.Group("/books")
	books.GET("", bookHandler.Index)
	books.GET("/search", bookHandler.Search)
	books.GET("/add", bookHandler.AddPage)
	books.POST("/add", bookHandler.Add)

	books.GET("/:id", bookHandler.Get)
	books.POST("/:id/delete", bookHandler.Delete)
	books.GET("/:id/edit", bookHandler.EditPage)
	books.POST("/:id/edit", bookHandler.Update)
	// books.GET("/:id/addloan", bookHandler.AddLoanPage)

	// memebers
	members := r.Group("/members")
	members.GET("", memberHandler.Index)
	members.GET("/add", memberHandler.AddPage)
	members.POST("/add", memberHandler.Add)
	members.GET("/search", memberHandler.Search)

	members.GET("/:id", memberHandler.GetById)
	members.POST("/:id/delete", memberHandler.DeleteById)
	members.GET("/:id/edit", memberHandler.EditPage)
	members.POST("/:id/edit", memberHandler.Update)
	members.GET("/:id/addloan", memberHandler.AddLoanPage)

	// loans
	loans := r.Group("/loans")
	loans.GET("", loanHandler.Index)
	loans.GET("/add", loanHandler.AddPage)
	loans.POST("/add", loanHandler.Add)
	loans.GET("/search")

	loans.GET("/:id", loanHandler.GetById)
	loans.POST("/:id/delete", loanHandler.DeleteById)
	loans.GET("/:id/edit", loanHandler.EditPage)
	loans.POST("/:id/edit", loanHandler.Update)
	loans.POST("/:id/return", loanHandler.ReturnLoan)
	loans.GET("/:id/return", loanHandler.ReturnPage)
	return r
}
