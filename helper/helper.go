package helper

func StringOrDefault(ptr *string, def string) string {
	if ptr != nil {
		return *ptr
	}
	return def
}

func IntOrDefault(ptr *int, def int) int {
	if ptr != nil {
		return *ptr
	}
	return def
}
