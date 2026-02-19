// # Domain Administration
//
// The Domains section provides extensive functionality for managing CommuniGate Pro
// domains and their internal objects. It handles the domain lifecycle, inheritance
// of settings, and mass-object discovery.
//
// Key capabilities include:
//   - Configuration Management: handling [Domain Settings], [Effective Settings],
//     and Domain-level [Aliases].
//   - Inheritance Control: managing [Account Defaults] and [Account Templates],
//     which dictate the initial state of new accounts within the domain.
//   - Object Discovery: using [Cli.ListDomainObjects] to efficiently search and
//     categorize accounts, aliases, and forwarders using filters and cookies for pagination.
//   - Operational State: suspending, resuming, and locating domain storage paths.
//   - Rules & Directory: managing Domain-level Mail/Signal rules and synchronizing
//     data with the Directory Service.
//
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
	Name    string     // unique identifier of the object (e.g., account name).
	Type    ObjectType // category: TypeAccount, TypeAlias, or TypeForwarder.
	Subtype string     // object-specific data.
	Targets []string   // destination addresses for aliases and forwarders.
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
//   - domain: the Domain name.
//   - storage: the "storage mount Point" name.
//
// This method executes the CREATEACCOUNTSTORAGE CLI command.
//
// Returns:
//   - error: an error if parameters are missing or the server fails to create the storage.
func (cli *Cli) CreateAccountStorage(domain, storage string) error {
	if domain == "" || storage == "" {
		return fmt.Errorf("domain name and storage path are required")
	}
	return cli.QueryNV("CREATEACCOUNTSTORAGE", domain, "PATH", storage)
}

// DeleteDirectoryRecords deletes Domain object records from the Directory.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the authenticated user Domain.
//
// This method executes the DELETEDIRECTORYRECORDS CLI command.
//
// Returns:
//   - error: an error if the operation fails or access is denied.
func (cli *Cli) DeleteDirectoryRecords(domain string) error {
	const cmd = "DELETEDIRECTORYRECORDS"
	if domain != "" {
		return cli.QueryNV(cmd, domain)
	} else {
		return cli.QueryNV(cmd)
	}
}

// GetAccountDefaultPrefs retrieves the Default Account Preferences for the specified Domain.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the administrator Domain.
//
// This method executes the GETACCOUNTDEFAULTPREFS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the default Preferences.
//   - error: an error if the retrieval fails.
func (cli *Cli) GetAccountDefaultPrefs(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTDEFAULTPREFS", Atom(domain))
}

// GetAccountDefaults retrieves the default Account settings for the specified Domain.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the administrator Domain.
//
// This method executes the GETACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the default settings.
//   - error: an error if the retrieval fails.
func (cli *Cli) GetAccountDefaults(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTDEFAULTS", Atom(domain))
}

// GetAccountTemplate retrieves the Account Template settings for the Domain.
//
// Parameters:
//   - domain: an optional Domain name. If empty, applies to the administrator Domain.
//
// This method executes the GETACCOUNTTEMPLATE CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Template settings.
//   - error: an error if the retrieval fails.
func (cli *Cli) GetAccountTemplate(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTTEMPLATE", Atom(domain))
}

// GetDomainAliases retrieves the list of Domain Aliases.
//
// Parameters:
//   - domain: the name of an existing Domain.
//
// This method executes the GETDOMAINALIASES CLI command.
//
// Returns:
//   - []string: an array with the Domain alias names.
//   - error: an error if the domain name is empty or the command fails.
func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("GETDOMAINALIASES", domain)
}

// GetDomainEffectiveSettings retrieves both explicitly set and default Domain settings.
//
// Parameters:
//   - domain: the name of an existing Domain. If empty, uses the current context.
//
// This method executes the GETDOMAINEFFECTIVESETTINGS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with all domain settings (explicit + default).
//   - error: an error if the command fails.
func (cli *Cli) GetDomainEffectiveSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINEFFECTIVESETTINGS", Atom(domain))
}

// GetDomainLocation retrieves the Domain file directory path relative to the Server base.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the administrator Domain.
//
// This method executes the GETDOMAINLOCATION CLI command.
//
// Returns:
//   - string: the Domain file path.
//   - error: an error if the domain name is empty or retrieval fails.
func (cli *Cli) GetDomainLocation(domain string) (string, error) {
	return cli.getString("GETDOMAINLOCATION", Atom(domain))
}

// GetDomainMailRules retrieves the list of Domain Queue Mail Rules.
//
// Parameters:
//   - domain: the name of an existing Domain.
//
// This method executes the GETDOMAINMAILRULES CLI command.
//
// Returns:
//   - [][MailRule]: an array of the  Mail Rules specified for the Domain.
//   - error: an error if the command fails.
func (cli *Cli) GetDomainMailRules(domain string) ([]MailRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMailRules("GETDOMAINMAILRULES", domain)
}

// GetDomainSettings retrieves explicitly set (non-default) Domain settings.
//
// Parameters:
//   - domain: the name of an existing Domain. If empty, uses the current context.
//
// This method executes the GETDOMAINSETTINGS CLI command.
//
// Returns:
//   - map[string]any: a dictionary containing only explicitly set settings.
//   - error: an error if the command fails.
func (cli *Cli) GetDomainSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINSETTINGS", Atom(domain))
}

// GetDomainSignalRules retrieves the list of Domain Signal Rules.
//
// Parameters:
//   - domain: the name of an existing Domain.
//
// This method executes the GETDOMAINSIGNALRULES CLI command.
//
// Returns:
//   - [][SignalRule]: an array of the  Signal Rules specified for the Domain.
//   - error: an error if the command fails.
func (cli *Cli) GetDomainSignalRules(domain string) ([]SignalRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSignalRules("GETDOMAINSIGNALRULES", domain)
}

// InsertDirectoryRecords inserts records for Domain objects into the Directory.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the authenticated user Domain.
//
// This method executes the INSERTDIRECTORYRECORDS CLI command.
//
// Returns:
//   - error: an error if the operation fails or access is denied.
func (cli *Cli) InsertDirectoryRecords(domain string) error {
	return cli.QueryNV("INSERTDIRECTORYRECORDS", Atom(domain))
}

// IsDomainExists checks if a domain exists by attempting to retrieve its location.
//
// Parameters:
//   - domain: the name of the domain to check.
//
// Returns:
//   - bool: true if the domain exists.
//   - error: an error if the check fails.
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
//   - domain: the Domain name.
//
// This method executes the LISTACCOUNTSTORAGE CLI command.
//
// Returns:
//   - []string: an array with "storage mount points" names.
//   - error: an error if the retrieval fails.
func (cli *Cli) ListAccountStorage(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("LISTACCOUNTSTORAGE", domain)
}

// ListAccounts retrieves the list of all Accounts in the Domain.
//
// Parameters:
//   - domain: an optional Domain name. If empty, applies to the administrator Domain.
//
// This method executes the LISTACCOUNTS CLI command.
//
// Returns:
//   - map[string]string: a dictionary where keys are Account names.
//   - error: an error if the command fails.
func (cli *Cli) ListAccounts(domain string) (map[string]string, error) {
	return cli.getMapString("LISTACCOUNTS", Atom(domain))
}

// ListAdminDomains gets the list of Domains that can be administered in the target Domain.
//
// Parameters:
//   - domain: an optional Domain name. If empty, applies to the authenticated user Domain.
//
// This method executes the LISTADMINDOMAINS CLI command.
//
// Returns:
//   - []string: an array with the Domain names.
//   - error: an error if the command fails.
func (cli *Cli) ListAdminDomains(domain string) ([]string, error) {
	return cli.getSliceString("LISTADMINDOMAINS", Atom(domain))
}

// ListDomainObjects retrieves a paginated and filtered list of Domain objects.
//
// Parameters:
//   - domain: the Domain name.
//   - limit: the maximum number of objects to list in one batch.
//   - filter: an optional substring to filter object names.
//   - what: bitmask [ObjectType] specifying which objects (Accounts, Aliases, Forwarders) to list.
//   - cookie: pagination token. Use an empty string for the first request.
//
// This method executes the LISTDOMAINOBJECTS CLI command.
//
// Returns:
//   - *[DomainObjectsResult]: a structure containing the objects found and the next cookie.
//   - error: an error if the retrieval or parsing fails.
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
//   - domain: the Domain name.
//   - filter: an optional substring to filter telnum numbers.
//   - limit: maximum number of numbers to return.
//
// This method executes the LISTDOMAINTELNUMS CLI command.
//
// Returns:
//   - map[string]string: a dictionary where each key is a Telnum, and the value is the Account.
//   - error: an error if the retrieval fails.
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
//   - domain: the name of the Domain to be resumed.
//
// This method executes the RESUMEDOMAIN CLI command.
//
// Returns:
//   - error: an error if the domain name is empty or the command fails.
func (cli *Cli) ResumeDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("RESUMEDOMAIN", domain)
}

// SetAccountDefaultPrefs replaces the Default Account Preferences for the Domain.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the authenticated user Domain.
//   - settings: a dictionary used to replace the Default Account Preferences.
//
// This method executes the SETACCOUNTDEFAULTPREFS CLI command.
//
// Returns:
//   - error: an error if settings are nil or the command fails.
func (cli *Cli) SetAccountDefaultPrefs(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETACCOUNTDEFAULTPREFS", Atom(domain), settings)
}

// SetAccountDefaults replaces the Domain Default Account settings.
//
// Parameters:
//   - domain: an optional Domain name.
//   - settings: a dictionary used to replace the Domain Default Account settings.
//
// This method executes the SETACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - error: an error if settings are nil or the command fails.
func (cli *Cli) SetAccountDefaults(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETACCOUNTDEFAULTS", Atom(domain), settings)
}

// SetAccountTemplate replaces the Account Template settings for the Domain.
//
// Parameters:
//   - domain: an optional Domain name. If empty, applies to the administrator Domain.
//   - settings: a dictionary used to replace the Domain Account Template.
//
// This method executes the SETACCOUNTTEMPLATE CLI command.
//
// Returns:
//   - error: an error if settings are nil or the command fails.
func (cli *Cli) SetAccountTemplate(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETACCOUNTTEMPLATE", Atom(domain), settings)
}

// SetDomainAliases sets the Domain aliases, replacing all existing ones.
//
// Parameters:
//   - domain: the name of an existing Domain.
//   - aliases: an array of Domain alias name strings.
//
// This method executes the SETDOMAINALIASES CLI command.
//
// Returns:
//   - error: an error if parameters are invalid or the command fails.
func (cli *Cli) SetDomainAliases(domain string, aliases []string) error {
	if domain == "" || aliases == nil {
		return fmt.Errorf("domain name and aliases array are required")
	}
	return cli.QueryNV("SETDOMAINALIASES", domain, aliases)
}

// SetDomainMailRules sets the Domain Queue Mail Rules, replacing all old ones.
//
// Parameters:
//   - domain: the name of an existing Domain.
//   - rules: an array containing the new Domain Queue Rules.
//
// This method executes the SETDOMAINMAILRULES CLI command.
//
// Returns:
//   - error: an error if parameters are invalid or access is denied.
func (cli *Cli) SetDomainMailRules(domain string, rules []MailRule) error {
	if domain == "" || rules == nil {
		return fmt.Errorf("domain name and rules array are required")
	}
	return setRulesGeneric(cli, "SETDOMAINMAILRULES", domain, rules)
}

// SetDomainSettings replaces the entire Domain settings dictionary.
//
// Parameters:
//   - domain: the name of an existing Domain.
//   - settings: a dictionary used to replace the Domain settings.
//
// This method executes the SETDOMAINSETTINGS CLI command.
//
// Returns:
//   - error: an error if parameters are invalid or the command fails.
func (cli *Cli) SetDomainSettings(domain string, settings map[string]any) error {
	if domain == "" || settings == nil {
		return fmt.Errorf("domain name and settings dictionary are required")
	}
	return cli.QueryNV("SETDOMAINSETTINGS", domain, settings)
}

// SetDomainSignalRules sets the Domain Signal Rules, replacing all old ones.
//
// Parameters:
//   - domain: the name of an existing Domain.
//   - rules: an array containing the new Domain Signal Rules.
//
// This method executes the SETDOMAINSIGNALRULES CLI command.
//
// Returns:
//   - error: an error if parameters are invalid or access is denied.
func (cli *Cli) SetDomainSignalRules(domain string, rules []SignalRule) error {
	if domain == "" || rules == nil {
		return fmt.Errorf("domain name and rules array are required")
	}
	return setRulesGeneric(cli, "SETDOMAINSIGNALRULES", domain, rules)
}

// SuspendDomain suspends a Domain, closing active accounts and preventing new ones.
//
// Parameters:
//   - domain: the name of the Domain to be suspended.
//
// This method executes the SUSPENDDOMAIN CLI command.
//
// Returns:
//   - error: an error if the domain name is empty or the command fails.
func (cli *Cli) SuspendDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("SUSPENDDOMAIN", domain)
}

// UpdateAccountDefaultPrefs modifies the Domain Default Account Preferences.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the authenticated user Domain.
//   - settings: a dictionary used to modify the Domain Default Account Preferences.
//
// This method executes the UPDATEACCOUNTDEFAULTPREFS CLI command.
//
// Returns:
//   - error: an error if settings are empty or the command fails.
func (cli *Cli) UpdateAccountDefaultPrefs(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEACCOUNTDEFAULTPREFS", Atom(domain), settings)
}

// UpdateAccountDefaults modifies the Default Account settings for the specified Domain.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the administrator Domain.
//   - settings: a dictionary used to modify the Domain Default Account settings.
//
// This method executes the UPDATEACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - error: an error if settings are empty or the command fails.
func (cli *Cli) UpdateAccountDefaults(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEACCOUNTDEFAULTS", Atom(domain), settings)
}

// UpdateAccountTemplate modifies the Account Template settings for the Domain.
//
// Parameters:
//   - domain: the Domain name. If empty, applies to the administrator Domain.
//   - settings: a dictionary used to modify the Domain Account Template.
//
// This method executes the UPDATEACCOUNTTEMPLATE CLI command.
//
// Returns:
//   - error: an error if settings are empty or the command fails.
func (cli *Cli) UpdateAccountTemplate(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEACCOUNTTEMPLATE", Atom(domain), settings)
}

// UpdateDomainSettings updates specific Domain settings without modifying others.
//
// Parameters:
//   - domain: the name of an existing Domain. If empty, uses the current context.
//   - settings: a dictionary containing the settings to update.
//
// This method executes the UPDATEDOMAINSETTINGS CLI command.
//
// Returns:
//   - error: an error if settings are empty or the command fails.
func (cli *Cli) UpdateDomainSettings(domain string, settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEDOMAINSETTINGS", Atom(domain), settings)
}
