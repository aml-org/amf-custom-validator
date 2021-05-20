package yaml

import (
	"errors"
	"github.com/kylelemons/go-gypsy/yaml"
	"strconv"
	"strings"
)

func Parse(profile string) (yaml.Node, error) {
	return yaml.Parse(strings.NewReader(profile))
}

func CleanYamlString(s string) string {
	if strings.Index(s, "\"") == 0 {
		s = s[1:len(s)]
	}
	if strings.Index(s, "\"") == len(s)-1 {
		s = s[0 : len(s)-1]
	}
	return s
}
func GetString(m yaml.Map, k string) (string, error) {
	switch v := m.Key(k); tc := v.(type) {
	case yaml.Scalar:
		s := CleanYamlString(tc.String())
		return s, nil
	default:
		return "", errors.New("Missing " + k + " property")
	}
}

func GetInt(m yaml.Map, k string) (int, error) {
	switch v := m.Key(k); tc := v.(type) {
	case yaml.Scalar:
		return strconv.Atoi(tc.String())
	default:
		return -1, errors.New("Missing " + k + " property")
	}
}

func GetMap(m yaml.Map, k string) (yaml.Map, error) {
	switch v := m.Key(k); tc := v.(type) {
	case yaml.Map:
		return tc, nil
	default:
		return nil, errors.New("Missing " + k + " property")
	}
}

func GetList(m yaml.Map, k string) (yaml.List, error) {
	switch v := m.Key(k); tc := v.(type) {
	case yaml.List:
		return tc, nil
	default:
		return nil, errors.New("Missing " + k + " property")
	}
}
