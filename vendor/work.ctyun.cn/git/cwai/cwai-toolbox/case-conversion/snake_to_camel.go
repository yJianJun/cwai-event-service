package case_conversion

import (
	"regexp"
	"strings"
)

var link = regexp.MustCompile("(^[A-Za-z0-9])|_([A-Za-z0-9])")
var firstLargeCase = regexp.MustCompile("^[A-Z]")

// SnakeToCamelBig converts `snake_case` to `BigCamelCase`.
//
// WARNING: `database_id` will be converted to `DatabaseId` instead of `DatabaseID`.
// Be aware of that before using.
//
// For `smallCamelCase`, see SnakeToCamelSmall.
func SnakeToCamelBig(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, "_", ""))
	})
}

// SnakeToCamelSmall converts `snake_case` to `smallCamelCase`.
//
// WARNING: `database_id` will be converted to `databaseId` instead of `databaseID`.
// Be aware of that before using.
//
// For `CamelCase`, see SnakeToCamel.
func SnakeToCamelSmall(str string) string {
	return firstLargeCase.ReplaceAllStringFunc(SnakeToCamelBig(str), func(s string) string {
		return strings.ToLower(s)
	})
}
