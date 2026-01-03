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
	TypeAccount   ObjectType = 1 << iota
	TypeAlias
	TypeForwarder
	TypeAll       = TypeAccount | TypeAlias | TypeForwarder
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

func (cli *Cli) GetDomainAliases(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required for GETDOMAINALIASES")
	}
	res, err := cli.query("GETDOMAINALIASES " + Marshal(domain))
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	aliases := make([]string, len(raw))
	for i, v := range raw {
		aliases[i] = fmt.Sprint(v)
	}
	return aliases, nil
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
	return fmt.Sprint(res), nil
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

func (cli *Cli) IsDomainExists(domain string) (bool, error) {
	loc, err := cli.GetDomainLocation(domain)
	if err != nil {
		return false, err
	}
	return loc != "", nil
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
		return make(map[string]string), nil
	}
	accounts := make(map[string]string, len(dict))
	for name, val := range dict {
		accounts[name] = fmt.Sprint(val)
	}
	return accounts, nil
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
		return nil, fmt.Errorf("cgpcli: ListDomainObjects unexpected response format")
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
		return nil, fmt.Errorf("domain name is required for LISTDOMAINTELNUMS")
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

func (cli *Cli) SetDomainAliases(domain string, aliases []string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required for SETDOMAINALIASES")
	}
	parts := []string{"SETDOMAINALIASES", Marshal(domain), Marshal(aliases)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) SetDomainSettings(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required for SETDOMAINSETTINGS")
	}
	parts := []string{"SETDOMAINSETTINGS", Marshal(domain), Marshal(settings)}
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
// Helpers
// -----------------------------------------------------------------------------

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
