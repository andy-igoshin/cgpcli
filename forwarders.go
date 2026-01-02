package cgpcli

import (
	"fmt"
	"strings"
)

// CreateForwarder создает новый Forwarder.
func (cli *Cli) CreateForwarder(forwarder, target string) error {
	if forwarder == "" || target == "" {
		return fmt.Errorf("forwarder name and target address are required")
	}
	parts := []string{"CREATEFORWARDER", Marshal(forwarder), "TO", Marshal(target)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// DeleteForwarder удаляет Forwarder.
func (cli *Cli) DeleteForwarder(forwarder string) error {
	if forwarder == "" {
		return fmt.Errorf("forwarder name is required")
	}
	_, err := cli.query("DELETEFORWARDER " + Marshal(forwarder))
	return err
}

// FindForwarders ищет все Forwarders в домене, указывающие на конкретный адрес.
func (cli *Cli) FindForwarders(domain, target string) ([]string, error) {
	if domain == "" || target == "" {
		return nil, fmt.Errorf("domain and target address are required")
	}
	parts := []string{"FINDFORWARDERS", Marshal(domain), "TO", Marshal(target)}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}

	list := make([]string, len(raw))
	for i, v := range raw {
		list[i] = fmt.Sprint(v)
	}
	return list, nil
}

// GetForwarder возвращает адрес, на который перенаправляется почта.
func (cli *Cli) GetForwarder(forwarder string) (string, error) {
	if forwarder == "" {
		return "", fmt.Errorf("forwarder name is required")
	}
	res, err := cli.query("GETFORWARDER " + Marshal(forwarder))
	if err != nil {
		return "", err
	}
	return fmt.Sprint(res), nil
}

// ListForwarders возвращает список всех Forwarders в указанном домене.
func (cli *Cli) ListForwarders(domain string) ([]string, error) {
	parts := []string{"LISTFORWARDERS"}
	if domain != "" {
		parts = append(parts, Marshal(domain))
	}

	var res any
	var err error
	if len(parts) > 1 {
		res, err = cli.query(strings.Join(parts, " "))
	} else {
		res, err = cli.query(parts[0])
	}

	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}

	list := make([]string, len(raw))
	for i, v := range raw {
		list[i] = fmt.Sprint(v)
	}
	return list, nil
}

// RenameForwarder переименовывает Forwarder.
func (cli *Cli) RenameForwarder(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new forwarder names are required")
	}
	parts := []string{"RENAMEFORWARDER", Marshal(oldName), "INTO", Marshal(newName)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}
