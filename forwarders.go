package cgpcli

import (
	"fmt"
)

func (cli *Cli) CreateForwarder(forwarder, target string) error {
	if forwarder == "" || target == "" {
		return fmt.Errorf("forwarder name and target address are required")
	}
	_, err := cli.Query("CREATEFORWARDER", forwarder, "TO", target)
	return err
}

func (cli *Cli) DeleteForwarder(forwarder string) error {
	if forwarder == "" {
		return fmt.Errorf("forwarder name is required")
	}
	_, err := cli.Query("DELETEFORWARDER", forwarder)
	return err
}

func (cli *Cli) FindForwarders(domain, target string) ([]string, error) {
	if domain == "" || target == "" {
		return nil, fmt.Errorf("domain and target address are required")
	}
	res, err := cli.Query("FINDFORWARDERS", domain, "TO", target)
	if err != nil {
		return nil, err
	}
	return cli.toStringSlice2(res)
}

func (cli *Cli) GetForwarder(forwarder string) (string, error) {
	if forwarder == "" {
		return "", fmt.Errorf("forwarder name is required")
	}
	res, err := cli.Query("GETFORWARDER", forwarder)
	if err != nil {
		return "", err
	}
	return cli.valToString("GETFORWARDER", res)
}

func (cli *Cli) ListForwarders(domain string) ([]string, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain name is required")
	}
	res, err := cli.Query("LISTFORWARDERS", domain)
	if err != nil {
		return nil, err
	}
	return cli.toStringSlice2(res)
}

func (cli *Cli) RenameForwarder(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new forwarder names are required")
	}
	_, err := cli.Query("RENAMEFORWARDER", oldName, "INTO", newName)
	return err
}
