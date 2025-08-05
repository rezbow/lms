package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	authViews "lms/internal/views/auth"
	staffView "lms/internal/views/staff"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StaffHandler struct {
	Repo      *repositories.StaffRepo
	Validator *validator.Validate
}

// get
func (sh *StaffHandler) Index(ctx *gin.Context) {
	pagination, err := readPagination(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	staff, err := sh.Repo.All(pagination)
	if err != nil {
		serverError(ctx)
		return
	}
	render(ctx, staffView.StaffSearch(staff), "staff")
}

// get
func (sh *StaffHandler) AddPage(ctx *gin.Context) {
	render(ctx, staffView.StaffForm(nil, views.Errors{}, "/staff/add"), "add staff")
}

// get
func (sh *StaffHandler) Get(ctx *gin.Context) {
	id, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	staff, err := sh.Repo.GetById(uint(id))
	if err != nil {
		switch err {
		case repositories.ErrInternal:
			serverError(ctx)
		case repositories.ErrNotFound:
			notfound(ctx)
		}
		return
	}

	render(ctx, staffView.Staff(staff), staff.FullName)
}

// post
func (sh *StaffHandler) Delete(ctx *gin.Context) {
	id, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	err = sh.Repo.DeleteById(uint(id))
	if err != nil {
		switch err {
		case repositories.ErrInternal:
			serverError(ctx)
		case repositories.ErrNotFound:
			notfound(ctx)
		default:
			// TODO:
		}
		return
	}
	redirect(ctx, "/staff")
}

// get
func (sh *StaffHandler) EditPage(ctx *gin.Context) {
	id, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	staff, err := sh.Repo.GetById(uint(id))
	if err != nil {
		switch err {
		case repositories.ErrInternal:
			serverError(ctx)
		case repositories.ErrNotFound:
			notfound(ctx)
		}
		return
	}
	render(ctx, staffView.StaffForm(staff, views.Errors{}, fmt.Sprintf("/staff/%d/edit", staff.ID)), staff.FullName)
}

// post
func (sh *StaffHandler) Edit(ctx *gin.Context) {
	var userInput struct {
		FullName    *string `form:"fullName" validate:"omitempty,min=1,max=100"`
		Username    *string `form:"username" validate:"omitempty,min=1,max=200"`
		PhoneNumber *string `form:"phoneNumber" validate:"omitempty,min=1,max=50"`
		Email       *string `form:"email" validate:"omitempty,email,min=1,max=100"`
		Password    *string `form:"password" validate:"omitempty,min=8,max=20"`
		Role        *string `form:"role" validate:"omitempty,oneof=admin librarian,min=1,max=50"`
		Status      *string `form:"status" validate:"omitempty,oneof=active suspended,min=1,max=50"`
	}

	id, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	staff, err := sh.Repo.GetById(id)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	endpoint := fmt.Sprintf("/staff/%d/edit", staff.ID)

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, staffView.StaffForm(
			staff,
			parseValidationErrors(err),
			endpoint,
		), staff.FullName)
		return
	}

	if err := sh.Validator.Struct(&userInput); err != nil {
		render(ctx, staffView.StaffForm(
			staff,
			parseValidationErrors(err),
			endpoint,
		), staff.FullName)
		return
	}

	if userInput.Email != nil {
		staff.Email = *userInput.Email
	}

	if userInput.FullName != nil {
		staff.FullName = *userInput.FullName
	}

	if userInput.PhoneNumber != nil {
		staff.PhoneNumber = *userInput.PhoneNumber
	}

	if userInput.Username != nil {
		staff.Username = *userInput.Username
	}

	if userInput.Password != nil {
		// TODO: hash and save
		hash, err := generateHash(*userInput.Password)
		if err != nil {
			serverError(ctx)
			return
		}
		staff.PasswordHash = hash
	}

	if userInput.Role != nil {
		staff.Role = *userInput.Role
	}

	if userInput.Status != nil {
		staff.Status = *userInput.Status
	}

	if err := sh.Repo.Update(staff); err != nil {
		switch err {
		case repositories.ErrNotFound:
			notfound(ctx)
		case repositories.ErrInternal:
			serverError(ctx)
		default:
			render(ctx, staffView.StaffForm(staff, sh.Repo.ConvertErrorsToFormErrors(err), endpoint), staff.FullName)
		}
		return
	}
	redirect(ctx, fmt.Sprintf("/staff/%d", staff.ID))
}

// post
func (sh *StaffHandler) Add(ctx *gin.Context) {
	var userInput struct {
		FullName    string `form:"fullName" binding:"required" validate:"required,min=1,max=100"`
		Username    string `form:"username" binding:"required" validate:"required,min=1,max=200"`
		PhoneNumber string `form:"phoneNumber" binding:"required" validate:"required,min=1,max=50"`
		Email       string `form:"email" binding:"required" validate:"required,email,min=1,max=50"`
		Role        string `form:"role" binding:"required" validate:"required,oneof=admin librarian,min=1,max=50"`
		Password    string `form:"password" binding:"required" validate:"required,min=8,max=20"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, staffView.StaffForm(
			&models.Staff{
				FullName:    userInput.FullName,
				Username:    userInput.Username,
				PhoneNumber: userInput.PhoneNumber,
				Email:       userInput.Email,
				Role:        userInput.Role,
			},
			parseValidationErrors(err),
			"/staff/add",
		), "add staff")
		return
	}

	if err := sh.Validator.Struct(&userInput); err != nil {
		render(ctx, staffView.StaffForm(
			&models.Staff{
				FullName:    userInput.FullName,
				Username:    userInput.Username,
				PhoneNumber: userInput.PhoneNumber,
				Email:       userInput.Email,
				Role:        userInput.Role,
			},
			parseValidationErrors(err),
			"/staff/add",
		), "add staff")
		return
	}

	// TODO: passwordHash
	hash, err := generateHash(userInput.Password)
	if err != nil {
		serverError(ctx)
		return
	}

	staff := models.Staff{
		FullName:     userInput.FullName,
		Username:     userInput.Username,
		PhoneNumber:  userInput.PhoneNumber,
		Email:        userInput.Email,
		Role:         userInput.Role,
		Status:       "active",
		PasswordHash: hash,
	}

	err = sh.Repo.Insert(&staff)
	if err != nil {
		switch err {
		case repositories.ErrInternal:
			serverError(ctx)
		case repositories.ErrNotFound:
			notfound(ctx)
		default:
			render(ctx, staffView.StaffForm(
				&models.Staff{
					FullName:    userInput.FullName,
					Username:    userInput.Username,
					PhoneNumber: userInput.PhoneNumber,
					Email:       userInput.Email,
					Role:        userInput.Role,
				},
				sh.Repo.ConvertErrorsToFormErrors(err),
				"/staff/add",
			), "add staff")
		}
		return
	}
	redirect(ctx, fmt.Sprintf("/staff/%d", staff.ID))
}

// get
func (sh *StaffHandler) Search(ctx *gin.Context) {
	fullName := ctx.Query("fullName")
	userName := ctx.Query("userName")
	role := ctx.Query("role")
	status := ctx.Query("status")

	pagination, err := readPagination(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	var total int64
	staff, err := sh.Repo.Filter(&models.StaffFilter{
		FullName: fullName,
		Username: userName,
		Role:     role,
		Status:   status,
	}, pagination, &total)

	if err != nil {
		serverError(ctx)
		return
	}

	render(ctx, staffView.StaffSearch(staff), "search staff")
}

func (sh *StaffHandler) LoginPage(ctx *gin.Context) {
	fmt.Println("lololo")
	render(ctx, authViews.LoginForm(views.Errors{}), "login")
}

func (sh *StaffHandler) Login(ctx *gin.Context) {
	var userInput struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, authViews.LoginForm(parseValidationErrors(err)), "login")
		return
	}

	staff, err := sh.Repo.GetByEmail(userInput.Email)
	if err != nil {
		if err == repositories.ErrNotFound {
			render(ctx, authViews.LoginForm(views.Errors{"login": "wrong user name or password"}), "login")
			return
		}
		serverError(ctx)
		return
	}

	if !compareHashAndPassoword(staff.PasswordHash, userInput.Password) {
		render(ctx, authViews.LoginForm(views.Errors{"login": "wrong user name or password"}), "login")
		return
	}

	session := sessions.Default(ctx)
	session.Set("staff", *staff)
	session.Save()

	redirect(ctx, "/")

}

func (sh *StaffHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	redirect(ctx, "/login")
}
