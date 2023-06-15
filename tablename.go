package neo4jx

import (
	"fmt"
	"reflect"
)

func ToLabelName(data interface{}) (string, error) {
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("Data type must be struct")
	}
	return t.Name(), nil
}
