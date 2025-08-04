package handlers

import (
	"net/url"
	"testing"
)

func TestSlugify(t *testing.T) {
	input := "رضا بوالخس"
	t.Log(slugify(input))
	t.Log(url.PathEscape(input))
}
