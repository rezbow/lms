package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	commonViews "lms/internal/views/common"
	memberViews "lms/internal/views/members"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MemberHandler struct {
	Repo      *repositories.MemberRepo
	Validator *validator.Validate
}

func (mh *MemberHandler) Index(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	members, err := mh.Repo.All(page, pageSize)
	if err != nil {
		render(ctx, commonViews.ServerError(err.Error()), "server error")
		return
	}
	render(ctx, memberViews.Index(members), "members")
}

func (mh *MemberHandler) GetById(ctx *gin.Context) {
	// hx := htmxHandler(ctx)
	memberId := ctx.Param("id")
	if _, err := strconv.Atoi(memberId); err != nil {
		render(ctx, commonViews.NotFound(), "404:((")
		return
	}
	member, err := mh.Repo.GetById(memberId)
	if err != nil {
		if err == repositories.ErrNotFound {
			render(ctx, commonViews.NotFound(), "404 :((")
			return
		}
		render(ctx, commonViews.ServerError(err.Error()), "internal server error")
		return
	}
	render(ctx, memberViews.MemberInfo(member), member.Name)
}

func (mh *MemberHandler) DeleteById(ctx *gin.Context) {
	memberId := ctx.Param("id")
	if err := mh.Repo.DeleteById(memberId); err != nil {
		if err == repositories.ErrNotFound {
			render(ctx, commonViews.NotFound(), "404 :((")
			return
		}
		render(ctx, commonViews.ServerError(err.Error()), "internal server error")
		return
	}
	redirect(ctx, "/members")
}

func (mh *MemberHandler) AddPage(ctx *gin.Context) {
	render(ctx, memberViews.MemberAddForm(), "add member")
}

func (mh *MemberHandler) Add(ctx *gin.Context) {
	var userInput struct {
		Name  string `form:"name" binding:"required" validate:"required,min=4,max=100"`
		Phone string `form:"phone" binding:"required" validate:"required,min=4,max=20"`
		Email string `form:"email" binding:"required" validate:"required,email,min=4,max=50"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
		return
	}

	if err := mh.Validator.Struct(&userInput); err != nil {
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
		return
	}

	member := models.Member{
		Name:  userInput.Name,
		Phone: userInput.Phone,
		Email: userInput.Email,
	}

	if err := mh.Repo.Insert(&member); err != nil {
		if err == repositories.ErrInternal {
			render(ctx, commonViews.ServerError(err.Error()), "internal server error")
			return
		}
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "form error")
		return
	}
	redirect(ctx, fmt.Sprintf("/members/%d", member.ID))
}

func (mh *MemberHandler) EditPage(ctx *gin.Context) {
	memberId := ctx.Param("id")
	member, err := mh.Repo.GetById(memberId)
	if err != nil {
		if err == repositories.ErrNotFound {
			render(ctx, commonViews.NotFound(), "404:((")
			return
		}
		render(ctx, commonViews.ServerError(""), "internal server error")
		return
	}
	render(ctx, memberViews.MemberEditForm(member), member.Name)
}

func (mh *MemberHandler) Update(ctx *gin.Context) {

	memberId := ctx.Param("id")
	id, err := strconv.Atoi(memberId)
	if err != nil {
		render(ctx, commonViews.NotFound(), "404:(")
		return
	}

	var memberUpdateForm struct {
		Name  *string `form:"name" validate:"omitempty,min=1,max=100"`
		Phone *string `form:"phone" validate:"omitempty,min=1,max=20"`
		Email *string `form:"email" validate:"omitempty,email,min=1,max=50"`
	}

	if err := ctx.ShouldBind(&memberUpdateForm); err != nil {
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
		return
	}

	if err := mh.Validator.Struct(&memberUpdateForm); err != nil {
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
		return
	}

	member := models.Member{
		ID: id,
	}

	if memberUpdateForm.Name != nil {
		member.Name = *memberUpdateForm.Name
	}
	if memberUpdateForm.Email != nil {
		member.Email = *memberUpdateForm.Email
	}
	if memberUpdateForm.Phone != nil {
		member.Phone = *memberUpdateForm.Phone
	}

	if err := mh.Repo.Update(&member); err != nil {
		if err == repositories.ErrInternal {
			render(ctx, commonViews.ServerError(err.Error()), "internal server error")
			return
		}
		render(ctx, commonViews.FormErrors([]string{err.Error()}), "form error")
		return
	}
	HXRedirect(ctx, "/members/"+memberId)
}

func (mh *MemberHandler) Search(ctx *gin.Context) {
	var total int64

	name := ctx.Query("name")
	phone := ctx.Query("phone")
	email := ctx.Query("email")

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	limit, err := strconv.Atoi(limitStr)

	members, err := mh.Repo.Filter(
		name,
		phone,
		email,
		limit,
		page,
		&total,
	)

	if err != nil {
		render(ctx, commonViews.ServerError(""), "server error")
		return
	}
	render(ctx, memberViews.MemberList(members), "members")
}
