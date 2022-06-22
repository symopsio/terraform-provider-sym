package utils

// ContainsStr checks if the given `str` is present in the given array of strings `list`.
func ContainsStr(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}
