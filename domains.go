// # Section: Domain Administration
package cgpcli

import (
	"fmt"
)

type ObjectType int

const (
	TypeAccount ObjectType = 1 << iota
	TypeAlias
	TypeForwarder
	TypeAll = TypeAccount | TypeAlias | TypeForwarder
)

type DomainObject struct {
	Name    string
	Type    ObjectType
	Subtype string
	Targets []string
}

type DomainObjectsResult struct {
	Objects map[string]DomainObject
	Cookie  string
}

func (cli *Cli) CreateAccountStorage(domain, storage string) error {
	if domain == "" || storage == "" {
		return fmt.Errorf("domain and storage path are required")
	}
	return cli.QueryNV("CREATEACCOUNTSTORAGE", domain, "PATH", storage)
}

func (cli *Cli) DeleteDirectoryRecords(domain string) (err error) {
	const cmd = "DELETEDIRECTORYRECORDS"
	if domain != "" {
		err = cli.QueryNV(cmd, domain)
	} else {
		err = cli.QueryNV(cmd)
	}
	return
}

func (cli *Cli) GetAccountDefaultPrefs(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTDEFAULTPREFS", Atom(domain))
}

func (cli *Cli) GetAccountDefaults(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTDEFAULTS", Atom(domain))
}

func (cli *Cli) GetAccountTemplate(domain string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTTEMPLATE", Atom(domain))
}

func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("GETDOMAINALIASES", domain)
}

func (cli *Cli) GetDomainEffectiveSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINEFFECTIVESETTINGS", Atom(domain))
}

func (cli *Cli) GetDomainLocation(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain name is required")
	}
	return cli.getString("GETDOMAINLOCATION", domain)
}

func (cli *Cli) GetDomainMailRules(domain string) ([]MailRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMailRules("GETDOMAINMAILRULES", domain)
}

func (cli *Cli) GetDomainSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINSETTINGS", Atom(domain))
}

func (cli *Cli) GetDomainSignalRules(domain string) ([]SignalRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSignalRules("GETDOMAINSIGNALRULES", domain)
}

func (cli *Cli) InsertDirectoryRecords(domain string) (err error) {
	const cmd = "INSERTDIRECTORYRECORDS"
	if domain != "" {
		err = cli.QueryNV(cmd, domain)
	} else {
		err = cli.QueryNV(cmd)
	}
	return
}

func (cli *Cli) IsDomainExists(domain string) (bool, error) {
	location, err := cli.GetDomainLocation(domain)
	if err != nil {
		return false, err
	}
	return location != "", nil
}

func (cli *Cli) ListAccountStorage(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("LISTACCOUNTSTORAGE", domain)
}

func (cli *Cli) ListAccounts(domain string) (map[string]string, error) {
	return cli.getMapString("LISTACCOUNTS", Atom(domain))
}

func (cli *Cli) ListAdminDomains(domain string) ([]string, error) {
	return cli.getSliceString("LISTADMINDOMAINS", domain)
}

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

func (cli *Cli) ResumeDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("RESUMEDOMAIN", domain)
}

func (cli *Cli) SetAccountDefaultPrefs(domain string, s map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("SETACCOUNTDEFAULTPREFS", domain, s)
}

func (cli *Cli) SetAccountDefaults(domain string, s map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("SETACCOUNTDEFAULTS", domain, s)
}

func (cli *Cli) SetAccountTemplate(domain string, s map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("SETACCOUNTTEMPLATE", domain, s)
}

func (cli *Cli) SetDomainAliases(domain string, aliases []string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("SETDOMAINALIASES", domain, aliases)
}

func (cli *Cli) SetDomainMailRules(domain string, rules []MailRule) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return setRulesGeneric(cli, "SETDOMAINMAILRULES", domain, rules)
}

func (cli *Cli) SetDomainSettings(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("SETDOMAINSETTINGS", domain, settings)
}

func (cli *Cli) SetDomainSignalRules(domain string, rules []SignalRule) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return setRulesGeneric(cli, "SETDOMAINSIGNALRULES", domain, rules)
}

func (cli *Cli) SuspendDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("SUSPENDDOMAIN", domain)
}

func (cli *Cli) UpdateAccountDefaultPrefs(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTDEFAULTPREFS", domain, settings)
}

func (cli *Cli) UpdateAccountDefaults(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTDEFAULTS", domain, settings)
}

func (cli *Cli) UpdateAccountTemplate(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTTEMPLATE", domain, settings)
}

func (cli *Cli) UpdateDomainSettings(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEDOMAINSETTINGS", domain, settings)
}
