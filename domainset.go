// # Section: Domain Set Administration
package cgpcli

import (
	"fmt"
	"strings"
)

// CreateDirectoryDomain создает новый домен на базе справочника.
func (cli *Cli) CreateDirectoryDomain(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain is required")
	}
	args := make([]any, 0, 3)
	args = append(args, "CREATEDIRECTORYDOMAIN", domain)
	if settings != nil {
		args = append(args, settings)
	}
	_, err := cli.Query(args...)
	return err
}

// CreateDomain создает новый домен.
func (cli *Cli) CreateDomain(domain string, shared bool, storage string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain is required")
	}
	args := make([]any, 0, 7)
	args = append(args, "CREATEDOMAIN", domain)
	if shared {
		args = append(args, "SHARED")
	}
	if storage != "" {
		args = append(args, "PATH", storage)
	}
	if settings != nil {
		args = append(args, settings)
	}
	_, err := cli.Query(args...)
	return err
}

// CreateDomainStorage создает точку монтирования хранилища.
func (cli *Cli) CreateDomainStorage(storage string, shared bool) error {
	if storage == "" {
		return fmt.Errorf("storage path is required")
	}
	args := make([]any, 0, 4)
	args = append(args, "CREATEDOMAINSTORAGE")
	if shared {
		args = append(args, "SHARED")
	}
	args = append(args, "PATH", storage)
	_, err := cli.Query(args...)
	return err
}

// DeleteDomain удаляет домен.
func (cli *Cli) DeleteDomain(domain string, force bool) error {
	if domain == "" {
		return fmt.Errorf("domain is required")
	}
	args := make([]any, 0, 3)
	args = append(args, "DELETEDOMAIN", domain)
	if force {
		args = append(args, "FORCE")
	}
	_, err := cli.Query(args...)
	return err
}

// GetClusterAccountDefaults возвращает настройки аккаунтов кластера по умолчанию.
func (cli *Cli) GetClusterAccountDefaults() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERACCOUNTDEFAULTS", "")
}

// GetClusterAccountPrefs возвращает предпочтения аккаунтов кластера по умолчанию.
func (cli *Cli) GetClusterAccountPrefs() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERACCOUNTPREFS", "")
}

// GetClusterDirectoryIntegration возвращает настройки интеграции справочника кластера.
func (cli *Cli) GetClusterDirectoryIntegration() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDIRECTORYINTEGRATION", "")
}

// GetClusterDomainDefaults возвращает настройки доменов кластера по умолчанию.
func (cli *Cli) GetClusterDomainDefaults() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDOMAINDEFAULTS", "")
}

// GetClusterTrustedCerts возвращает список доверенных сертификатов кластера.
func (cli *Cli) GetClusterTrustedCerts() ([]string, error) {
	return cli.getTrustedCerts("GETCLUSTERTRUSTEDCERTS")
}

// GetDirectoryIntegration возвращает настройки интеграции справочника сервера.
func (cli *Cli) GetDirectoryIntegration() (map[string]any, error) {
	return cli.getMapAny("GETDIRECTORYINTEGRATION", "")
}

// GetDomainDefaults возвращает настройки доменов сервера по умолчанию.
func (cli *Cli) GetDomainDefaults() (map[string]any, error) {
	return cli.getMapAny("GETDOMAINDEFAULTS", "")
}

// GetServerAccountDefaults возвращает настройки аккаунтов сервера по умолчанию.
func (cli *Cli) GetServerAccountDefaults() (map[string]any, error) {
	return cli.getMapAny("GETSERVERACCOUNTDEFAULTS", "")
}

// GetServerAccountPrefs возвращает предпочтения аккаунтов сервера по умолчанию.
func (cli *Cli) GetServerAccountPrefs() (map[string]any, error) {
	return cli.getMapAny("GETSERVERACCOUNTPREFS", "")
}

// GetServerTrustedCerts возвращает список доверенных сертификатов сервера.
func (cli *Cli) GetServerTrustedCerts() ([]string, error) {
	return cli.getTrustedCerts("GETSERVERTRUSTEDCERTS")
}

// ListClusterTelnums возвращает телефонные номера кластера.
func (cli *Cli) ListClusterTelnums(f string, l int) (map[string]string, error) {
	return cli.listTelnums("LISTCLUSTERTELNUMS", f, l)
}

// ListDomainStorage возвращает список путей хранилища.
func (cli *Cli) ListDomainStorage(shared bool) ([]string, error) {
	args := make([]any, 0, 2)
	args = append(args, "LISTDOMAINSTORAGE")
	if shared {
		args = append(args, "SHARED")
	}
	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}
	return cli.toStringSlice2(res)
}

// ListDomains возвращает список всех доменов.
func (cli *Cli) ListDomains() ([]string, error) { return cli.getSliceString("LISTDOMAINS", "") }

// ListServerTelnums возвращает телефонные номера сервера.
func (cli *Cli) ListServerTelnums(f string, l int) (map[string]string, error) {
	return cli.listTelnums("LISTSERVERTELNUMS", f, l)
}

// MainDomainName возвращает имя основного домена.
func (cli *Cli) MainDomainName() (string, error) {
	res, err := cli.Query("MAINDOMAINNAME")
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("MAINDOMAINNAME: expected string, got %T", res)
	}
	return s, nil
}

// ReloadDirectoryDomains перезагружает домены из справочника.
func (cli *Cli) ReloadDirectoryDomains() error {
	_, err := cli.Query("RELOADDIRECTORYDOMAINS")
	return err
}

// RenameDomain переименовывает домен.
func (cli *Cli) RenameDomain(oldName, newName string, storage string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new names are required")
	}
	args := make([]any, 0, 6)
	args = append(args, "RENAMEDOMAIN", oldName, "INTO", newName)
	if storage != "" {
		args = append(args, "PATH", storage)
	}
	_, err := cli.Query(args...)
	return err
}

// SetClusterAccountDefaults устанавливает настройки аккаунтов кластера.
func (cli *Cli) SetClusterAccountDefaults(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERACCOUNTDEFAULTS", "", s)
}

// SetClusterAccountPrefs устанавливает предпочтения аккаунтов кластера.
func (cli *Cli) SetClusterAccountPrefs(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERACCOUNTPREFS", "", s)
}

// SetClusterDirectoryIntegration устанавливает настройки интеграции справочника кластера.
func (cli *Cli) SetClusterDirectoryIntegration(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERDIRECTORYINTEGRATION", "", s)
}

// SetClusterDomainDefaults устанавливает настройки доменов кластера.
func (cli *Cli) SetClusterDomainDefaults(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERDOMAINDEFAULTS", "", s)
}

// SetClusterTrustedCerts устанавливает доверенные сертификаты кластера.
func (cli *Cli) SetClusterTrustedCerts(certs []string) error {
	return cli.setTrustedCerts("SETCLUSTERTRUSTEDCERTS", certs)
}

// SetDirectoryIntegration устанавливает настройки интеграции справочника сервера.
func (cli *Cli) SetDirectoryIntegration(s map[string]any) error {
	return cli.setMapAny("SETDIRECTORYINTEGRATION", "", s)
}

// SetDomainDefaults устанавливает настройки доменов сервера.
func (cli *Cli) SetDomainDefaults(s map[string]any) error {
	return cli.setMapAny("SETDOMAINDEFAULTS", "", s)
}

// SetServerAccountDefaults устанавливает настройки аккаунтов сервера.
func (cli *Cli) SetServerAccountDefaults(s map[string]any) error {
	return cli.setMapAny("SETSERVERACCOUNTDEFAULTS", "", s)
}

// SetServerAccountPrefs устанавливает предпочтения аккаунтов сервера.
func (cli *Cli) SetServerAccountPrefs(s map[string]any) error {
	return cli.setMapAny("SETSERVERACCOUNTPREFS", "", s)
}

// SetServerTrustedCerts устанавливает доверенные сертификаты сервера.
func (cli *Cli) SetServerTrustedCerts(certs []string) error {
	return cli.setTrustedCerts("SETSERVERTRUSTEDCERTS", certs)
}

// UpdateClusterAccountDefaults обновляет настройки аккаунтов кластера.
func (cli *Cli) UpdateClusterAccountDefaults(s map[string]any) error {
	return cli.setMapAny("UPDATECLUSTERACCOUNTDEFAULTS", "", s)
}

// UpdateClusterAccountPrefs обновляет предпочтения аккаунтов кластера.
func (cli *Cli) UpdateClusterAccountPrefs(s map[string]any) error {
	return cli.setMapAny("UPDATECLUSTERACCOUNTPREFS", "", s)
}

// UpdateClusterDomainDefaults обновляет настройки доменов кластера.
func (cli *Cli) UpdateClusterDomainDefaults(s map[string]any) error {
	return cli.setMapAny("UPDATECLUSTERDOMAINDEFAULTS", "", s)
}

// UpdateDomainDefaults обновляет настройки доменов сервера.
func (cli *Cli) UpdateDomainDefaults(s map[string]any) error {
	return cli.setMapAny("UPDATEDOMAINDEFAULTS", "", s)
}

// UpdateServerAccountDefaults обновляет настройки аккаунтов сервера.
func (cli *Cli) UpdateServerAccountDefaults(s map[string]any) error {
	return cli.setMapAny("UPDATESERVERACCOUNTDEFAULTS", "", s)
}

// UpdateServerAccountPrefs обновляет предпочтения аккаунтов сервера.
func (cli *Cli) UpdateServerAccountPrefs(s map[string]any) error {
	return cli.setMapAny("UPDATESERVERACCOUNTPREFS", "", s)
}

// --- Private Methods ---

func (cli *Cli) getTrustedCerts(cmd string) ([]string, error) {
	res, err := cli.getMapAny(cmd, "")
	if err != nil {
		return nil, err
	}
	val, ok := res["TrustedCertificates"]
	if !ok || val == nil {
		return []string{}, nil
	}
	raw, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected []any for TrustedCertificates, got %T", cmd, val)
	}
	certs := make([]string, 0, len(raw))
	for i, item := range raw {
		b, ok := item.([]byte)
		if !ok {
			return nil, fmt.Errorf("%s: expected []byte at index %d, got %T", cmd, i, item)
		}
		certs = append(certs, string(b))
	}
	return certs, nil
}

func (cli *Cli) listTelnums(cmd, filter string, limit int) (map[string]string, error) {
	args := make([]any, 0, 4)
	args = append(args, cmd)
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if limit > 0 {
		args = append(args, limit)
	}
	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	delete(dict, "")
	return cli.toStringMap(cmd, dict)
}

func (cli *Cli) setTrustedCerts(cmd string, certs []string) error {
	if certs == nil {
		return fmt.Errorf("%s: certs slice is nil", cmd)
	}
	replacer := strings.NewReplacer(
		"-----BEGIN CERTIFICATE-----", "",
		"-----END CERTIFICATE-----", "",
		"\n", "", "\r", "", " ", "",
	)
	prepared := make([]any, 0, len(certs))
	for _, c := range certs {
		clean := replacer.Replace(c)
		if clean != "" {
			prepared = append(prepared, DataBlock(clean))
		}
	}
	_, err := cli.Query(cmd, map[string]any{"TrustedCertificates": prepared})
	return err
}
