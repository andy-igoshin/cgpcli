// # Section: Alert Administration
package cgpcli

import (
	"fmt"
)

func (cli *Cli) GetAccountAlerts(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapString("GETACCOUNTALERTS", Atom(account))
}

func (cli *Cli) GetClusterAlerts() (map[string]string, error) {
	return cli.getMapString("GETCLUSTERALERTS")
}

func (cli *Cli) GetDomainAlerts(domain string) (map[string]string, error) {
	return cli.getMapString("GETDOMAINALERTS", Atom(domain))
}

func (cli *Cli) GetServerAlerts() (map[string]string, error) {
	return cli.getMapString("GETSERVERALERTS")
}

func (cli *Cli) PostAccountAlert(account, alert string) error {
	if account == "" || alert == "" {
		return fmt.Errorf("account name and alert text are required")
	}
	return cli.QueryNV("POSTACCOUNTALERT", Atom(account), "ALERT", alert)
}

func (cli *Cli) PostClusterAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	return cli.QueryNV("POSTCLUSTERALERT", alert)
}

func (cli *Cli) PostDomainAlert(domain, alert string) error {
	if domain == "" || alert == "" {
		return fmt.Errorf("domain name and alert text are required")
	}
	return cli.QueryNV("POSTDOMAINALERT", domain, "ALERT", alert)
}

func (cli *Cli) PostServerAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	return cli.QueryNV("POSTSERVERALERT", alert)
}

func (cli *Cli) RemoveAccountAlert(account, timeStamp string) error {
	if account == "" || timeStamp == "" {
		return fmt.Errorf("account name and timeStamp are required")
	}
	return cli.QueryNV("REMOVEACCOUNTALERT", Atom(account), "ALERT", timeStamp)
}

func (cli *Cli) RemoveClusterAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	return cli.QueryNV("REMOVECLUSTERALERT", timeStamp)
}

func (cli *Cli) RemoveDomainAlert(domain, timeStamp string) error {
	if domain == "" || timeStamp == "" {
		return fmt.Errorf("domain name and timeStamp are required")
	}
	return cli.QueryNV("REMOVEDOMAINALERT", domain, "ALERT", timeStamp)
}

func (cli *Cli) RemoveServerAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	return cli.QueryNV("REMOVESERVERALERT", timeStamp)
}

func (cli *Cli) SetAccountAlerts(account string, alerts map[string]any) error {
	if account == "" || alerts == nil {
		return fmt.Errorf("account name and alerts dictionary are required")
	}
	return cli.QueryNV("SETACCOUNTALERTS", Atom(account), alerts)
}

func (cli *Cli) SetClusterAlerts(alerts map[string]any) error {
	return cli.QueryNV("SETCLUSTERALERTS", alerts)
}

func (cli *Cli) SetDomainAlerts(domain string, alerts map[string]any) error {
	return cli.QueryNV("SETDOMAINALERTS", Atom(domain), alerts)
}

func (cli *Cli) SetServerAlerts(alerts map[string]any) error {
	return cli.QueryNV("SETSERVERALERTS", alerts)
}
