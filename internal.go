package cgpcli

import (
	"fmt"
	"strconv"
	"time"
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
		s, err := toString(cmd, v)
		if err != nil {
			return nil, err
		}
		result[i] = s
	}
	return result, nil
}

func (cli *Cli) getString(cmd string, args ...string) (string, error) {
	// Создаем общий слайс: [0] - команда, остальные - аргументы
	params := make([]any, len(args)+1)
	params[0] = cmd
	for i, v := range args {
		params[i+1] = v
	}

	res, err := cli.Query(params...)
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
		// Вместо "тихого" приведения через fmt.Sprintf, возвращаем ошибку
		return "", fmt.Errorf("%s: unexpected value type %T (value: %v)", cmd, v, v)
	}
}

// toStringMap — конвертирует map[string]any в map[string]string (для ListAccounts и др.).
func (cli *Cli) toStringMap(cmd string, data map[string]any) (map[string]string, error) {
	result := make(map[string]string, len(data))
	for k, v := range data {
		s, err := toString(cmd, v)
		if err != nil {
			return nil, err
		}
		result[k] = s
	}
	return result, nil
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
