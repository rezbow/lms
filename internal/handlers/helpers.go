package handlers

import (
	"lms/internal/views"
	commonViews "lms/internal/views/common"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/donseba/go-htmx"
	"github.com/gin-gonic/gin"
)

var htmxApp = htmx.New()

func isHtmx(ctx *gin.Context) bool {
	return ctx.GetHeader("HX-Request") == "true"
}

func redirect(ctx *gin.Context, location string) {
	HXRedirect(ctx, location)
}

func HXRedirect(ctx *gin.Context, location string) {
	if isHtmx(ctx) {
		ctx.Header("HX-Redirect", location)
		ctx.Status(http.StatusNoContent)
		return
	}
	ctx.Redirect(http.StatusSeeOther, location)
}

func render(ctx *gin.Context, html templ.Component, title string) {
	if isHtmx(ctx) {
		html.Render(ctx.Request.Context(), ctx.Writer)
	} else {
		views.LayoutNew(
			html,
			title,
		).Render(ctx.Request.Context(), ctx.Writer)
	}
	ctx.Status(http.StatusOK)
}

func htmxHandler(ctx *gin.Context) *htmx.Handler {
	return htmxApp.NewHandler(ctx.Writer, ctx.Request)
}

func readID(ctx *gin.Context) (int, error) {
	id := ctx.Param("id")
	if idInt, err := strconv.Atoi(id); err != nil {
		return 0, err
	} else {
		if idInt < 1 {
			return 0, err
		}
		return idInt, nil
	}
}

func serverError(ctx *gin.Context) {
	render(ctx, commonViews.ServerError(""), "server error")
}

func notfound(ctx *gin.Context) {
	render(ctx, commonViews.NotFound(), "404:((")
}

func formError(ctx *gin.Context, err error) {
	render(ctx, commonViews.FormErrors([]string{err.Error()}), "error")
}
