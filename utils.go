// # Utilities
//
// The Utilities section provides developer-focused functions for data visualization
// and debugging. When working with the CommuniGate Pro CLI, responses often
// arrive as deeply nested structures of maps and slices; these utilities
// simplify the inspection of such data.
//
// Key capabilities include:
//   - Recursive Visualization: [PrettyPrint] allows for a clean, indented
//     tree-view output of any data structure returned by the CLI.
//   - Sorted Inspection: maps are automatically sorted by key during output,
//     ensuring consistent and predictable debugging sessions.
//   - Type Metadata: each value in the pretty-printed output is annotated
//     with its Go type information, making it easier to map CLI responses
//     to application logic.
//   - Formatted Time: date and time objects are rendered in a human-readable
//     format, bypassing the default verbose [time.Time] representation.
//
// These functions are intended for use during the development and
// troubleshooting phases to verify the integrity of data received
// from the server.
package cgpcli

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// PrettyPrint prints a data structure to the standard output in a
// human-readable format.
//
// It recursively traverses maps, slices, and basic types, providing
// indentation and type annotations. This is especially useful for
// inspecting complex dictionaries returned by the CommuniGate Pro CLI.
//
// Parameters:
//   - v: the data structure to be printed.
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
