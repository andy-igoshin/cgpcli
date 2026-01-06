package cgpcli

import (
	"fmt"
)

// CreateForwarder создает новый Forwarder.
func (cli *Cli) CreateForwarder(forwarder, target string) error {
	if forwarder == "" || target == "" {
		return fmt.Errorf("forwarder name and target address are required")
	}
	_, err := cli.Query("CREATEFORWARDER", forwarder, "TO", target)
	return err
}

// DeleteForwarder удаляет Forwarder.
func (cli *Cli) DeleteForwarder(forwarder string) error {
	if forwarder == "" {
		return fmt.Errorf("forwarder name is required")
	}
	_, err := cli.Query("DELETEFORWARDER", forwarder)
	return err
}

// FindForwarders ищет все Forwarders в домене, указывающие на конкретный адрес.
func (cli *Cli) FindForwarders(domain, target string) ([]string, error) {
	if domain == "" || target == "" {
		return nil, fmt.Errorf("domain and target address are required")
	}
	res, err := cli.Query("FINDFORWARDERS", domain, "TO", target)
	if err != nil {
		return nil, err
	}
	return cli.toStringSlice(res), nil
}

// GetForwarder возвращает адрес, на который перенаправляется почта.
func (cli *Cli) GetForwarder(forwarder string) (string, error) {
	if forwarder == "" {
		return "", fmt.Errorf("forwarder name is required")
	}
	res, err := cli.Query("GETFORWARDER", forwarder)
	if err != nil {
		return "", err
	}

	if s, ok := res.(string); ok {
		return s, nil
	}
	return fmt.Sprint(res), nil
}

// ListForwarders возвращает список всех Forwarders в указанном домене.
func (cli *Cli) ListForwarders(domain string) ([]string, error) {
	// Используем make для сохранения линейной структуры кода
	args := make([]any, 0, 2)
	args = append(args, "LISTFORWARDERS")
	if domain != "" {
		args = append(args, domain)
	}

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	return cli.toStringSlice(res), nil
}

// RenameForwarder переименовывает Forwarder.
func (cli *Cli) RenameForwarder(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new forwarder names are required")
	}
	_, err := cli.Query("RENAMEFORWARDER", oldName, "INTO", newName)
	return err
}

// toStringSlice — вспомогательный метод для эффективного приведения типов
func (cli *Cli) toStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return []string{}
	}

	list := make([]string, len(raw))
	for i, item := range raw {
		if s, ok := item.(string); ok {
			list[i] = s
		} else {
			list[i] = fmt.Sprint(item)
		}
	}
	return list
}
