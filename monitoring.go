// # Section: Monitoring
package cgpcli

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// StatsCapacity определяет начальную емкость слайса для списка статистики.
// Значение 2048 покрывает стандартный набор OID большинства инсталляций без переаллокаций.
const statsCapacity = 2048

// ErrStatNotFound — внутренняя ошибка для идентификации конца списка или отсутствия элемента.
var ErrStatNotFound = errors.New("statistics element is not found")

// GetDialogInfo возвращает данные о состоянии Signal Dialog по его ID.
func (cli *Cli) GetDialogInfo(dialogID int64) (map[string]any, error) {
	return cli.getMapAny("GETDIALOGINFO", dialogID)
}

// GetNextStatName возвращает имя следующего элемента статистики.
func (cli *Cli) GetNextStatName(objectID string) (string, error) {
	const cmd = "GETNEXTSTATNAME"
	res, err := cli.Query(cmd, objectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", io.EOF
		}
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
	}
	return s, nil
}

// GetStatElement возвращает текущее значение элемента статистики (SNMP).
func (cli *Cli) GetStatElement(objectID string) (any, error) {
	if objectID == "" {
		return nil, fmt.Errorf("objectID is required")
	}
	return cli.Query("GETSTATELEMENT", objectID)
}

// ListAllStats собирает все доступные имена элементов статистики в один слайс.
func (cli *Cli) ListAllStats() ([]string, error) {
	stats := make([]string, 0, statsCapacity)
	var current string

	for {
		next, err := cli.GetNextStatName(current)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		stats = append(stats, next)
		current = next
	}
	return stats, nil
}

// SetStatElement обновляет значение "Custom" элемента статистики.
func (cli *Cli) SetStatElement(objectID, op string, val int64) error {
	if objectID == "" {
		return fmt.Errorf("objectID is required")
	}
	if op != "INC" && op != "SET" {
		return fmt.Errorf("operation must be INC or SET")
	}
	return cli.QueryNV("SETSTATELEMENT", objectID, op, val)
}

// Shutdown инициирует остановку сервера CommuniGate Pro.
func (cli *Cli) Shutdown() error {
	return cli.QueryNV("SHUTDOWN")
}
