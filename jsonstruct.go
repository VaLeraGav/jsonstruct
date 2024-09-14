package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"unicode"
)

const (
	stringType     = "string"
	float64Type    = "int"
	boolType       = "bool"
	mapIntType     = "map[string]interface{}"
	arrIntType     = "[]interface{}"
	defaultIntType = "interface{}"
)

var (
	ErrConvertInterfaceToMap   = errors.New("error format : []interface{} -> map[string]interface{}")
	ErrConvertInterfaceToArray = errors.New("error format : interface{}  ->  []interface{}")
)

func Convert(jsonStr string, nameStruct string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}
	acc, err := generateStruct(result, nameStruct)
	if err != nil {
		return "", fmt.Errorf("generateStruct: %w", err)
	}
	return acc, nil
}

func WriteFile(filename string, strStruct string) error {
	fileNameAbs, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileNameAbs, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	packageName := getPackageName(fileNameAbs)

	strStruct = fmt.Sprintf("package %s \n\n %s", packageName, strStruct)

	fmt.Fprintf(f, "%s\n", strStruct)
	return nil
}

func getPackageName(fileNameAbs string) string {
	dir := filepath.Dir(fileNameAbs)
	return filepath.Base(dir)
}

func generateStruct(data map[string]interface{}, nameStruct string) (string, error) {
	var acc string
	var typeValue string
	var nameField string
	var mapInterface map[string]interface{}

	nameStructToUpper := toUpperFirstLetter(nameStruct)

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	acc += fmt.Sprintf("type %s struct {\n", nameStructToUpper)

	for _, key := range keys {
		value := data[key]
		nameField = toUpperFirstLetter(key)
		typeValue = getType(value)

		if typeValue == mapIntType || typeValue == arrIntType {
			if typeValue == mapIntType {
				var ok bool
				mapInterface, ok = value.(map[string]interface{})
				if !ok {
					return "", fmt.Errorf("mapInterface: %w", ErrConvertInterfaceToMap)
				}

				acc += buildLine(nameField, nameField, key)
			}

			if typeValue == arrIntType {
				arrInterface, ok := value.([]interface{})
				if !ok {
					return "", fmt.Errorf("arrInterface: %w", ErrConvertInterfaceToArray)
				}

				if len(arrInterface) == 0 {
					acc += buildLine(nameField, arrIntType, key)
					continue
				}

				firstElement := getType(arrInterface[0])
				if firstElement != mapIntType {
					acc += buildLine(nameField, arrIntType, key)
					continue
				}

				// TODO:ошибка если длинный массив с повторами
				idLongestArray := longestMap(arrInterface)

				firstElem := arrInterface[idLongestArray]

				mapInterface, ok = firstElem.(map[string]interface{})
				if !ok {
					return "", fmt.Errorf("mapInterface from firstElem: %w", ErrConvertInterfaceToMap)
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

func longestMap(maps []interface{}) int {
	longestLength := 0
	idLongestMap := 0

	for id, item := range maps {
		if m, ok := item.(map[string]interface{}); ok {
			if len(m) > longestLength {
				longestLength = len(m)
				idLongestMap = id
			}
		}
	}

	return idLongestMap
}

func buildLine(nameField string, typeValue string, key string) string {
	return fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\n", nameField, typeValue, key)
}

func getType(value interface{}) string {
	switch value.(type) {
	case string:
		return stringType
	case float64:
		return float64Type
	case bool:
		return boolType
	case map[string]interface{}:
		return mapIntType
	case []interface{}:
		return arrIntType
	default:
		return defaultIntType
	}
}

func toUpperFirstLetter(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
