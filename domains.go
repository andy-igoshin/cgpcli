package cgpcli

import (
	"fmt"
	"strconv"
	"strings"
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
	Total   int
	Cookie  string
}

type DomainSetResult struct {
	Items map[string]string
	Total int
}

// -----------------------------------------------------------------------------
// Methods (Alphabetical Order)
// -----------------------------------------------------------------------------

func (cli *Cli) CreateAccountStorage(domain, storage string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	if storage == "" {
		return fmt.Errorf("storage path is required")
	}
	parts := []string{"CREATEACCOUNTSTORAGE", Marshal(domain), "PATH", Marshal(storage)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) DeleteDirectoryRecords(domain string) error {
	parts := []string{"DELETEDIRECTORYRECORDS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) GetAccountDefaults(domain string) (map[string]any, error) {
	parts := []string{"GETACCOUNTDEFAULTS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	defaults, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTDEFAULTS: expected map, got %T", res)
	}
	return defaults, nil
}

func (cli *Cli) GetAccountDefaultPrefs(domain string) (map[string]any, error) {
	parts := []string{"GETACCOUNTDEFAULTPREFS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	prefs, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTDEFAULTPREFS: expected map, got %T", res)
	}
	return prefs, nil
}

func (cli *Cli) GetAccountTemplate(domain string) (map[string]any, error) {
	parts := []string{"GETACCOUNTTEMPLATE"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	template, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTTEMPLATE: expected map, got %T", res)
	}
	return template, nil
}

func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	res, err := cli.query("GETDOMAINALIASES " + Marshal(domain))
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	aliases := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprint(v)
		}
		aliases = append(aliases, s)
	}
	return aliases, nil
}

func (cli *Cli) GetDomainEffectiveSettings(domain string) (map[string]any, error) {
	parts := []string{"GETDOMAINEFFECTIVESETTINGS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	settings, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETDOMAINEFFECTIVESETTINGS: expected map, got %T", res)
	}
	return settings, nil
}

func (cli *Cli) GetDomainLocation(domain string) (string, error) {
	parts := []string{"GETDOMAINLOCATION"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return "", err
	}
	location, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETDOMAINLOCATION unexpected response type: %T", res)
	}
	return location, nil
}

func (cli *Cli) GetDomainMailRules(domain string) ([]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	res, err := cli.query("GETDOMAINMAILRULES " + Marshal(domain))
	if err != nil {
		return nil, err
	}
	rules, ok := res.([]any)
	if !ok {
		return make([]any, 0), nil
	}
	return rules, nil
}

func (cli *Cli) GetDomainSettings(domain string) (map[string]any, error) {
	parts := []string{"GETDOMAINSETTINGS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	settings, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETDOMAINSETTINGS: expected map, got %T", res)
	}
	return settings, nil
}

func (cli *Cli) GetDomainSignalRules(domain string) ([]any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	res, err := cli.query("GETDOMAINSIGNALRULES " + Marshal(domain))
	if err != nil {
		return nil, err
	}
	rules, ok := res.([]any)
	if !ok {
		return make([]any, 0), nil
	}
	return rules, nil
}

func (cli *Cli) InsertDirectoryRecords(domain string) error {
	parts := []string{"INSERTDIRECTORYRECORDS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) ListAccounts(domain string) (map[string]string, error) {
	parts := []string{"LISTACCOUNTS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		if res == nil {
			return make(map[string]string), nil
		}
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

	res, err := cli.query("LISTACCOUNTSTORAGE " + Marshal(domain))
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok {
		if res == nil {
			return []string{}, nil
		}
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
	parts := []string{"LISTADMINDOMAINS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}

	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}

	domains := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprint(v)
		}
		domains = append(domains, s)
	}

	return domains, nil
}

func (cli *Cli) ListDomainObjects(domainName string, limit int, filter string, what ObjectType, cookie string) (*DomainObjectsResult, error) {
	parts := []string{"LISTDOMAINOBJECTS", Marshal(domainName)}
	if filter != "" {
		parts = append(parts, "FILTER", Marshal(filter))
	}
	parts = append(parts, strconv.Itoa(limit))
	if what != 0 {
		parts = append(parts, what.ToCGPString())
	}
	parts = append(parts, "COOKIE", Marshal(cookie))

	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok || len(raw) < 5 {
		return nil, fmt.Errorf("LISTDOMAINOBJECTS: unexpected response format")
	}

	result := &DomainObjectsResult{
		Objects: make(map[string]DomainObject),
		Cookie:  fmt.Sprint(raw[4]),
	}

	switch v := raw[0].(type) {
	case int64:
		result.Total = int(v)
	case string:
		result.Total, _ = strconv.Atoi(v)
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
				if first, isSlice := v[0].([]any); isSlice {
					obj.Type = TypeForwarder
					obj.Subtype = "forwarder"
					for _, addr := range first {
						obj.Targets = append(obj.Targets, fmt.Sprint(addr))
					}
				} else {
					obj.Type = TypeAlias
					obj.Subtype = "alias"
					for _, addr := range v {
						obj.Targets = append(obj.Targets, fmt.Sprint(addr))
					}
				}
			}
		}
		result.Objects[name] = obj
	}
	return result, nil
}

func (cli *Cli) ListDomainTelnums(domain, filter string, limit int) (*DomainSetResult, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	parts := []string{"LISTDOMAINTELNUMS", Marshal(domain)}
	if filter != "" {
		parts = append(parts, "FILTER", Marshal(filter))
	}
	parts = append(parts, strconv.Itoa(limit))
	return cli.parseDomainSet(strings.Join(parts, " "))
}

func (cli *Cli) ResumeDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	_, err := cli.query("RESUMEDOMAIN " + Marshal(domain))
	return err
}

func (cli *Cli) SetAccountDefaults(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	parts := []string{"SETACCOUNTDEFAULTS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetAccountDefaultPrefs(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings map is required")
	}
	parts := []string{"SETACCOUNTDEFAULTPREFS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetAccountTemplate(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	parts := []string{"SETACCOUNTTEMPLATE"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetDomainAliases(domain string, aliases []string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	parts := []string{"SETDOMAINALIASES", Marshal(domain), Marshal(aliases)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetDomainMailRules(domain string, rules []any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	parts := []string{"SETDOMAINMAILRULES", Marshal(domain), Marshal(rules)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetDomainSettings(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	parts := []string{"SETDOMAINSETTINGS", Marshal(domain), Marshal(settings)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetDomainSignalRules(domain string, rules []any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	parts := []string{"SETDOMAINSIGNALRULES", Marshal(domain), Marshal(rules)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SuspendDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	_, err := cli.query("SUSPENDDOMAIN " + Marshal(domain))
	return err
}

func (cli *Cli) UpdateAccountDefaults(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings map is required")
	}
	parts := []string{"UPDATEACCOUNTDEFAULTS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) UpdateAccountDefaultPrefs(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings map is required")
	}
	parts := []string{"UPDATEACCOUNTDEFAULTPREFS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) UpdateAccountTemplate(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings map is required")
	}
	parts := []string{"UPDATEACCOUNTTEMPLATE"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) UpdateDomainSettings(domain string, settings map[string]any) error {
	if settings == nil {
		return nil
	}
	parts := []string{"UPDATEDOMAINSETTINGS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}
	parts = append(parts, Marshal(settings))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// -----------------------------------------------------------------------------
// Helpers / Extended API
// -----------------------------------------------------------------------------

func (cli *Cli) IsDomainExists(domain string) (bool, error) {
	location, err := cli.GetDomainLocation(domain)
	if err != nil {
		return false, err
	}
	return location != "", nil
}

func (cli *Cli) parseDomainSet(cmd string) (*DomainSetResult, error) {
	res, err := cli.query(cmd)
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("parseDomainSet: expected dictionary, got %T", res)
	}
	result := &DomainSetResult{Items: make(map[string]string)}
	for k, v := range dict {
		if k == "" {
			result.Total, _ = strconv.Atoi(fmt.Sprint(v))
			continue
		}
		result.Items[k] = fmt.Sprint(v)
	}
	return result, nil
}

func (t ObjectType) ToCGPString() string {
	var parts []string
	if t&TypeAccount != 0 {
		parts = append(parts, "ACCOUNTS")
	}
	if t&TypeAlias != 0 {
		parts = append(parts, "ALIASES")
	}
	if t&TypeForwarder != 0 {
		parts = append(parts, "FORWARDERS")
	}
	return strings.Join(parts, " ")
}
