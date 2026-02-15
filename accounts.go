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
)

var accountTypeNames = [...]string{
	"", "MultiMailbox", "TextMailbox", "MailDirMailbox",
	"SlicedMailbox", "AGrade", "BGrade", "CGrade", "DGrade",
}

func (t AccountType) String() string {
	if int(t) < 0 || int(t) >= len(accountTypeNames) {
		return ""
	}
	return accountTypeNames[t]
}

// CheckAccountPassword verifies if the provided password matches the Account's password.
//
// Parameters:
//   - account: This parameter specifies the name of the Account to check.
//   - password: The plain-text password to verify against the server record.
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
//   - account: This parameter specifies the name for the new Account.
//     The name can contain the @ symbol followed by the Domain name;
//     otherwise, the Account is created in the specified Domain (or
//     the administrator Domain if not specified).
//   - settings: This optional parameter specifies the initial Account settings.
//     If specified, it is used to modify the Template settings for the target Domain.
//   - accType: This optional parameter specifies the type of the Account to create
//     (e.g., MultiMailbox, MailDirMailbox). If [TypeDefault] is used,
//     a MultiMailbox-type Account is created.
//   - storage: This optional parameter specifies the "storage mount point"
//     directory for the Account data (without the .mnt suffix).
//   - legacy: This optional flag tells the system to create an Account
//     with a Legacy (visible for legacy mailers) INBOX.
//
// This method executes the CREATEACCOUNT CLI command.
//
// Returns an error if the account name is empty or if the server
// fails to execute the command.
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
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// This method executes the DELETEACCOUNT CLI command.
//
// Returns an error if the account name is empty or if the server
// fails to execute the command.
func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("DELETEACCOUNT", account)
}

// DeleteAccountMailRule removes a Mail Rule from an existing Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rule: This string parameter specifies the name of the [MailRule] to be removed.
//
// If the specified Mail Rule does not exist, the command does nothing and
// does not return an error.
//
// This method executes the UPDATEACCOUNTMAILRULE CLI command with the DELETE keyword.
//
// Returns an error if the account name or rule name is empty, or if the
// server fails to execute the command.
func (cli *Cli) DeleteAccountMailRule(account, rule string) error {
	if account == "" || rule == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	return cli.QueryNV("UPDATEACCOUNTMAILRULE", account, "DELETE", rule)
}

// DeleteAccountSignalRule removes a Signal Rule from an existing Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rule: This string parameter specifies the name of the [SignalRule] to be removed.
//
// If the specified Signal Rule does not exist, the command does nothing and
// does not return an error.
//
// This method executes the UPDATEACCOUNTSIGNALRULE CLI command with the DELETE keyword.
//
// Returns an error if the account name or rule name is empty, or if the
// server fails to execute the command.
func (cli *Cli) DeleteAccountSignalRule(account, rule string) error {
	if account == "" || rule == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	return cli.QueryNV("UPDATEACCOUNTSIGNALRULE", account, "DELETE", rule)
}

// GetAccountACL retrieves the Access Control Lists (ACLs) for an Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account (target Account).
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - authAccount: This optional parameter specifies the name of an Account on whose
//     behalf the operation should be executed.
//
// If the authAccount is specified, the ACL info is returned only if the specified
// Account has the Admin access right for the target Account.
//
// The method produces a dictionary with the ACL elements.
//
// This method executes the GETACCOUNTACL CLI command.
//
// Returns a map[string]any containing the ACL elements, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountACL(account, authAccount string) (res map[string]any, err error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTACL"
	if authAccount != "" {
		res, err = cli.getMapAny(cmd, Atom(account), "AUTH", authAccount)
	} else {
		res, err = cli.getMapAny(cmd, Atom(account))
	}
	return
}

// GetAccountACLRights retrieves the effective access rights for a specific Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account (target Account).
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - authAccount: This parameter specifies the name of an Account whose effective
//     access rights for the target Account should be retrieved.
//
// The method produces a string containing the effective access rights.
//
// This method executes the GETACCOUNTACLRIGHTS CLI command.
//
// Returns a string with the effective rights, or an error if target account
// or auth account names are empty, or if the server fails to execute the command.
func (cli *Cli) GetAccountACLRights(account, authAccount string) (string, error) {
	if account == "" || authAccount == "" {
		return "", fmt.Errorf("target account and auth account names are required")
	}
	return cli.getString("GETACCOUNTACLRIGHTS", Atom(account), "AUTH", authAccount)
}

// GetAccountAccessModes retrieves the enabled services (AccessModes) for an Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
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
// Returns a slice of strings containing the names of enabled services, or an
// error if the account name is empty or settings cannot be retrieved.
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
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces an array containing the Account alias names.
//
// This method executes the GETACCOUNTALIASES CLI command.
//
// Returns a slice of strings with the alias names, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountAliases(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTALIASES", account)
}

// GetAccountEffectivePrefs retrieves the effective Account Preferences.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces a dictionary containing the Account Preferences.
// Both the explicitly set and the default settings are included in the result.
//
// Note: All users can use this command to retrieve their own effective Preferences.
//
// This method executes the GETACCOUNTEFFECTIVEPREFS CLI command.
//
// Returns a map[string]any with the preferences, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountEffectivePrefs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTEFFECTIVEPREFS", account)
}

// GetAccountEffectiveSettings retrieves the effective Account settings.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name. The single asterisk (*) symbol
//     can be used to specify the current authenticated Account.
//
// The method produces a dictionary containing the Account settings.
// Both the explicitly set and the default Account settings are included
// in the dictionary.
//
// Note: All users can send this command for their own Accounts.
//
// This method executes the GETACCOUNTEFFECTIVESETTINGS CLI command.
//
// Returns a map[string]any with the settings, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountEffectiveSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTEFFECTIVESETTINGS", Atom(account))
}

// GetAccountInfo retrieves multiple elements from the Account "info" dictionary.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name. The single asterisk (*) symbol
//     can be used to specify the current authenticated Account.
//   - keys: An optional array specifying the names of the info keys to retrieve
//     (e.g., "LastLogin", "LastAddress").
//
// Note: Key names are case-sensitive. If Account "info" data is stored in .info
// dictionary files, do NOT include the hash (#) symbol in the key names.
//
// The method produces a dictionary containing all those "info" elements that
// exist and are specified in the keys array. If the keys array is empty,
// the entire "info" dictionary is returned.
//
// Note: All users can use this command to retrieve elements from their own
// Account "info" data.
//
// This method executes the GETACCOUNTINFO CLI command.
//
// Returns a map[string]any with the requested info elements, or an error if
// the account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountInfo(account string, keys []string) (res map[string]any, err error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTINFO"
	if len(keys) > 0 {
		res, err = cli.getMapAny(cmd, Atom(account), keys)
	} else {
		res, err = cli.getMapAny(cmd, Atom(account))
	}
	return
}

// GetAccountLocation retrieves the storage path for an Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces a string containing the Account file directory path
// (for multi-mailbox Accounts) or the Account INBOX Mailbox path
// (for single-mailbox Accounts).
//
// The path is relative to the file directory of the Account Domain.
//
// This method executes the GETACCOUNTLOCATION CLI command.
//
// Returns a string with the Account file path, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountLocation(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	return cli.getString("GETACCOUNTLOCATION", account)
}

// GetAccountMailRules retrieves the list of Account Queue Rules.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces an array of the Queue Rules specified for the Account.
//
// This method executes the GETACCOUNTMAILRULES CLI command.
//
// Returns a slice of [MailRule] structures, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountMailRules(account string) ([]MailRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMailRules("GETACCOUNTMAILRULES", account)
}

// GetAccountOneInfo retrieves a single specific element from the Account "info" dictionary.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name. The single asterisk (*) symbol
//     can be used to specify the current authenticated Account.
//   - key: This parameter specifies the name of the requested "info" element.
//
// Note: The key name is case-sensitive. Do NOT include the hash (#) symbol
// in the key name.
//
// The method produces the specified "info" element. If the element is not found,
// the output is an empty string.
//
// Note: All users can use this command to retrieve elements from their own
// Account "info" data.
//
// This method executes the GETACCOUNTINFO CLI command with the Key keyword.
//
// Returns the specified info element (as any), or an error if the account
// or key names are empty, or if the server fails to execute the command.
func (cli *Cli) GetAccountOneInfo(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTINFO", Atom(account), "Key", key)
}

// GetAccountOneSetting retrieves a single setting from the effective Account settings list.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name. The single asterisk (*) symbol
//     can be used to specify the current authenticated Account.
//   - key: This parameter specifies the name of the setting to read.
//
// The method produces an object which can be a string, an array, a dictionary,
// or a null-object representing the requested setting.
//
// Note: All users can send the GETACCOUNTONESETTING command for their own Accounts.
//
// This method executes the GETACCOUNTONESETTING CLI command.
//
// Returns the setting value (as any), or an error if the account or key
// names are empty, or if the server fails to execute the command.
func (cli *Cli) GetAccountOneSetting(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTONESETTING", Atom(account), key)
}

// GetAccountPlainPassword retrieves the clear-text password for an Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//   - method: This optional parameter specifies the authentication method.
//     Supported values include "SIP" and "RADIUS".
//
// If the method is not specified (empty string), the command retrieves the
// main Account password. If a method is specified, the server retrieves
// the password specifically configured for that authentication context.
//
// Note: This command only works if the password is stored using a reversible
// encryption method. It will not return the plain-text password if a
// non-reversible hash algorithm is used for storage.
//
// This method executes the GETACCOUNTPLAINPASSWORD CLI command.
//
// Returns a string containing the plain-text password, or an error if
// the account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountPlainPassword(account string, method string) (s string, err error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTPLAINPASSWORD"
	if method == "SIP" || method == "RADIUS" {
		s, err = cli.getString(cmd, account, "METHOD", method)
	} else {
		s, err = cli.getString(cmd, account)
	}
	return
}

// GetAccountPrefs retrieves the Account Preferences.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces a dictionary containing the Account Preferences.
// Unlike effective preferences, this returns only the explicitly set values.
//
// Note: All users can use the GETACCOUNTPREFS command to retrieve their
// own Account Preferences.
//
// This method executes the GETACCOUNTPREFS CLI command.
//
// Returns a map[string]any with the preferences, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountPrefs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTPREFS", account)
}

// GetAccountPresence retrieves the current "presence" status of an Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method returns the presence information in one of the following formats:
//   - []any (array of two strings): The first element is the presence status
//     (e.g., "Available", "Away"), and the second is the custom status message.
//   - string: The presence status (returned if no custom status message is set).
//   - nil: Returned if the presence status is not set at all.
//
// This method executes the GETACCOUNTPRESENCE CLI command.
//
// Returns the presence data as any, or an error if the account name is empty
// or the server fails to execute the command.
func (cli *Cli) GetAccountPresence(account string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.Query("GETACCOUNTPRESENCE", account)
}

// GetAccountRIMAPs retrieves the list of Account RIMAP (Remote IMAP) records.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces a dictionary where each element represents a RIMAP record
// configured for the Account. RIMAP (Remote IMAP) records are used to poll
// external mailboxes and fetch messages into the local Account.
//
// This method executes the GETACCOUNTRIMAPS CLI command.
//
// Returns a map[string]any containing the RIMAP records, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountRIMAPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRIMAPS", account)
}

// GetAccountRPOPs retrieves the list of Account RPOP (Remote POP) records.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces a dictionary where each element represents an RPOP record
// configured for the Account. RPOP (Remote POP) records are used to collect
// mail from external POP3 servers into the local Account.
//
// This method executes the GETACCOUNTRPOPS CLI command.
//
// Returns a map[string]any containing the RPOP records, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountRPOPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRPOPS", account)
}

// GetAccountRSIPs retrieves the list of Account RSIP (Remote SIP) records.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces a dictionary where each element represents an RSIP record
// configured for the Account. RSIP records are typically used for registering
// on remote SIP servers or managing remote SIP-related configurations for
// the account.
//
// This method executes the GETACCOUNTRSIPS CLI command.
//
// Returns a map[string]any containing the RSIP records, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountRSIPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRSIPS", account)
}

// GetAccountRights retrieves the access rights granted to the specified user.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces an array listing all Account Server or Domain Access rights.
// These rights determine what administrative or service-level actions the user
// is permitted to perform within the server or their domain.
//
// This method executes the GETACCOUNTRIGHTS CLI command.
//
// Returns a slice of strings containing the access rights, or an error if
// the account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountRights(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTRIGHTS", account)
}

// GetAccountSettings retrieves the explicitly set Account settings.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name. The single asterisk (*) symbol
//     can be used to specify the current authenticated Account.
//
// The method produces a dictionary containing the Account settings.
// Only the explicitly set (non-default) settings are included in the result.
// To get the combined set of explicit and default settings, use
// GetAccountEffectiveSettings instead.
//
// Note: All users can send the GETACCOUNTSETTINGS command for their own Accounts.
//
// This method executes the GETACCOUNTSETTINGS CLI command.
//
// Returns a map[string]any with the settings, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTSETTINGS", Atom(account))
}

// GetAccountSignalRules retrieves the list of Account Signal Rules.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces an array of the Signal Rules specified for the Account.
// Signal Rules are used to manage real-time communication events.
//
// This method executes the GETACCOUNTSIGNALRULES CLI command.
//
// Returns a slice of [SignalRule] structures, or an error if the
// account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountSignalRules(account string) ([]SignalRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSignalRules("GETACCOUNTSIGNALRULES", account)
}

// GetAccountTelnums retrieves the list of telephone numbers assigned to the Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// The method produces an array containing all telephone numbers assigned to
// the specified Account. These numbers are used for telephony services,
// such as PBX routing or SMS.
//
// This method executes the GETACCOUNTTELNUMS CLI command.
//
// Returns a slice of strings containing the telephone numbers, or an error
// if the account name is empty or the server fails to execute the command.
func (cli *Cli) GetAccountTelnums(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTTELNUMS", account)
}

// IsAccountExists checks for the existence of an account on the server.
//
// Parameters:
//   - account: This parameter specifies the name of the account to check.
//     The name can include the Domain name.
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
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//
// This command immediately terminates all active connections and sessions
// associated with the account, including POP, IMAP, FTP, WebUser,
// AirSync, and others.
//
// This method executes the KILLACCOUNTSESSIONS CLI command.
//
// Returns an error if the account name is empty or if the server
// fails to terminate the sessions.
func (cli *Cli) KillAccountSessions(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("KILLACCOUNTSESSIONS", account)
}

// ModifyAccountTelnums changes the telephone numbers assigned to the Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - params: A dictionary specifying the operation to perform.
//
// The params dictionary must contain an "op" key with one of the following values:
//   - "add": Requires a "telnum" key. Atomically adds the number. Returns an
//     error if the number already exists.
//   - "del": Requires a "telnum" key. Atomically removes the number. Returns
//     an error if the number is not found.
//   - "pop": Removes the first assigned number atomically.
//
// This method executes the MODIFYACCOUNTTELNUMS CLI command.
//
// Returns a map[string]any which, for the "pop" operation, may contain the
// "telnum" key with the removed number. Returns an error if parameters are
// invalid or the server operation fails.
func (cli *Cli) ModifyAccountTelnums(account string, params map[string]any) (map[string]any, error) {
	if account == "" || params == nil {
		return nil, fmt.Errorf("account name and parameters dictionary are required")
	}
	return cli.getMapAny("MODIFYACCOUNTTELNUMS", account, params)
}

// RenameAccount renames an existing account and optionally moves its data.
//
// Parameters:
//   - oldName: The current name of the account (can include domain).
//   - newName: The new name for the account (can include domain).
//   - storage: (Optional) Specifies the "storage mount Point" directory for the
//     moved account data. The name should be specified without the .mnt suffix.
//
// This command is used to rename accounts within the same domain or across
// domains, and can also be used to relocate account storage folders.
//
// This method executes the RENAMEACCOUNT CLI command.
//
// Returns an error if names are empty or if the server fails to complete
// the rename/move operation.
func (cli *Cli) RenameAccount(oldName, newName, storage string) (err error) {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new names are required")
	}

	const cmd = "RENAMEACCOUNT"
	if storage != "" {
		err = cli.QueryNV(cmd, oldName, "INTO", newName, "PATH", storage)
	} else {
		err = cli.QueryNV(cmd, oldName, "INTO", newName)
	}
	return
}

// SetAccountACL modifies the Access Control List (ACL) for Account Access Rights.
//
// Parameters:
//   - account: Name of the target Account. Use "*" for the current authenticated user.
//   - authAccount: (Optional) Name of an Account on whose behalf the operation
//     is executed. Requires 'Admin' rights for the target account.
//   - acl: A dictionary where keys are identifiers (e.g., "anyone" or "user@domain")
//     and values are strings with access right symbols.
//
// Modification Rules (Value strings):
//   - Starts with "+": Rights are added to the existing set.
//   - Starts with "-": Rights are removed from the existing set.
//   - Otherwise: The string replaces the current set of rights.
//   - Empty string result: The ACL entry for that identifier is removed.
//
// This method executes the SETACCOUNTACL CLI command.
//
// Returns an error if account or acl are missing, or if the server operation fails.
func (cli *Cli) SetAccountACL(account, authAccount string, acl map[string]any) (err error) {
	if account == "" || acl == nil {
		return fmt.Errorf("account name and acl dictionary are required")
	}

	const cmd = "SETACCOUNTACL"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "AUTH", authAccount, acl)
	} else {
		err = cli.QueryNV(cmd, Atom(account), acl)
	}
	return
}

// SetAccountAccessModes updates the allowed access protocols/modes for an Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//   - modes: A slice of strings representing allowed protocols (e.g., "IMAP", "POP", "HTTP").
//
// Behavior:
//   - If modes is nil: The setting is removed (inherited from Domain/Server defaults).
//   - If modes is empty: Access is restricted to "None" (explicitly blocking all protocols).
//   - If modes contains values: An array is constructed with a prefix (30) followed
//     by the specified modes to satisfy the CommuniGate Pro CLI dictionary format.
//
// This is a high-level helper that internally uses UpdateAccountSettings to
// modify the "AccessModes" key.
//
// Returns an error if the account name is empty or the update operation fails.
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
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - aliases: A slice of strings containing the new Account alias names.
//
// Behavior:
//
//	This is an "all-or-nothing" operation. All existing aliases for the
//	Account are removed and replaced with the ones provided in the aliases
//	array. To clear all aliases, pass an empty (but non-nil) slice.
//
// This method executes the SETACCOUNTALIASES CLI command.
//
// Returns an error if the account name is empty, the aliases slice is nil,
// or the server fails to update the alias list.
func (cli *Cli) SetAccountAliases(account string, aliases []string) error {
	if account == "" || aliases == nil {
		return fmt.Errorf("account name and aliases array are required")
	}
	return cli.QueryNV("SETACCOUNTALIASES", account, aliases)
}

// SetAccountMailRules sets the list of Account Queue Mail Rules.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rules: A slice of [MailRule] structures containing the new rules.
//
// Behavior:
//
//	This command replaces all existing Mail Rules for the account with the
//	provided set. To clear all rules, pass an empty (but non-nil) slice.
//	Mail Rules are processed by the Enqueuer component when a message
//	is delivered to the account's INBOX.
//
// Note: Any user can use this method to modify rules for their own account.
//
// This method executes the SETACCOUNTMAILRULES CLI command.
//
// Returns an error if the account name is empty, the rules slice is nil,
// or the server fails to update the rules.
func (cli *Cli) SetAccountMailRules(account string, rules []MailRule) error {
	if account == "" || rules == nil {
		return fmt.Errorf("account name and rules array are required")
	}
	return setRulesGeneric(cli, "SETACCOUNTMAILRULES", account, rules)
}

// SetAccountPassword updates the password for the specified Account.
//
// Parameters:
//   - account: Name of an existing Account (can include the Domain name).
//   - password: The new password string. If an empty string is provided
//     along with a tag, the corresponding application-specific password is removed.
//   - method: (Optional) Specifies the Account Access Mode.
//     "SIP" modifies the Alternative SIP Password.
//     "RADIUS" modifies the Alternative RADIUS Password.
//     Otherwise, the primary CommuniGate Password is modified.
//   - tag: (Optional) Specifies the tag for an application-specific password.
//   - check: If true, the operation succeeds only if the new password
//     matches the server's size and complexity restrictions.
//
// Any user can modify their own password (complexity check is always enforced for self-service updates).
//
// This method executes the SETACCOUNTPASSWORD CLI command.
//
// Returns an error if account/password are missing or if the server
// rejects the password change.
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

// SetAccountPrefs sets the Preferences for the specified account.
//
// Parameters:
//   - account: The name of an existing account (can include the domain name).
//   - settings: A dictionary (map[string]any) containing the new preference settings.
//
// Features:
//
//	This command completely replaces all old account preferences with the new ones provided.
//	Preferences typically include user interface settings (WebUser), language
//	parameters, time zones, and client-specific configurations.
//
// This method executes the SETACCOUNTPREFS CLI command.
//
// Returns an error if the account name is empty, the settings dictionary is nil,
// or if the server rejects the operation.
func (cli *Cli) SetAccountPrefs(account string, settings map[string]any) error {
	if account == "" || settings == nil {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTPREFS", account, settings)
}

// SetAccountRIMAPs sets the Account RIMAP records.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - records: A dictionary containing the Account RIMAP records.
//
// Behavior:
//
//	All old Account RIMAP records are removed. The new records are set
//	from the provided dictionary.
//
// This method executes the SETACCOUNTRIMAPS CLI command.
//
// Returns an error if the account name is empty, the records dictionary is nil,
// or if the server operation fails.
func (cli *Cli) SetAccountRIMAPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTRIMAPS", account, records)
}

// SetAccountRPOPs sets the Account RPOP records.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - records: A dictionary containing the Account RPOP records.
//
// Behavior:
//
//	All old Account RPOP records are removed. The new records are set
//	from the provided dictionary.
//
// This method executes the SETACCOUNTRPOPS CLI command.
//
// Returns an error if the account name is empty, the records dictionary is nil,
// or if the server operation fails.
func (cli *Cli) SetAccountRPOPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTRPOPS", account, records)
}

// SetAccountRSIPs sets the Account RSIP records.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - records: A dictionary containing the Account RSIP records.
//
// Behavior:
//
//	All old Account RSIP records are removed. The new records are set
//	from the provided dictionary.
//
// This method executes the SETACCOUNTRSIPS CLI command.
//
// Returns an error if the account name is empty, the records dictionary is nil,
// or if the server operation fails.
func (cli *Cli) SetAccountRSIPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTRSIPS", account, records)
}

// SetAccountSettings changes the Account settings.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//   - settings: A dictionary used to replace the Account settings dictionary.
//
// Behavior:
//
//	All old Account settings are removed. The settings are replaced with
//	the ones provided in the settings dictionary.
//
// This method executes the SETACCOUNTSETTINGS CLI command.
//
// Returns an error if the account name is empty, the settings dictionary
// is nil, or the server fails to update the settings.
func (cli *Cli) SetAccountSettings(account string, settings map[string]any) error {
	if account == "" || settings == nil {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTSETTINGS", account, settings)
}

// SetAccountSignalRules sets the Account Signal Rules.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rules: An array of [SignalRule] structures containing the new Signal Rules.
//
// Behavior:
//
//	All old Account Signal Rules are removed. The new rules are set
//	from the provided array.
//
// This method executes the SETACCOUNTSIGNALRULES CLI command.
//
// Returns an error if the account name is empty, the rules slice is nil,
// or if the server operation fails.
func (cli *Cli) SetAccountSignalRules(account string, rules []SignalRule) error {
	if account == "" || rules == nil {
		return fmt.Errorf("account name and rules array are required")
	}
	return setRulesGeneric(cli, "SETACCOUNTSIGNALRULES", account, rules)
}

// SetAccountTelnums assigns telephone numbers to the Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - telnums: An array containing the telephone number strings.
//
// Behavior:
//
//	All old numbers assigned to the Account are removed. The new numbers
//	are set from the provided array.
//
// This method executes the SETACCOUNTTELNUMS CLI command.
//
// Returns an error if the account name is empty, the telnums slice is nil,
// or if the server operation fails.
func (cli *Cli) SetAccountTelnums(account string, telnums []string) error {
	if account == "" || telnums == nil {
		return fmt.Errorf("account name and telnums array are required")
	}
	return cli.QueryNV("SETACCOUNTTELNUMS", account, telnums)
}

// SetAccountType changes the Account type.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - accType: This parameter specifies the new Account type (e.g., MultiMailbox, AGrade, etc.).
//
// Behavior:
//
//	The command updates the storage and management type of the account.
//	Note that the current Account type must also belong to the target
//	type set for the operation to succeed.
//
// This method executes the SETACCOUNTTYPE CLI command.
//
// Returns an error if the account name is empty, the account type is
// TypeDefault, or if the server operation fails.
func (cli *Cli) SetAccountType(account string, accType AccountType) error {
	if account == "" || accType == TypeDefault {
		return fmt.Errorf("account name and account type are required")
	}
	return cli.QueryNV("SETACCOUNTTYPE", account, accType.String())
}

// UpdateAccountInfo updates the Account "info" fields.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - info: A dictionary containing the Account "info" fields to be
//     updated or added.
//
// Behavior:
//
//	This method performs a partial update of the Account "info" dictionary.
//	It modifies only the keys provided in the info dictionary,
//	while leaving other existing keys in the Account Information dictionary
//	unchanged.
//
// This method executes the UPDATEACCOUNTINFO CLI command.
//
// Returns an error if the account name is empty, the info dictionary
// is empty, or if the server operation fails.
func (cli *Cli) UpdateAccountInfo(account string, info map[string]any) error {
	if account == "" || len(info) == 0 {
		return fmt.Errorf("account name and info are required")
	}
	return cli.QueryNV("UPDATEACCOUNTINFO", account, info)
}

// UpdateAccountMailRule updates or creates a specific Mail Rule for the Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rule: A [MailRule] structure containing the rule definition.
//
// Behavior:
//
//	The behavior depends on the content of the rule structure:
//	- If only priority and name are provided (fewer than 4 rule elements):
//	  The method updates the priority of an existing rule with the given name.
//	  If the rule does not exist, an error is returned.
//	- If 4 or more elements are provided (conditions and actions included):
//	  The entire rule is stored. If a rule with the same name already exists,
//	  it is replaced by this new definition.
//
// This method executes the UPDATEACCOUNTMAILRULE CLI command.
//
// Returns an error if the account name is empty or if the server operation fails.
func (cli *Cli) UpdateAccountMailRule(account string, rule MailRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return updateRuleGeneric(cli, "UPDATEACCOUNTMAILRULE", account, rule)
}

// UpdateAccountPrefs modifies the Account Preferences.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - settings: A dictionary containing the preference settings to update.
//
// Behavior:
//
//	This method performs a partial update of the Account Preferences.
//	Omitted elements remain unmodified. If a value is set to the string
//	"default", that specific preference is removed, causing the account
//	to inherit the default value from its Domain or the System.
//
// This method executes the UPDATEACCOUNTPREFS CLI command.
//
// Returns an error if the account name is empty, the settings dictionary
// is empty, or if the server operation fails.
func (cli *Cli) UpdateAccountPrefs(account string, settings map[string]any) error {
	if account == "" || len(settings) == 0 {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("UPDATEACCOUNTPREFS", account, settings)
}

// UpdateAccountSettings updates the Account settings.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - settings: A dictionary containing the settings to be updated.
//
// Behavior:
//
//	This method performs a partial update of the Account settings dictionary.
//	Omitted settings remain unmodified. If a setting value is specified
//	as the string "default", the specific Account setting is removed,
//	allowing the account to use the default value (typically inherited
//	from the Domain or Class).
//
// This method executes the UPDATEACCOUNTSETTINGS CLI command.
//
// Returns an error if the account name is empty, the settings dictionary
// is empty, or if the server operation fails.
func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {
	if account == "" || len(settings) == 0 {
		return fmt.Errorf("account name and settings dictionary are required")
	}
	return cli.QueryNV("UPDATEACCOUNTSETTINGS", account, settings)
}

// UpdateAccountSignalRule updates or creates a specific Signal Rule for the Account.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - rule: A [SignalRule] structure containing the rule definition.
//
// Behavior:
//
//	The behavior depends on the content of the rule structure:
//	- If only priority and name are provided (fewer than 4 rule elements):
//	  The method updates the priority of an existing Signal Rule with the
//	  given name. If the rule does not exist, an error is returned.
//	- If 4 or more elements are provided (conditions and actions included):
//	  The entire rule is stored as a new Signal Rule. If a rule with the
//	  same name already exists, it is replaced.
//
// This method executes the UPDATEACCOUNTSIGNALRULE CLI command.
//
// Returns an error if the account name is empty or if the server operation fails.
func (cli *Cli) UpdateAccountSignalRule(account string, rule SignalRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return updateRuleGeneric(cli, "UPDATEACCOUNTSIGNALRULE", account, rule)
}

// UpdateScheduledTask sets the Account Scheduled Task record.
//
// Parameters:
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - task: A dictionary containing the Scheduled Task data (id, program/script, etc.).
//
// Behavior:
//   - If the task 'id' does not exist, a new task is created.
//   - If the 'program' or 'script' element is missing, an existing task
//     with the given 'id' is deleted.
//   - Tasks can be one-time (no 'period') or recurring (period is set to
//     'day', 'week', 'month', 'year', or a number of seconds).
//   - The task executes either a Real-Time Application or a Synchronous Script.
//
// This method executes the UPDATESCHEDULEDTASK CLI command.
//
// Returns an error if the account name is empty, the task dictionary
// is empty, or if the server operation fails.
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
//   - account: This parameter specifies the name of an existing Account.
//     The name can include the Domain name.
//   - identity: The 'From:' header value to verify (e.g., "Name <user@domain.com>").
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
//   - account: This parameter specifies the name of an existing Account.
//   - password: The password string to check (clear text).
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
