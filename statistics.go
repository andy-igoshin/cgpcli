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
	args := make([]any, 0, 4)
	args = append(args, "GETACCOUNTSTAT", Atom(account))
	if key != "" {
		args = append(args, "KEY", key)
	}
	return cli.Query(args...)
}

// GetDomainStat возвращает статистику домена.
func (cli *Cli) GetDomainStat(domain string, key string) (any, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	args := make([]any, 0, 4)
	args = append(args, "GETDOMAINSTAT", Atom(domain))
	if key != "" {
		args = append(args, "KEY", key)
	}
	return cli.Query(args...)
}

// ResetAccountStat сбрасывает статистику аккаунта.
func (cli *Cli) ResetAccountStat(account string, key string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	args := make([]any, 0, 4)
	args = append(args, "RESETACCOUNTSTAT", Atom(account))
	if key != "" {
		args = append(args, "KEY", key)
	}
	_, err := cli.Query(args...)
	return err
}

// ResetDomainStat сбрасывает статистику домена.
func (cli *Cli) ResetDomainStat(domain string, key string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	args := make([]any, 0, 4)
	args = append(args, "RESETDOMAINSTAT", Atom(domain))
	if key != "" {
		args = append(args, "KEY", key)
	}
	_, err := cli.Query(args...)
	return err
}
