package cgpcli

import (
	"fmt"
	"strconv"
	"strings"
)

// CreateDirectoryDomain создает новый домен на базе справочника (Directory-based).
func (cli *Cli) CreateDirectoryDomain(domainName string, settings map[string]any) error {
	if domainName == "" {
		return fmt.Errorf("domainName is required")
	}
	parts := []string{"CREATEDIRECTORYDOMAIN", domainName}
	if settings != nil {
		parts = append(parts, Marshal(settings))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

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

// CreateDomainStorage создает точку монтирования хранилища для новых доменов.
func (cli *Cli) CreateDomainStorage(storage string, shared bool) error {
	if storage == "" {
		return fmt.Errorf("storage path is required")
	}
	parts := []string{"CREATEDOMAINSTORAGE"}
	if shared {
		parts = append(parts, "SHARED")
	}
	parts = append(parts, "PATH", storage)
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

// GetClusterAccountDefaults возвращает настройки аккаунтов по умолчанию для всего кластера.
func (cli *Cli) GetClusterAccountDefaults() (map[string]any, error) {
	res, err := cli.query("GETCLUSTERACCOUNTDEFAULTS")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETCLUSTERACCOUNTDEFAULTS: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetClusterAccountPrefs возвращает предпочтения аккаунтов по умолчанию для всего кластера.
func (cli *Cli) GetClusterAccountPrefs() (map[string]any, error) {
	res, err := cli.query("GETCLUSTERACCOUNTPREFS")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETCLUSTERACCOUNTPREFS: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetClusterDirectoryIntegration возвращает общекластерные настройки интеграции со справочником.
func (cli *Cli) GetClusterDirectoryIntegration() (map[string]any, error) {
	res, err := cli.query("GETCLUSTERDIRECTORYINTEGRATION")
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETCLUSTERDIRECTORYINTEGRATION: unexpected response type: %T", res)
	}

	return dict, nil
}

// GetClusterDomainDefaults возвращает настройки доменов по умолчанию для всего кластера.
func (cli *Cli) GetClusterDomainDefaults() (map[string]any, error) {
	res, err := cli.query("GETCLUSTERDOMAINDEFAULTS")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETCLUSTERDOMAINDEFAULTS: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetClusterTrustedCerts возвращает список доверенных сертификатов кластера в формате строк.
func (cli *Cli) GetClusterTrustedCerts() ([]string, error) {
	res, err := cli.query("GETCLUSTERTRUSTEDCERTS")
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETCLUSTERTRUSTEDCERTS: expected dictionary, got %T", res)
	}

	rawList, ok := dict["TrustedCertificates"].([]any)
	if !ok {
		return []string{}, nil
	}

	certs := make([]string, 0, len(rawList))
	for _, v := range rawList {
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprint(v)
		}
		certs = append(certs, strings.Trim(s, "[]"))
	}
	return certs, nil
}

// GetDirectoryIntegration возвращает общесерверные настройки интеграции со справочником.
func (cli *Cli) GetDirectoryIntegration() (map[string]any, error) {
	res, err := cli.query("GETDIRECTORYINTEGRATION")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETDIRECTORYINTEGRATION: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetDomainDefaults возвращает общесерверные настройки доменов по умолчанию.
func (cli *Cli) GetDomainDefaults() (map[string]any, error) {
	res, err := cli.query("GETDOMAINDEFAULTS")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETDOMAINDEFAULTS: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetServerAccountDefaults возвращает настройки аккаунтов по умолчанию для сервера.
func (cli *Cli) GetServerAccountDefaults() (map[string]any, error) {
	res, err := cli.query("GETSERVERACCOUNTDEFAULTS")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETSERVERACCOUNTDEFAULTS: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetServerAccountPrefs возвращает предпочтения аккаунтов по умолчанию для сервера.
func (cli *Cli) GetServerAccountPrefs() (map[string]any, error) {
	res, err := cli.query("GETSERVERACCOUNTPREFS")
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETSERVERACCOUNTPREFS: unexpected response type: %T", res)
	}
	return dict, nil
}

// GetServerTrustedCerts возвращает список доверенных сертификатов сервера в формате строк (PEM/Base64).
func (cli *Cli) GetServerTrustedCerts() ([]string, error) {
	res, err := cli.query("GETSERVERTRUSTEDCERTS")
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETSERVERTRUSTEDCERTS: expected dictionary, got %T", res)
	}

	rawList, ok := dict["TrustedCertificates"].([]any)
	if !ok {
		return []string{}, nil
	}

	certs := make([]string, 0, len(rawList))
	for _, v := range rawList {
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprint(v)
		}
		certs = append(certs, strings.Trim(s, "[]"))
	}
	return certs, nil
}

// ListDomains возвращает список имен всех доменов на сервере.
func (cli *Cli) ListDomains() ([]string, error) {
	res, err := cli.query("LISTDOMAINS")
	if err != nil {
		return nil, err
	}
	rawList, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	domains := make([]string, 0, len(rawList))
	for _, v := range rawList {
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprint(v)
		}
		domains = append(domains, s)
	}
	return domains, nil
}

// ListDomainStorage возвращает список точек монтирования хранилища.
func (cli *Cli) ListDomainStorage(shared bool) ([]string, error) {
	parts := []string{"LISTDOMAINSTORAGE"}
	if shared {
		parts = append(parts, "SHARED")
	}

	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	rawList, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("LISTDOMAINSTORAGE: unexpected response type: %T", res)
	}

	result := make([]string, 0, len(rawList))
	for _, v := range rawList {
		s, ok := v.(string)
		if !ok {
			s = fmt.Sprint(v)
		}
		result = append(result, s)
	}

	return result, nil
}

// ListServerTelnums возвращает список телефонных номеров в некластерных доменах.
func (cli *Cli) ListServerTelnums(filter string, limit int) (map[string]string, error) {
	parts := []string{"LISTSERVERTELNUMS"}
	if filter != "" {
		parts = append(parts, "FILTER", filter)
	}
	parts = append(parts, strconv.Itoa(limit))

	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	rawDict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("LISTSERVERTELNUMS: unexpected response type: %T", res)
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

// ListClusterTelNums возвращает список телефонных номеров в общих доменах кластера.
func (cli *Cli) ListClusterTelNums(filter string, limit int) (map[string]string, error) {
	parts := []string{"LISTCLUSTERTELNUMS"}
	if filter != "" {
		parts = append(parts, "FILTER", filter)
	}
	parts = append(parts, strconv.Itoa(limit))

	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}

	rawDict, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("LISTCLUSTERTELNUMS: unexpected response type: %T", res)
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

// MainDomainName возвращает имя главного (основного) домена сервера.
func (cli *Cli) MainDomainName() (string, error) {
	res, err := cli.query("MAINDOMAINNAME")
	if err != nil {
		return "", err
	}

	name, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("MAINDOMAINNAME: unexpected response type: %T", res)
	}

	return name, nil
}

// ReloadDirectoryDomains инициирует сканирование дерева справочника.
func (cli *Cli) ReloadDirectoryDomains() error {
	_, err := cli.query("RELOADDIRECTORYDOMAINS")
	return err
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

// SetClusterAccountDefaults полностью заменяет настройки аккаунтов по умолчанию для кластера.
func (cli *Cli) SetClusterAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETCLUSTERACCOUNTDEFAULTS " + Marshal(settings))
	return err
}

// SetClusterAccountPrefs полностью заменяет предпочтения аккаунтов по умолчанию для кластера.
func (cli *Cli) SetClusterAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETCLUSTERACCOUNTPREFS " + Marshal(settings))
	return err
}

// SetClusterDirectoryIntegration заменяет общекластерные настройки интеграции со справочником.
// Эти команды доступны только в Динамическом Кластере.
func (cli *Cli) SetClusterDirectoryIntegration(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETCLUSTERDIRECTORYINTEGRATION " + Marshal(settings))
	return err
}

// SetClusterDomainDefaults полностью заменяет настройки доменов по умолчанию для кластера.
func (cli *Cli) SetClusterDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETCLUSTERDOMAINDEFAULTS " + Marshal(settings))
	return err
}

// SetServerTrustedCerts заменяет список доверенных сертификатов сервера.
func (cli *Cli) SetClusterTrustedCerts(certs []string) error {
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
		prepared = append(prepared, DataBlock(clean))
	}

	data := map[string]any{
		"TrustedCertificates": prepared,
	}

	_, err := cli.query("SETCLUSTERTRUSTEDCERTS " + Marshal(data))
	return err
}

// SetDirectoryIntegration заменяет общесерверные настройки интеграции со справочником.
func (cli *Cli) SetDirectoryIntegration(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETDIRECTORYINTEGRATION " + Marshal(settings))
	return err
}

// SetDomainDefaults полностью заменяет общесерверные настройки доменов по умолчанию.
func (cli *Cli) SetDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETDOMAINDEFAULTS " + Marshal(settings))
	return err
}

// SetServerAccountDefaults заменяет настройки аккаунтов по умолчанию для сервера.
func (cli *Cli) SetServerAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETSERVERACCOUNTDEFAULTS " + Marshal(settings))
	return err
}

// SetServerAccountPrefs заменяет предпочтения аккаунтов по умолчанию для сервера.
func (cli *Cli) SetServerAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("SETSERVERACCOUNTPREFS " + Marshal(settings))
	return err
}

// SetServerTrustedCerts заменяет список доверенных сертификатов сервера.
func (cli *Cli) SetServerTrustedCerts(certs []string) error {
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
		prepared = append(prepared, DataBlock(clean))
	}

	data := map[string]any{
		"TrustedCertificates": prepared,
	}

	_, err := cli.query("SETSERVERTRUSTEDCERTS " + Marshal(data))
	return err
}

// UpdateClusterAccountDefaults инкрементально обновляет настройки аккаунтов для кластера.
func (cli *Cli) UpdateClusterAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("UPDATECLUSTERACCOUNTDEFAULTS " + Marshal(settings))
	return err
}

// UpdateClusterAccountPrefs инкрементально обновляет предпочтения аккаунтов для кластера.
func (cli *Cli) UpdateClusterAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("UPDATECLUSTERACCOUNTPREFS " + Marshal(settings))
	return err
}

// UpdateClusterDomainDefaults инкрементально обновляет настройки доменов для кластера.
func (cli *Cli) UpdateClusterDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("UPDATECLUSTERDOMAINDEFAULTS " + Marshal(settings))
	return err
}

// UpdateDomainDefaults инкрементально обновляет общесерверные настройки доменов по умолчанию.
func (cli *Cli) UpdateDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("UPDATEDOMAINDEFAULTS " + Marshal(settings))
	return err
}

// UpdateServerAccountDefaults инкрементально обновляет настройки аккаунтов для сервера.
func (cli *Cli) UpdateServerAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("UPDATESERVERACCOUNTDEFAULTS " + Marshal(settings))
	return err
}

// UpdateServerAccountPrefs инкрементально обновляет предпочтения аккаунтов для сервера.
func (cli *Cli) UpdateServerAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	_, err := cli.query("UPDATESERVERACCOUNTPREFS " + Marshal(settings))
	return err
}
