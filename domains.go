package cgpcli

import (
	"fmt"
	"strconv"
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Methods
// -----------------------------------------------------------------------------

func (cli *Cli) CreateAccountStorage(domain, storage string) error {
	if domain == "" || storage == "" {
		return fmt.Errorf("domain and storage path are required")
	}
	_, err := cli.Query("CREATEACCOUNTSTORAGE", domain, "PATH", storage)
	return err
}

func (cli *Cli) DeleteDirectoryRecords(domain string) error {
	args := make([]any, 0, 2)
	args = append(args, "DELETEDIRECTORYRECORDS")
	if domain != "" {
		args = append(args, domain)
	}
	_, err := cli.Query(args...)
	return err
}

func (cli *Cli) GetAccountDefaults(domain string) (map[string]any, error) {
	return cli.getSettingsMapAlt("GETACCOUNTDEFAULTS", domain)
}

func (cli *Cli) GetAccountDefaultPrefs(domain string) (map[string]any, error) {
	return cli.getSettingsMapAlt("GETACCOUNTDEFAULTPREFS", domain)
}

func (cli *Cli) GetAccountTemplate(domain string) (map[string]any, error) {
	return cli.getSettingsMapAlt("GETACCOUNTTEMPLATE", domain)
}

func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	res, err := cli.Query("GETDOMAINALIASES", domain)
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}

	aliases := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := cli.asString(v)
		if !ok {
			return nil, fmt.Errorf("GETDOMAINALIASES: unexpected element type %T", v)
		}
		aliases = append(aliases, s)
	}
	return aliases, nil
}

func (cli *Cli) GetDomainEffectiveSettings(domain string) (map[string]any, error) {
	return cli.getSettingsMapAlt("GETDOMAINEFFECTIVESETTINGS", domain)
}

func (cli *Cli) GetDomainLocation(domain string) (string, error) {
	res, err := cli.getDomainData("GETDOMAINLOCATION", domain)
	if err != nil {
		return "", err
	}
	location, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETDOMAINLOCATION unexpected response type: %T", res)
	}
	return location, nil
}

func (cli *Cli) GetDomainMailRules(domain string) ([]MailRule, error) {
	return cli.getMailRules("GETDOMAINMAILRULES", domain)
}

func (cli *Cli) GetDomainSettings(domain string) (map[string]any, error) {
	return cli.getSettingsMapAlt("GETDOMAINSETTINGS", domain)
}

func (cli *Cli) GetDomainSignalRules(domain string) ([]SignalRule, error) {
	return cli.getSignalRules("GETDOMAINSIGNALRULES", domain)
}

func (cli *Cli) InsertDirectoryRecords(domain string) error {
	args := make([]any, 0, 2)
	args = append(args, "INSERTDIRECTORYRECORDS")
	if domain != "" {
		args = append(args, domain)
	}
	_, err := cli.Query(args...)
	return err
}

func (cli *Cli) ListAccounts(domain string) (map[string]string, error) {
	res, err := cli.getDomainData("LISTACCOUNTS", domain)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return make(map[string]string), nil
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("LISTACCOUNTS unexpected response type: %T", res)
	}
	accounts := make(map[string]string, len(dict))
	for name, val := range dict {
		s, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("LISTACCOUNTS account %q has unexpected value type %T: %v", name, val, val)
		}
		accounts[name] = s
	}
	return accounts, nil
}

func (cli *Cli) ListAccountStorage(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	res, err := cli.Query("LISTACCOUNTSTORAGE", domain)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return []string{}, nil
	}
	raw, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("LISTACCOUNTSTORAGE unexpected response type: %T", res)
	}
	storages := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("LISTACCOUNTSTORAGE contains non-string element: %T", v)
		}
		storages = append(storages, s)
	}
	return storages, nil
}

func (cli *Cli) ListAdminDomains(domain string) ([]string, error) {
	res, err := cli.getDomainData("LISTADMINDOMAINS", domain)
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}

	domains := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := cli.asString(v)
		if !ok {
			return nil, fmt.Errorf("LISTADMINDOMAINS: unexpected element type %T", v)
		}
		domains = append(domains, s)
	}
	return domains, nil
}

func (cli *Cli) ListDomainObjects(domain string, limit int, filter string, what ObjectType, cookie string) (*DomainObjectsResult, error) {
	args := make([]any, 0, 8)
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

	// Всегда передаем COOKIE для фиксации структуры ответа (5 элементов)
	args = append(args, "COOKIE", cookie)

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok || len(raw) < 4 {
		return nil, fmt.Errorf("LISTDOMAINOBJECTS: unexpected response format (len=%d)", len(raw))
	}

	result := &DomainObjectsResult{
		Objects: make(map[string]DomainObject),
	}

	// Извлекаем Cookie (5-й элемент), только если он физически присутствует
	if len(raw) >= 5 {
		if c, ok := raw[4].(string); ok {
			result.Cookie = c
		}
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
				// Forwarder: [[address]]
				if first, isSlice := v[0].([]any); isSlice {
					obj.Type = TypeForwarder
					obj.Subtype = "forwarder"
					obj.Targets = cli.toStrictStringSlice(first)
				} else {
					// Alias: [owner]
					obj.Type = TypeAlias
					obj.Subtype = "alias"
					obj.Targets = cli.toStrictStringSlice(v)
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

	items := make(map[string]string, len(dict))
	for k, v := range dict {
		if k == "" {
			continue // Пропускаем счетчик CGP
		}
		if s, ok := cli.asString(v); ok {
			items[k] = s
		}
	}
	return items, nil
}

func (cli *Cli) ResumeDomain(domain string) error {
	_, err := cli.Query("RESUMEDOMAIN", domain)
	return err
}

func (cli *Cli) SetAccountDefaults(domain string, s map[string]any) error {
	return cli.domSet("SETACCOUNTDEFAULTS", domain, s)
}

func (cli *Cli) SetAccountDefaultPrefs(domain string, s map[string]any) error {
	return cli.domSet("SETACCOUNTDEFAULTPREFS", domain, s)
}

func (cli *Cli) SetAccountTemplate(domain string, s map[string]any) error {
	return cli.domSet("SETACCOUNTTEMPLATE", domain, s)
}

func (cli *Cli) SetDomainAliases(domain string, aliases []string) error {
	_, err := cli.Query("SETDOMAINALIASES", domain, aliases)
	return err
}

func (cli *Cli) SetDomainMailRules(domain string, rules []MailRule) error {
	return setRulesGeneric(cli, "SETDOMAINMAILRULES", domain, rules)
}

func (cli *Cli) SetDomainSettings(domain string, settings map[string]any) error {
	_, err := cli.Query("SETDOMAINSETTINGS", domain, settings)
	return err
}

func (cli *Cli) SetDomainSignalRules(domain string, rules []SignalRule) error {
	return setRulesGeneric(cli, "SETDOMAINSIGNALRULES", domain, rules)
}

func (cli *Cli) SuspendDomain(domain string) error {
	_, err := cli.Query("SUSPENDDOMAIN", domain)
	return err
}

func (cli *Cli) UpdateAccountDefaults(domain string, s map[string]any) error {
	return cli.domSet("UPDATEACCOUNTDEFAULTS", domain, s)
}

func (cli *Cli) UpdateAccountDefaultPrefs(domain string, s map[string]any) error {
	return cli.domSet("UPDATEACCOUNTDEFAULTPREFS", domain, s)
}

func (cli *Cli) UpdateAccountTemplate(domain string, s map[string]any) error {
	return cli.domSet("UPDATEACCOUNTTEMPLATE", domain, s)
}

func (cli *Cli) UpdateDomainSettings(domain string, settings map[string]any) error {
	if settings == nil {
		return nil
	}
	args := make([]any, 0, 3)
	args = append(args, "UPDATEDOMAINSETTINGS")
	if domain != "" {
		args = append(args, domain)
	}
	args = append(args, settings)
	_, err := cli.Query(args...)
	return err
}

// -----------------------------------------------------------------------------
// Helpers / Internal
// -----------------------------------------------------------------------------

func (cli *Cli) asString(v any) (string, bool) {
	switch val := v.(type) {
	case string:
		return val, true
	case int64:
		return strconv.FormatInt(val, 10), true
	case int:
		return strconv.Itoa(val), true
	default:
		return "", false
	}
}

func (cli *Cli) toStrictStringSlice(raw []any) []string {
	res := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := cli.asString(v); ok {
			res = append(res, s)
		}
	}
	return res
}

func (cli *Cli) getDomainData(cmd, domain string) (any, error) {
	args := make([]any, 0, 2)
	args = append(args, cmd)
	if domain != "" {
		args = append(args, Atom(domain))
	}
	return cli.Query(args...)
}

func (cli *Cli) getSettingsMapAlt(cmd, domain string) (map[string]any, error) {
	res, err := cli.getDomainData(cmd, domain)
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return dict, nil
}

func (cli *Cli) domSet(cmd, domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	args := make([]any, 0, 3)
	args = append(args, cmd)
	if domain != "" {
		args = append(args, domain)
	}
	args = append(args, settings)
	_, err := cli.Query(args...)
	return err
}

func (cli *Cli) IsDomainExists(domain string) (bool, error) {
	location, err := cli.GetDomainLocation(domain)
	if err != nil {
		return false, err
	}
	return location != "", nil
}

func (cli *Cli) toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}
