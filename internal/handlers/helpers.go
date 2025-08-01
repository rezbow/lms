package handlers

import (
	"errors"
	"lms/internal/utils"
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

func readPagination(ctx *gin.Context) (*utils.Pagination, error) {
	var page int
	var limit int
	var err error

	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, err
		}
	}

	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}
	}

	if page <= 0 || limit <= 0 {
		return nil, errors.New("needs positive integer")
	}
	return utils.NewPagination(page, limit), nil
}

func readIntFromQuery(str string) (int, error) {
	if str == "" {
		return -1, nil
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, errors.New("needs positive integer")
	}
	return i, nil

}
