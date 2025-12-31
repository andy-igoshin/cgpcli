package cgpcli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Marshal(v any) string {
	switch val := v.(type) {
	case string:
		if isCGPLiteral(val) {
			return val
		}
		return fmt.Sprintf("%q", val)

	case int, int64, int32:
		var n int64
		switch i := val.(type) {
		case int:
			n = int64(i)
		case int32:
			n = int64(i)
		case int64:
			n = i
		}
		return formatCGPSize(n)

	case bool:
		if val {
			return "YES"
		}
		return "NO"

	case time.Time:
		if val.IsZero() {
			return `""`
		}
		return fmt.Sprintf("%q", val.Format("02-Jan-2006 15:04:05 -0700"))

	case []byte:
		return fmt.Sprintf("#[%X]", val)

	case map[string]any:
		var b strings.Builder
		b.WriteString("{ ")
		for k, v := range val {
			b.WriteString(k)
			b.WriteString(" = ")
			b.WriteString(Marshal(v))
			b.WriteString("; ")
		}
		b.WriteString("}")
		return b.String()

	case []any:
		var b strings.Builder
		b.WriteString("( ")
		for i, v := range val {
			b.WriteString(Marshal(v))
			if i < len(val)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(" )")
		return b.String()

	default:
		if val == nil {
			return `""`
		}
		return fmt.Sprintf("%v", val)
	}
}

func isCGPLiteral(s string) bool {
	if s == "unlimited" || s == "default" || s == "none" {
		return true
	}
	matchQuota := true
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		if i == len(s)-1 && strings.ContainsRune("KMGTkmgt", r) {
			continue
		}
		matchQuota = false
		break
	}
	return matchQuota && len(s) > 0
}

func formatCGPSize(v int64) string {
	if v == -1 {
		return "unlimited"
	}
	if v == 0 {
		return "0"
	}
	units := []struct {
		suffix string
		value  int64
	}{
		{"T", 1024 * 1024 * 1024 * 1024},
		{"G", 1024 * 1024 * 1024},
		{"M", 1024 * 1024},
		{"K", 1024},
	}
	for _, u := range units {
		if v%u.value == 0 {
			return fmt.Sprintf("%d%s", v/u.value, u.suffix)
		}
	}
	return strconv.FormatInt(v, 10)
}
