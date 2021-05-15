package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Fixture struct {
	Profile string
	Parsed string
}

const root string = "./data/"

func Fixtures() []Fixture {
	fixtures := make([]Fixture, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Index(path, ".yaml") > -1 {
			parsed := strings.ReplaceAll(path, ".yaml",".parsed")
			fixture := Fixture{
				Profile: path,
				Parsed: parsed,
			}
			fixtures = append(fixtures, fixture)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return fixtures
}

func (f Fixture) ReadParsed() string {
	bytes, err := ioutil.ReadFile(f.Parsed)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
