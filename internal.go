package cgpcli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// getMapAny — получение словаря (map[string]any) от сервера.
func (cli *Cli) getMapAny(cmd string, args ...any) (map[string]any, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// getSliceAny — получение произвольного слайса ([]any) от сервера.
func (cli *Cli) getSliceAny(cmd string, args ...any) ([]any, error) {
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	s, ok := res.([]any)
	if !ok {
		if res == nil {
			return []any{}, nil
		}
		return nil, fmt.Errorf("%s: expected []any, got %T", cmd, res)
	}

	return s, nil
}

// getSliceString — получение списка строк.
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

// setMapAny — строгая отправка словаря данных (SET, UPDATE, MODIFY).
func (cli *Cli) setMapAny(cmd string, subject string, data map[string]any) (err error) {
	if data == nil {
		return fmt.Errorf("%s: data map is nil", cmd)
	}

	if subject != "" {
		err = cli.QueryNV(cmd, Atom(subject), data)
	} else {
		err = cli.QueryNV(cmd, data)
	}
	return
}

// toAnySlice — безопасное приведение к слайсу для вложенных данных (Rules).
func toAnySlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return []any{}
}

func toCGPText(s string) string {
	return strings.ReplaceAll(s, "\n", `\e`)
}

func toGoText(s string) string {
	return strings.ReplaceAll(s, `\e`, "\n")
}

// toInt — безопасное приведение к int для числовых параметров.
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

// toInt64 безопасно приводит различные числовые типы к int64.
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

// toStringMap — конвертирует map[string]any в map[string]string (для ListAccounts и др.).
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
