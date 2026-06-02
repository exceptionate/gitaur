package lib

import (
	"reflect"
	"strings"
)

func HasField[T any](fieldName string) bool {
	t := reflect.TypeFor[T]()
	// If a pointer type is passed, deref it.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for field := range t.Fields() {
		if strings.EqualFold(field.Name, fieldName) {
			return true
		}
	}
	return false
}
