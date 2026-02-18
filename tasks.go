// # Named Task Administration
//
// The Named Tasks section provides management for the CommuniGate Pro
// internal scheduler. Named Tasks allow administrators and users to
// automate actions (such as executing CG/PL scripts) based on specific
// time intervals or calendar schedules.
//
// Key capabilities include:
//   - Lifecycle Management: creating, renaming, and deleting tasks associated
//     with specific accounts using [Cli.CreateNamedTask] and [Cli.DeleteNamedTask].
//   - Configuration: retrieving and updating task parameters.
//   - Discovery: listing tasks at both account and domain levels to
//     monitor automated workflows across the server.
//   - Owner Resolution: automatic handling of account names, ensuring
//     tasks are correctly linked to the local part of the account address.
//
package cgpcli

import (
	"fmt"
	"strings"
)

// CreateNamedTask creates a new Named Task.
//
// Parameters:
//   - task: the name for the new Named Task. Can include the @ symbol followed by the Domain name.
//     If the Domain name is not specified, the command applies to the administrator Domain.
//   - account: the owner Account name.
//
// This method executes the CREATENAMEDTASK CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateNamedTask(task, account string) error {
	if task == "" || account == "" {
		return fmt.Errorf("task name and account name are required")
	}

	if idx := strings.IndexByte(account, '@'); idx != -1 {
		account = account[:idx]
	}
	return cli.QueryNV("CREATENAMEDTASK", task, "FOR", account)
}

// DeleteNamedTask removes an existing Named Task.
//
// Parameters:
//   - task: the name of an existing Named Task. The name can include the Domain name.
//
// This method executes the DELETENAMEDTASK CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteNamedTask(task string) error {
	if task == "" {
		return fmt.Errorf("task name is required")
	}
	return cli.QueryNV("DELETENAMEDTASK", task)
}

// GetNamedTask retrieves the Named Task settings.
//
// Parameters:
//   - task: the name of an existing Named Task. The name can include the Domain name.
//
// This method executes the GETNAMEDTASK CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Named Task settings.
//   - error: an error if the command fails.
func (cli *Cli) GetNamedTask(task string) (map[string]any, error) {
	if task == "" {
		return nil, fmt.Errorf("task name is required")
	}
	return cli.getMapAny("GETNAMEDTASK", task)
}

// ListAccountNamedTasks retrieves the list of all Named Tasks owned by
// the specified Account.
//
// Parameters:
//   - account: the owner Account name.
//
// This method executes the LISTACCOUNTNAMEDTASKS CLI command.
//
// Returns:
//   - map[string]any: a dictionary where keys are Task names and values
//     are dictionaries containing Task owner, Real Name, and Application program.
//   - error: an error if the command fails.
func (cli *Cli) ListAccountNamedTasks(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("LISTACCOUNTNAMEDTASKS", account)
}

// ListDomainNamedTasks retrieves the list of all Named Tasks in the Domain.
//
// Parameters:
//   - domain: an optional Domain name. If the Domain name is not specified,
//     the command applies to the administrator Domain.
//
// This method executes the LISTDOMAINNAMEDTASKS CLI command.
//
// Returns:
//   - map[string]any: a dictionary where keys are Task names and values
//     are dictionaries containing Task owner, Real Name, and Application program.
//   - error: an error if the command fails.
func (cli *Cli) ListDomainNamedTasks(domain string) (map[string]any, error) {
	return cli.getMapAny("LISTDOMAINNAMEDTASKS", Atom(domain))
}

// RenameNamedTask renames an existing Named Task.
//
// Parameters:
//   - oldTask: the name of an existing Named Task. The name can include the Domain name.
//   - newTask: the new Named Task name.
//
// This method executes the RENAMENAMEDTASK CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameNamedTask(oldTask, newTask string) error {
	if oldTask == "" || newTask == "" {
		return fmt.Errorf("oldName and newName are required")
	}
	return cli.QueryNV("RENAMENAMEDTASK", oldTask, "INTO", newTask)
}

// UpdateNamedTask sets the Named Task settings.
//
// Parameters:
//   - task: the name of an existing Named Task. The name can include the Domain name.
//   - settings: a dictionary used to update the Named Task settings dictionary.
//
// This method executes the UPDATENAMEDTASK CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateNamedTask(task string, settings map[string]any) error {
	if task == "" || len(settings) == 0 {
		return fmt.Errorf("task name and settings are required")
	}
	return cli.QueryNV("UPDATENAMEDTASK", task, settings)
}
