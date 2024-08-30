package struct_convert_test

import (
	"strings"
	"testing"

	"github.com/VaLeraGav/struct_convert"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

func TestStructConvert(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected string
	}{
		{
			name: "Check Array",
			jsonData: `{
  "query": "Виктор Иван",
  "count": 7
}`,
			expected: `
type Menuitem struct {
	Value string ` + "\t`json:\"value\"`" + `
	Name  string ` + "\t`json:\"name\"`" + `
}
type Test struct {
	Id      string      ` + "\t`json:\"id\"`" + `
	Value   string      ` + "\t`json:\"value\"`" + `
	Menuitem` + "\t" + `[]Menuitem ` + "\t`json:\"menuitem\"`" + `
}`,
		},
		{
			name: "Check Array",
			jsonData: `{
				"id": "file",
				"value": "File",
				"menuitem": [
						{"value": "New", "name": "New"},
						{"value": "Open", "name": "New"}
				]
			}`,
			expected: `
type Menuitem struct {
	Value string ` + "\t`json:\"value\"`" + `
	Name  string ` + "\t`json:\"name\"`" + `
}
type Test struct {
	Id      string      ` + "\t`json:\"id\"`" + `
	Value   string      ` + "\t`json:\"value\"`" + `
	Menuitem` + "\t" + `[]Menuitem ` + "\t`json:\"menuitem\"`" + `
}`,
		},
		{
			name: "Check Array",
			jsonData: `{
	"query": "Виктор Иван",
	"count": 7,
	"parts": ["NAME", "SURNAME"]
	}`,
			expected: `
type Generated struct {
        Query   string ` + "\t`json:\"query\"`" + `
        Count   int     ` + "\t`json:\"сount\"`" + `
        Parts   map[string]interface{}  ` + "\t`json:\"parts\"`" + `
`,
		},
	}

	for _, tc := range testCases {
		result, err := struct_convert.StructConvert(tc.jsonData, "Test")
		if err != nil {
			t.Errorf("StructConvert returned an error:\t %v", err)
			continue
		}

		outputPrepare := prepare(result)
		expectedPrepare := prepare(tc.expected)

		if outputPrepare != expectedPrepare {
			t.Errorf("\nInput: \n%s\nStructConvert: \n%s\n Expected:\n%s\n", tc.jsonData, Green+outputPrepare+Reset, Red+expectedPrepare+Reset)
		}
	}
}

func prepare(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "\t", ""), "\n", "")
}
