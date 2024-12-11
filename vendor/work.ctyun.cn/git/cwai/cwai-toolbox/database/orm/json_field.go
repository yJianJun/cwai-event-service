package orm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

// Value struct to json
func Value(field interface{}) (driver.Value, error) {
	b, err := json.Marshal(field)
	return string(b), err
}

// Scan json to struct
func Scan(field interface{}, value interface{}) error {
	switch value.(type) {
	case string:
		return json.Unmarshal([]byte(value.(string)), field)
	case []byte:
		return json.Unmarshal(value.([]byte), field)
	}
	return fmt.Errorf("Scan unknown type: %s", reflect.TypeOf(value))
}
