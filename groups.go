// # Group Administration
//
// The Groups section provides functionality for managing CommuniGate Pro groups.
// Groups serve as versatile objects for mail distribution and signal routing.
//
// Key capabilities include:
//   - Lifecycle Management: Creating, deleting, and renaming groups across
//     the environment.
//   - Membership & Configuration: Using [Cli.GetGroup] and [Cli.SetGroup] to manage
//     both group-specific settings and the list of group members simultaneously.
//   - Discovery: Listing all groups within a specific domain and checking
//     for the existence of individual groups.
//   - Unified Data Access: The [GroupResult] structure simplifies interaction
//     by separating administrative settings from the member list.
//
// Use these methods to automate the management of distribution lists and
// departmental communication channels.
package cgpcli

import (
	"fmt"
)

// GroupResult represents the combined data of a CommuniGate Pro group,
// separating operational settings from the actual member list.
type GroupResult struct {
	Settings map[string]any // Administrative settings and attributes of the group.
	Members  []string       // List of addresses or object names belonging to the group.
}

// CreateGroup создает новую группу.
func (cli *Cli) CreateGroup(group string, settings map[string]any) (err error) {
	if group == "" {
		return fmt.Errorf("group name is required")
	}

	const cmd = "CREATEGROUP"
	if settings != nil {
		err = cli.QueryNV(cmd, group, settings)
	} else {
		err = cli.QueryNV(cmd, group)
	}
	return
}

// DeleteGroup удаляет группу.
func (cli *Cli) DeleteGroup(group string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	return cli.QueryNV("DELETEGROUP", group)
}

// GetGroup возвращает настройки группы и список её участников (Members).
func (cli *Cli) GetGroup(group string) (*GroupResult, error) {
	if group == "" {
		return nil, fmt.Errorf("group name is required")
	}

	res, err := cli.getMapAny("GETGROUP", group)
	if err != nil {
		return nil, err
	}

	result := &GroupResult{
		Settings: make(map[string]any),
		Members:  []string{},
	}

	for k, v := range res {
		if k == "Members" {
			m, err := toStringSlice("GETGROUP", v)
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
		return false, fmt.Errorf("group name is required")
	}

	err := cli.QueryNV("GETGROUP", group)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListGroups возвращает список имен всех групп в указанном домене.
func (cli *Cli) ListGroups(domain string) ([]string, error) {
	return cli.getSliceString("LISTGROUPS", Atom(domain))
}

// RenameGroup переименовывает группу или переносит её в другой домен.
func (cli *Cli) RenameGroup(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new group names are required")
	}
	return cli.QueryNV("RENAMEGROUP", oldName, "INTO", newName)
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

	return cli.QueryNV("SETGROUP", group, fullData)
}
