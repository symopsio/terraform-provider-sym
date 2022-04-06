package utils

import (
	"encoding/base64"
	"io/ioutil"
)

// ParseImpl takes in an impl in base64, text, or filename format,
// and consistently returns the impl in text format, if possible.
func ParseImpl(impl string) string {
	return parseImpl(impl, true)
}

// ParseRemoteImpl takes in an impl in base64 or text format,
// and consistently returns the impl in text format, if possible.
func ParseRemoteImpl(impl string) string {
	return parseImpl(impl, false)
}

// parseImpl decodes the given implementation and parses it to
// text format whenever possible.
//
// Args:
//		impl: the string to parse
//		isFilePath: whether the impl is a file path or the full contents of the impl
func parseImpl(impl string, isFilePath bool) string {
	contents, err := base64.StdEncoding.DecodeString(impl)
	if err != nil && isFilePath {
		contents, err = ioutil.ReadFile(impl)
	}

	if err != nil {
		return impl
	}
	return string(contents)
}
