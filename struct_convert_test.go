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

// TODO: так как используется range порядок каждый раз меняется

func TestStructConvert(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		expected string
	}{
		{
			name: "Easy",
			jsonData: `
					{
						"query": "Test",
						"count": 7
					}`,
			expected: `type Test struct {
			Query string ` + getJson("query") + `
			Count int    ` + getJson("count") + `
		}`,
		},
		{
			name: "Multidimensional object",
			jsonData: `{
						"id": "file",
						"value": "File",
						"menuitem": [
								{"value": "New"},
								{"value": "Open", "name": "New"}
						]
					}`,
			expected: `
		type Menuitem struct {
			Value string ` + getJson("value") + `
			Name  string ` + getJson("name") + `
		}
		type Test struct {
			Id      string     ` + getJson("id") + `
			Value   string     ` + getJson("value") + `
			Menuitem` + "\t" + `[]Menuitem ` + getJson("menuitem") + `
		}`,
		},
		{
			name: "Three levels of nesting an object",
			jsonData: `
			{
				"id": "file",
				"menuitem":	{
						"value": "New",
						"menuitem1":	{
							"value": "New"
						},
						"menuitem2":	{
							"value": "New"
						}
					}
			}`,
			expected: `
type Menuitem1 struct {
  Value   string  ` + getJson("value") + `
}
type Menuitem2 struct {
  Value   string  ` + getJson("value") + `
}
type Menuitem struct {
  Menuitem2       Menuitem2     ` + getJson("menuitem2") + `
  Value   string  							` + getJson("value") + `
  Menuitem1       Menuitem1     ` + getJson("menuitem1") + `
}
type Test struct {
  Id      string  				  ` + getJson("id") + `
  Menuitem        Menuitem  ` + getJson("menuitem") + `
}
			`,
		},
		{
			name: "Check Array",
			jsonData: `
			{
				"query": "Виктор Иван",
				"count": 7,
				"parts": ["NAME", "SURNAME"]
			}`,
			expected: `
type Test struct {
	Query   string  ` + getJson("query") + `
	Count   int     ` + getJson("count") + `
	Parts   []interface{}   ` + getJson("parts") + `
}`,
		},
		{
			name: "Check Array",
			jsonData: `
			{
				"query": "Виктор Иван",
				"count": 7,
				"parts": []
			}`,
			expected: `
type Test struct {
	Query   string  ` + getJson("query") + `
	Count   int     ` + getJson("count") + `
	Parts   []interface{}   ` + getJson("parts") + `
}`,
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
			t.Errorf("\nName: \n%s\nInput: \n%s\nStructConvert: \n%s\n Expected:\n%s\n \n==================================",
				tc.name, tc.jsonData, Green+result+Reset, Red+tc.expected+Reset,
			)
		}
	}
}

func prepare(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "\t", ""), "\n", "")
}

func getJson(value string) string {
	return "\t`json:\"" + value + "\"`"
}
