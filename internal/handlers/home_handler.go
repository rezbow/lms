package handlers

import (
	"lms/internal/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
	session := sessions.Default(ctx)
	render(ctx, views.Home(session), "Home")
}
