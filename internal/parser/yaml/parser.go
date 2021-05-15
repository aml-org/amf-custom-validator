package yaml

import (
	"errors"
	"github.com/kylelemons/go-gypsy/yaml"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func readFile(path string) (yaml.Node, error) {
	f,err := os.Open(path)
	defer func() {
		err :=  f.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	return yaml.Parse(f)
}

func readURL(path string) (yaml.Node, error) {
	response, err := http.Get(path)
	if err != nil {
		return nil,err
	}
	defer func() {
		err :=  response.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}

	return yaml.Parse(response.Body)
}

func Parse(path string) (yaml.Node, error) {
	var node yaml.Node
	var err error

	if strings.Index(path, "file://") == 0 || strings.Index(path, "://") == -1 {
		node, err = readFile(strings.ReplaceAll(path, "file://", ""))
		if err != nil {
			return node, errors.New("Error reading provided YAML file " + path)
		}
	} else {
		node, err = readURL(path)
		if err != nil {
			return node, errors.New("Error reading provided YAML URL " + path)
		}
	}

	return node, nil
}

func GetString(m yaml.Map, k string) (string, error) {
	switch v := m.Key(k); tc := v.(type) {
	case yaml.Scalar:
		return tc.String(), nil
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