package internal

import (
	"lms/internal/models"
	"lms/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		staff := utils.ExtractStaffFromSession(session)
		if staff == nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		staff := utils.ExtractStaffFromSession(session)
		if staff == nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}
		if staff.Role != models.RoleAdmin {
			ctx.Redirect(http.StatusSeeOther, "/")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
