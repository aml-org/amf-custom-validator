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
	Generated string
}

func Fixtures(root string) []Fixture {
	fixtures := make([]Fixture, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic("error reading data directory")
		}
		if strings.Index(path, ".yaml") > -1 {
			parsed := strings.ReplaceAll(path, ".yaml",".parsed")
			generated := strings.ReplaceAll(path, ".yaml",".rego")
			fixture := Fixture{
				Profile: path,
				Parsed: parsed,
				Generated: generated,
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

func (f Fixture) ReadGenerated() string {
	bytes, err := ioutil.ReadFile(f.Generated)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func ForceWrite(f string, data string) {
	ioutil.WriteFile(f, []byte(data),0644)
}