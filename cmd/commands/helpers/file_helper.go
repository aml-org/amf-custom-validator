package helpers

import (
	"errors"
	"os"
)

func ExistsFile(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func OpenFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	CheckError(err)
	return file
}

func CreateFile(path string) *os.File {
	file, err := os.Create(path)
	CheckError(err)
	return file
}

func OpenOrCreateFile(path string) *os.File {
	if ExistsFile(path) {
		return OpenFile(path)
	} else {
		return CreateFile(path)
	}
}

func WriteString(file *os.File, content string) {
	_, wErr := file.WriteString(content)
	CheckError(wErr)
	sErr := file.Sync()
	CheckError(sErr)
}
