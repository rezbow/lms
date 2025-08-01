package utils

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
