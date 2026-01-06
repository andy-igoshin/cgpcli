package cgpcli

import (
	"fmt"
	"strings"
)

// CreateDirectoryDomain создает новый домен на базе справочника (Directory-based).
func (cli *Cli) CreateDirectoryDomain(domainName string, settings map[string]any) error {
	if domainName == "" {
		return fmt.Errorf("domainName is required")
	}

	args := make([]any, 0, 3)
	args = append(args, "CREATEDIRECTORYDOMAIN", domainName)
	if settings != nil {
		args = append(args, settings)
	}
	_, err := cli.Query(args...)
	return err
}

// CreateDomain создает новый домен.
func (cli *Cli) CreateDomain(domain string, shared bool, storage string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
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

// CreateDomainStorage создает точку монтирования хранилища для новых доменов.
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
		return fmt.Errorf("domain name is required")
	}

	args := make([]any, 0, 3)
	args = append(args, "DELETEDOMAIN", domain)
	if force {
		args = append(args, "FORCE")
	}

	_, err := cli.Query(args...)
	return err
}

// Вспомогательная функция для получения map[string]any из Query
func (cli *Cli) getSettingsMap(command string) (map[string]any, error) {
	res, err := cli.Query(command)
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type: %T", command, res)
	}
	return dict, nil
}

func (cli *Cli) GetClusterAccountDefaults() (map[string]any, error) {
	return cli.getSettingsMap("GETCLUSTERACCOUNTDEFAULTS")
}
func (cli *Cli) GetClusterAccountPrefs() (map[string]any, error) {
	return cli.getSettingsMap("GETCLUSTERACCOUNTPREFS")
}
func (cli *Cli) GetClusterDirectoryIntegration() (map[string]any, error) {
	return cli.getSettingsMap("GETCLUSTERDIRECTORYINTEGRATION")
}
func (cli *Cli) GetClusterDomainDefaults() (map[string]any, error) {
	return cli.getSettingsMap("GETCLUSTERDOMAINDEFAULTS")
}
func (cli *Cli) GetDirectoryIntegration() (map[string]any, error) {
	return cli.getSettingsMap("GETDIRECTORYINTEGRATION")
}
func (cli *Cli) GetDomainDefaults() (map[string]any, error) {
	return cli.getSettingsMap("GETDOMAINDEFAULTS")
}
func (cli *Cli) GetServerAccountDefaults() (map[string]any, error) {
	return cli.getSettingsMap("GETSERVERACCOUNTDEFAULTS")
}
func (cli *Cli) GetServerAccountPrefs() (map[string]any, error) {
	return cli.getSettingsMap("GETSERVERACCOUNTPREFS")
}

// GetClusterTrustedCerts возвращает список доверенных сертификатов кластера.
func (cli *Cli) GetClusterTrustedCerts() ([]string, error) {
	res, err := cli.Query("GETCLUSTERTRUSTEDCERTS")
	if err != nil {
		return nil, err
	}
	return cli.parseTrustedCerts(res)
}

// GetServerTrustedCerts возвращает список доверенных сертификатов сервера.
func (cli *Cli) GetServerTrustedCerts() ([]string, error) {
	res, err := cli.Query("GETSERVERTRUSTEDCERTS")
	if err != nil {
		return nil, err
	}
	return cli.parseTrustedCerts(res)
}

func (cli *Cli) parseTrustedCerts(res any) ([]string, error) {
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected dictionary, got %T", res)
	}

	rawList, ok := dict["TrustedCertificates"].([]any)
	if !ok {
		return []string{}, nil
	}

	certs := make([]string, 0, len(rawList))
	for _, v := range rawList {
		s := fmt.Sprint(v)
		// Очищаем скобки, если парсер оставил их в сыром виде
		certs = append(certs, strings.Trim(s, "[]"))
	}
	return certs, nil
}

// ListDomains возвращает список имен всех доменов на сервере.
func (cli *Cli) ListDomains() ([]string, error) {
	res, err := cli.Query("LISTDOMAINS")
	if err != nil {
		return nil, err
	}
	return cli.toStringSlice(res), nil
}

// ListDomainStorage возвращает список точек монтирования хранилища.
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
	return cli.toStringSlice(res), nil
}

// Вспомогательный метод для LIST...TELNUMS
func (cli *Cli) listTelNums(cmd, filter string, limit int) (map[string]string, error) {
	args := make([]any, 0, 4)
	args = append(args, cmd)
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	args = append(args, limit)

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	rawDict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type: %T", cmd, res)
	}

	result := make(map[string]string, len(rawDict))
	for k, v := range rawDict {
		if k == "" {
			continue
		}
		if s, ok := v.(string); ok {
			result[k] = s
		} else {
			result[k] = fmt.Sprint(v)
		}
	}
	return result, nil
}

func (cli *Cli) ListServerTelnums(f string, l int) (map[string]string, error) {
	return cli.listTelNums("LISTSERVERTELNUMS", f, l)
}
func (cli *Cli) ListClusterTelNums(f string, l int) (map[string]string, error) {
	return cli.listTelNums("LISTCLUSTERTELNUMS", f, l)
}

// MainDomainName возвращает имя основного домена.
func (cli *Cli) MainDomainName() (string, error) {
	res, err := cli.Query("MAINDOMAINNAME")
	if err != nil {
		return "", err
	}
	if s, ok := res.(string); ok {
		return s, nil
	}
	return fmt.Sprint(res), nil
}

func (cli *Cli) ReloadDirectoryDomains() error {
	_, err := cli.Query("RELOADDIRECTORYDOMAINS")
	return err
}

// RenameDomain переименовывает домен.
func (cli *Cli) RenameDomain(oldName, newName string, storage string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new domain names are required")
	}

	args := make([]any, 0, 6)
	args = append(args, "RENAMEDOMAIN", oldName, "INTO", newName)
	if storage != "" {
		args = append(args, "PATH", storage)
	}

	_, err := cli.Query(args...)
	return err
}

// Сеттеры и Обновления (Set/Update)
func (cli *Cli) SetClusterAccountDefaults(s map[string]any) error {
	return cli.sendSettings("SETCLUSTERACCOUNTDEFAULTS", s)
}
func (cli *Cli) SetClusterAccountPrefs(s map[string]any) error {
	return cli.sendSettings("SETCLUSTERACCOUNTPREFS", s)
}
func (cli *Cli) SetClusterDirectoryIntegration(s map[string]any) error {
	return cli.sendSettings("SETCLUSTERDIRECTORYINTEGRATION", s)
}
func (cli *Cli) SetClusterDomainDefaults(s map[string]any) error {
	return cli.sendSettings("SETCLUSTERDOMAINDEFAULTS", s)
}
func (cli *Cli) SetDirectoryIntegration(s map[string]any) error {
	return cli.sendSettings("SETDIRECTORYINTEGRATION", s)
}
func (cli *Cli) SetDomainDefaults(s map[string]any) error {
	return cli.sendSettings("SETDOMAINDEFAULTS", s)
}
func (cli *Cli) SetServerAccountDefaults(s map[string]any) error {
	return cli.sendSettings("SETSERVERACCOUNTDEFAULTS", s)
}
func (cli *Cli) SetServerAccountPrefs(s map[string]any) error {
	return cli.sendSettings("SETSERVERACCOUNTPREFS", s)
}

func (cli *Cli) UpdateClusterAccountDefaults(s map[string]any) error {
	return cli.sendSettings("UPDATECLUSTERACCOUNTDEFAULTS", s)
}
func (cli *Cli) UpdateClusterAccountPrefs(s map[string]any) error {
	return cli.sendSettings("UPDATECLUSTERACCOUNTPREFS", s)
}
func (cli *Cli) UpdateClusterDomainDefaults(s map[string]any) error {
	return cli.sendSettings("UPDATECLUSTERDOMAINDEFAULTS", s)
}
func (cli *Cli) UpdateDomainDefaults(s map[string]any) error {
	return cli.sendSettings("UPDATEDOMAINDEFAULTS", s)
}
func (cli *Cli) UpdateServerAccountDefaults(s map[string]any) error {
	return cli.sendSettings("UPDATESERVERACCOUNTDEFAULTS", s)
}
func (cli *Cli) UpdateServerAccountPrefs(s map[string]any) error {
	return cli.sendSettings("UPDATESERVERACCOUNTPREFS", s)
}

func (cli *Cli) sendSettings(cmd string, s map[string]any) error {
	if s == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.Query(cmd, s)
	return err
}

// Работа с сертификатами
func (cli *Cli) SetClusterTrustedCerts(certs []string) error {
	return cli.setTrustedCerts("SETCLUSTERTRUSTEDCERTS", certs)
}
func (cli *Cli) SetServerTrustedCerts(certs []string) error {
	return cli.setTrustedCerts("SETSERVERTRUSTEDCERTS", certs)
}

func (cli *Cli) setTrustedCerts(cmd string, certs []string) error {
	if certs == nil {
		return fmt.Errorf("certs array is required")
	}

	prepared := make([]any, 0, len(certs))
	replacer := strings.NewReplacer(
		"-----BEGIN CERTIFICATE-----", "",
		"-----END CERTIFICATE-----", "",
		"\n", "", "\r", "", " ", "",
	)

	for _, c := range certs {
		clean := replacer.Replace(c)
		if clean == "" {
			continue
		}
		// Используем наш тип DataBlock, чтобы маршалер выдал [base64]
		prepared = append(prepared, DataBlock(clean))
	}

	_, err := cli.Query(cmd, map[string]any{"TrustedCertificates": prepared})
	return err
}
