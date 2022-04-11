package commands

import (
	"fmt"
	"os"
)

func Fallback(cmd string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s is not a recognized command. Available commands are: validate, normalize, generate, help\n", cmd)
	os.Exit(1)
}
