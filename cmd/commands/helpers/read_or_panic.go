package helpers

import "io/ioutil"

func ReadOrPanic(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}
