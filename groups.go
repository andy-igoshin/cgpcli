// # Section: Group Administration
package cgpcli

import (
	"fmt"
)

// GroupResult объединяет настройки и список участников
type GroupResult struct {
	Settings map[string]any
	Members  []string
}

// CreateGroup создает новую группу.
func (cli *Cli) CreateGroup(group string, settings map[string]any) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	args := make([]any, 0, 3)
	args = append(args, "CREATEGROUP", group)
	if settings != nil {
		args = append(args, settings)
	}
	_, err := cli.Query(args...)
	return err
}

// DeleteGroup удаляет группу.
func (cli *Cli) DeleteGroup(group string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	_, err := cli.Query("DELETEGROUP", group)
	return err
}

// GetGroup возвращает настройки группы и список её участников (Members).
func (cli *Cli) GetGroup(group string) (*GroupResult, error) {
	if group == "" {
		return nil, fmt.Errorf("group name is required")
	}
	res, err := cli.Query("GETGROUP", group)
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETGROUP: expected map, got %T", res)
	}

	result := &GroupResult{
		Settings: make(map[string]any),
		Members:  []string{},
	}

	for k, v := range dict {
		if k == "Members" {
			m, err := cli.toStringSlice2(v)
			if err == nil {
				result.Members = m
			}
			continue
		}
		result.Settings[k] = v
	}
	return result, nil
}

// IsGroupExists проверяет, существует ли группа.
func (cli *Cli) IsGroupExists(group string) (bool, error) {
	if group == "" {
		return false, nil
	}
	_, err := cli.Query("GETGROUP", group)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ListGroups возвращает список имен всех групп в указанном домене.
func (cli *Cli) ListGroups(domain string) ([]string, error) {
	args := make([]any, 0, 2)
	args = append(args, "LISTGROUPS")
	if domain != "" {
		args = append(args, domain)
	}
	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}
	return cli.toStringSlice2(res)
}

// RenameGroup переименовывает группу или переносит её в другой домен.
func (cli *Cli) RenameGroup(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new group names are required")
	}
	_, err := cli.Query("RENAMEGROUP", oldName, "INTO", newName)
	return err
}

// SetGroup заменяет все настройки и список участников группы.
func (cli *Cli) SetGroup(group string, settings map[string]any, members []string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}

	fullData := make(map[string]any)
	if settings != nil {
		for k, v := range settings {
			fullData[k] = v
		}
	}

	if members != nil {
		anyMembers := make([]any, len(members))
		for i, m := range members {
			anyMembers[i] = m
		}
		fullData["Members"] = anyMembers
	}

	_, err := cli.Query("SETGROUP", group, fullData)
	return err
}
