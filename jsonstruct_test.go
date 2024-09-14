package jsonstruct_test

import (
	"strings"
	"testing"

	"github.com/VaLeraGav/jsonstruct"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

func TestJsonStruct(t *testing.T) {
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
			expected: `
			type Test struct {
				Count int    ` + getJson("count") + `
				Query string ` + getJson("query") + `
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
				Name    string  ` + getJson("name") + `
				Value string 	` + getJson("value") + `
			}
			type Test struct {
				Id      	string     	` + getJson("id") + `
				Menuitem 	[]Menuitem 	` + getJson("menuitem") + `
				Value  		 string     ` + getJson("value") + `
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
			type Menuitem2 struct {
				Value   string  		` + getJson("value") + `
			}
			type Menuitem1 struct {
				Value   string  		` + getJson("value") + `
			}
			type Menuitem struct {
				Menuitem1   Menuitem1 	` + getJson("menuitem1") + `
				Menuitem2   Menuitem2 	` + getJson("menuitem2") + `
				Value   	string  	` + getJson("value") + `
			}
			type Test struct {
				Id      	string  	` + getJson("id") + `
				Menuitem  	Menuitem  	` + getJson("menuitem") + `
			}`,
		},
		{
			name: "Check Array 1",
			jsonData: `
			{
				"query": "Виктор Иван",
				"count": 7,
				"parts": ["NAME", "SURNAME"]
			}`,
			expected: `
			type Test struct {
				Count   int     		` + getJson("count") + `
				Parts   []interface{}   ` + getJson("parts") + `
				Query   string  		` + getJson("query") + `
			}`,
		},
		{
			name: "Check Array 2",
			jsonData: `
			{
				"query": "Виктор Иван",
				"count": 7,
				"parts": []
			}`,
			expected: `
			type Test struct {
				Count   int            	` + getJson("count") + `
				Parts   []interface{}  	` + getJson("parts") + `
				Query   string  		` + getJson("query") + `
			}`,
		},
	}

	for _, tc := range testCases {
		result, err := jsonstruct.Convert(tc.jsonData, "Test")
		if err != nil {
			t.Errorf("Convert returned an error:\t %v", err)
			continue
		}

		outputPrepare := prepare(result)
		expectedPrepare := prepare(tc.expected)

		if outputPrepare != expectedPrepare {
			t.Errorf("\nName: \n%s\nInput: \n%s\nConvert: \n%s\n Expected:\n%s\n \n==================================",
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
