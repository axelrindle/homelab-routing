package util

import (
	"strings"
)

// Tests whether a string s contains any of the given substrings.
func ContainsAny(s string, substrs []string) bool {
	for _, needle := range substrs {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}
