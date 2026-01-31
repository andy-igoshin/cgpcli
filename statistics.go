// # Section: Statistics
package cgpcli

import (
	"fmt"
)

// GetAccountStat возвращает статистику аккаунта.
func (cli *Cli) GetAccountStat(account string, key string) (res any, err error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTSTAT"
	if key != "" {
		res, err = cli.Query(cmd, Atom(account), "KEY", key)
	} else {
		res, err = cli.Query(cmd, Atom(account))
	}
	return
}

// GetDomainStat возвращает статистику домена.
func (cli *Cli) GetDomainStat(domain string, key string) (res any, err error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}

	const cmd = "GETDOMAINSTAT"
	if key != "" {
		res, err = cli.Query(cmd, Atom(domain), "KEY", key)
	} else {
		res, err = cli.Query(cmd, Atom(domain))
	}
	return
}

// ResetAccountStat сбрасывает статистику аккаунта.
func (cli *Cli) ResetAccountStat(account string, key string) (err error) {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	const cmd = "RESETACCOUNTSTAT"
	if key != "" {
		err = cli.QueryNV(cmd, Atom(account), "KEY", key)
	} else {
		err = cli.QueryNV(cmd, Atom(account))
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
		err = cli.QueryNV(cmd, Atom(domain), "KEY", key)
	} else {
		err = cli.QueryNV(cmd, Atom(domain))
	}
	return
}
