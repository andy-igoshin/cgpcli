// # Section: Alert Administration
package cgpcli

import (
	"fmt"
)

func (cli *Cli) GetAccountAlerts(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account is required")
	}
	const cmd = "GETACCOUNTALERTS"
	res, err := cli.getMapAny(cmd, Atom(account))
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

func (cli *Cli) GetClusterAlerts() (map[string]string, error) {
	const cmd = "GETCLUSTERALERTS"
	res, err := cli.getMapAny(cmd)
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

func (cli *Cli) GetDomainAlerts(domain string) (map[string]string, error) {
	const cmd = "GETDOMAINALERTS"
	res, err := cli.getMapAny(cmd, Atom(domain))
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

func (cli *Cli) GetServerAlerts() (map[string]string, error) {
	const cmd = "GETSERVERALERTS"
	res, err := cli.getMapAny(cmd)
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

func (cli *Cli) PostAccountAlert(account, alert string) error {
	if account == "" || alert == "" {
		return fmt.Errorf("account and alert text are required")
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
		return fmt.Errorf("domain and alert text are required")
	}
	return cli.QueryNV("POSTDOMAINALERT", Atom(domain), "ALERT", alert)
}

func (cli *Cli) PostServerAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	return cli.QueryNV("POSTSERVERALERT", alert)
}

func (cli *Cli) RemoveAccountAlert(account, timeStamp string) error {
	if account == "" || timeStamp == "" {
		return fmt.Errorf("account and timeStamp are required")
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
		return fmt.Errorf("domain and timeStamp are required")
	}
	return cli.QueryNV("REMOVEDOMAINALERT", Atom(domain), "ALERT", timeStamp)
}

func (cli *Cli) RemoveServerAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	return cli.QueryNV("REMOVESERVERALERT", timeStamp)
}

func (cli *Cli) SetAccountAlerts(account string, alerts map[string]any) error {
	if account == "" {
		return fmt.Errorf("account is required")
	}
	return cli.setMapAny("SETACCOUNTALERTS", account, alerts)
}

func (cli *Cli) SetClusterAlerts(alerts map[string]any) error {
	return cli.setMapAny("SETCLUSTERALERTS", "", alerts)
}

func (cli *Cli) SetDomainAlerts(domain string, alerts map[string]any) error {
	return cli.setMapAny("SETDOMAINALERTS", domain, alerts)
}

func (cli *Cli) SetServerAlerts(alerts map[string]any) error {
	return cli.setMapAny("SETSERVERALERTS", "", alerts)
}
