// # Access Rights Administration
//
// The Access Rights section handles the assignment of administrative privileges to accounts.
// In CommuniGate Pro, "Rights" represent a set of permissions that grant an account
// authority over server-wide or domain-specific administrative tasks.
//
// Key capabilities include:
//   - Privilege Management: using [Cli.SetAccountRights] to define a specific set
//     of administrative permissions for any given account.
//   - Granular Control: assigning predefined CommuniGate Pro rights such as "CanModifySettings",
//     "CanCreateAccounts", or "CanManageQueues".
package cgpcli

import (
	"fmt"
)

// SetAccountRights sets the Account Server Access rights.
//
// Parameters:
//   - account: the name of an existing Account.
//     The name can include the Domain name.
//   - rights: this array should contain the Access Right codes.
//
// Note: This operation is not additive. All old Account access rights
// are removed and replaced by the rights array.
//
// Returns:
//   - error: an error if the account does not exist or the command fails.
func (cli *Cli) SetAccountRights(account string, rights []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("SETACCOUNTRIGHTS", account, rights)
}
