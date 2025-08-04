package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	commonViews "lms/internal/views/common"
	loanViews "lms/internal/views/loans"
	memberViews "lms/internal/views/members"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MemberHandler struct {
	Repo      *repositories.MemberRepo
	Validator *validator.Validate
}

func (mh *MemberHandler) Index(ctx *gin.Context) {
	pagination, err := readPagination(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	members, err := mh.Repo.All(pagination)
	if err != nil {
		render(ctx, commonViews.ServerError(err.Error()), "server error")
		return
	}
	render(ctx, memberViews.MemberSearch(members), "members")
}

func (mh *MemberHandler) GetById(ctx *gin.Context) {

	memberId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
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
	render(ctx, memberViews.MemberInfo(member), member.FullName)
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
	render(ctx, memberViews.MemberForm(nil, views.Errors{}, "/members/add"), "add member")
}

func (mh *MemberHandler) Add(ctx *gin.Context) {
	var userInput struct {
		FullName    string `form:"fullName" binding:"required" validate:"required,min=4,max=100"`
		PhoneNumber string `form:"phoneNumber" binding:"required" validate:"required,min=4,max=20"`
		Email       string `form:"email" binding:"required" validate:"required,email,min=4,max=50"`
		NationalId  string `form:"nationalId" binding:"required" validate:"required,min=4,max=50"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, memberViews.MemberForm(nil, parseValidationErrors(err), "/members/add"), "add member")
		return
	}

	if err := mh.Validator.Struct(&userInput); err != nil {
		render(ctx, memberViews.MemberForm(nil, parseValidationErrors(err), "/members/add"), "add member")
		return
	}

	member := models.Member{
		FullName:    userInput.FullName,
		PhoneNumber: userInput.PhoneNumber,
		NationalId:  userInput.NationalId,
		Email:       userInput.Email,
	}

	if err := mh.Repo.Insert(&member); err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
		} else if err == repositories.ErrInternal {
			serverError(ctx)
		} else {
			render(ctx, memberViews.MemberForm(nil, mh.Repo.ConvertErrorsToFormErrors(err), "/members/add"), member.FullName)
		}
		return
	}

	redirect(ctx, fmt.Sprintf("/members/%d", member.ID))
}

func (mh *MemberHandler) EditPage(ctx *gin.Context) {
	memberId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	member, err := mh.Repo.GetById(memberId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, memberViews.MemberForm(member, views.Errors{}, fmt.Sprintf("/members/%d/edit", member.ID)), member.FullName)
}

func (mh *MemberHandler) Update(ctx *gin.Context) {

	memberId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	var memberUpdateForm struct {
		FullName    *string `form:"fullName" validate:"omitempty,min=1,max=100"`
		PhoneNumber *string `form:"phoneNumber" validate:"omitempty,min=1,max=20"`
		NationalId  *string `form:"nationalId" validate:"omitempty,min=1,max=20"`
		Email       *string `form:"email" validate:"omitempty,email,min=1,max=50"`
		Status      *string `form:"status" validate:"omitempty,oneof=active suspended,min=1,max=50"`
	}

	member, err := mh.Repo.GetById(memberId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	endpoint := fmt.Sprintf("/members/%d/edit", member.ID)

	if err := ctx.ShouldBind(&memberUpdateForm); err != nil {
		render(ctx, memberViews.MemberForm(member, parseValidationErrors(err), endpoint), member.FullName)
		return
	}

	if err := mh.Validator.Struct(&memberUpdateForm); err != nil {
		render(ctx, memberViews.MemberForm(member, parseValidationErrors(err), endpoint), member.FullName)
		return
	}

	if memberUpdateForm.FullName != nil {
		member.FullName = *memberUpdateForm.FullName
	}
	if memberUpdateForm.Email != nil {
		member.Email = *memberUpdateForm.Email
	}
	if memberUpdateForm.PhoneNumber != nil {
		member.PhoneNumber = *memberUpdateForm.PhoneNumber
	}

	if memberUpdateForm.NationalId != nil {
		member.NationalId = *memberUpdateForm.NationalId
	}

	if err := mh.Repo.Update(member); err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
		} else if err == repositories.ErrInternal {
			serverError(ctx)
		} else {
			render(ctx, memberViews.MemberForm(member, mh.Repo.ConvertErrorsToFormErrors(err), endpoint), member.FullName)
		}
		return
	}
	redirect(ctx, fmt.Sprintf("/members/%d", memberId))
}

func (mh *MemberHandler) Search(ctx *gin.Context) {
	var total int64

	name := ctx.Query("fullName")
	phone := ctx.Query("phoneNumber")
	email := ctx.Query("email")
	nationalId := ctx.Query("nationalId")

	pagination, err := readPagination(ctx)

	if err != nil {
		notfound(ctx)
		return
	}

	members, err := mh.Repo.Filter(
		&models.MemberFilter{
			FullName:    name,
			PhoneNumber: phone,
			Email:       email,
			NationalId:  nationalId,
		},
		pagination,
		&total,
	)

	if err != nil {
		serverError(ctx)
		return
	}
	render(ctx, memberViews.MemberSearch(members), "members")
}

func (mh *MemberHandler) AddLoanPage(ctx *gin.Context) {
	memberId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	member, err := mh.Repo.GetById(memberId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, loanViews.LoanForm(&models.Loan{MemberId: member.ID}, views.Errors{}, "/loans/add"), "add loan")
}

func parseRepoError(err error) {

}
