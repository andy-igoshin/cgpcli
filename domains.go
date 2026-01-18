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
	_, err := cli.Query("CREATEACCOUNTSTORAGE", domain, "PATH", storage)
	return err
}

func (cli *Cli) DeleteDirectoryRecords(domain string) error {
	if domain == "" {
		_, err := cli.Query("DELETEDIRECTORYRECORDS")
		return err
	}
	_, err := cli.Query("DELETEDIRECTORYRECORDS", domain)
	return err
}

func (cli *Cli) GetAccountDefaultPrefs(domain string) (map[string]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMapAny("GETACCOUNTDEFAULTPREFS", domain)
}

func (cli *Cli) GetAccountDefaults(domain string) (map[string]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMapAny("GETACCOUNTDEFAULTS", domain)
}

func (cli *Cli) GetAccountTemplate(domain string) (map[string]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMapAny("GETACCOUNTTEMPLATE", domain)
}

func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("GETDOMAINALIASES", domain)
}

func (cli *Cli) GetDomainEffectiveSettings(domain string) (map[string]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMapAny("GETDOMAINEFFECTIVESETTINGS", domain)
}

func (cli *Cli) GetDomainLocation(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain name is required")
	}
	res, err := cli.Query("GETDOMAINLOCATION", domain)
	if err != nil {
		return "", err
	}
	location, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETDOMAINLOCATION: expected string, got %T", res)
	}
	return location, nil
}

func (cli *Cli) GetDomainMailRules(domain string) ([]MailRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMailRules("GETDOMAINMAILRULES", domain)
}

func (cli *Cli) GetDomainSettings(domain string) (map[string]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getMapAny("GETDOMAINSETTINGS", domain)
}

func (cli *Cli) GetDomainSignalRules(domain string) ([]SignalRule, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSignalRules("GETDOMAINSIGNALRULES", domain)
}

func (cli *Cli) InsertDirectoryRecords(domain string) error {
	if domain == "" {
		_, err := cli.Query("INSERTDIRECTORYRECORDS")
		return err
	}
	_, err := cli.Query("INSERTDIRECTORYRECORDS", domain)
	return err
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
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	raw, err := cli.getMapAny("LISTACCOUNTS", domain)
	if err != nil {
		return nil, err
	}
	return cli.toStringMap("LISTACCOUNTS", raw)
}

func (cli *Cli) ListAdminDomains(domain string) ([]string, error) {
	return cli.getSliceString("LISTADMINDOMAINS", domain)
}

func (cli *Cli) ListDomainObjects(domain string, limit int, filter string, what ObjectType, cookie string) (*DomainObjectsResult, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	args := make([]any, 0, 10)
	args = append(args, "LISTDOMAINOBJECTS", domain)
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

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok || len(raw) < 5 {
		return nil, fmt.Errorf("LISTDOMAINOBJECTS: unexpected response format")
	}

	result := &DomainObjectsResult{Objects: make(map[string]DomainObject)}
	if c, ok := raw[4].(string); ok {
		result.Cookie = c
	}

	items, ok := raw[1].(map[string]any)
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
				// Используем toStringSlice2 для корректной обработки заголовков CGP
				if first, isSlice := v[0].([]any); isSlice {
					obj.Type = TypeForwarder
					obj.Subtype = "forwarder"
					obj.Targets, targetErr = toStringSlice("LISTDOMAINOBJECTS", first)
				} else {
					obj.Type = TypeAlias
					obj.Subtype = "alias"
					obj.Targets, targetErr = toStringSlice("LISTDOMAINOBJECTS", v)
				}
				if targetErr != nil {
					return nil, fmt.Errorf("LISTDOMAINOBJECTS: invalid target format for %s: %w", name, targetErr)
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
	args := make([]any, 0, 5)
	args = append(args, "LISTDOMAINTELNUMS", domain)
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	args = append(args, limit)

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("LISTDOMAINTELNUMS: expected map, got %T", res)
	}
	delete(dict, "")
	return cli.toStringMap("LISTDOMAINTELNUMS", dict)
}

func (cli *Cli) ResumeDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	_, err := cli.Query("RESUMEDOMAIN", domain)
	return err
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
	_, err := cli.Query("SETDOMAINALIASES", domain, aliases)
	return err
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
	_, err := cli.Query("SUSPENDDOMAIN", domain)
	return err
}

func (cli *Cli) UpdateAccountDefaultPrefs(domain string, s map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTDEFAULTPREFS", domain, s)
}

func (cli *Cli) UpdateAccountDefaults(domain string, s map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTDEFAULTS", domain, s)
}

func (cli *Cli) UpdateAccountTemplate(domain string, s map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTTEMPLATE", domain, s)
}

func (cli *Cli) UpdateDomainSettings(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.setMapAny("UPDATEDOMAINSETTINGS", domain, settings)
}
