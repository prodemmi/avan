package table

import (
	"reflect"
	"regexp"
	"strings"
)

func GetProviderConfigType(input int64) string {
	switch input {
	case 0:
		return "SMS"
	case 1:
		return "EMAIL"
	case 2:
		return "TELEGRAM"
	case 3:
		return "SLACK"
	case 4:
		return "WEBHOOK"
	}

	return "UNKNOWN"
}

func ToTitleCase(input string) string {
	// Insert space before all caps (like 'camelCase' -> 'camel Case')
	re := regexp.MustCompile("([a-z])([A-Z])")
	title := re.ReplaceAllString(input, "$1 $2")
	// Capitalize the first letter of each word
	return strings.Title(title)
}

func ToTableMapList[T any](inputs []T) []map[string]any {
	var result []map[string]any
	for _, item := range inputs {
		// Create a new map for each struct item
		itemMap := make(map[string]any)

		// Use reflection to get the fields and values of the struct
		val := reflect.ValueOf(item)
		typ := reflect.TypeOf(item)

		// Ensure it's a struct
		if val.Kind() == reflect.Struct {
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				fieldName := typ.Field(i).Name
				itemMap[fieldName] = field.Interface()
			}
		}

		// Add the map to the result slice
		result = append(result, itemMap)
	}
	return result
}
