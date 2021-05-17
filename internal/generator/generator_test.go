package generator


import (
	"github.com/aml-org/amfopa/internal/parser"
	"github.com/aml-org/amfopa/test"
	"strings"
	"testing"
)

func TestGenerated(t *testing.T) {

	for _,fix := range test.Fixtures("../../test/data") {
		profile, err := parser.Parse(fix.Profile)
		if err != nil {
			panic(err)
		}
		generated := Generate(*profile)
		//test.ForceWrite(fix.Generated, generated.Code)
		actual := strings.TrimSpace(generated.Code)
		expected := strings.TrimSpace(fix.ReadGenerated())

		if actual != expected {
			t.Errorf("Error in expected profile %s\n\nActual:\n%s\n----\nExpected:\n%s", fix.Profile, actual, expected)
		}
	}
}