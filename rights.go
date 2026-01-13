// # Section: Access Rights Administration
package cgpcli

import (
	"fmt"
)

// SetAccountRights устанавливает права доступа к серверу для аккаунта.
// Внимание: все старые права будут удалены и заменены новыми.
func (cli *Cli) SetAccountRights(account string, rights []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETACCOUNTRIGHTS", account, rights)
	return err
}
