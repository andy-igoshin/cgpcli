// # Domain Administration
//
// The Domains section provides extensive functionality for managing CommuniGate Pro
// domains and their internal objects. It handles the domain lifecycle, inheritance
// of settings, and mass-object discovery.
//
// Key capabilities include:
//   - Configuration Management: Handling [Domain Settings], [Effective Settings],
//     and Domain-level [Aliases].
//   - Inheritance Control: Managing [Account Defaults] and [Account Templates],
//     which dictate the initial state of new accounts within the domain.
//   - Object Discovery: Using [Cli.ListDomainObjects] to efficiently search and
//     categorize accounts, aliases, and forwarders using filters and cookies for pagination.
//   - Operational State: Suspending, resuming, and locating domain storage paths.
//   - Rules & Directory: Managing Domain-level Mail/Signal rules and synchronizing
//     data with the Directory Service.
//
// This section is vital for multi-tenant environments where administrators need
// to enforce domain-wide policies and manage large sets of objects.
package cgpcli

import (
	"fmt"
)

// ObjectType defines the bitmask for filtering domain objects (accounts, aliases, forwarders).
type ObjectType int

const (
	TypeAccount ObjectType = 1 << iota
	TypeAlias
	TypeForwarder
	TypeAll = TypeAccount | TypeAlias | TypeForwarder
)

// DomainObject represents a generic object within a domain,
// holding its name, type, and associated targets.
type DomainObject struct {
	Name    string     // The unique identifier of the object (e.g., account name).
	Type    ObjectType // The category: TypeAccount, TypeAlias, or TypeForwarder.
	Subtype string     // Object-specific data (e.g., file extension for accounts).
	Targets []string   // Destination addresses for Aliases and Forwarders.
}

// DomainObjectsResult contains a batch of domain objects and a cookie
// for paginated retrieval.
type DomainObjectsResult struct {
	Objects map[string]DomainObject
	Cookie  string
}

// CreateAccountStorage creates a "storage mount point" for new Accounts in the Domain.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - storage: This parameter specifies the "storage mount Point" name.
//
// Returns:
//   - error: An error if parameters are missing or the server fails to create the storage.
func (cli *Cli) CreateAccountStorage(domain, storage string) error {
	if domain == "" || storage == "" {
		return fmt.Errorf("domain name and storage path are required")
	}
	return cli.QueryNV("CREATEACCOUNTSTORAGE", domain, "PATH", storage)
}

// DeleteDirectoryRecords deletes Domain object records from the Directory.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the authenticated user Domain.
//
// Returns:
//   - error: An error if the operation fails or access is denied.
func (cli *Cli) DeleteDirectoryRecords(domain string) (err error) {
	const cmd = "DELETEDIRECTORYRECORDS"
	if domain != "" {
		err = cli.QueryNV(cmd, domain)
	} else {
		err = cli.QueryNV(cmd)
	}
	return
}

// GetAccountDefaultPrefs retrieves the Default Account Preferences for the specified Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//
// Returns:
//   - map[string]any: A dictionary with the default Preferences.
//   - error: An error if the retrieval fails.
func (cli *Cli) GetAccountDefaultPrefs(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTDEFAULTPREFS", Atom(domain))
}

// GetAccountDefaults retrieves the default Account settings for the specified Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//
// Returns:
//   - map[string]any: A dictionary with the default settings.
//   - error: An error if the retrieval fails.
func (cli *Cli) GetAccountDefaults(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTDEFAULTS", Atom(domain))
}

// GetAccountTemplate retrieves the Account Template settings for the Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//
// Returns:
//   - map[string]any: A dictionary with the Template settings.
//   - error: An error if the retrieval fails.
func (cli *Cli) GetAccountTemplate(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTTEMPLATE", Atom(domain))
}

// GetDomainAliases retrieves the list of Domain Aliases.
//
// Parameters:
//   - domain: This parameter specifies the name of an existing Domain.
//
// Returns:
//   - []string: An array with the Domain alias names.
//   - error: An error if the domain name is empty or the command fails.
func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("GETDOMAINALIASES", domain)
}

// GetDomainEffectiveSettings retrieves both explicitly set and default Domain settings.
//
// Parameters:
//   - domain: The name of an existing Domain. If empty, uses the current context.
//
// Returns:
//   - map[string]any: A dictionary with all domain settings (explicit + default).
//   - error: An error if the command fails.
func (cli *Cli) GetDomainEffectiveSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINEFFECTIVESETTINGS", Atom(domain))
}

// GetDomainLocation retrieves the Domain file directory path relative to the Server base.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//
// Returns:
//   - string: The Domain file path.
//   - error: An error if the domain name is empty or retrieval fails.
func (cli *Cli) GetDomainLocation(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain name is required")
	}
	return cli.getString("GETDOMAINLOCATION", domain)
}

// GetDomainMailRules retrieves the list of Domain Queue Mail Rules.
//
// Parameters:
//   - domain: This parameter specifies the name of an existing Domain.
//
// Returns:
//   - [][MailRule]: An array of the  Mail Rules specified for the Domain.
//   - error: An error if the command fails.
func (cli *Cli) GetDomainMailRules(domain string) ([]MailRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMailRules("GETDOMAINMAILRULES", domain)
}

// GetDomainSettings retrieves explicitly set (non-default) Domain settings.
//
// Parameters:
//   - domain: The name of an existing Domain. If empty, uses the current context.
//
// Returns:
//   - map[string]any: A dictionary containing only explicitly set settings.
//   - error: An error if the command fails.
func (cli *Cli) GetDomainSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINSETTINGS", Atom(domain))
}

// GetDomainSignalRules retrieves the list of Domain Signal Rules.
//
// Parameters:
//   - domain: This parameter specifies the name of an existing Domain.
//
// Returns:
//   - [][SignalRule]: An array of the  Signal Rules specified for the Domain.
//   - error: An error if the command fails.
func (cli *Cli) GetDomainSignalRules(domain string) ([]SignalRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSignalRules("GETDOMAINSIGNALRULES", domain)
}

// InsertDirectoryRecords inserts records for Domain objects into the Directory.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the authenticated user Domain.
//
// Returns:
//   - error: An error if the operation fails or access is denied.
func (cli *Cli) InsertDirectoryRecords(domain string) error {
	return cli.QueryNV("INSERTDIRECTORYRECORDS", Atom(domain))
}

// IsDomainExists checks if a domain exists by attempting to retrieve its location.
//
// Parameters:
//   - domain: The name of the domain to check.
//
// Returns:
//   - bool: True if the domain exists.
//   - error: An error if the check fails.
func (cli *Cli) IsDomainExists(domain string) (bool, error) {
	location, err := cli.GetDomainLocation(domain)
	if err != nil {
		return false, err
	}
	return location != "", nil
}

// ListAccountStorage lists Account "storage mount points" in the specified Domain.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//
// Returns:
//   - []string: An array with "storage mount points" names.
//   - error: An error if the retrieval fails.
func (cli *Cli) ListAccountStorage(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("LISTACCOUNTSTORAGE", domain)
}

// ListAccounts retrieves the list of all Accounts in the Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//
// Returns:
//   - map[string]string: A dictionary where keys are Account names.
//   - error: An error if the command fails.
func (cli *Cli) ListAccounts(domain string) (map[string]string, error) {
	return cli.getMapString("LISTACCOUNTS", Atom(domain))
}

// ListAdminDomains gets the list of Domains that can be administered in the target Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the authenticated user Domain.
//
// Returns:
//   - []string: An array with the Domain names.
//   - error: An error if the command fails.
func (cli *Cli) ListAdminDomains(domain string) ([]string, error) {
	return cli.getSliceString("LISTADMINDOMAINS", domain)
}

// ListDomainObjects retrieves a paginated and filtered list of Domain objects.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - limit: The maximum number of objects to list in one batch.
//   - filter: Optional substring to filter object names.
//   - what: Bitmask [ObjectType] specifying which objects (Accounts, Aliases, Forwarders) to list.
//   - cookie: Pagination token. Use an empty string for the first request.
//
// Returns:
//   - *[DomainObjectsResult]: A structure containing the objects found and the next cookie.
//   - error: An error if the retrieval or parsing fails.
func (cli *Cli) ListDomainObjects(domain string, limit int, filter string, what ObjectType, cookie string) (*DomainObjectsResult, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}

	const cmd = "LISTDOMAINOBJECTS"
	args := make([]any, 0, 8)
	args = append(args, domain)
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	args = append(args, limit)
	if what&TypeAccount != 0 {
		args = append(args, "ACCOUNTS")
	}
	if what&TypeAlias != 0 {
		args = append(args, "ALIASES")
	}
	if what&TypeForwarder != 0 {
		args = append(args, "FORWARDERS")
	}
	args = append(args, "COOKIE", cookie)

	res, err := cli.getSliceAny(cmd, args...)
	if err != nil {
		return nil, err
	}

	if len(res) < 5 {
		return nil, fmt.Errorf("%s: unexpected response format", cmd)
	}

	result := &DomainObjectsResult{Objects: make(map[string]DomainObject)}
	if c, ok := res[4].(string); ok {
		result.Cookie = c
	}

	items, ok := res[1].(map[string]any)
	if !ok {
		return result, nil
	}

	for name, val := range items {
		obj := DomainObject{Name: name}
		switch v := val.(type) {
		case string:
			obj.Type = TypeAccount
			obj.Subtype = v
		case []any:
			if len(v) > 0 {
				var targetErr error
				if first, isSlice := v[0].([]any); isSlice {
					obj.Type = TypeForwarder
					obj.Subtype = "forwarder"
					obj.Targets, targetErr = toStringSlice(cmd, first)
				} else {
					obj.Type = TypeAlias
					obj.Subtype = "alias"
					obj.Targets, targetErr = toStringSlice(cmd, v)
				}
				if targetErr != nil {
					return nil, fmt.Errorf("%s: invalid target format for %s: %w", cmd, name, targetErr)
				}
			}
		}
		result.Objects[name] = obj
	}
	return result, nil
}

// ListDomainTelnums reads Telnum numbers created in the specified Domain.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - filter: Optional substring to filter telnum numbers.
//   - limit: Maximum number of numbers to return.
//
// Returns:
//   - map[string]string: A dictionary where each key is a Telnum, and the value is the Account.
//   - error: An error if the retrieval fails.
func (cli *Cli) ListDomainTelnums(domain, filter string, limit int) (map[string]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}

	const cmd = "LISTDOMAINTELNUMS"
	args := make([]any, 0, 4)
	args = append(args, domain)
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	args = append(args, limit)

	res, err := cli.getMapAny(cmd, args...)
	if err != nil {
		return nil, err
	}
	delete(res, "")
	return toStringMap(cmd, res)
}

// ResumeDomain resumes a suspended Domain so Accounts can be opened.
//
// Parameters:
//   - domain: This parameter specifies the name of the Domain to be resumed.
//
// Returns:
//   - error: An error if the domain name is empty or the command fails.
func (cli *Cli) ResumeDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("RESUMEDOMAIN", domain)
}

// SetAccountDefaultPrefs replaces the Default Account Preferences for the Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the authenticated user Domain.
//   - settings: A dictionary used to replace the Default Account Preferences.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetAccountDefaultPrefs(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETACCOUNTDEFAULTPREFS", Atom(domain), settings)
}

// SetAccountDefaults replaces the Domain Default Account settings.
//
// Parameters:
//   - domain: The Domain name.
//   - settings: A dictionary used to replace the Domain Default Account settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetAccountDefaults(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETACCOUNTDEFAULTS", Atom(domain), settings)
}

// SetAccountTemplate replaces the Account Template settings for the Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//   - settings: A dictionary used to replace the Domain Account Template.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetAccountTemplate(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETACCOUNTTEMPLATE", Atom(domain), settings)
}

// SetDomainAliases sets the Domain aliases, replacing all existing ones.
//
// Parameters:
//   - domain: The name of an existing Domain.
//   - aliases: An array of Domain alias name strings.
//
// Returns:
//   - error: An error if parameters are invalid or the command fails.
func (cli *Cli) SetDomainAliases(domain string, aliases []string) error {
	if domain == "" || aliases == nil {
		return fmt.Errorf("domain name and aliases array are required")
	}
	return cli.QueryNV("SETDOMAINALIASES", domain, aliases)
}

// SetDomainMailRules sets the Domain Queue Mail Rules, replacing all old ones.
//
// Parameters:
//   - domain: The name of an existing Domain.
//   - rules: An array containing the new Domain Queue Rules.
//
// Returns:
//   - error: An error if parameters are invalid or access is denied.
func (cli *Cli) SetDomainMailRules(domain string, rules []MailRule) error {
	if domain == "" || rules == nil {
		return fmt.Errorf("domain name and rules array are required")
	}
	return setRulesGeneric(cli, "SETDOMAINMAILRULES", domain, rules)
}

// SetDomainSettings replaces the entire Domain settings dictionary.
//
// Parameters:
//   - domain: The name of an existing Domain.
//   - settings: A dictionary used to replace the Domain settings.
//
// Returns:
//   - error: An error if parameters are invalid or the command fails.
func (cli *Cli) SetDomainSettings(domain string, settings map[string]any) error {
	if domain == "" || settings == nil {
		return fmt.Errorf("domain name and settings dictionary are required")
	}
	return cli.QueryNV("SETDOMAINSETTINGS", domain, settings)
}

// SetDomainSignalRules sets the Domain Signal Rules, replacing all old ones.
//
// Parameters:
//   - domain: The name of an existing Domain.
//   - rules: An array containing the new Domain Signal Rules.
//
// Returns:
//   - error: An error if parameters are invalid or access is denied.
func (cli *Cli) SetDomainSignalRules(domain string, rules []SignalRule) error {
	if domain == "" || rules == nil {
		return fmt.Errorf("domain name and rules array are required")
	}
	return setRulesGeneric(cli, "SETDOMAINSIGNALRULES", domain, rules)
}

// SuspendDomain suspends a Domain, closing active accounts and preventing new ones.
//
// Parameters:
//   - domain: This parameter specifies the name of the Domain to be suspended.
//
// Returns:
//   - error: An error if the domain name is empty or the command fails.
func (cli *Cli) SuspendDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("SUSPENDDOMAIN", domain)
}

// UpdateAccountDefaultPrefs modifies the Domain Default Account Preferences.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the authenticated user Domain.
//   - settings: A dictionary used to modify the Domain Default Account Preferences.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateAccountDefaultPrefs(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEACCOUNTDEFAULTPREFS", Atom(domain), settings)
}

// UpdateAccountDefaults modifies the Default Account settings for the specified Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//   - settings: A dictionary used to modify the Domain Default Account settings.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateAccountDefaults(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEACCOUNTDEFAULTS", Atom(domain), settings)
}

// UpdateAccountTemplate modifies the Account Template settings for the Domain.
//
// Parameters:
//   - domain: The Domain name. If empty, applies to the administrator Domain.
//   - settings: A dictionary used to modify the Domain Account Template.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateAccountTemplate(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEACCOUNTTEMPLATE", Atom(domain), settings)
}

// UpdateDomainSettings updates specific Domain settings without modifying others.
//
// Parameters:
//   - domain: The name of an existing Domain. If empty, uses the current context.
//   - settings: A dictionary containing the settings to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateDomainSettings(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEDOMAINSETTINGS", Atom(domain), settings)
}
