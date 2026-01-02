package cgpcli

import (
	"fmt"
	"strings"
)

// SETACCOUNTRIGHTS устанавливает права доступа к серверу для аккаунта.
// Внимание: все старые права будут удалены и заменены новыми.
func (cli *Cli) SetAccountRights(account string, rights []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	anyRights := make([]any, len(rights))
	for i, r := range rights {
		anyRights[i] = r
	}

	// Marshal(account) гарантирует, что имя будет в кавычках, если это нужно
	parts := []string{"SETACCOUNTRIGHTS", Marshal(account), Marshal(anyRights)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}
