// # Mailbox Administration
//
// The Mailboxes section provides methods for managing account folder structures,
// access control lists (ACL), and mailbox-specific metadata. It supports
// operations for both private user folders and shared resources.
//
// Key capabilities include:
//   - Lifecycle Management: Creating, renaming, and deleting mailboxes.
//     Supports recursive operations for entire folder trees using the "MAILBOXES" keyword.
//   - Access Control: Granular management of permissions via [Cli.GetMailboxACL]
//     and [Cli.SetMailboxACL], including support for effective rights calculation.
//   - Metadata & Organization: Managing mailbox classes (e.g., for MAPI/Groupware
//     compatibility) and mailbox-level aliases.
//   - Discovery: Listing mailbox trees with filters and retrieving detailed
//     information such as message counts and storage size.
//   - On-behalf Operations: Most methods support the "AUTH" parameter,
//     enabling administrative actions on behalf of a specific user.
//
// Use these methods to manage IMAP folder structures, shared calendars,
// and complex permission schemes in multi-user environments.
package cgpcli

import (
	"fmt"
)

// CreateMailbox creates a Mailbox in the specified Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name for the new Mailbox.
//   - class: (Optional) This parameter specifies the Mailbox class for the new Mailbox.
//   - authAccount: (Optional) This parameter specifies the name of an Account on whose behalf this operation should be executed.
//
// Returns:
//   - error: An error if the command fails.
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
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name of the Mailbox to be deleted.
//   - recursive: If true, the keyword MAILBOXES is used, and all nested Mailboxes (submailboxes) are deleted, too.
//   - authAccount: (Optional) This parameter specifies the name of an Account on whose behalf the operation should be executed.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteMailbox(account, mailbox string, recursive bool, authAccount string) (err error) {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account name and mailbox name are required")
	}

	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}
	if authAccount != "" {
		err = cli.QueryNV("DELETEMAILBOX", Atom(account), keyword, mailbox, "AUTH", authAccount)
	} else {
		err = cli.QueryNV("DELETEMAILBOX", Atom(account), keyword, mailbox)
	}
	return
}

// GetMailboxACL retrieves the access control list for the Account Mailbox.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name of an existing Mailbox in the specified Account.
//   - authAccount: (Optional) This parameter specifies the name of an Account on whose behalf the operation should be executed.
//
// Returns:
//   - map[string]string: A dictionary with the Mailbox access elements.
//   - error: An error if the command fails.
func (cli *Cli) GetMailboxACL(account, mailbox, authAccount string) (res map[string]string, err error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account name and mailbox name are required")
	}

	const cmd = "GETMAILBOXACL"
	if authAccount != "" {
		res, err = cli.getMapString(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	} else {
		res, err = cli.getMapString(cmd, Atom(account), "MAILBOX", mailbox)
	}
	return
}

// GetMailboxAliases retrieves the Account Mailbox aliases.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// Returns:
//   - map[string]string: A dictionary where each key is the name of an existing Mailbox alias, and the value is the name of the Mailbox it points to.
//   - error: An error if the command fails.
func (cli *Cli) GetMailboxAliases(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapString("GETMAILBOXALIASES", Atom(account))
}

// GetMailboxInfo retrieves the internal information about the Account Mailbox.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name of an existing Mailbox in the specified Account.
//   - authAccount: (Optional) This parameter specifies the name of an Account on whose behalf the operation should be executed.
//
// Returns:
//   - map[string]any: A dictionary with the Mailbox internal information.
//   - error: An error if the command fails.
func (cli *Cli) GetMailboxInfo(account, mailbox, authAccount string) (res map[string]any, err error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}

	const cmd = "GETMAILBOXINFO"
	if authAccount != "" {
		res, err = cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	} else {
		res, err = cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox)
	}
	return
}

// GetMailboxRights retrieves the effective Mailbox access rights based on the Mailbox ACL for the given authAccountName.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name of an existing Mailbox in the specified Account.
//   - authAccount: This parameter specifies the name of an Account whose access rights should be retrieved.
//
// Returns:
//   - string: A string with the effective Mailbox access rights.
//   - error: An error if the command fails.
func (cli *Cli) GetMailboxRights(account, mailbox, authAccount string) (string, error) {
	if account == "" || mailbox == "" || authAccount == "" {
		return "", fmt.Errorf("account, mailbox and authAccount names are required")
	}
	return cli.getString("GETMAILBOXRIGHTS", Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
}

// GetMailboxSubscription retrieves the list of Account "subscribed Mailboxes".
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// Returns:
//   - []string: An array with the list of Account "subscribed Mailboxes".
//   - error: An error if the command fails.
func (cli *Cli) GetMailboxSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETMAILBOXSUBSCRIPTION", Atom(account))
}

// ListMailboxes retrieves the list of Account Mailboxes.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - filter: (Optional) This parameter specifies the filter string to apply to Account Mailbox names (wildcards "*" and "%" are supported).
//   - authAccount: (Optional) This parameter specifies the name of an Account on whose behalf the LIST operation should be executed.
//
// Returns:
//   - map[string]any: A dictionary where each key is a Mailbox name and the value contains Mailbox information or nested structures.
//   - error: An error if the command fails.
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
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - oldName: This parameter specifies the name of the Mailbox to be renamed.
//   - newName: This parameter specifies the new name for the Mailbox.
//   - recursive: If true, the keyword MAILBOXES is used, and all nested Mailboxes (submailboxes) are renamed, too.
//   - authAccount: (Optional) This optional parameter specifies the name of an Account on whose behalf the operation should be executed.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) RenameMailbox(account, oldName, newName string, recursive bool, authAccount string) (err error) {
	if account == "" || oldName == "" || newName == "" {
		return fmt.Errorf("account, old and new names are required")
	}

	const cmd = "RENAMEMAILBOX"
	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), keyword, oldName, "INTO", newName, "AUTH", authAccount)
	} else {
		err = cli.QueryNV(cmd, Atom(account), keyword, oldName, "INTO", newName)
	}
	return
}

// SetMailboxACL modifies the access control list for the Account Mailbox.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name of an existing Mailbox in the specified Account.
//   - acl: This parameter specifies the access right elements to be modified. Use "+" or "-" prefixes to add or remove rights.
//   - authAccount: (Optional) This parameter specifies the name of an Account on whose behalf the operation should be executed.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetMailboxACL(account, mailbox string, acl map[string]string, authAccount string) (err error) {
	if account == "" || mailbox == "" || acl == nil {
		return fmt.Errorf("account, mailbox and non-empty acl are required")
	}

	const cmd = "SETMAILBOXACL"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount, acl)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, acl)
	}
	return
}

// SetMailboxAliases sets the Account Mailbox aliases.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - aliases: The set of new Mailbox aliases.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetMailboxAliases(account string, aliases map[string]string) error {
	if account == "" || aliases == nil {
		return fmt.Errorf("account name and aliases dictionary are required")
	}
	return cli.QueryNV("SETMAILBOXALIASES", Atom(account), aliases)
}

// SetMailboxClass sets the "class" of an Account Mailbox.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - mailbox: This parameter specifies the name of an existing Mailbox in the specified Account.
//   - class: The Mailbox class.
//   - authAccount: (Optional) This parameter specifies the name of an Account whose Mailbox access rights should be used.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetMailboxClass(account, mailbox, class, authAccount string) (err error) {
	if account == "" || mailbox == "" || class == "" {
		return fmt.Errorf("account, mailbox and class are required")
	}

	const cmd = "SETMAILBOXCLASS"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount, "CLASS", class)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "CLASS", class)
	}
	return
}

// SetMailboxSubscription sets the Account "subscribed Mailboxes" list.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account. The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - subs: The list of subscribed Mailboxes. Each array element should be a string with a Mailbox name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetMailboxSubscription(account string, subs []string) error {
	if account == "" || subs == nil {
		return fmt.Errorf("account name and subs array are required")
	}
	return cli.QueryNV("SETMAILBOXSUBSCRIPTION", Atom(account), subs)
}
