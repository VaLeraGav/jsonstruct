# jsonstruct

Library convert json to struct.

```go
func main() {
	jsonStr := `{
	"query": "Виктор Иван",
	"count": 7,
	"parts": [{"sdfsdf" : 12}]
	}`

	acc, err := jsonstruct.Convert(jsonStr, "Name")
	if err != nil {
		return "", fmt.Errorf("Convert: %w", err)
	}

	err = jsonstruct.WriteFile("./name.go", acc)
	if err != nil {
		return "", fmt.Errorf("WriteFile: %w", err)
	}
}
```

`./name` file will be created.go to the root

```go
type Parts struct {
        Count int `json:"count"`
}
type Name struct {
        Count   int     `json:"count"`
        Parts   []Parts `json:"parts"`
        Query   string  `json:"query"`
}
```
