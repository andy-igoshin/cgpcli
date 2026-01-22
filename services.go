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
	const cmd = "BALANCE"
	res, err := cli.Query(cmd, Atom(account), params)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// Dataset выполняет операции с наборами данных аккаунта (адресные книги и т.д.).
func (cli *Cli) Dataset(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	const cmd = "DATASET"
	res, err := cli.Query(cmd, Atom(account), params)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// RemoveAccountSubset удаляет конкретный "dataset" аккаунта (например, 'RepliedAddresses').
func (cli *Cli) RemoveAccountSubset(account, subset string) (err error) {
	if account == "" || subset == "" {
		return fmt.Errorf("account name and subset name are required")
	}
	_, err = cli.Query("REMOVEACCOUNTSUBSET", Atom(account), "SUBSET", subset)
	return err
}

// Roster выполняет операции с ростером (списком контактов) аккаунта.
func (cli *Cli) Roster(account string, params map[string]any) (map[string]any, error) {
	if account == "" || len(params) == 0 {
		return nil, fmt.Errorf("account name and params are required")
	}
	const cmd = "ROSTER"
	res, err := cli.Query(cmd, Atom(account), params)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}
