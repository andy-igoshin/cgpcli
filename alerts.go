// # Section: Alert Administration
package cgpcli

import (
	"fmt"
)

func (cli *Cli) GetAccountAlerts(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account is required")
	}
	res, err := cli.getMapAny("GETACCOUNTALERTS", account)
	if err != nil {
		return nil, err
	}
	return toStringMap("GETACCOUNTALERTS", res)
}

func (cli *Cli) GetClusterAlerts() (map[string]string, error) {
	res, err := cli.getMapAny("GETCLUSTERALERTS", "")
	if err != nil {
		return nil, err
	}
	return toStringMap("GETCLUSTERALERTS", res)
}

func (cli *Cli) GetDomainAlerts(domain string) (map[string]string, error) {
	res, err := cli.getMapAny("GETDOMAINALERTS", domain)
	if err != nil {
		return nil, err
	}
	return toStringMap("GETDOMAINALERTS", res)
}

func (cli *Cli) GetServerAlerts() (map[string]string, error) {
	res, err := cli.getMapAny("GETSERVERALERTS", "")
	if err != nil {
		return nil, err
	}
	return toStringMap("GETSERVERALERTS", res)
}

func (cli *Cli) PostAccountAlert(account, alert string) error {
	if account == "" || alert == "" {
		return fmt.Errorf("account and alert text are required")
	}
	_, err := cli.Query("POSTACCOUNTALERT", Atom(account), "ALERT", alert)
	return err
}

func (cli *Cli) PostClusterAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	_, err := cli.Query("POSTCLUSTERALERT", alert)
	return err
}

func (cli *Cli) PostDomainAlert(domain, alert string) error {
	if domain == "" || alert == "" {
		return fmt.Errorf("domain and alert text are required")
	}
	_, err := cli.Query("POSTDOMAINALERT", Atom(domain), "ALERT", alert)
	return err
}

func (cli *Cli) PostServerAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert text is required")
	}
	_, err := cli.Query("POSTSERVERALERT", alert)
	return err
}

func (cli *Cli) RemoveAccountAlert(account, timeStamp string) error {
	if account == "" || timeStamp == "" {
		return fmt.Errorf("account and timeStamp are required")
	}
	_, err := cli.Query("REMOVEACCOUNTALERT", Atom(account), "ALERT", timeStamp)
	return err
}

func (cli *Cli) RemoveClusterAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	_, err := cli.Query("REMOVECLUSTERALERT", timeStamp)
	return err
}

func (cli *Cli) RemoveDomainAlert(domain, timeStamp string) error {
	if domain == "" || timeStamp == "" {
		return fmt.Errorf("domain and timeStamp are required")
	}
	_, err := cli.Query("REMOVEDOMAINALERT", Atom(domain), "ALERT", timeStamp)
	return err
}

func (cli *Cli) RemoveServerAlert(timeStamp string) error {
	if timeStamp == "" {
		return fmt.Errorf("timeStamp is required")
	}
	_, err := cli.Query("REMOVESERVERALERT", timeStamp)
	return err
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
