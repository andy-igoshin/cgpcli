// # Account Services
//
// The Account Services section provides access to user-specific application
// data and billing operations. These methods interface with the high-level
// subsystems of CommuniGate Pro that handle real-time user interactions
// and resource accounting.
//
// Key capabilities include:
//   - Billing & Quotas: using the [Cli.Balance] method to perform complex
//     financial operations, including balance inquiries, credit reservations,
//     and service charging.
//   - Data Management: interacting with account-level datasets (such as
//     collected addresses) via [Cli.Dataset] and [Cli.RemoveAccountSubset].
//   - Instant Messaging: managing the XMPP/SIP [Cli.Roster], allowing for
//     programmatic contact list manipulation, subscription management,
//     and presence monitoring.
package cgpcli

// Balance manages Account Billing Balances.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - params: a dictionary containing the "op" string element (list, reserve, release,
//     charge, credit, read, readAll, history, remove) and operation-specific data.
//
// This method executes the BALANCE CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the billing operation results.
//   - error: an error if the command fails.
func (cli *Cli) Balance(account string, params map[string]any) (map[string]any, error) {
	if account == "" {
		return nil, ErrAccountRequired
	}
	if len(params) == 0 {
		return nil, ErrParamsRequired
	}
	return cli.getMapAny("BALANCE", Atom(account), params)
}

// Dataset manages Account "datasets".
//
// Parameters:
//
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - params: a dictionary containing:
//     subsetName: (string) Target dataset or dataset subset;
//     what: (string) The operation to apply:
//     listSubsets: Lists all subsets of the specified dataset;
//     createSet: Creates the specified dataset;
//     removeSet: Removes the specified dataset;
//     listEntries: Lists subset entries;
//     setEntry: Creates or updates an entry;
//     deleteEntry: Removes the specified entry;
//     addRandomEntry: Adds a new entry with a generated name;
//     findAddress: Finds an entry by "addressbook.Email" attribute.
//
// This method executes the DATASET CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the operation results.
//   - error: an error if the command fails.
func (cli *Cli) Dataset(account string, params map[string]any) (map[string]any, error) {
	if account == "" {
		return nil, ErrAccountRequired
	}
	if len(params) == 0 {
		return nil, ErrParamsRequired
	}
	return cli.getMapAny("DATASET", Atom(account), params)
}

// RemoveAccountSubset removes a specific Account "dataset" by name.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - subset: the name of an existing data subset in the specified Account.
//
// This method executes the REMOVEACCOUNTSUBSET CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RemoveAccountSubset(account, subset string) error {
	if account == "" {
		return ErrAccountRequired
	}
	if subset == "" {
		return ErrSubsetRequired
	}
	return cli.QueryNV("REMOVEACCOUNTSUBSET", Atom(account), "SUBSET", subset)
}

// Roster manages the Account contact list (IM/Presence).
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - params: a dictionary containing the "what" string element (List, Update,
//     remove, Presence, probe) and operation-specific elements.
//
// This method executes the ROSTER CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the roster operation results.
//   - error: an error if the command fails.
func (cli *Cli) Roster(account string, params map[string]any) (map[string]any, error) {
	if account == "" {
		return nil, ErrAccountRequired
	}
	if len(params) == 0 {
		return nil, ErrParamsRequired
	}
	return cli.getMapAny("ROSTER", Atom(account), params)
}
