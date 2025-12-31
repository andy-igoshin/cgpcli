package cgpcli

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (cli *Cli) CreateAccount(account string, settings map[string]any) error {

	if account == "" {
		return fmt.Errorf("account name is required")
	}

	// Формат: CREATEACCOUNT name { settings }
	parts := make([]string, 0, 3)
	parts = append(parts, "CREATEACCOUNT", account)

	if settings != nil {
		parts = append(parts, Marshal(settings))
	}

	_, err := cli.query(strings.Join(parts, " "))

	return err
}

func (cli *Cli) DeleteAccount(account string) error {

	if account == "" {
		return fmt.Errorf("account name is required")
	}

	parts := make([]string, 0, 2)
	parts = append(parts, "DELETEACCOUNT", account)

	_, err := cli.query(strings.Join(parts, " "))

	return err
}

func (cli *Cli) GetAccountInfo(account string, key any) (any, error) {

	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	// Минимальная емкость 2 (команда + имя), может вырасти до 4
	parts := make([]string, 0, 4)
	parts = append(parts, "GETACCOUNTINFO", account)

	if key != nil {
		switch v := key.(type) {
		case string:
			// Формат: GETACCOUNTINFO account Key "name"
			parts = append(parts, "Key", Marshal(v))
		case []string:
			// Превращаем []string в []any для Marshal, чтобы получить (val1, val2)
			anyList := make([]any, len(v))
			for i, s := range v {
				anyList[i] = s
			}
			parts = append(parts, Marshal(anyList))
		case []any:
			parts = append(parts, Marshal(v))
		}
	}

	return cli.query(strings.Join(parts, " "))
}

func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {

	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	// Команда GETACCOUNT возвращает структуру настроек
	res, err := cli.query("GETACCOUNTSETTINGS " + account)
	if err != nil {
		return nil, err
	}

	settings, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTSETTINGS: expected map, got %T", res)
	}

	return settings, nil
}

func (cli *Cli) ListAccounts(domain string) ([]string, error) {

	parts := make([]string, 0, 2)
	parts = append(parts, "LISTACCOUNTS")
	if domain != "" {
		parts = append(parts, domain)
	}

	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("LISTACCOUNTS: expected map, got %T", res)
	}

	accounts := make([]string, 0, len(dict))
	for name := range dict {
		accounts = append(accounts, name)
	}
	sort.Strings(accounts)

	return accounts, nil
}

func (cli *Cli) ListDomainObjects(domainName string, limit int, filter string, what string, cookie string) ([]any, error) {
	var parts []string

	parts = append(parts, "LISTDOMAINOBJECTS", Marshal(domainName))

	if filter != "" {
		parts = append(parts, "FILTER", Marshal(filter))
	}

	parts = append(parts, strconv.Itoa(limit))

	if what != "" {
		parts = append(parts, what)
	}

	parts = append(parts, "COOKIE", Marshal(cookie))

	// Вызываем query, который сам делает Send и Unmarshal
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	result, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("cgpcli: ListDomainObjects unexpected response type: %T", res)
	}

	return result, nil
}

func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {

	if account == "" {
		return fmt.Errorf("account name is required")
	}

	parts := make([]string, 0, 3)
	parts = append(parts, "UPDATEACCOUNTSETTINGS", account, Marshal(settings))

	_, err := cli.query(strings.Join(parts, " "))

	return err
}
