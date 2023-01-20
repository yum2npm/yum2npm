package utils

import (
	"regexp"
)

func NamedMatches(pattern string, value string) map[string]string {
	var regex = regexp.MustCompile(pattern)

	match := regex.FindStringSubmatch(value)
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}
