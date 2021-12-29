package yaml

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
)

type Yaml struct {
	data *yaml.Node
}

// NewYaml returns a pointer to a new `Yaml` object after unmarshaling `body` bytes
func NewYaml(body []byte) (*Yaml, error) {
	var val yaml.Node
	err := yaml.Unmarshal(body, &val)
	if err != nil {
		return nil, errors.New("unmarshal []byte to yaml failed: " + err.Error())
	}
	root := val.Content[0]
	return &Yaml{root}, nil
}

// Check if the given branch was found
func (y *Yaml) IsFound() bool {
	if y.data == nil {
		return false
	}
	return true
}

// Get returns a pointer to a new `Yaml` object for `key` in its `map` representation
//
// Example:
//      y.Get("xx").Get("yy").Int()
func (y *Yaml) Get(key interface{}) *Yaml {
	res, _ := y.GetOrError(key)
	return res // always returns yaml node, if key not present yaml node with nil value is returned
}

func (y *Yaml) GetOrError(key interface{}) (*Yaml, error)  {
	found := false
	for _, n := range y.data.Content {
		if found {
			return &Yaml{n}, nil
		}
		if n.Kind == yaml.ScalarNode && n.Value == key {
			found = true
		}
	}
	return &Yaml{nil}, errors.New("Key not found")
}

// GetPath searches for the item as specified by the branch
//
// Example:
//      y.GetPath("bb", "cc").Int()
func (y *Yaml) GetPath(branch ...interface{}) *Yaml {
	yin := y
	for _, p := range branch {
		yin = yin.Get(p)
	}
	return yin
}

// Array type asserts to an `array`
func (y *Yaml) Array() ([]*Yaml, error) {
	if !y.IsFound() {
		return nil, errors.New("not node found")
	}
	if y.data.Kind == yaml.SequenceNode {
		var acc []*Yaml
		for _, n := range y.data.Content {
			acc = append(acc, &Yaml{n})
		}
		return acc, nil
	}
	return nil, errors.New(fmt.Sprintf("cannot transform into array non-scalar YAML node at [%d,%d]", y.data.Line, y.data.Column))
}

func (y *Yaml) IsArray() bool {
	return y.data.Kind == yaml.SequenceNode
}

// return the size of array
func (y *Yaml) GetArraySize() (int, error) {
	if y.data.Kind == yaml.SequenceNode {
		return len(y.data.Content), nil
	}
	return -1, errors.New(fmt.Sprintf("error getting size of non-sequence yaml node at [%d,%d]", y.data.Line, y.data.Column))
}

// GetIndex returns a pointer to a new `Yaml` object.
// for `index` in its `array` representation
//
// Example:
//      y.Get("xx").GetIndex(1).String()
func (y *Yaml) GetIndex(index int) *Yaml {
	a, err := y.Array()
	if err == nil {
		if len(a) > index {
			return a[index]
		}
	}
	return &Yaml{nil}
}

func (y *Yaml) Pos() (int, int) {
	return y.data.Line, y.data.Column
}

// Int type asserts to `int`
func (y *Yaml) Int() (int, error) {
	if !y.IsFound() {
		return -1, errors.New("not node found")
	}
	if y.data.Kind == yaml.ScalarNode && y.data.Tag == "!!int" {
		i, e := strconv.Atoi(y.data.Value)
		if e == nil {
			return i, nil
		}
		return -1, errors.New(fmt.Sprintf("error in int conversion of parsed int at [%d,%d]", y.data.Line, y.data.Column))
	}
	return 0, errors.New(fmt.Sprintf("type assertion to int failed for scalar at [%d,%d]", y.data.Line, y.data.Column))
}

// Bool type asserts to `bool`
func (y *Yaml) Bool() (bool, error) {
	if !y.IsFound() {
		return false, errors.New("not node found")
	}
	if y.data.Kind == yaml.ScalarNode && y.data.Tag == "!!bool" {
		i, e := strconv.ParseBool(y.data.Value)
		if e == nil {
			return i, nil
		}
		return false, errors.New(fmt.Sprintf("error in bool conversion of parsed bool at [%d,%d]", y.data.Line, y.data.Column))
	}
	return false, errors.New(fmt.Sprintf("type assertion to bool failed for scalar at [%d,%d]", y.data.Line, y.data.Column))
}

// String type asserts to `string`
func (y *Yaml) String() (string, error) {
	if !y.IsFound() {
		return "", errors.New("not node found")
	}
	if y.data.Kind == yaml.ScalarNode && y.data.Tag == "!!str" {
		return y.data.Value, nil
	}
	return "", errors.New("type assertion to string failed")
}

func (y *Yaml) Float() (float64, error) {
	if !y.IsFound() {
		return -1, errors.New("not node found")
	}
	if y.data.Kind == yaml.ScalarNode && y.data.Tag == "!!float" {
		f, e := strconv.ParseFloat(y.data.Value, 64)
		if e == nil {
			return f, nil
		}
		return -1, errors.New(fmt.Sprintf("error in float conversion of parsed float at [%d,%d]", y.data.Line, y.data.Column))
	}
	return -1, errors.New(fmt.Sprintf("type assertion to float failed for scalar at [%d,%d]", y.data.Line, y.data.Column))
}

// Map type asserts to `map`
func (y *Yaml) Map() (map[string]*Yaml, error) {
	if !y.IsFound() {
		return nil, errors.New("not node found")
	}
	if y.data.Kind == yaml.MappingNode {
		acc := make(map[string]*Yaml)
		k := ""
		for i, n := range y.data.Content {
			if i%2 == 0 {
				k = n.Value
			} else {
				acc[k] = &Yaml{n}
			}
		}
		return acc, nil
	}
	return nil, errors.New(fmt.Sprintf("type assertion to map[interface]*yaml.Node failed at [%d,%d]", y.data.Line, y.data.Column))
}

// Check if it is a map
func (y *Yaml) IsMap() bool {
	return y.data.Kind == yaml.MappingNode
}

// Get all the keys of the map
func (y *Yaml) GetMapKeys() ([]string, error) {
	m, err := y.Map()

	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)

	}
	return keys, nil
}
