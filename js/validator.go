//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"syscall/js"

	"github.com/aml-org/amf-custom-validator/internal/validator"
	c "github.com/aml-org/amf-custom-validator/pkg/config"
)

func validateWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 3 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		debug := args[2].Bool()

		res, err := validator.ValidateRegoV1(profileString, dataString, debug, nil)
		if err != nil {
			return err.Error()
		}
		return res

	})
	return jsonFunc
}

func validateWrapperV0() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 3 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		debug := args[2].Bool()

		res, err := validator.ValidateRegoV0(profileString, dataString, debug, nil)
		if err != nil {
			return err.Error()
		}
		return res

	})
	return jsonFunc
}

func validateWithConfigurationWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 5 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		debug := args[2].Bool()
		if !args[3].IsUndefined() {
			return "ValidationConfiguration not yet supported, please set this parameter as 'undefined'. ValidationConfiguration needs to import functions from JS which is not supported by Go 1.19"
		}
		reportConfig := buildReportConfiguration(args[4])
		res, err := validator.ValidateWithConfigurationV1(profileString, dataString, debug, nil, c.DefaultValidationConfiguration{}, reportConfig)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func validateWithConfigurationWrapperV0() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 5 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		debug := args[2].Bool()
		if !args[3].IsUndefined() {
			return "ValidationConfiguration not yet supported, please set this parameter as 'undefined'. ValidationConfiguration needs to import functions from JS which is not supported by Go 1.19"
		}
		reportConfig := buildReportConfiguration(args[4])
		res, err := validator.ValidateWithConfigurationV0(profileString, dataString, debug, nil, c.DefaultValidationConfiguration{}, reportConfig)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func buildReportConfiguration(value js.Value) c.ReportConfiguration {
	includeReportCreationTime := value.Get("IncludeReportCreationTime").Bool()
	reportSchemaIri := value.Get("ReportSchemaIri").String()
	lexicalSchemaIri := value.Get("LexicalSchemaIri").String()
	return c.ReportConfiguration{
		IncludeReportCreationTime: includeReportCreationTime,
		ReportSchemaIri:           reportSchemaIri,
		LexicalSchemaIri:          lexicalSchemaIri,
	}
}

func genRegoWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		res, err := validator.GenerateRegoV1(profileString, false, nil)
		if err != nil {
			return err.Error()
		}
		return res.Code
	})
	return jsonFunc
}

func genRegoWrapperV0() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		res, err := validator.GenerateRegoV0(profileString, false, nil)
		if err != nil {
			return err.Error()
		}
		return res.Code
	})
	return jsonFunc
}

func normalizeInputWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		dataString := args[0].String()
		normalizedInput, err := validator.ProcessInput(dataString, false, nil)
		if err != nil {
			return err.Error()
		}
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "  ")
		err = enc.Encode(normalizedInput)
		if err != nil {
			return err.Error()
		}

		return b.String()
	})
	return jsonFunc
}

func exitWrapper(c chan bool) js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		c <- true
		return nil
	})
	return jsonFunc
}

func main() {
	c := make(chan bool)
	// validate
	f := validateWrapper()
	js.Global().Set("__AMF__validateCustomProfile", f)
	// validate with configuration
	f = validateWithConfigurationWrapper()
	js.Global().Set("__AMF__validateCustomProfileWithConfiguration", f)
	// gen rego
	f = genRegoWrapper()
	js.Global().Set("__AMF__generateRego", f)
	// normalizeInput
	f = normalizeInputWrapper()
	js.Global().Set("__AMF__normalizeInput", f)

	// validate legacy
	f = validateWrapperV0()
	js.Global().Set("__AMF__validateCustomProfileV0", f)
	// validate with configuration legacy
	f = validateWithConfigurationWrapperV0()
	js.Global().Set("__AMF__validateCustomProfileWithConfigurationV0", f)
	// gen rego legacy
	f = genRegoWrapperV0()
	js.Global().Set("__AMF__generateRegoV0", f)

	// exit
	f = exitWrapper(c)
	js.Global().Set("__AMF__terminateValidator", f)
	<-c
}
