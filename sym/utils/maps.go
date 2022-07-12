package utils

// GetStr retrieves a string safely from a given map `m` at the given `key`, if it exists.
// If the value does not exist or is not a string, returns `("", false)`.
func GetStr(m map[string]interface{}, key string) (string, bool) {
	value, found := m[key]
	if !found {
		// No value set in map, return failure.
		return "", false
	}

	if _, ok := value.(string); !ok {
		// Value set in map is not a string, return failure.
		return "", false
	}

	// Value is set in map and is a string, so return it.
	return value.(string), true
}
