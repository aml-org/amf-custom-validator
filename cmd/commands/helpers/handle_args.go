package helpers

import (
	"fmt"
	"os"
)

func ValidateNsArgs(nsArgs []int, usage string) {
	for _, nArgs := range nsArgs {
		if len(os.Args) == nArgs {
			return
		}
	}
	errorArgs(usage)
}

func ValidateNArgs(nArgs int, usage string) {
	if len(os.Args) != nArgs {
		errorArgs(usage)
	}
}

func errorArgs(usage string) {
	_, _ = fmt.Fprintf(os.Stderr, "Wrong number of arguments. Usage %s\n", usage)
	os.Exit(1)
}
