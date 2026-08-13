// # Account Administration
//
// The Accounts section provides comprehensive methods for managing CommuniGate Pro accounts
// across their entire lifecycle. It supports core administrative tasks such as
// creation, deletion, and renaming, as well as fine-grained configuration of
// Settings (admin-level limits) and Preferences (user-level options).
//
// Key capabilities include:
//   - Password management and identity verification.
//   - Full support for Mail and Signal Rules.
//   - Access control via ACLs and Rights management.
//   - Alias and Telephony number (Telnums) configuration.
package cgpcli

import (
	"fmt"
	"strings"
)

// AccountType represents the storage and management type of a CommuniGate Pro account.
type AccountType int

const (
	// TypeDefault uses the default storage type defined for the Domain.
	TypeDefault AccountType = iota

	// MultiMailbox represents a standard multi-mailbox account.
	MultiMailbox

	// TextMailbox represents a legacy text-based single mailbox account.
	TextMailbox

	// MailDirMailbox represents an account using MailDir storage format.
	MailDirMailbox

	// SlicedMailbox represents an account using Sliced storage format.
	SlicedMailbox

	// AGrade, BGrade, CGrade, DGrade represent special account types
	AGrade
	BGrade
	CGrade
	DGrade

	// ResourceMailbox is a special type of account
	ResourceMailbox
)

var accountTypeNames = [...]string{
	"", "MultiMailbox", "TextMailbox", "MailDirMailbox", "SlicedMailbox",
	"AGrade", "BGrade", "CGrade", "DGrade", "ResourceMailbox",
}

func (t AccountType) String() string {
	if int(t) < 0 || int(t) >= len(accountTypeNames) {
		return ""
	}
	return accountTypeNames[t]
}

// BlockMode represents the blocking mode of a CommuniGate Pro account.
type BlockMode int

const (
	// BlockFull specifies full account blocking.
	BlockFull BlockMode = iota + 1
	// BlockPart specifies partial account blocking with specific restrictions.
	BlockPart
)

var blockModeNames = [...]string{"", "Full", "Part"}

func (m BlockMode) String() string {
	if int(m) < 1 || int(m) >= len(blockModeNames) {
		return ""
	}
	return blockModeNames[m]
}

// BlockFlags represents optional restriction flags applied during partial account blocking.
type BlockFlags int

const (
	PreventNone BlockFlags = 0
	// PreventLogin prevents user authorization in the system by any means.
	PreventLogin BlockFlags = 1 << (iota - 1)
	// PreventDel prevents deletion or movement of emails from the user's mailboxes.
	PreventDel
	// PreventShowFolders prevents showing delegated folders of this account to other users.
	PreventShowFolders
)

// BlockAccount blocks an existing Account using full or partial restrictions.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - mode: the blocking mode ([BlockFull] or [BlockPart]).
//   - flags: a bitmask of optional restriction flags ([PreventLogin], [PreventDel], [PreventShowFolders]).
//     Important: These flags are applied to [BlockPart] mode only. For partial blocking,
//     you must specify at least one flag for the command to execute successfully.
//
// This method executes the BLOCKACCOUNT CLI command.
//
// Returns:
//   - error: an error if the command fails or if validation constraints are violated.
func (cli *Cli) BlockAccount(account string, mode BlockMode, flags BlockFlags) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	modeStr := mode.String()
	if modeStr == "" {
		return fmt.Errorf("invalid block mode")
	}

	args := make([]any, 0, 6)
	args = append(args, account, "MODE", modeStr)

	if mode == BlockPart {
		if flags == 0 {
			return fmt.Errorf("at least one block flag must be specified for partial blocking")
		}
		if flags&PreventLogin != 0 {
			args = append(args, "PREVENTLOGIN")
		}
		if flags&PreventDel != 0 {
			args = append(args, "PREVENTDEL")
		}
		if flags&PreventShowFolders != 0 {
			args = append(args, "PREVENTSHOWFOLDERS")
		}
	}

	return cli.QueryNV("BLOCKACCOUNT", args...)
}

// CheckAccountPassword verifies if the provided password matches the Account's password.
//
// Parameters:
//   - account: the Account name.
//   - password: the plain-text password to verify against the server record.
//
// This is a helper method that performs a client-side comparison by retrieving
// the plain-text password from the server.
//
// Important: Unlike VerifyAccountPassword(), this method does NOT trigger
// the "External Authentication" helper. This makes it safe to use within
// the External Authentication helper itself to verify credentials against
// the local CommuniGate Pro storage without causing recursive loops.
//
// Note: This method only works if the Account's password is stored using a
// reversible encryption method.
//
// Returns:
//   - (true, nil) if the passwords match.
//   - (false, nil) if the passwords do not match.
//   - (false, error) if the password cannot be retrieved or verified.
func (cli *Cli) CheckAccountPassword(account, password string) (bool, error) {
	if password == "" {
		return false, fmt.Errorf("password is required")
	}
	plain, err := cli.GetAccountPlainPassword(account, "")
	if err != nil {
		return false, err
	}
	return plain == password, nil
}

// CreateAccount creates a new Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - settings: an optional dictionary of initial Account settings.
//   - accType: an optional Account type (e.g., MultiMailbox, MailDirMailbox).
//     If [TypeDefault] is used, a MultiMailbox-type Account is created.
//   - storage: an optional "storage mount point" directory for the Account data
//     (without the .mnt suffix).
//   - legacy: if true, creates the Account with a Legacy (visible for legacy mailers) INBOX.
//
// This method executes the CREATEACCOUNT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateAccount(account string, settings map[string]any, accType AccountType, storage string, legacy bool) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	args := make([]any, 0, 6)
	args = append(args, account)
	if accType != TypeDefault {
		args = append(args, accType.String())
	}
	if storage != "" {
		args = append(args, "PATH", storage)
	}
	if legacy {
		args = append(args, "LEGACY")
	}
	if settings != nil {
		args = append(args, settings)
	}
	return cli.QueryNV("CREATEACCOUNT", args...)
}

// DeleteAccount removes an existing Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// This method executes the DELETEACCOUNT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("DELETEACCOUNT", account)
}

// DeleteAccountMailRule removes a Mail Rule from an existing Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - rule: the name of the [MailRule] to be removed.
//
// If the specified Mail Rule does not exist, the command does nothing and
// does not return an error.
//
// This method executes the UPDATEACCOUNTMAILRULE CLI command with the DELETE keyword.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteAccountMailRule(account, rule string) error {
	if account == "" || rule == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	return cli.QueryNV("UPDATEACCOUNTMAILRULE", account, "DELETE", rule)
}

// DeleteAccountSignalRule removes a Signal Rule from an existing Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - rule: the name of the [SignalRule] to be removed.
//
// If the specified Signal Rule does not exist, the command does nothing and
// does not return an error.
//
// This method executes the UPDATEACCOUNTSIGNALRULE CLI command with the DELETE keyword.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteAccountSignalRule(account, rule string) error {
	if account == "" || rule == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	return cli.QueryNV("UPDATEACCOUNTSIGNALRULE", account, "DELETE", rule)
}

// GetAccountACL retrieves the Access Control Lists (ACLs) for an Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - authAccount: an optional name of an Account on whose behalf the operation
//     is executed. If specified, the ACL info is returned only if that Account
//     has the Admin access right for the target Account.
//
// This method executes the GETACCOUNTACL CLI command.
//
// Returns:
//   - map[string]any: a dictionary containing the ACL elements.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountACL(account, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTACL"
	if authAccount != "" {
		return cli.getMapAny(cmd, Atom(account), "AUTH", authAccount)
	} else {
		return cli.getMapAny(cmd, Atom(account))
	}
}

// GetAccountACLRights retrieves the effective access rights for a specific Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - authAccount: the name of an Account whose effective access rights for the
//     target Account should be retrieved.
//
// This method executes the GETACCOUNTACLRIGHTS CLI command.
//
// Returns:
//   - string: a string with the effective access rights.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountACLRights(account, authAccount string) (string, error) {
	if account == "" || authAccount == "" {
		return "", fmt.Errorf("target account and auth account names are required")
	}
	return cli.getString("GETACCOUNTACLRIGHTS", Atom(account), "AUTH", authAccount)
}

// GetAccountAccessModes retrieves the enabled services (AccessModes) for an Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// The AccessMode setting specifies Enabled Services. The returned slice of strings
// represents the names of the enabled services.
//
// Special values:
//   - If the setting is "All", the method returns []string{"All"}.
//   - If the setting is "None" or empty, the method returns an empty slice.
//   - If the setting is an array, the method skips the first element (the numeric
//     threshold) and returns the names of the explicitly listed services.
//
// Supported services include: Mail, POP, IMAP, WebMail, PWD, Agent, WebSite,
// Relay, Roaming, FTP, MAPI, TLS, S/MIME, LDAP, WebCAL, RADIUS, SIP, PBX,
// XMPP, XIMSS, Signal, AirSync, HTTP, MobilePBX, XMedia, YMedia, MobileClient,
// ClientMail, ClientIM, ClientVoIP.
//
// Returns:
//   - []string: a slice containing the names of enabled services.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountAccessModes(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	settings, err := cli.GetAccountSettings(account)
	if err != nil {
		return nil, err
	}
	raw, ok := settings["AccessModes"]
	if !ok || raw == nil {
		return nil, nil
	}
	switch val := raw.(type) {
	case string:
		if strings.EqualFold(val, "All") {
			return []string{"All"}, nil
		}
		return []string{}, nil
	case []any:
		if len(val) <= 1 {
			return []string{}, nil
		}
		return toStringSlice("GetAccountAccessModes", val[1:])
	default:
		return []string{}, nil
	}
}

// GetAccountAliases retrieves the list of Account aliases.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// This method executes the GETACCOUNTALIASES CLI command.
//
// Returns:
//   - []string: a slice of alias name strings.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountAliases(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTALIASES", account)
}

// GetAccountEffectivePrefs retrieves the effective Account Preferences.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// Both the explicitly set and the default preferences are included in the result.
//
// Note: All users can use this command to retrieve their own effective Preferences.
//
// This method executes the GETACCOUNTEFFECTIVEPREFS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the effective Account Preferences.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountEffectivePrefs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTEFFECTIVEPREFS", account)
}

// GetAccountEffectiveSettings retrieves the effective Account settings.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// Both the explicitly set and the default Account settings are included in the result.
//
// Note: All users can send this command for their own Accounts.
//
// This method executes the GETACCOUNTEFFECTIVESETTINGS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the effective Account settings.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountEffectiveSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTEFFECTIVESETTINGS", Atom(account))
}

// GetAccountInfo retrieves multiple elements from the Account "info" dictionary.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - keys: an optional slice of info key names to retrieve (e.g., "LastLogin", "LastAddress").
//     If empty, the entire "info" dictionary is returned.
//
// Note: Key names are case-sensitive. If Account "info" data is stored in .info
// dictionary files, do NOT include the hash (#) symbol in the key names.
//
// Note: All users can use this command to retrieve elements from their own Account "info" data.
//
// This method executes the GETACCOUNTINFO CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the requested info elements.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountInfo(account string, keys []string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTINFO"
	if len(keys) > 0 {
		return cli.getMapAny(cmd, Atom(account), keys)
	} else {
		return cli.getMapAny(cmd, Atom(account))
	}
}

// GetAccountLocation retrieves the storage path for an Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// The path is relative to the file directory of the Account Domain. For multi-mailbox
// Accounts it points to the Account file directory; for single-mailbox Accounts it
// points to the INBOX Mailbox path.
//
// This method executes the GETACCOUNTLOCATION CLI command.
//
// Returns:
//   - string: the Account file path.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountLocation(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	return cli.getString("GETACCOUNTLOCATION", account)
}

// GetAccountMailRules retrieves the list of Account Queue Mail Rules.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// This method executes the GETACCOUNTMAILRULES CLI command.
//
// Returns:
//   - [][MailRule]: a slice of MailRule structures.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountMailRules(account string) ([]MailRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMailRules("GETACCOUNTMAILRULES", account)
}

// GetAccountOneInfo retrieves a single element from the Account "info" dictionary.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - key: the name of the requested "info" element. The key name is case-sensitive.
//     Do NOT include the hash (#) symbol.
//
// If the element is not found, the output is an empty string.
//
// Note: All users can use this command to retrieve elements from their own Account "info" data.
//
// This method executes the GETACCOUNTINFO CLI command with the Key keyword.
//
// Returns:
//   - any: the specified info element.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountOneInfo(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTINFO", Atom(account), "Key", key)
}

// GetAccountOneSetting retrieves a single setting from the effective Account settings list.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - key: the name of the setting to read.
//
// The command produces an object which can be a string, an array, a dictionary,
// or a null-object representing the requested setting.
//
// Note: All users can send the GETACCOUNTONESETTING command for their own Accounts.
//
// This method executes the GETACCOUNTONESETTING CLI command.
//
// Returns:
//   - any: the setting value (string, array, dictionary, or nil).
//   - error: an error if the command fails.
func (cli *Cli) GetAccountOneSetting(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTONESETTING", Atom(account), key)
}

// GetAccountPlainPassword retrieves the clear-text password for an Account.
//
// Parameters:
//   - account: the Account name.
//   - method: an optional authentication method. Supported values are "SIP" and "RADIUS".
//     If empty, the main Account password is retrieved.
//
// Note: This command only works if the password is stored using a reversible
// encryption method. It will not return the plain-text password if a
// non-reversible hash algorithm is used for storage.
//
// This method executes the GETACCOUNTPLAINPASSWORD CLI command.
//
// Returns:
//   - string: the plain-text password.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountPlainPassword(account string, method string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTPLAINPASSWORD"
	if method == "SIP" || method == "RADIUS" {
		return cli.getString(cmd, account, "METHOD", method)
	} else {
		return cli.getString(cmd, account)
	}
}

// GetAccountPrefs retrieves the Account Preferences.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// Only the explicitly set (non-default) preferences are included in the result.
// To get the combined set of explicit and default preferences, use
// GetAccountEffectivePrefs instead.
//
// Note: All users can use the GETACCOUNTPREFS command to retrieve their own Account Preferences.
//
// This method executes the GETACCOUNTPREFS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Account Preferences.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountPrefs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTPREFS", account)
}

// GetAccountPresence retrieves the current "presence" status of an Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// The method returns the presence information in one of the following formats:
//   - []any (array of two strings): the first element is the presence status
//     (e.g., "Available", "Away"), and the second is the custom status message.
//   - string: the presence status (returned if no custom status message is set).
//   - nil: returned if the presence status is not set at all.
//
// This method executes the GETACCOUNTPRESENCE CLI command.
//
// Returns:
//   - any: the presence data.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountPresence(account string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.Query("GETACCOUNTPRESENCE", account)
}

// GetAccountRIMAPs retrieves the list of Account RIMAP (Remote IMAP) records.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// RIMAP records are used to poll external mailboxes and fetch messages into the local Account.
//
// This method executes the GETACCOUNTRIMAPS CLI command.
//
// Returns:
//   - map[string]any: a dictionary containing the RIMAP records.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountRIMAPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRIMAPS", account)
}

// GetAccountRPOPs retrieves the list of Account RPOP (Remote POP) records.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// RPOP records are used to collect mail from external POP3 servers into the local Account.
//
// This method executes the GETACCOUNTRPOPS CLI command.
//
// Returns:
//   - map[string]any: a dictionary containing the RPOP records.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountRPOPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRPOPS", account)
}

// GetAccountRSIPs retrieves the list of Account RSIP (Remote SIP) records.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// RSIP records are used for registering on remote SIP servers or managing remote
// SIP-related configurations for the Account.
//
// This method executes the GETACCOUNTRSIPS CLI command.
//
// Returns:
//   - map[string]any: a dictionary containing the RSIP records.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountRSIPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRSIPS", account)
}

// GetAccountRights retrieves the access rights granted to the specified Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// The method produces an array listing all Account Server or Domain Access rights.
// These rights determine what administrative or service-level actions the user
// is permitted to perform within the server or their domain.
//
// This method executes the GETACCOUNTRIGHTS CLI command.
//
// Returns:
//   - []string: a slice of access right strings.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountRights(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTRIGHTS", account)
}

// GetAccountSettings retrieves the explicitly set Account settings.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// Only the explicitly set (non-default) settings are included in the result.
// To get the combined set of explicit and default settings, use
// GetAccountEffectiveSettings instead.
//
// Note: All users can send the GETACCOUNTSETTINGS command for their own Accounts.
//
// This method executes the GETACCOUNTSETTINGS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Account settings.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTSETTINGS", Atom(account))
}

// GetAccountSignalRules retrieves the list of Account Signal Rules.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// Signal Rules are used to manage real-time communication events.
//
// This method executes the GETACCOUNTSIGNALRULES CLI command.
//
// Returns:
//   - [][SignalRule]: a slice of SignalRule structures.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountSignalRules(account string) ([]SignalRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSignalRules("GETACCOUNTSIGNALRULES", account)
}

// GetAccountTelnums retrieves the list of telephone numbers assigned to the Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// These numbers are used for telephony services such as PBX routing or SMS.
//
// This method executes the GETACCOUNTTELNUMS CLI command.
//
// Returns:
//   - []string: a slice of telephone number strings.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountTelnums(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTTELNUMS", account)
}

// IsAccountExists checks for the existence of an account on the server.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// This is a helper method that does not correspond to a single CLI command.
// It internally calls GetAccountLocation to verify existence. If the account
// exists, the server returns its storage path; otherwise, it returns an error.
//
// Returns:
//   - (true, nil) if the account is found.
//   - (false, nil) if the account is confirmed to not exist.
//   - (false, error) if the check fails due to communication or permission issues.
func (cli *Cli) IsAccountExists(account string) (bool, error) {
	if account == "" {
		return false, fmt.Errorf("account name is required")
	}

	location, err := cli.GetAccountLocation(account)
	if err != nil {
		return false, err
	}
	return location != "", nil
}

// KillAccountSessions interrupts all active sessions for the specified Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// This command immediately terminates all active connections and sessions
// associated with the account, including POP, IMAP, FTP, WebUser,
// AirSync, and others.
//
// This method executes the KILLACCOUNTSESSIONS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) KillAccountSessions(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("KILLACCOUNTSESSIONS", account)
}

// ModifyAccountTelnums changes the telephone numbers assigned to the Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - params: a dictionary specifying the operation to perform.
//
// The params dictionary must contain an "op" key with one of the following values:
//   - "add": requires a "telnum" key. Atomically adds the number. Returns an
//     error if the number already exists.
//   - "del": requires a "telnum" key. Atomically removes the number. Returns
//     an error if the number is not found.
//   - "pop": removes the first assigned number atomically.
//
// This method executes the MODIFYACCOUNTTELNUMS CLI command.
//
// Returns:
//   - map[string]any: for the "pop" operation, may contain the "telnum" key with the removed number.
//   - error: an error if the command fails.
func (cli *Cli) ModifyAccountTelnums(account string, params map[string]any) (map[string]any, error) {
	if account == "" || params == nil {
		return nil, fmt.Errorf("account name and parameters dictionary are required")
	}
	return cli.getMapAny("MODIFYACCOUNTTELNUMS", account, params)
}

// RenameAccount renames an existing account and optionally moves its data.
//
// Parameters:
//   - oldName: the current name of the account. The name can include the Domain name.
//   - newName: the new name for the account. The name can include the Domain name.
//   - storage: an optional "storage mount point" directory for the moved account data
//     (without the .mnt suffix).
//
// This method executes the RENAMEACCOUNT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameAccount(oldName, newName, storage string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new names are required")
	}

	const cmd = "RENAMEACCOUNT"
	if storage != "" {
		return cli.QueryNV(cmd, oldName, "INTO", newName, "PATH", storage)
	} else {
		return cli.QueryNV(cmd, oldName, "INTO", newName)
	}
}

// SetAccountACL modifies the Access Control List (ACL) for Account Access Rights.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - acl: a dictionary where keys are identifiers (e.g., "anyone" or "user@domain")
//     and values are strings with access right symbols.
//   - authAccount: an optional name of an Account on whose behalf the operation
//     is executed. Requires Admin rights for the target account.
//
// Modification rules for the value strings:
//   - starts with "+": rights are added to the existing set.
//   - starts with "-": rights are removed from the existing set.
//   - otherwise: the string replaces the current set of rights.
//   - empty string result: the ACL entry for that identifier is removed.
//
// This method executes the SETACCOUNTACL CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountACL(account string, acl map[string]any, authAccount string) error {
	if account == "" || acl == nil {
		return fmt.Errorf("account name and acl dictionary are required")
	}

	const cmd = "SETACCOUNTACL"
	if authAccount != "" {
		return cli.QueryNV(cmd, Atom(account), "AUTH", authAccount, acl)
	} else {
		return cli.QueryNV(cmd, Atom(account), acl)
	}
}

// SetAccountAccessModes updates the allowed access protocols/modes for an Account.
//
// Parameters:
//   - account: the Account name.
//   - modes: a slice of strings representing allowed protocols (e.g., "IMAP", "POP", "HTTP").
//
// Behavior:
//   - If modes is nil: the setting is removed (inherited from Domain/Server defaults).
//   - If modes is empty: access is restricted to "None" (explicitly blocking all protocols).
//   - If modes contains values: an array is constructed with a prefix (30) followed
//     by the specified modes to satisfy the CommuniGate Pro CLI dictionary format.
//
// This is a high-level helper that internally uses UpdateAccountSettings to
// modify the "AccessModes" key.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountAccessModes(account string, modes []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	var value any
	switch {
	case modes == nil:
		value = nil
	case len(modes) == 0:
		value = "None"
	default:
		raw := make([]any, 0, len(modes)+1)
		raw = append(raw, 30)
		for _, m := range modes {
			raw = append(raw, m)
		}
		value = raw
	}
	return cli.UpdateAccountSettings(account, map[string]any{"AccessModes": value})
}

// SetAccountAliases sets the list of aliases for the specified Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - aliases: a slice of strings containing the new Account alias names.
//
// This is an "all-or-nothing" operation: all existing aliases for the Account are
// removed and replaced with the ones provided. To clear all aliases, pass an
// empty (but non-nil) slice.
//
// This method executes the SETACCOUNTALIASES CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountAliases(account string, aliases []string) error {
	if account == "" || aliases == nil {
		return fmt.Errorf("account name and aliases array are required")
	}
	return cli.QueryNV("SETACCOUNTALIASES", account, aliases)
}

// SetAccountMailRules sets the list of Account Queue Mail Rules.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - rules: a slice of [MailRule] structures containing the new rules.
//
// This command replaces all existing Mail Rules for the account with the provided set.
// To clear all rules, pass an empty (but non-nil) slice. Mail Rules are processed by
// the Enqueuer component when a message is delivered to the account's INBOX.
//
// Note: Any user can use this method to modify rules for their own account.
//
// This method executes the SETACCOUNTMAILRULES CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountMailRules(account string, rules []MailRule) error {
	if account == "" || rules == nil {
		return fmt.Errorf("account name and rules array are required")
	}
	return setRulesGeneric(cli, "SETACCOUNTMAILRULES", account, rules)
}

// SetAccountPassword updates the password for the specified Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - password: the new password string. If empty and a tag is provided, the
//     corresponding application-specific password is removed.
//   - method: an optional Account Access Mode. "SIP" modifies the Alternative SIP Password;
//     "RADIUS" modifies the Alternative RADIUS Password. Otherwise, the primary
//     CommuniGate Password is modified.
//   - tag: an optional tag for an application-specific password.
//   - check: if true, the operation succeeds only if the new password matches the
//     server's size and complexity restrictions.
//
// Any user can modify their own password (complexity check is always enforced for self-service updates).
//
// This method executes the SETACCOUNTPASSWORD CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountPassword(account, password, method, tag string, check bool) error {
	if account == "" || password == "" {
		return fmt.Errorf("account name and password are required")
	}

	args := make([]any, 0, 8)
	args = append(args, account, "PASSWORD", password)
	if method != "" {
		args = append(args, "METHOD", method)
	}
	if tag != "" {
		args = append(args, "NAME", tag)
	}
	if check {
		args = append(args, "CHECK")
	}
	return cli.QueryNV("SETACCOUNTPASSWORD", args...)
}

// SetAccountPrefs sets the Preferences for the specified Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - settings: a dictionary of Account Preferences.
//
// This command completely replaces all existing Account Preferences with the provided ones.
// Preferences typically include user interface settings (WebUser), language parameters,
// time zones, and client-specific configurations.
//
// This method executes the SETACCOUNTPREFS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountPrefs(account string, settings map[string]any) error {
	if account == "" || settings == nil {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTPREFS", account, settings)
}

// SetAccountRIMAPs sets the Account RIMAP records.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - records: a dictionary containing the Account RIMAP records.
//
// All existing Account RIMAP records are removed and replaced with the provided dictionary.
//
// This method executes the SETACCOUNTRIMAPS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountRIMAPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTRIMAPS", account, records)
}

// SetAccountRPOPs sets the Account RPOP records.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - records: a dictionary containing the Account RPOP records.
//
// All existing Account RPOP records are removed and replaced with the provided dictionary.
//
// This method executes the SETACCOUNTRPOPS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountRPOPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTRPOPS", account, records)
}

// SetAccountRSIPs sets the Account RSIP records.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - records: a dictionary containing the Account RSIP records.
//
// All existing Account RSIP records are removed and replaced with the provided dictionary.
//
// This method executes the SETACCOUNTRSIPS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountRSIPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTRSIPS", account, records)
}

// SetAccountSettings changes the Account settings.
//
// Parameters:
//   - account: the Account name.
//   - settings: a dictionary of Account settings.
//
// All existing Account settings are removed and replaced with the provided dictionary.
//
// This method executes the SETACCOUNTSETTINGS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountSettings(account string, settings map[string]any) error {
	if account == "" || settings == nil {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTSETTINGS", account, settings)
}

// SetAccountSignalRules sets the Account Signal Rules.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - rules: a slice of [SignalRule] structures containing the new Signal Rules.
//
// All existing Account Signal Rules are removed and replaced with the provided slice.
//
// This method executes the SETACCOUNTSIGNALRULES CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountSignalRules(account string, rules []SignalRule) error {
	if account == "" || rules == nil {
		return fmt.Errorf("account name and rules array are required")
	}
	return setRulesGeneric(cli, "SETACCOUNTSIGNALRULES", account, rules)
}

// SetAccountTelnums assigns telephone numbers to the Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - telnums: a slice containing the telephone number strings.
//
// All existing numbers assigned to the Account are removed and replaced with the provided slice.
//
// This method executes the SETACCOUNTTELNUMS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountTelnums(account string, telnums []string) error {
	if account == "" || telnums == nil {
		return fmt.Errorf("account name and telnums array are required")
	}
	return cli.QueryNV("SETACCOUNTTELNUMS", account, telnums)
}

// SetAccountType changes the Account type.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - accType: the new Account type (e.g., MultiMailbox, AGrade).
//
// Note: the current Account type must also belong to the target type set for the
// operation to succeed.
//
// This method executes the SETACCOUNTTYPE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetAccountType(account string, accType AccountType) error {
	if account == "" || accType == TypeDefault {
		return fmt.Errorf("account name and account type are required")
	}
	return cli.QueryNV("SETACCOUNTTYPE", account, accType.String())
}

// UnblockAccount removes all blocking restrictions from an existing Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//
// This method executes the UNBLOCKACCOUNT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UnblockAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("UNBLOCKACCOUNT", account)
}

// UpdateAccountInfo updates the Account "info" fields.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - info: a dictionary containing the Account "info" fields to be updated or added.
//
// This method performs a partial update of the Account "info" dictionary.
// Only the keys provided in the info dictionary are modified; all other existing
// keys in the Account Information dictionary remain unchanged.
//
// This method executes the UPDATEACCOUNTINFO CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateAccountInfo(account string, info map[string]any) error {
	if account == "" || len(info) == 0 {
		return fmt.Errorf("account name and info are required")
	}
	return cli.QueryNV("UPDATEACCOUNTINFO", account, info)
}

// UpdateAccountMailRule updates or creates a specific Mail Rule for the Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - rule: a [MailRule] structure containing the rule definition.
//
// Behavior depends on the content of the rule structure:
//   - if only priority and name are provided (fewer than 4 rule elements): the method
//     updates the priority of an existing rule with the given name. If the rule does
//     not exist, an error is returned.
//   - if 4 or more elements are provided (conditions and actions included): the entire
//     rule is stored. If a rule with the same name already exists, it is replaced.
//
// This method executes the UPDATEACCOUNTMAILRULE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateAccountMailRule(account string, rule MailRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return updateRuleGeneric(cli, "UPDATEACCOUNTMAILRULE", account, rule)
}

// UpdateAccountPrefs modifies the Account Preferences.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - settings: a dictionary of Account Preferences.
//
// This method performs a partial update of the Account Preferences. Omitted elements
// remain unmodified. If a value is set to the string "default", that specific
// preference is removed, causing the account to inherit the default value from
// its Domain or the System.
//
// This method executes the UPDATEACCOUNTPREFS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateAccountPrefs(account string, settings map[string]any) error {
	if account == "" || len(settings) == 0 {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("UPDATEACCOUNTPREFS", account, settings)
}

// UpdateAccountSettings updates the Account settings.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - settings: a dictionary of Account settings.
//
// This method performs a partial update of the Account settings dictionary. Omitted
// settings remain unmodified. If a setting value is specified as the string "default",
// the specific Account setting is removed, allowing the account to use the default
// value (typically inherited from the Domain or Class).
//
// This method executes the UPDATEACCOUNTSETTINGS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {
	if account == "" || len(settings) == 0 {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("UPDATEACCOUNTSETTINGS", account, settings)
}

// UpdateAccountSignalRule updates or creates a specific Signal Rule for the Account.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - rule: a [SignalRule] structure containing the rule definition.
//
// Behavior depends on the content of the rule structure:
//   - if only priority and name are provided (fewer than 4 rule elements): the method
//     updates the priority of an existing Signal Rule with the given name. If the rule
//     does not exist, an error is returned.
//   - if 4 or more elements are provided (conditions and actions included): the entire
//     rule is stored as a new Signal Rule. If a rule with the same name already exists,
//     it is replaced.
//
// This method executes the UPDATEACCOUNTSIGNALRULE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateAccountSignalRule(account string, rule SignalRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return updateRuleGeneric(cli, "UPDATEACCOUNTSIGNALRULE", account, rule)
}

// UpdateScheduledTask sets the Account Scheduled Task record.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - task: a dictionary containing the Scheduled Task data (id, program/script, etc.).
//
// Behavior:
//   - if the task "id" does not exist, a new task is created.
//   - if the "program" or "script" element is missing, an existing task with the given "id" is deleted.
//   - tasks can be one-time (no "period") or recurring (period is set to "day", "week",
//     "month", "year", or a number of seconds).
//   - the task executes either a Real-Time Application or a Synchronous Script.
//
// This method executes the UPDATESCHEDULEDTASK CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateScheduledTask(account string, task map[string]any) error {
	if account == "" || len(task) == 0 {
		return fmt.Errorf("account name and task dictionary are required")
	}
	return cli.QueryNV("UPDATESCHEDULEDTASK", account, task)
}

// VerifyAccountIdentity checks if the Account is authorized to use a specific
// 'From:' header value.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - identity: the 'From:' header value to verify (e.g., "Name <user@domain.com>").
//
// This method executes the VERIFYACCOUNTIDENTITY CLI command.
//
// Returns:
//   - (true, nil) if the identity is verified and allowed.
//   - (false, error) if the verification fails or the server returns an error.
func (cli *Cli) VerifyAccountIdentity(account, identity string) (bool, error) {
	if account == "" || identity == "" {
		return false, fmt.Errorf("account name and identity are required")
	}
	err := cli.QueryNV("VERIFYACCOUNTIDENTITY", account, "FOR", identity)
	return err == nil, err
}

// VerifyAccountPassword verifies the Account password.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//   - password: the password to check (clear text).
//
// This method executes the VERIFYACCOUNTPASSWORD CLI command.
//
// Returns:
//   - (true, nil) if the password is correct.
//   - (false, error) if the password is incorrect or a server error occurred.
func (cli *Cli) VerifyAccountPassword(account, password string) (bool, error) {
	if account == "" || password == "" {
		return false, fmt.Errorf("account name and password are required")
	}
	err := cli.QueryNV("VERIFYACCOUNTPASSWORD", account, "PASSWORD", password)
	return err == nil, err
}
