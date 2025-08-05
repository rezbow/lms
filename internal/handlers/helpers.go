package handlers

import (
	"errors"
	"lms/internal/utils"
	"lms/internal/views"
	commonViews "lms/internal/views/common"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func redirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusSeeOther, location)
}

func render(ctx *gin.Context, html templ.Component, title string) {
	session := sessions.Default(ctx)

	views.Layout(
		html,
		title,
		session,
	).Render(ctx.Request.Context(), ctx.Writer)
	ctx.Status(http.StatusOK)
}

func readID(ctx *gin.Context) (uint, error) {
	id := ctx.Param("id")
	if idInt, err := strconv.Atoi(id); err != nil {
		return 0, err
	} else {
		if idInt < 1 {
			return 0, err
		}
		return uint(idInt), nil
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
		return 0, nil
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

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func parseValidationErrors(err error) views.Errors {
	errors := make(views.Errors)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			errors[lowerFirst(fe.Field())] = fe.Error()
		}
	} else {
		errors["_"] = err.Error()
	}
	return errors
}

func slugify(input string) string {
	var builder strings.Builder
	for _, c := range input {
		switch {
		case unicode.Is(unicode.Latin, c) && unicode.IsLetter(c):
			builder.WriteRune(unicode.ToLower(c))
		case unicode.Is(unicode.Arabic, c):
			builder.WriteRune(c)
		case unicode.IsDigit(c):
			builder.WriteRune(c)
		case c == ' ' || c == '-' || c == '_':
			builder.WriteRune('-')
		}
	}

	slug := builder.String()

	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return strings.Trim(slug, "-")
}

func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func compareHashAndPassoword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
