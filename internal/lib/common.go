package lib

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/exceptionate/gitaur/internal/ui"
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

func PromptWithDefault(label string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	if defaultValue != "" {
		fmt.Print(ui.Label.Render(fmt.Sprintf("%s [%s]: ", label, defaultValue)))
	} else {
		fmt.Print(ui.Label.Render(fmt.Sprintf("%s: ", label)))
	}

	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	return input
}
