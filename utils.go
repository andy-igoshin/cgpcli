package cgpcli

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// PrettyPrint выводит структуру данных в удобном для чтения виде.
func PrettyPrint(v any) {
	prettyPrintRecursive(v, 0)
}

func prettyPrintRecursive(v any, indent int) {
	prefix := strings.Repeat("  ", indent)
	val := reflect.ValueOf(v)

	if !val.IsValid() {
		fmt.Print("<invalid>\n")
		return
	}

	// Разрываем указатели и интерфейсы
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			fmt.Print("<nil>\n")
			return
		}
		val = val.Elem()
	}

	// Спец-обработка для типов, реализующих Stringer (например, time.Time)
	if stringer, ok := v.(fmt.Stringer); ok {
		fmt.Printf("%s\n", stringer.String())
		return
	}

	// Обработка слайса байт как строки (наш DataBlock из парсера)
	if val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8 {
		fmt.Printf("[]byte: %q\n", val.Bytes())
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		fmt.Printf("(struct %s):\n", val.Type().Name())
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			fmt.Printf("%s  %s: ", prefix, field.Name)
			prettyPrintRecursive(val.Field(i).Interface(), indent+2)
		}

	case reflect.Map:
		if val.Len() == 0 {
			fmt.Print("{}\n")
			return
		}
		fmt.Printf("(map len=%d):\n", val.Len())

		// Сбор и сортировка ключей
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
		})

		for _, key := range keys {
			fmt.Printf("%s  %v: ", prefix, key)
			prettyPrintRecursive(val.MapIndex(key).Interface(), indent+2)
		}

	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			fmt.Print("[]\n")
			return
		}
		fmt.Printf("(slice len=%d):\n", val.Len())
		for i := 0; i < val.Len(); i++ {
			fmt.Printf("%s  - ", prefix)
			prettyPrintRecursive(val.Index(i).Interface(), indent+2)
		}

	default:
		fmt.Printf("%v\n", val.Interface())
	}
}
