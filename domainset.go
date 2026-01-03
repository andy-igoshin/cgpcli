package cgpcli

import (
	"fmt"
	"strings"
)

// CreateDomain создает новый домен.
func (cli *Cli) CreateDomain(domain string, shared bool, storage string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	parts := []string{"CREATEDOMAIN", Marshal(domain)}
	if shared {
		parts = append(parts, "SHARED")
	}
	if storage != "" {
		parts = append(parts, "PATH", Marshal(storage))
	}
	if settings != nil {
		parts = append(parts, Marshal(settings))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// DeleteDomain удаляет домен.
func (cli *Cli) DeleteDomain(domain string, force bool) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	parts := []string{"DELETEDOMAIN", Marshal(domain)}
	if force {
		parts = append(parts, "FORCE")
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// ListDomains возвращает список имен всех доменов на сервере.
func (cli *Cli) ListDomains() ([]string, error) {
	res, err := cli.query("LISTDOMAINS")
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	domains := make([]string, len(raw))
	for i, v := range raw {
		domains[i] = fmt.Sprint(v)
	}
	return domains, nil // Исправлено: добавлен nil для ошибки
}

// ListServerTelnums возвращает список телефонных номеров сервера.
func (cli *Cli) ListServerTelnums(filter string, limit int) (*DomainSetResult, error) {
	parts := []string{"LISTSERVERTELNUMS"}
	if filter != "" {
		parts = append(parts, "FILTER", Marshal(filter))
	}
	parts = append(parts, fmt.Sprint(limit))
	return cli.parseDomainSet(strings.Join(parts, " "))
}

// RenameDomain переименовывает домен.
func (cli *Cli) RenameDomain(oldName, newName string, storage string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new domain names are required")
	}
	parts := []string{"RENAMEDOMAIN", Marshal(oldName), "INTO", Marshal(newName)}
	if storage != "" {
		parts = append(parts, "PATH", Marshal(storage))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}
