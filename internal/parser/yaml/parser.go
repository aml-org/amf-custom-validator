package yaml

import (
	"strings"
)

func CleanYamlString(s string) string {
	if strings.Index(s, "\"") == 0 {
		s = s[1:len(s)]
	}
	if strings.Index(s, "\"") == len(s)-1 {
		s = s[0 : len(s)-1]
	}
	return s
}
