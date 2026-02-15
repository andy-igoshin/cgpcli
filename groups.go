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

// CreateGroup creates new Groups.
//
// Parameters:
//   - group: This parameter specifies the name for the new Group. The name can contain the @ symbol followed by the Domain name, in this case the Group is created in the specified Domain. If the Domain name is not specified, the command applies to the administrator Domain.
//   - settings: (Optional) This optional parameter specifies the initial Group settings and the members list.
//
// Returns:
//   - error: An error if the command fails.
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

// DeleteGroup removes Groups.
//
// Parameters:
//   - group: This parameter specifies the name of an existing Group. The name can include the Domain name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteGroup(group string) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}
	return cli.QueryNV("DELETEGROUP", group)
}

// GetGroup retrieves the Group settings.
//
// Parameters:
//   - group: This parameter specifies the name of an existing Group. The name can include the Domain name.
//
// Returns:
//   - *[GroupResult]: A structure containing the Group settings and members.
//   - error: An error if the command fails.
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
//   - group: This parameter specifies the name of an existing Group. The name can include the Domain name.
//
// Returns:
//   - bool: True if the group exists.
//   - error: An error if the check fails (excluding 'not found' errors).
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
//   - domain: (Optional) This optional parameter specifies the Domain name. If the Domain name is not specified, the command applies to the administrator Domain.
//
// Returns:
//   - []string: An array with the names of all Groups in the specified (or default) Domain.
//   - error: An error if the command fails.
func (cli *Cli) ListGroups(domain string) ([]string, error) {
	return cli.getSliceString("LISTGROUPS", Atom(domain))
}

// RenameGroup renames Groups.
//
// Parameters:
//   - oldName: This parameter specifies the name of an existing Group. The name can include the Domain name.
//   - newName: This parameter specifies the new Group name. The name can include the Domain name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) RenameGroup(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new group names are required")
	}
	return cli.QueryNV("RENAMEGROUP", oldName, "INTO", newName)
}

// SetGroup sets the Group settings.
//
// Parameters:
//   - group: This parameter specifies the name of an existing Group. The name can include the Domain name.
//   - settings: This dictionary is used to replace the Group settings dictionary.
//   - members: (Optional) The list of group members to be included in the settings dictionary.
//
// Returns:
//   - error: An error if the command fails.
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
