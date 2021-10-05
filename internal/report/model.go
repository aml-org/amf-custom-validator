package report

type Severity = string

const (
	Warning   Severity = "http://www.w3.org/ns/shacl#Warning"
	Violation Severity = "http://www.w3.org/ns/shacl#Violation"
	Info      Severity = "http://www.w3.org/ns/shacl#Info"
)

type any = interface{}

type object = map[string]any

type Ref struct {
	Id string `json:"@id"`
}

type Typed struct {
	Type []string `json:"@type"`
}

type Result struct {
	Typed
	SourceShapeName string  `json:"sourceShapeName"`
	FocusNode       Ref     `json:"focusNode"`
	ResultSeverity  *Ref    `json:"resultSeverity,omitempty"`
	ResultMessage   string  `json:"resultMessage"`
	Trace           []Trace `json:"trace"`
}

type Trace struct {
	Typed
	Component  string     `json:"component"`
	ResultPath string     `json:"resultPath"`
	TraceValue TraceValue `json:"traceValue"`
	Location   *Location  `json:"location,omitempty"`
}

type TraceValue struct {
	Typed
	Actual    any       `json:"actual,omitempty"`
	Condition *string   `json:"condition,omitempty"`
	Expected  any       `json:"expected,omitempty"`
	Negated   bool      `json:"negated"`
	Argument  *string   `json:"argument,omitempty"`
	SubResult *[]Result `json:"subResult,omitempty"`
}

type Location struct {
	Typed
	Uri   *Ref  `json:"uri,omitempty"`
	Range Range `json:"range"`
}

type Range struct {
	Typed
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Typed
	Line   int `json:"line"`
	Column int `json:"column"`
}
