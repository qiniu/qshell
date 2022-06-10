package utils

import "strings"

func IsHostUnavailableError(err error) bool {
	if err == nil {
		return false
	}

	info := err.Error()
	return strings.Contains(info, "dial tcp: lookup") && strings.Contains(info, ": no such host")
}
