// # Account Services
//
// The Account Services section provides access to user-specific application
// data and billing operations. These methods interface with the high-level
// subsystems of CommuniGate Pro that handle real-time user interactions
// and resource accounting.
//
// Key capabilities include:
//   - Billing & Quotas: Using the [Cli.Balance] method to perform complex
//     financial operations, including balance inquiries, credit reservations,
//     and service charging.
//   - Data Management: Interacting with account-level datasets (such as
//     collected addresses) via [Cli.Dataset] and [Cli.RemoveAccountSubset].
//   - Instant Messaging: Managing the XMPP/SIP [Cli.Roster], allowing for
//     programmatic contact list manipulation, subscription management,
//     and presence monitoring.
//
// These services are essential for integrating CommuniGate Pro with external
// billing platforms, CRM systems, and custom messaging clients.
package cgpcli

import (
	"fmt"
)

// Balance manages Account Billing Balances.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - params: A dictionary containing the "op" string element (list, reserve, release,
//     charge, credit, read, readAll, history, remove) and operation-specific data.
//
// Returns:
//   - map[string]any: A dictionary with the billing operation results.
//   - error: An error if the command fails or permissions are insufficient.
func (cli *Cli) Balance(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	return cli.getMapAny("BALANCE", Atom(account), params)
}

// Dataset manages Account "datasets".
//
// Parameters:
//   - account: The name of an existing Account (use "*" for current user).
//   - params: A dictionary containing:
//     subsetName - (string) Target dataset or dataset subset;
//     what - (string) The operation to apply:
//     listSubsets - Lists all subsets of the specified dataset;
//     createSet - Creates the specified dataset;
//     removeSet - Removes the specified dataset;
//     listEntries - Lists subset entries;
//     setEntry - Creates or updates an entry;
//     deleteEntry - Removes the specified entry;
//     addRandomEntry - Adds a new entry with a generated name;
//     findAddress - Finds an entry by "addressbook.Email" attribute.
//
// Returns:
//   - map[string]any: A dictionary with the operation results.
//   - error: An error if the command fails.
func (cli *Cli) Dataset(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	return cli.getMapAny("DATASET", Atom(account), params)
}

// RemoveAccountSubset removes a specific Account "dataset" by name.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - subset: This parameter specifies the name of an existing data subset in the specified Account.
//
// Returns:
//   - error: An error if the account or subset is not found.
func (cli *Cli) RemoveAccountSubset(account, subset string) error {
	if account == "" || subset == "" {
		return fmt.Errorf("account name and subset name are required")
	}
	return cli.QueryNV("REMOVEACCOUNTSUBSET", Atom(account), "SUBSET", subset)
}

// Roster manages the Account contact list (IM/Presence).
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - params: A dictionary containing the "what" string element (List, Update,
//     remove, Presence, probe) and operation-specific elements.
//
// Returns:
//   - map[string]any: A dictionary with the roster operation results.
//   - error: An error if the command fails.
func (cli *Cli) Roster(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	return cli.getMapAny("ROSTER", Atom(account), params)
}
