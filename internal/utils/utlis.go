package utils

import (
	"github.com/gin-contrib/sessions"
	"lms/internal/models"
)

func DefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func GetIntFromMap(m map[string]any, key string) (int, bool) {
	i, ok := m[key].(int)
	if !ok {
		return 0, ok
	}
	return i, ok
}

func ExtractStaffFromSession(session sessions.Session) *models.Staff {
	s := session.Get("staff")
	if session == nil {
		return nil
	}
	staff, ok := s.(models.Staff)
	if !ok {
		return nil
	}
	return &staff
}
