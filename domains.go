package cgpcli

import (
	"fmt"
	"strconv"
	"strings"
)

// ObjectType определяет битовую маску для категорий объектов домена
type ObjectType int

const (
	TypeAccount   ObjectType = 1 << iota // 1
	TypeAlias                            // 2
	TypeForwarder                        // 4
	TypeAll       = TypeAccount | TypeAlias | TypeForwarder
)

// DomainObject представляет детализированный объект домена
type DomainObject struct {
	Name    string
	Type    ObjectType
	Subtype string
	Targets []string
}

// DomainObjectsResult структура ответа с поддержкой пагинации
type DomainObjectsResult struct {
	Objects map[string]DomainObject
	Total   int
	Cookie  string
}

// DomainSetResult представляет результат глобального поиска по всем доменам
type DomainSetResult struct {
	Items map[string]string // Карта Объект -> Аккаунт-владелец
	Total int               // Общее количество найденных записей
}

// -----------------------------------------------------------------------------
// Cli Methods (Alphabetical Order)
// -----------------------------------------------------------------------------

func (cli *Cli) CreateDomain(domain string, shared bool, storage string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	parts := []string{"CREATEDOMAIN", Marshal(domain)}
	if shared {
		parts = append(parts, "SHARED")
	}
	if storage != "" {
		parts = append(parts, "PATH", Marshal(storage))
	}
	if settings != nil {
		parts = append(parts, Marshal(settings))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

func (cli *Cli) DeleteDomain(domain string, force bool) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	parts := []string{"DELETEDOMAIN", Marshal(domain)}
	if force {
		parts = append(parts, "FORCE")
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

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
		if _, isStr := res.(string); isStr {
			return make(map[string]string), nil
		}
		return nil, fmt.Errorf("LISTACCOUNTS: expected dictionary response, got %T", res)
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

func (cli *Cli) ListDomains() ([]string, error) {
	res, err := cli.query("LISTDOMAINS")
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	domains := make([]string, len(raw))
	for i, v := range raw {
		domains[i] = fmt.Sprint(v)
	}
	return domains, nil
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

func (cli *Cli) ListServerTelnums(filter string, limit int) (*DomainSetResult, error) {
	parts := []string{"LISTSERVERTELNUMS"}
	if filter != "" {
		parts = append(parts, "FILTER", Marshal(filter))
	}
	parts = append(parts, strconv.Itoa(limit))
	return cli.parseDomainSet(strings.Join(parts, " "))
}

func (cli *Cli) RenameDomain(oldName, newName string, storage string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new domain names are required")
	}
	parts := []string{"RENAMEDOMAIN", Marshal(oldName), "INTO", Marshal(newName)}
	if storage != "" {
		parts = append(parts, "PATH", Marshal(storage))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
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
	anyAliases := make([]any, len(aliases))
	for i, a := range aliases {
		anyAliases[i] = a
	}
	parts := []string{"SETDOMAINALIASES", Marshal(domain), Marshal(anyAliases)}
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
// Type & Parser Helpers
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

func (t ObjectType) String() string {
	switch t {
	case TypeAccount:
		return "Account"
	case TypeAlias:
		return "Alias"
	case TypeForwarder:
		return "Forwarder"
	default:
		return "Unknown"
	}
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

// -----------------------------------------------------------------------------
// Extended API
// -----------------------------------------------------------------------------

// IsDomainExists проверяет существование домена.
func (cli *Cli) IsDomainExists(domain string) (bool, error) {
	loc, err := cli.GetDomainLocation(domain)
	if err != nil {
		return false, err
	}
	return loc != "", nil
}
