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

// GetStrArray retrieves an array of strings safely from a given map `m` at the given `key`, if it exists.
// If the value does not exist or is not an array of strings, returns `([]string{}, bool)`.
func GetStrArray(m map[string]interface{}, key string) ([]string, bool) {
	value, found := m[key]
	if !found {
		// No value set in map, return failure.
		return []string{}, false
	}

	if value, ok := value.([]string); ok {
		// If we got a value, and it's already an array of strings, just return it.
		return value, true
	}

	// If we don't already have an array of strings, try casting.
	if _, ok := value.([]interface{}); !ok {
		// We either don't have an array or we don't have an array of items that might be strings,
		// so return the failure case.
		return []string{}, false
	}

	// Since we know we have an array of interfaces, try casting all the items to strings.
	var stringValues []string
	for _, item := range value.([]interface{}) {
		_, ok := item.(string)

		if ok {
			stringValues = append(stringValues, item.(string))
		} else {
			// If any item in the array is not a string, return failure.
			return []string{}, false
		}
	}

	// Value is set in map and is an array of strings, so return it.
	return stringValues, true
}
