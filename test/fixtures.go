package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Fixture struct {
	Profile   string
	Parsed    string
	Generated string
}

type IntegrationFixture string

func Fixtures(root string) []Fixture {
	fixtures := make([]Fixture, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic("error reading data directory")
		}
		if strings.Index(path, ".yaml") > -1 && strings.Index(path, "integration") == -1 {
			parsed := strings.ReplaceAll(path, ".yaml", ".parsed")
			generated := strings.ReplaceAll(path, ".yaml", ".rego")
			fixture := Fixture{
				Profile:   path,
				Parsed:    parsed,
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

func IntegrationFixtures(root string) []IntegrationFixture {
	var fixtures []IntegrationFixture
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.Index(path, "profile") > -1 {
			fixtures = append(fixtures, IntegrationFixture(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return fixtures
}

func (f IntegrationFixture) ReadProfile() string {
	bytes, err := ioutil.ReadFile(string(f) + "/profile.yaml")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (f IntegrationFixture) ReadFixturePositiveData() string {
	bytes, err := ioutil.ReadFile(string(f) + "/positive.data.jsonld")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (f IntegrationFixture) ReadFixtureNegativeData() string {
	bytes, err := ioutil.ReadFile(string(f) + "/negative.data.jsonld")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (f IntegrationFixture) ReadFixturePositiveReport() string {
	bytes, err := ioutil.ReadFile(string(f) + "/positive.report.jsonld")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (f IntegrationFixture) ReadFixtureNegativeReport() string {
	bytes, err := ioutil.ReadFile(string(f) + "/negative.report.jsonld")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (f Fixture) ReadProfile() string {
	bytes, err := ioutil.ReadFile(f.Profile)
	if err != nil {
		panic(err)
	}
	return string(bytes)
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

// Only for fixing tests
func ForceWrite(f string, data string) {
	ioutil.WriteFile(f, []byte(data), 0644)
}
