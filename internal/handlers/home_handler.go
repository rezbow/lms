package handlers

import (
	"lms/internal/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", views.Home())
}
