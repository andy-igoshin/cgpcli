package cgpcli

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// PrettyPrint выводит структуру данных в удобном для чтения виде.
func PrettyPrint(v any) {
	prettyPrintRecursive(v, 0)
}

func prettyPrintRecursive(v any, indent int) {
	prefix := strings.Repeat("  ", indent)

	if v == nil {
		fmt.Printf("nil\n")
		return
	}

	switch val := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Заменяем "interface {}" на "any" для красоты
		tStr := strings.ReplaceAll(fmt.Sprintf("%T", val), "interface {}", "any")
		fmt.Printf("(%s len=%d):\n", tStr, len(val))

		for _, k := range keys {
			fmt.Printf("%s  %s: ", prefix, k)
			prettyPrintRecursive(val[k], indent+1)
		}

	case []any:
		tStr := strings.ReplaceAll(fmt.Sprintf("%T", val), "interface {}", "any")
		fmt.Printf("(%s len=%d):\n", tStr, len(val))

		for _, item := range val {
			fmt.Printf("%s  - ", prefix)
			prettyPrintRecursive(item, indent+1)
		}

	case string:
		if val == "" {
			fmt.Printf("(empty) (string)\n")
		} else {
			fmt.Printf("%s (string)\n", val)
		}

	case time.Time:
		fmt.Printf("%s (time.Time)\n", val.Format("2006-01-02 15:04:05"))

	default:
		tStr := strings.ReplaceAll(fmt.Sprintf("%T", v), "interface {}", "any")
		fmt.Printf("%v (%s)\n", v, tStr)
	}
}
