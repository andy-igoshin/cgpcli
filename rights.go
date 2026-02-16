// # Access Rights Administration
//
// The Access Rights section handles the assignment of administrative privileges to accounts.
// In CommuniGate Pro, "Rights" represent a set of permissions that grant an account
// authority over server-wide or domain-specific administrative tasks.
//
// Key capabilities include:
//   - Privilege Management: Using [Cli.SetAccountRights] to define a specific set
//     of administrative permissions for any given account.
//   - Granular Control: Assigning predefined CommuniGate Pro rights such as "CanModifySettings",
//     "CanCreateAccounts", or "CanManageQueues".
//
// Use these methods to automate the provisioning of delegated administration
// and to enforce security policies by restricting access to sensitive
// server operations.
package cgpcli

import (
	"fmt"
)

// SetAccountRights sets the Account Server Access rights.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rights: This array should contain the Access Right codes.
//
// Note: This operation is not additive. All old Account access rights
// are removed and replaced by the rights array.
//
// Returns:
//   - error: An error if the account does not exist or the command fails.
func (cli *Cli) SetAccountRights(account string, rights []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("SETACCOUNTRIGHTS", account, rights)
}
