package helpers

import (
	"fmt"
	"os"
)

func ValidateNArgs(nArgs int, usage string) {
	if len(os.Args) != nArgs {
		_, _ = fmt.Fprintf(os.Stderr, "Wrong number of arguments. Usage %s\n", usage)
		os.Exit(1)
	}
}
