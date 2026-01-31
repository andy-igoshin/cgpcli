// # Section: Access Rights Administration
package cgpcli

import (
	"fmt"
)

// SetAccountRights устанавливает права доступа к серверу для аккаунта.
func (cli *Cli) SetAccountRights(account string, rights []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("SETACCOUNTRIGHTS", account, rights)
}
