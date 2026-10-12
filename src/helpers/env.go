package helpers

import (
	"os"
	"regexp"
)

var windows_env_pattern = regexp.MustCompile(`%([^%]+)%`)

func ExpandEnv(value string) string {
	return os.ExpandEnv(windows_env_pattern.ReplaceAllString(value, "${$1}"))
}
