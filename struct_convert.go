package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"unicode"
)

func main() {
	// jsonStr := `{
	// "query": "Виктор Иван",
	// "count": 7,
	// "parts": [12, 12]
	// }`
	jsonStr := `{
				"id": "file",
				"value": "File",
				"menuitem": [
						{"value": "New", "name": "New"},
						{"value": "Open", "name": "New"}
				]
			}`

	acc, err := StructConvert(jsonStr, "Generated")
	fmt.Println(err)
	fmt.Println(acc)

}

func StructConvert(jsonStr string, nameStruct string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return "", err
	}
	acc, err := generateStruct(result, nameStruct)
	if err != nil {
		return "", err
	}
	return acc, nil
}

func generateStruct(data map[string]interface{}, nameStruct string) (string, error) {
	var acc string
	var typeValue string
	var nameField string
	var mapInterface map[string]interface{}
	var firstElem interface{}

	nameStructToUpper := toUpperFirstLetter(nameStruct)

	acc += fmt.Sprintf("type %s struct {\n", nameStructToUpper)

	for key, value := range data {
		nameField = toUpperFirstLetter(key)
		typeValue = goType(value)

		if typeValue == "map[string]interface{}" || typeValue == "[]interface{}" {
			if typeValue == "map[string]interface{}" {
				var ok bool
				mapInterface, ok = value.(map[string]interface{})
				if !ok {
					return "", errors.New("Error: []interface{} -> map[string]interface{}")
				}

				acc += buildLine(nameField, nameField, key)
			}

			if typeValue == "[]interface{}" {
				arrInterface, ok := value.([]interface{})
				if !ok {
					return "", errors.New("Error: interface{}  ->  []interface{}")
				}

				if len(arrInterface) == 0 {
					firstElem = arrInterface
				} else {
					firstElem = arrInterface[0]
				}

				fmt.Println(goType(arrInterface))

				// firstElemType := goType(firstElem)
				// if goType(firstElem) == "string" {
				// 	acc += buildLine(nameField, "map[string]interface{}", key)
				// 	continue
				// }

				mapInterface, ok = firstElem.(map[string]interface{})
				if !ok {
					return "", errors.New("Error: []interface{} -> map[string]interface{}")
				}

				arrayNameField := fmt.Sprintf("[]%s", nameField)
				acc += buildLine(nameField, arrayNameField, key)
			}

			structStr, err := generateStruct(mapInterface, key)
			if err != nil {
				return "", err
			}
			acc = structStr + acc
			continue
		}
		acc += buildLine(nameField, typeValue, key)
	}

	acc += "}\n"
	return acc, nil
}

func buildLine(nameField string, typeValue string, key string) string {
	return fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\n", nameField, typeValue, key)
}

func goType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case float64:
		return "int"
	case bool:
		return "bool"
	case map[string]interface{}:
		return "map[string]interface{}"
	case []interface{}:
		return "[]interface{}"
	default:
		return "interface{}"
	}
}

func toUpperFirstLetter(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
