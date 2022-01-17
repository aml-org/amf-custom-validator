package path

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"strings"
)

type PropertyPath interface {
	// Source Returns the raw source path expression
	Source() string

	// Expanded returns the expanded version of Source using an IriExpander (which wraps a context). It discards inverse path '^' characters
	Expanded(iriExpander *misc.IriExpander) (string, error)

	// Trace returns the expanded version of Source using an IriExpander (which wraps a context). It includes inverse path '^' characters for tracing purposes
	Trace(iriExpander *misc.IriExpander) (string, error)
}

type BasePath struct {
	source string
}

// Common code to all paths

func expandAndJoinParts(iriExpander *misc.IriExpander, parts []PropertyPath, separator string, keepInversePaths bool) (string, error) {
	expandedParts, err := expandParts(iriExpander, parts, keepInversePaths)
	if err != nil {
		return "", err
	} else {
		return strings.Join(expandedParts, separator), nil
	}
}

func expandParts(iriExpander *misc.IriExpander, parts []PropertyPath, keepInversePaths bool) ([]string, error) {
	expandedParts := make([]string, len(parts))
	for i, path := range parts {
		expanded, err := expandPart(iriExpander, keepInversePaths, path)
		if err != nil {
			return nil, err
		}
		expandedParts[i] = expanded
	}
	expandedParts = addExpressionNesting(parts, expandedParts)
	return expandedParts, nil
}

func expandPart(iriExpander *misc.IriExpander, keepInversePaths bool, path PropertyPath) (string, error) {
	// Not the prettiest if Go does not allow using partially applied methods as arguments
	if keepInversePaths {
		return path.Trace(iriExpander)
	} else {
		return path.Expanded(iriExpander)
	}
}

func addExpressionNesting(original []PropertyPath, expanded []string) []string {
	for i, path := range original {
		switch path.(type) {
		case AndPath, OrPath:
			expanded[i] = fmt.Sprintf("(%s)", expanded[i])
		}
	}
	return expanded
}
