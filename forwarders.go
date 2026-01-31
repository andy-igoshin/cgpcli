// # Section: Forwarder Administration
package cgpcli

import (
	"fmt"
)

func (cli *Cli) CreateForwarder(forwarder, target string) error {
	if forwarder == "" || target == "" {
		return fmt.Errorf("forwarder name and target address are required")
	}
	return cli.QueryNV("CREATEFORWARDER", forwarder, "TO", target)
}

func (cli *Cli) DeleteForwarder(forwarder string) error {
	if forwarder == "" {
		return fmt.Errorf("forwarder name is required")
	}
	return cli.QueryNV("DELETEFORWARDER", forwarder)
}

func (cli *Cli) FindForwarders(domain, target string) ([]string, error) {
	if domain == "" || target == "" {
		return nil, fmt.Errorf("domain and target address are required")
	}

	const cmd = "FINDFORWARDERS"
	res, err := cli.Query(cmd, domain, "TO", target)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

func (cli *Cli) GetForwarder(forwarder string) (string, error) {
	if forwarder == "" {
		return "", fmt.Errorf("forwarder name is required")
	}
	return cli.getString("GETFORWARDER", forwarder)
}

func (cli *Cli) ListForwarders(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	return cli.getSliceString("LISTFORWARDERS", domain)
}

func (cli *Cli) RenameForwarder(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old name and new name are required")
	}
	return cli.QueryNV("RENAMEFORWARDER", oldName, "INTO", newName)
}
