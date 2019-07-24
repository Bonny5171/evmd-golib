package utils

import (
	"regexp"
	"strings"
)

func RemoveSpaces(s string) string {
	n := strings.Split(strings.TrimSpace(s), " ")
	return strings.Join(n, "")
}

func RemoveSpacesAndLower(s string) string {
	return strings.ToLower(RemoveSpaces(s))
}

func SnakeCase(s string) string {
	re := regexp.MustCompile(`([[:lower:]])([[:upper:]])`)
	return strings.ReplaceAll(strings.ToLower(re.ReplaceAllString(s, "${1}_${2}")), " ", "_")
}
