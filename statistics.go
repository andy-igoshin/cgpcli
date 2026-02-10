// # Statistics
//
// The Statistics section provides methods for monitoring the activity and
// resource consumption of specific domains and individual accounts.
// These metrics are invaluable for usage auditing, billing verification,
// and identifying unusual activity patterns.
//
// Key capabilities include:
//   - Granular Monitoring: Retrieving specific counters or the entire statistics
//     dictionary for an account or domain.
//   - Reset Mechanisms: Clearing accumulated counters via [Cli.ResetAccountStat]
//     and [Cli.ResetDomainStat], allowing for periodic (e.g., monthly or daily)
//     reporting cycles.
//   - Resource Tracking: Accessing data related to protocol usage, login
//     frequency, and storage interaction at the entity level.
//
// These methods are essential for building per-customer reporting tools,
// monitoring SLA compliance, and automating resource-based policy enforcement.
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
