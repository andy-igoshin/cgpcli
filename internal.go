package cgpcli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// getInt — получение целого числа (int64) от сервера.
func (cli *Cli) getInt(cmd string, args ...any) (int64, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return 0, err
	}

	d, ok := res.(int64)
	if !ok {
		return 0, fmt.Errorf("%s: response is %T, expected string", cmd, res)
	}
	return d, nil
}

// getMapAny — получение словаря (map[string]any) от сервера.
func (cli *Cli) getMapAny(cmd string, args ...any) (map[string]any, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return make(map[string]any), nil
	}

	switch v := res.(type) {
	case map[string]any:
		return v, nil
	case string:
		if v == "" || v == "OK" {
			return make(map[string]any), nil
		}
		return nil, fmt.Errorf("%s: expected map, got non-empty string %q", cmd, v)
	default:
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
}

// getMapString — получение словаря (map[string]string) от сервера.
func (cli *Cli) getMapString(cmd string, args ...any) (map[string]string, error) {
	res, err := cli.getMapAny(cmd, args...)
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

// getSliceAny — получение слайса ([]any) от сервера.
func (cli *Cli) getSliceAny(cmd string, args ...any) ([]any, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	switch v := res.(type) {
	case []any:
		return v, nil
	case string:
		if v == "" || v == "OK" {
			return []any{}, nil
		}
		return nil, fmt.Errorf("%s: expected []any, got string %q", cmd, v)
	case nil:
		return []any{}, nil
	default:
		return nil, fmt.Errorf("%s: expected []any, got %T", cmd, res)
	}
}

// getSliceByte — получение слайса ([]byte) от сервера.
func (cli *Cli) getSliceByte(cmd string, args ...any) ([]byte, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	s, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}

	return s, nil
}

// getSliceString — получение слайса ([]string) от сервера.
func (cli *Cli) getSliceString(cmd string, args ...any) ([]string, error) {
	raw, err := cli.getSliceAny(cmd, args...)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(raw))
	for i, v := range raw {
		s, err := toString(cmd, v)
		if err != nil {
			return nil, err
		}
		res[i] = s
	}
	return res, nil
}

// getString — получение строки (string) от сервера.
func (cli *Cli) getString(cmd string, args ...any) (string, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return "", err
	}

	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: response is %T, expected string", cmd, res)
	}
	return s, nil
}

// parseFileInfoMap converts a raw map from CGP (with ST* keys) into a map of FileInfo structures.
func parseFileInfoMap(raw map[string]any) map[string]*FileInfo {
	res := make(map[string]*FileInfo, len(raw))
	for name, attrs := range raw {
		attrMap, ok := attrs.(map[string]any)
		if !ok {
			continue
		}

		res[name] = &FileInfo{
			Name:     name,
			Size:     toInt64(attrMap["STFileSize"]),
			Created:  toTime(attrMap["STCreated"]),
			Modified: toTime(attrMap["STModified"]),
		}
	}
	return res
}

func replaceSpecChars(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 64)

	for _, symbol := range s {
		switch symbol {
		case '\\':
			// replace \ -> \\   (CGP backslash)
			sb.WriteString(`\\`)

		case '"':
			// replace \" -> \\" (CGP quote)
			sb.WriteString(`\"`)

		case '\n':
			// replace \n -> \\e (CGP End-of-Line)
			sb.WriteString(`\e`)

		case '\t':
			// replace \t -> \\t (CGP Tab)
			sb.WriteString(`\t`)

		case '\r', 0x0000:
			continue

		default:
			sb.WriteRune(symbol)
		}
	}

	return sb.String()
}

// toAnySlice — приведение к слайсу []any.
func toAnySlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return []any{}
}

// toInt — приведение различных числовых типов к int.
func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

// toInt64 приведение различных числовых типов к int64.
func toInt64(v any) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	default:
		return 0
	}
}

func toString(cmd string, v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case time.Time:
		return val.Format("02-Jan-2006 15:04:05"), nil
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(val), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("%s: unexpected value type %T (value: %v)", cmd, v, v)
	}
}

// toStringMap — конвертирует map[string]any в map[string]string.
func toStringMap(cmd string, data map[string]any) (map[string]string, error) {
	res := make(map[string]string, len(data))
	for k, v := range data {
		s, err := toString(cmd, v)
		if err != nil {
			return nil, err
		}
		res[k] = s
	}
	return res, nil
}

// toStringSlice - конвертирует []any в []string
func toStringSlice(cmd string, v any) ([]string, error) {
	arr, ok := v.([]any)
	if !ok {
		if v == nil {
			return []string{}, nil
		}
		return nil, fmt.Errorf("%s: expected []any, got %T", cmd, v)
	}
	res := make([]string, len(arr))
	for i, val := range arr {
		s, err := toString(cmd, val)
		if err != nil {
			return nil, err
		}
		res[i] = s
	}
	return res, nil
}

func toTime(v any) time.Time {
	if t, ok := v.(time.Time); ok {
		return t
	}
	return time.Time{}
}
