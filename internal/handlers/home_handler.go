package handlers

import (
	"lms/internal/views"

	"github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
	render(ctx, views.Home(), "Home")
}
