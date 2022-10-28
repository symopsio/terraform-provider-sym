package utils

// ContainsString returns true if the lookup string is present in the given slice.
func ContainsString(slice []string, lookup string) bool {
	for _, item := range slice {
		if item == lookup {
			return true
		}
	}
	return false
}
