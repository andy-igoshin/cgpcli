// # Mailbox Administration
//
// The Mailboxes section provides methods for managing account folder structures,
// access control lists (ACL), and mailbox-specific metadata. It supports
// operations for both private user folders and shared resources.
//
// Key capabilities include:
//   - Lifecycle Management: creating, renaming, and deleting mailboxes.
//     Supports recursive operations for entire folder trees using the "MAILBOXES" keyword.
//   - Access Control: granular management of permissions via [Cli.GetMailboxACL]
//     and [Cli.SetMailboxACL], including support for effective rights calculation.
//   - Metadata & Organization: managing mailbox classes (e.g., for MAPI/Groupware
//     compatibility) and mailbox-level aliases.
//   - Discovery: listing mailbox trees with filters and retrieving detailed
//     information such as message counts and storage size.
//   - On-behalf Operations: most methods support the "AUTH" parameter,
//     enabling administrative actions on behalf of a specific user.
package cgpcli

import (
	"fmt"
)

// CreateMailbox creates a Mailbox in the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name for the new Mailbox.
//   - class: mailbox class for the new Mailbox.
//   - authAccount: an optional name of an Account on whose behalf the operation is executed.
//
// This method executes the MAILBOX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateMailbox(account, mailbox, class, authAccount string) error {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account name and mailbox name are required")
	}

	args := make([]any, 0, 7)
	args = append(args, Atom(account), "MAILBOX", mailbox)
	if class != "" {
		args = append(args, "CLASS", class)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	return cli.QueryNV("CREATEMAILBOX", args...)
}

// DeleteMailbox removes a Mailbox from the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name of the Mailbox to be deleted.
//   - recursive: if true, the keyword MAILBOXES is used, and all nested Mailboxes (submailboxes) are deleted, too.
//   - authAccount: an optional name of an Account on whose behalf the operation is executed.
//
// This method executes the MAILBOX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteMailbox(account, mailbox string, recursive bool, authAccount string) error {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account name and mailbox name are required")
	}

	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}
	if authAccount != "" {
		return cli.QueryNV("DELETEMAILBOX", Atom(account), keyword, mailbox, "AUTH", authAccount)
	} else {
		return cli.QueryNV("DELETEMAILBOX", Atom(account), keyword, mailbox)
	}
}

// GetMailboxACL retrieves the access control list for the Account Mailbox.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name of an existing Mailbox in the specified Account.
//   - authAccount: an optional name of an Account on whose behalf the operation is executed.
//
// This method executes the GETMAILBOXACL CLI command.
//
// Returns:
//   - map[string]string: a dictionary with the Mailbox access elements.
//   - error: an error if the command fails.
func (cli *Cli) GetMailboxACL(account, mailbox, authAccount string) (map[string]string, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account name and mailbox name are required")
	}

	const cmd = "GETMAILBOXACL"
	if authAccount != "" {
		return cli.getMapString(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	} else {
		return cli.getMapString(cmd, Atom(account), "MAILBOX", mailbox)
	}
}

// GetMailboxAliases retrieves the Account Mailbox aliases.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// This method executes the GETMAILBOXALIASES CLI command.
//
// Returns:
//   - map[string]string: a dictionary where each key is the name of an existing Mailbox alias, and the value is the name of the Mailbox it points to.
//   - error: an error if the command fails.
func (cli *Cli) GetMailboxAliases(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapString("GETMAILBOXALIASES", Atom(account))
}

// GetMailboxInfo retrieves the internal information about the Account Mailbox.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name of an existing Mailbox in the specified Account.
//   - authAccount: an optional name of an Account on whose behalf the operation is executed.
//
// This method executes the GETMAILBOXINFO CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Mailbox internal information.
//   - error: an error if the command fails.
func (cli *Cli) GetMailboxInfo(account, mailbox, authAccount string) (map[string]any, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}

	const cmd = "GETMAILBOXINFO"
	if authAccount != "" {
		return cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	} else {
		return cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox)
	}
}

// GetMailboxRights retrieves the effective Mailbox access rights based on the Mailbox ACL for the given authAccountName.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name of an existing Mailbox in the specified Account.
//   - authAccount: the name of an Account whose access rights should be retrieved.
//
// This method executes the GETMAILBOXRIGHTS CLI command.
//
// Returns:
//   - string: a string with the effective Mailbox access rights.
//   - error: an error if the command fails.
func (cli *Cli) GetMailboxRights(account, mailbox, authAccount string) (string, error) {
	if account == "" || mailbox == "" || authAccount == "" {
		return "", fmt.Errorf("account, mailbox and authAccount names are required")
	}
	return cli.getString("GETMAILBOXRIGHTS", Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
}

// GetMailboxSubscription retrieves the list of Account "subscribed Mailboxes".
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// This method executes the GETMAILBOXSUBSCRIPTION CLI command.
//
// Returns:
//   - []string: an array with the list of Account "subscribed Mailboxes".
//   - error: an error if the command fails.
func (cli *Cli) GetMailboxSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETMAILBOXSUBSCRIPTION", Atom(account))
}

// ListMailboxes retrieves the list of Account Mailboxes.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - filter: an optional filter string for Account Mailbox names (wildcards "*" and "%" are supported).
//   - authAccount: an optional name of an Account on whose behalf the LIST operation is executed.
//
// This method executes the LISTMAILBOXES CLI command.
//
// Returns:
//   - map[string]any: a dictionary where each key is a Mailbox name and the value contains Mailbox information or nested structures.
//   - error: an error if the command fails.
func (cli *Cli) ListMailboxes(account, filter, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	args := make([]any, 0, 5)
	args = append(args, Atom(account))
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	return cli.getMapAny("LISTMAILBOXES", args...)
}

// RenameMailbox renames a Mailbox in the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - oldName: the name of the Mailbox to be renamed.
//   - newName: the new name for the Mailbox.
//   - recursive: if true, the keyword MAILBOXES is used, and all nested Mailboxes (submailboxes) are renamed, too.
//   - authAccount: an optional name of an Account on whose behalf the operation is executed.
//
// This method executes the RENAMEMAILBOX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameMailbox(account, oldName, newName string, recursive bool, authAccount string) error {
	if account == "" || oldName == "" || newName == "" {
		return fmt.Errorf("account, old and new names are required")
	}

	const cmd = "RENAMEMAILBOX"
	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}
	if authAccount != "" {
		return cli.QueryNV(cmd, Atom(account), keyword, oldName, "INTO", newName, "AUTH", authAccount)
	} else {
		return cli.QueryNV(cmd, Atom(account), keyword, oldName, "INTO", newName)
	}
}

// SetMailboxACL modifies the access control list for the Account Mailbox.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name of an existing Mailbox in the specified Account.
//   - acl: the access right elements to be modified. Use "+" or "-" prefixes to add or remove rights.
//   - authAccount: an optional name of an Account on whose behalf the operation is executed.
//
// This method executes the SETMAILBOXACL CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetMailboxACL(account, mailbox string, acl map[string]string, authAccount string) error {
	if account == "" || mailbox == "" || acl == nil {
		return fmt.Errorf("account, mailbox and non-empty acl are required")
	}

	const cmd = "SETMAILBOXACL"
	if authAccount != "" {
		return cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount, acl)
	} else {
		return cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, acl)
	}
}

// SetMailboxAliases sets the Account Mailbox aliases.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - aliases: the set of new Mailbox aliases.
//
// This method executes the SETMAILBOXALIASES CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetMailboxAliases(account string, aliases map[string]string) error {
	if account == "" || aliases == nil {
		return fmt.Errorf("account name and aliases dictionary are required")
	}
	return cli.QueryNV("SETMAILBOXALIASES", Atom(account), aliases)
}

// SetMailboxClass sets the "class" of an Account Mailbox.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: the name of an existing Mailbox in the specified Account.
//   - class: the Mailbox class.
//   - authAccount: an optional name of an Account whose Mailbox access rights should be used.
//
// This method executes the SETMAILBOXCLASS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetMailboxClass(account, mailbox, class, authAccount string) error {
	if account == "" || mailbox == "" || class == "" {
		return fmt.Errorf("account, mailbox and class are required")
	}

	const cmd = "SETMAILBOXCLASS"
	if authAccount != "" {
		return cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount, "CLASS", class)
	} else {
		return cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "CLASS", class)
	}
}

// SetMailboxSubscription sets the Account "subscribed Mailboxes" list.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - subs: the list of subscribed Mailboxes. Each array element should be a string with a Mailbox name.
//
// This method executes the SETMAILBOXSUBSCRIPTION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetMailboxSubscription(account string, subs []string) error {
	if account == "" || subs == nil {
		return fmt.Errorf("account name and subs array are required")
	}
	return cli.QueryNV("SETMAILBOXSUBSCRIPTION", Atom(account), subs)
}
