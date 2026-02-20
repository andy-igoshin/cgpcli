// # Group Administration
//
// The Groups section provides functionality for managing CommuniGate Pro groups.
// Groups serve as versatile objects for mail distribution and signal routing.
//
// Key capabilities include:
//   - Lifecycle Management: creating, deleting, and renaming groups across
//     the environment.
//   - Membership & Configuration: using [Cli.GetGroup] and [Cli.SetGroup] to manage
//     both group-specific settings and the list of group members simultaneously.
//   - Discovery: listing all groups within a specific domain and checking
//     for the existence of individual groups.
//   - Unified Data Access: the [GroupResult] structure simplifies interaction
//     by separating administrative settings from the member list.
package cgpcli

import (
	"fmt"
)

// GroupResult represents the combined data of a CommuniGate Pro group,
// separating operational settings from the actual member list.
type GroupResult struct {
	Settings map[string]any // administrative settings and attributes of the group.
	Members  []string       // list of addresses or object names belonging to the group.
}

// CreateGroup creates new Groups.
//
// Parameters:
//   - group: the name for the new Group. The name can include the Domain name, in this case the Group is created in the specified Domain.
//     If the Domain name is not specified, the command applies to the administrator Domain.
//   - settings: an optional initial Group settings and the members list.
//
// This method executes the CREATEGROUP CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateGroup(group string, settings map[string]any) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}

	const cmd = "CREATEGROUP"
	if settings != nil {
		return cli.QueryNV(cmd, group, settings)
	} else {
		return cli.QueryNV(cmd, group)
	}
}

// DeleteGroup removes Groups.
//
// Parameters:
//   - group: the name of an existing Group. The name can include the Domain name.
//
// This method executes the DELETEGROUP CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteGroup(group string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	return cli.QueryNV("DELETEGROUP", group)
}

// GetGroup retrieves the Group settings.
//
// Parameters:
//   - group: the name of an existing Group. The name can include the Domain name.
//
// This method executes the GETGROUP CLI command.
//
// Returns:
//   - *[GroupResult]: a structure containing the Group settings and members.
//   - error: an error if the command fails.
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

// IsGroupExists checks if a group exists by attempting to retrieve its settings.
//
// Parameters:
//   - group: the name of an existing Group. The name can include the Domain name.
//
// This method executes the GETGROUP CLI command.
//
// Returns:
//   - bool: true if the group exists.
//   - error: an error if the command fails.
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

// ListGroups retrieves the list of all Groups in the Domain.
//
// Parameters:
//   - domain: an optional Domain name. If empty, applies to the administrator Domain.
//
// This method executes the LISTGROUPS CLI command.
//
// Returns:
//   - []string: an array with the names of all Groups in the specified (or default) Domain.
//   - error: an error if the command fails.
func (cli *Cli) ListGroups(domain string) ([]string, error) {
	return cli.getSliceString("LISTGROUPS", Atom(domain))
}

// RenameGroup renames Groups.
//
// Parameters:
//   - oldName: the name of an existing Group. The name can include the Domain name.
//   - newName: the new Group name. The name can include the Domain name.
//
// This method executes the RENAMEGROUP CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameGroup(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new group names are required")
	}
	return cli.QueryNV("RENAMEGROUP", oldName, "INTO", newName)
}

// SetGroup sets the Group settings.
//
// Parameters:
//   - group: the name of an existing Group. The name can include the Domain name.
//   - settings: this dictionary is used to replace the Group settings dictionary.
//   - members: an optional list of group members to be included in the settings dictionary.
//
// This method executes the SETGROUP CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
