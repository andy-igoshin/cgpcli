// # Section: Statistics
package cgpcli

import (
	"fmt"
)

// GetAccountStat возвращает статистику аккаунта.
func (cli *Cli) GetAccountStat(account string, key string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	args := make([]any, 0, 3)
	args = append(args, Atom(account))
	if key != "" {
		args = append(args, "KEY", key)
	}
	return cli.Query("GETACCOUNTSTAT", args...)
}

// GetDomainStat возвращает статистику домена.
func (cli *Cli) GetDomainStat(domain string, key string) (any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}

	const cmd = "GETDOMAINSTAT"

	if key != "" {
		return cli.Query(cmd, Atom(domain), "KEY", key)
	} else {
		return cli.Query(cmd, Atom(domain))
	}
}

// ResetAccountStat сбрасывает статистику аккаунта.
func (cli *Cli) ResetAccountStat(account string, key string) (err error) {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	const cmd = "RESETACCOUNTSTAT"

	if key != "" {
		_, err = cli.Query(cmd, Atom(account), "KEY", key)
	} else {
		_, err = cli.Query(cmd, Atom(account))
	}
	return
}

// ResetDomainStat сбрасывает статистику домена.
func (cli *Cli) ResetDomainStat(domain string, key string) (err error) {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "RESETDOMAINSTAT"

	if key != "" {
		_, err = cli.Query(cmd, Atom(domain), "KEY", key)
	} else {
		_, err = cli.Query(cmd, Atom(domain))
	}
	return
}
