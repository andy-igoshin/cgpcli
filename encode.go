package cgpcli

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DataBlock []byte

// Marshal преобразует типы Go в формат строк CommuniGate Pro CLI.
// По умолчанию вызывается в контексте значения.
func Marshal(v any) string {
	return marshal(v, false)
}

// marshal — внутренняя функция, учитывающая контекст (ключ словаря или значение).
func marshal(v any, asKey bool) string {
	switch val := v.(type) {
	case DataBlock:
		// Формат Datablock: [BASE64]
		return "[" + string(val) + "]"

	case []byte:
		// Бинарные данные передаются как #[HEX]
		return fmt.Sprintf("#[%X]", val)

	case string:
		// 1. Если это КЛЮЧ: разрешаем без кавычек только простые идентификаторы (буквы, цифры, _)
		if asKey && isCGPIdentifier(val) {
			return val
		}
		// 2. Если это ЗНАЧЕНИЕ: без кавычек только системные литералы (unlimited, YES) и квоты (10M)
		if !asKey && isCGPSystemLiteral(val) {
			return val
		}
		// Во всех остальных случаях (и для строк-значений) — используем кавычки
		return fmt.Sprintf("%q", val)

	case int, int32, int64:
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
		// Формат времени для CLI команд
		return fmt.Sprintf("%q", val.Format("02-Jan-2006 15:04:05 -0700"))

	case map[string]any:
		var b strings.Builder
		b.WriteString("{ ")

		// Сортируем ключи для детерминированного вывода (важно для тестов и логов)
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			b.WriteString(marshal(k, true)) // Ключ
			b.WriteString(" = ")
			b.WriteString(marshal(val[k], false)) // Значение
			b.WriteString("; ")
		}
		b.WriteString("}")
		return b.String()

	case []any:
		var b strings.Builder
		b.WriteString("( ")
		for i, v := range val {
			b.WriteString(marshal(v, false)) // Элементы массива — это значения
			if i < len(val)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(" )")
		return b.String()

	case map[string]string:
		// Приводим к map[string]any для переиспользования логики
		m := make(map[string]any, len(val))
		for k, v := range val {
			m[k] = v
		}
		return marshal(m, false)

	case []string:
		// Приводим к []any
		s := make([]any, len(val))
		for i, v := range val {
			s[i] = v
		}
		return marshal(s, false)

	default:
		if val == nil {
			return `""`
		}
		return fmt.Sprintf("%v", val)
	}
}

// isCGPIdentifier проверяет, может ли строка быть КЛЮЧОМ без кавычек.
// Разрешены только буквы, цифры и подчеркивание. Точки и тире требуют кавычек.
func isCGPIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

// isCGPSystemLiteral проверяет, может ли строка быть ЗНАЧЕНИЕМ без кавычек.
// Разрешены системные константы и числовые квоты с суффиксами.
func isCGPSystemLiteral(s string) bool {
	if s == "" {
		return false
	}

	switch strings.ToLower(s) {
    case "default", "unlimited", "yes", "no", "none", "any":
        return true
    }

	// Проверка на число или квоту (например, 10G)
	isNumeric := true
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		// Суффикс допустим только в конце строки
		if i == len(s)-1 && strings.ContainsRune("KMGTkmgt", r) {
			continue
		}
		isNumeric = false
		break
	}
	return isNumeric
}

// formatCGPSize преобразует байты в формат CGP (1K, 5M, 10G) или "unlimited".
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
