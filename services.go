// # Section: Account Services
package cgpcli

import (
	"fmt"
)

// Balance выполняет операции с биллингом аккаунта (list, reserve, charge, credit и др.).
func (cli *Cli) Balance(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	return cli.getMapAny("BALANCE", Atom(account), params)
}

// Dataset выполняет операции с наборами данных аккаунта (адресные книги и т.д.).
func (cli *Cli) Dataset(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	return cli.getMapAny("DATASET", Atom(account), params)
}

// RemoveAccountSubset удаляет конкретный "dataset" аккаунта (например, 'RepliedAddresses').
func (cli *Cli) RemoveAccountSubset(account, subset string) error {
	if account == "" || subset == "" {
		return fmt.Errorf("account name and subset name are required")
	}
	return cli.QueryNV("REMOVEACCOUNTSUBSET", Atom(account), "SUBSET", subset)
}

// Roster выполняет операции с ростером (списком контактов) аккаунта.
func (cli *Cli) Roster(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	return cli.getMapAny("ROSTER", Atom(account), params)
}
