package internal

import (
	"lms/internal/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	bookHandler *handlers.BookHandler,
	authorHandler *handlers.AuthorHandler,
	memberHandler *handlers.MemberHandler,
	loanHandler *handlers.LoanHandler,
	categoryHandler *handlers.CategoryHandler,
	staffHandler *handlers.StaffHandler,
	dashboradHandler *handlers.DashboardHanlder,
) *gin.Engine {
	r := gin.Default()

	//r.HTMLRender = &TemplRender{}

	store := cookie.NewStore([]byte("cfaa7e52"))
	r.Use(sessions.Sessions("session", store))
	r.GET("/login", staffHandler.LoginPage)
	r.POST("/login", staffHandler.Login)

	r.Use(AuthRequired())
	r.GET("/", dashboradHandler.Dashboard)
	r.POST("/logout", staffHandler.Logout)

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
	books.GET("/:id/addloan", bookHandler.AddLoanPage)

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
	loans.GET("/search", loanHandler.Search)

	loans.GET("/:id", loanHandler.GetById)
	loans.POST("/:id/delete", loanHandler.DeleteById)
	loans.GET("/:id/edit", loanHandler.EditPage)
	loans.POST("/:id/edit", loanHandler.Update)
	loans.POST("/:id/return", loanHandler.ReturnLoan)
	loans.GET("/:id/return", loanHandler.ReturnPage)

	// categories
	categories := r.Group("/categories")
	categories.GET("/add", categoryHandler.AddPage)
	categories.POST("/add", categoryHandler.Add)

	categories.GET("/:slug", categoryHandler.Get)
	categories.GET("/:slug/edit", categoryHandler.EditPage)
	categories.POST("/:slug/edit", categoryHandler.Edit)
	categories.POST("/:slug/delete", categoryHandler.Delete)

	// staff
	staff := r.Group("/staff", AdminRequired())
	staff.GET("", staffHandler.Index)
	staff.GET("/add", staffHandler.AddPage)
	staff.POST("/add", staffHandler.Add)
	staff.GET("/search", staffHandler.Search)

	staff.GET("/:id", staffHandler.Get)
	staff.POST("/:id/delete", staffHandler.Delete)
	staff.GET("/:id/edit", staffHandler.EditPage)
	staff.POST("/:id/edit", staffHandler.Edit)

	// staff
	authors := r.Group("/authors")
	authors.GET("", authorHandler.Index)
	authors.GET("/add", authorHandler.AddPage)
	authors.POST("/add", authorHandler.Add)
	authors.GET("/search", authorHandler.Search)

	authors.GET("/:id", authorHandler.Get)
	authors.POST("/:id/delete", authorHandler.Delete)
	authors.GET("/:id/edit", authorHandler.EditPage)
	authors.POST("/:id/edit", authorHandler.Edit)
	return r
}
