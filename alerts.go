// # Alert Administration
//
// The Alerts section provides functionality for managing system and user notifications
// across all administrative levels: Server, Cluster, Domain, and Account.
//
// Key capabilities include:
//   - Retrieving active alerts using "GET...ALERTS" commands.
//   - Posting new custom alerts for administrators or users to see.
//   - Removing specific alerts using their unique timestamps.
//   - Batch updating alert dictionaries.
package cgpcli

import (
	"fmt"
)

// GetAccountAlerts retrieves the Account Alerts.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// This method executes the GETACCOUNTALERTS CLI command.
//
// Returns a map[string]string containing the alerts and timestamps,
// or an error if the account name is empty or the server fails
// to execute the command.
//
// This method executes the GETACCOUNTALERTS CLI command.
func (cli *Cli) GetAccountAlerts(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapString("GETACCOUNTALERTS", Atom(account))
}

// GetClusterAlerts retrieves the cluster-wide Alerts.
//
// This method executes the GETCLUSTERALERTS CLI command.
//
// Returns a map[string]string containing the cluster alerts and timestamps,
// or an error if the server fails to execute the command.
//
// This method executes the GETCLUSTERALERTS CLI command.
func (cli *Cli) GetClusterAlerts() (map[string]string, error) {
	return cli.getMapString("GETCLUSTERALERTS")
}

// GetDomainAlerts retrieves the domain-wide Alerts.
//
// Parameters:
//   - domain: an optional Domain name.
//     If an empty string is provided, the command typically targets the
//     current domain context.
//
// This method executes the GETDOMAINALERTS CLI command.
//
// Returns a map[string]string containing the domain alerts and timestamps,
// or an error if the server fails to execute the command.
//
// This method executes the GETDOMAINALERTS CLI command.
func (cli *Cli) GetDomainAlerts(domain string) (map[string]string, error) {
	return cli.getMapString("GETDOMAINALERTS", Atom(domain))
}

// GetServerAlerts retrieves the server-wide Alerts.
//
// This method executes the GETSERVERALERTS CLI command.
//
// Returns a map[string]string containing the server alerts and timestamps,
// or an error if the server fails to execute the command.
//
// This method executes the GETSERVERALERTS CLI command.
func (cli *Cli) GetServerAlerts() (map[string]string, error) {
	return cli.getMapString("GETSERVERALERTS")
}

// PostAccountAlert posts an alert message to a specific Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - alert: the text of the alert to be posted.
//
// This method executes the POSTACCOUNTALERT CLI command.
//
// Returns an error if the account name or alert text is empty, or if
// the server fails to execute the command.
//
// This method executes the POSTACCOUNTALERT CLI command.
func (cli *Cli) PostAccountAlert(account, alert string) error {
	if account == "" || alert == "" {
		return fmt.Errorf("account name and alert text are required")
	}
	return cli.QueryNV("POSTACCOUNTALERT", Atom(account), "ALERT", alert)
}

// PostClusterAlert posts a cluster-wide Alert message.
//
// Parameters:
//   - alert: the text of the alert to be posted.
//
// This method executes the POSTCLUSTERALERT CLI command.
//
// Returns an error if the alert text is empty or if the server fails
// to execute the command.
//
// This method executes the POSTCLUSTERALERT CLI command.
func (cli *Cli) PostClusterAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	return cli.QueryNV("POSTCLUSTERALERT", alert)
}

// PostDomainAlert posts a domain-wide Alert message.
//
// Parameters:
//   - domain: the Domain name.
//   - alert: the text of the alert to be posted.
//
// This method executes the POSTDOMAINALERT CLI command.
//
// Returns an error if the domain name or alert text is empty, or if
// the server fails to execute the command.
//
// This method executes the POSTDOMAINALERT CLI command.
func (cli *Cli) PostDomainAlert(domain, alert string) error {
	if domain == "" || alert == "" {
		return fmt.Errorf("domain name and alert text are required")
	}
	return cli.QueryNV("POSTDOMAINALERT", domain, "ALERT", alert)
}

// PostServerAlert posts a server-wide Alert message.
//
// Parameters:
//   - alert: the text of the alert to be posted.
//
// This method executes the POSTSERVERALERT CLI command.
//
// Returns an error if the alert text is empty or if the server fails
// to execute the command.
//
// This method executes the POSTSERVERALERT CLI command.
func (cli *Cli) PostServerAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	return cli.QueryNV("POSTSERVERALERT", alert)
}

// RemoveAccountAlert removes a specific alert message from an Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - timeStamp: the exact timestamp string of the alert to be removed,
//     as retrieved by GetAccountAlerts.
//
// This method executes the REMOVEACCOUNTALERT CLI command.
//
// Returns an error if the account name or timeStamp is empty, or if
// the server fails to execute the command.
//
// This method executes the REMOVEACCOUNTALERT CLI command.
func (cli *Cli) RemoveAccountAlert(account, timeStamp string) error {
	if account == "" || timeStamp == "" {
		return fmt.Errorf("account name and timeStamp are required")
	}
	return cli.QueryNV("REMOVEACCOUNTALERT", Atom(account), "ALERT", timeStamp)
}

// RemoveClusterAlert removes a specific cluster-wide alert message.
//
// Parameters:
//   - timeStamp: the exact timestamp string of the cluster-wide alert
//     to be removed, as retrieved by GetClusterAlerts.
//
// This method executes the REMOVECLUSTERALERT CLI command.
//
// Returns an error if the timeStamp is empty or if the server
// fails to execute the command.
//
// This method executes the REMOVECLUSTERALERT CLI command.
func (cli *Cli) RemoveClusterAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	return cli.QueryNV("REMOVECLUSTERALERT", timeStamp)
}

// RemoveDomainAlert removes a specific domain-wide alert message.
//
// Parameters:
//   - domain: the Domain name.
//   - timeStamp: the exact timestamp string of the domain alert
//     to be removed, as retrieved by GetDomainAlerts.
//
// This method executes the REMOVEDOMAINALERT CLI command.
//
// Returns an error if the domain name or timeStamp is empty, or if
// the server fails to execute the command.
//
// This method executes the REMOVEDOMAINALERT CLI command.
func (cli *Cli) RemoveDomainAlert(domain, timeStamp string) error {
	if domain == "" || timeStamp == "" {
		return fmt.Errorf("domain name and timeStamp are required")
	}
	return cli.QueryNV("REMOVEDOMAINALERT", domain, "ALERT", timeStamp)
}

// RemoveServerAlert removes a specific server-wide alert message.
//
// Parameters:
//   - timeStamp: the exact timestamp string of the server-wide alert
//     to be removed, as retrieved by GetServerAlerts.
//
// This method executes the REMOVESERVERALERT CLI command.
//
// Returns an error if the timeStamp is empty or if the server
// fails to execute the command.
//
// This method executes the REMOVESERVERALERT CLI command.
func (cli *Cli) RemoveServerAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	return cli.QueryNV("REMOVESERVERALERT", timeStamp)
}

// SetAccountAlerts replaces the entire Account alert dictionary.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - alerts: a dictionary used to replace the current Account alert dictionary.
//
// This method executes the SETACCOUNTALERTS CLI command.
//
// Returns an error if the account name is empty, the alerts map is nil,
// or if the server fails to execute the command.
//
// This method executes the SETACCOUNTALERTS CLI command.
func (cli *Cli) SetAccountAlerts(account string, alerts map[string]any) error {
	if account == "" || alerts == nil {
		return fmt.Errorf("account name and alerts dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTALERTS", Atom(account), alerts)
}

// SetClusterAlerts replaces the entire cluster-wide Alert dictionary.
//
// Parameters:
//   - alerts: a dictionary used to replace the current Cluster alert dictionary.
//
// This method executes the SETCLUSTERALERTS CLI command.
//
// Returns an error if the server fails to execute the command.
//
// This method executes the SETCLUSTERALERTS CLI command.
func (cli *Cli) SetClusterAlerts(alerts map[string]any) error {
	return cli.QueryNV("SETCLUSTERALERTS", alerts)
}

// SetDomainAlerts replaces the entire Domain alert dictionary.
//
// Parameters:
//   - domain: an optional Domain name.
//     If an empty string is provided, the command targets the current domain context.
//   - alerts: a dictionary used to replace the current Domain alert dictionary.
//
// This method executes the SETDOMAINALERTS CLI command.
//
// Returns an error if the server fails to execute the command.
//
// This method executes the SETDOMAINALERTS CLI command.
func (cli *Cli) SetDomainAlerts(domain string, alerts map[string]any) error {
	return cli.QueryNV("SETDOMAINALERTS", Atom(domain), alerts)
}

// SetServerAlerts replaces the entire server-wide Alert dictionary.
//
// Parameters:
//   - alerts: a dictionary used to replace the current Server alert dictionary.
//
// This method executes the SETSERVERALERTS CLI command.
//
// Returns an error if the server fails to execute the command.
//
// This method executes the SETSERVERALERTS CLI command.
func (cli *Cli) SetServerAlerts(alerts map[string]any) error {
	return cli.QueryNV("SETSERVERALERTS", alerts)
}
