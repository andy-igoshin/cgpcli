package cgpcli

import (
	"fmt"
	"strings"
)

// GroupResult объединяет настройки и список участников
type GroupResult struct {
	Settings map[string]any
	Members  []string
}

// -----------------------------------------------------------------------------
// CLI Methods
// -----------------------------------------------------------------------------

// CreateGroup создает новую группу.
func (cli *Cli) CreateGroup(group string, settings map[string]any) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	parts := []string{"CREATEGROUP", Marshal(group)}
	if settings != nil {
		parts = append(parts, Marshal(settings))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// DeleteGroup удаляет группу.
func (cli *Cli) DeleteGroup(group string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	_, err := cli.query("DELETEGROUP " + Marshal(group))
	return err
}

// GetGroup возвращает настройки группы и список её участников (Members).
func (cli *Cli) GetGroup(group string) (*GroupResult, error) {
	if group == "" {
		return nil, fmt.Errorf("group name is required")
	}
	res, err := cli.query("GETGROUP " + Marshal(group))
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETGROUP: expected map, got %T", res)
	}

	result := &GroupResult{
		Settings: make(map[string]any),
		Members:  make([]string, 0),
	}

	for k, v := range dict {
		if k == "Members" {
			if mRaw, ok := v.([]any); ok {
				for _, m := range mRaw {
					result.Members = append(result.Members, fmt.Sprint(m))
				}
			}
			continue
		}
		result.Settings[k] = v
	}
	return result, nil
}

// ListGroups возвращает список имен всех групп в указанном домене.
func (cli *Cli) ListGroups(domain string) ([]string, error) {
	parts := []string{"LISTGROUPS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}

	var res any
	var err error
	if len(parts) > 1 {
		res, err = cli.query(strings.Join(parts, " "))
	} else {
		res, err = cli.query(parts[0])
	}

	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}

	list := make([]string, len(raw))
	for i, v := range raw {
		list[i] = fmt.Sprint(v)
	}
	return list, nil
}

// RenameGroup переименовывает группу или переносит её в другой домен.
func (cli *Cli) RenameGroup(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new group names are required")
	}
	parts := []string{"RENAMEGROUP", Marshal(oldName), "INTO", Marshal(newName)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetGroup заменяет все настройки и список участников группы.
func (cli *Cli) SetGroup(group string, settings map[string]any, members []string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}

	fullData := make(map[string]any)
	for k, v := range settings {
		fullData[k] = v
	}

	if members != nil {
		anyMembers := make([]any, len(members))
		for i, m := range members {
			anyMembers[i] = m
		}
		fullData["Members"] = anyMembers
	}

	parts := []string{"SETGROUP", Marshal(group), Marshal(fullData)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// -----------------------------------------------------------------------------
// Extended API / Helpers
// -----------------------------------------------------------------------------

// IsGroupExists проверяет, существует ли группа.
func (cli *Cli) IsGroupExists(group string) (bool, error) {
	if group == "" {
		return false, nil
	}
	_, err := cli.query("GETGROUP " + Marshal(group))
	if err != nil {
		return false, nil
	}
	return true, nil
}
