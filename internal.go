package cgpcli

import (
	"fmt"
	"strconv"
)

// getMapAny — строгое получение словаря (map[string]any) от сервера.
func (cli *Cli) getMapAny(cmd string, subject string) (map[string]any, error) {
	var res any
	var err error

	if subject != "" {
		res, err = cli.Query(cmd, Atom(subject))
	} else {
		res, err = cli.Query(cmd)
	}

	if err != nil {
		return nil, err
	}

	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map[string]any, got %T", cmd, res)
	}

	return m, nil
}

// setMapAny — строгая отправка словаря данных (SET, UPDATE, MODIFY).
func (cli *Cli) setMapAny(cmd string, subject string, data map[string]any) error {
	if data == nil {
		return fmt.Errorf("%s: data map is nil", cmd)
	}

	var err error
	if subject != "" {
		_, err = cli.Query(cmd, Atom(subject), data)
	} else {
		_, err = cli.Query(cmd, data)
	}
	return err
}

// getSliceAny — строгое получение списка ([]any) от сервера.
func (cli *Cli) getSliceAny(cmd string, subject string) ([]any, error) {
	var res any
	var err error

	if subject != "" {
		res, err = cli.Query(cmd, Atom(subject))
	} else {
		res, err = cli.Query(cmd)
	}

	if err != nil {
		return nil, err
	}

	s, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected []any, got %T", cmd, res)
	}

	return s, nil
}

// getSliceString — получение списка строк с поддержкой приведения чисел через strconv.
func (cli *Cli) getSliceString(cmd string, subject string) ([]string, error) {
	raw, err := cli.getSliceAny(cmd, subject)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(raw))
	for i, v := range raw {
		s, err := cli.valToString(cmd, v)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}
	return result, nil
}

func parseCode(b []byte) CliCode {
	code := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			break
		}
		code = code*10 + int(c-'0')
	}
	return CliCode(code)
}

// toAnySlice — безопасное приведение к слайсу для вложенных данных (Rules).
func (cli *Cli) toAnySlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return []any{}
}

// toInt — безопасное приведение к int для числовых параметров.
func (cli *Cli) toInt(v any) int {
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

// toStringMap — конвертирует map[string]any в map[string]string (для ListAccounts и др.).
func (cli *Cli) toStringMap(cmd string, data map[string]any) (map[string]string, error) {
	result := make(map[string]string, len(data))
	for k, v := range data {
		s, err := cli.valToString(cmd, v)
		if err != nil {
			return nil, err
		}
		result[k] = s
	}
	return result, nil
}

// toStringSlice — вспомогательный метод для эффективного приведения типов
func (cli *Cli) toStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return []string{}
	}

	list := make([]string, len(raw))
	for i, item := range raw {
		if s, ok := item.(string); ok {
			list[i] = s
		} else {
			list[i] = fmt.Sprint(item)
		}
	}
	return list
}

// toStringSlice2 — временный строгий конвертер для точечного рефакторинга.
func (cli *Cli) toStringSlice2(v any) ([]string, error) {
	raw, ok := v.([]any)
	if !ok {
		if v == nil {
			return []string{}, nil
		}
		return nil, fmt.Errorf("expected []any, got %T", v)
	}

	result := make([]string, len(raw))
	for i, item := range raw {
		s, err := cli.valToString("", item)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}
	return result, nil
}

// valToString — системное приведение значения к строке (string или числа через strconv).
func (cli *Cli) valToString(cmd string, v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("%s: unexpected value type %T (value: %v)", cmd, v, v)
	}
}
