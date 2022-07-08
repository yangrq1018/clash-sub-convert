package util

// FirstString returns the first non-empty string
func FirstString(ss ...string) string {
	for i := range ss {
		if ss[i] != "" {
			return ss[i]
		}
	}
	return ""
}
