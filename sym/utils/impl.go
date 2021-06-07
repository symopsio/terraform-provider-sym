package utils

import (
	"encoding/base64"
	"io/ioutil"
)

// Takes in an impl in base64, text, or filename format,
// and consistently returns the impl in text format, if possible.
func ParseImpl(impl string) string {
	return parseImpl(impl, true)
}

// Takes in an impl in base64 or text format,
// and consistently returns the impl in text format, if possible.
func ParseRemoteImpl(impl string) string {
	return parseImpl(impl, false)
}

func parseImpl(impl string, readFile bool) string {
	contents, err := base64.StdEncoding.DecodeString(impl)
	if err != nil && readFile {
		contents, err = ioutil.ReadFile(impl)
	}

	if err != nil {
		return impl
	}
	return string(contents)
}
