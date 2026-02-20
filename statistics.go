// # Statistics
//
// The Statistics section provides methods for monitoring the activity and
// resource consumption of specific domains and individual accounts.
// These metrics are invaluable for usage auditing, billing verification,
// and identifying unusual activity patterns.
//
// Key capabilities include:
//   - Granular Monitoring: retrieving specific counters or the entire statistics
//     dictionary for an account or domain.
//   - Reset Mechanisms: clearing accumulated counters via [Cli.ResetAccountStat]
//     and [Cli.ResetDomainStat], allowing for periodic (e.g., monthly or daily)
//     reporting cycles.
//   - Resource Tracking: accessing data related to protocol usage, login
//     frequency, and storage interaction at the entity level.
package cgpcli

import (
	"fmt"
)

// GetAccountStat retrieves statistics data about the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - key: an optional name of the statistical entry to retrieve.
//
// This method executes the GETACCOUNTSTAT CLI command.
//
// Returns:
//   - any: a number or a timeStamp with the specified statistical information,
//     or (if the key is not specified) a dictionary with all available statistical data.
//     If the statistical data for the specified key does not exist, an empty string is returned.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountStat(account string, key string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTSTAT"
	if key != "" {
		return cli.Query(cmd, Atom(account), "KEY", key)
	} else {
		return cli.Query(cmd, Atom(account))
	}
}

// GetDomainStat retrieves statistics data about the specified Domain.
//
// Parameters:
//   - domain: the Domain name. The asterisk (*) symbol can
//     be used to specify the Domain of the current authenticated Account.
//   - key: an optional name of the statistical entry to retrieve.
//
// This method executes the GETDOMAINSTAT CLI command.
//
// Returns:
//   - any: a string with the specified statistical information, or (if the
//     key is not specified) a dictionary with all available statistical data.
//   - error: an error if the command fails.
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

// ResetAccountStat resets statistics data about the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - key: an optional name of the statistical entry to reset.
//     If the key is not specified, all Account statistical entries are reset.
//
// This method executes the RESETACCOUNTSTAT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ResetAccountStat(account string, key string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	const cmd = "RESETACCOUNTSTAT"
	if key != "" {
		return cli.QueryNV(cmd, Atom(account), "KEY", key)
	} else {
		return cli.QueryNV(cmd, Atom(account))
	}
}

// ResetDomainStat resets statistics data about the specified Domain.
//
// Parameters:
//   - domain: the Domain name. The asterisk (*) symbol can
//     be used to specify the Domain of the current authenticated Account.
//   - key: an optional name of the statistical entry to reset.
//     If the key is not specified, all Domain statistical entries are reset.
//
// This method executes the RESETDOMAINSTAT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ResetDomainStat(domain string, key string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "RESETDOMAINSTAT"
	if key != "" {
		return cli.QueryNV(cmd, Atom(domain), "KEY", key)
	} else {
		return cli.QueryNV(cmd, Atom(domain))
	}
}
