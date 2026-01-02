package cgpcli

import (
	"fmt"
	"strings"
)

// CreateAccount создает новый аккаунт с указанными настройками.
func (cli *Cli) CreateAccount(account string, settings map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"CREATEACCOUNT", account}
	if settings != nil {
		parts = append(parts, Marshal(settings))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// DeleteAccount удаляет указанный аккаунт.
func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.query("DELETEACCOUNT " + account)
	return err
}

// GetAccountAliases возвращает список алиасов (псевдонимов) указанного аккаунта.
func (cli *Cli) GetAccountAliases(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETACCOUNTALIASES " + account)
	if err != nil {
		return nil, err
	}
	list, ok := res.([]any)
	if !ok {
		if res == "" {
			return []string{}, nil
		}
		return nil, fmt.Errorf("GETACCOUNTALIASES: expected array, got %T", res)
	}
	aliases := make([]string, len(list))
	for i, v := range list {
		aliases[i] = fmt.Sprint(v)
	}
	return aliases, nil
}

// GetAccountEffectiveSettings возвращает итоговые настройки (явные + унаследованные).
// Принимает имя аккаунта или "*" для текущего пользователя.
func (cli *Cli) GetAccountEffectiveSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETACCOUNTEFFECTIVESETTINGS " + account)
	if err != nil {
		return nil, err
	}
	settings, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTEFFECTIVESETTINGS: expected map, got %T", res)
	}
	return settings, nil
}

// GetAccountInfo возвращает срез данных из словаря "info" (LastLogin, и т.д.).
// Если keys пуст, запрашиваются все доступные данные.
func (cli *Cli) GetAccountInfo(account string, keys []string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	parts := []string{"GETACCOUNTINFO", Marshal(account)}
	if len(keys) > 0 {
		anyKeys := make([]any, len(keys))
		for i, k := range keys {
			anyKeys[i] = k
		}
		parts = append(parts, Marshal(anyKeys))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	settings, ok := res.(map[string]any)
	if !ok {
		if res == "" {
			return make(map[string]any), nil
		}
		return nil, fmt.Errorf("GETACCOUNTINFO: expected map, got %T", res)
	}
	return settings, nil
}

// GetAccountLocation возвращает путь к директории аккаунта или его INBOX.
func (cli *Cli) GetAccountLocation(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETACCOUNTLOCATION " + account)
	if err != nil {
		return "", err
	}
	location, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETACCOUNTLOCATION: expected string, got %T", res)
	}
	return location, nil
}

// GetAccountOneSetting возвращает значение одной конкретной эффективной настройки.
func (cli *Cli) GetAccountOneSetting(account string, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	parts := []string{"GETACCOUNTONESETTING", account, key}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetAccountPlainPassword возвращает пароль аккаунта в открытом виде (недокументировано).
// Параметр method может быть "SIP" или "RADIUS".
func (cli *Cli) GetAccountPlainPassword(account string, method string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	parts := []string{"GetAccountPlainPassword", account}
	if method == "SIP" || method == "RADIUS" {
		parts = append(parts, "METHOD", method)
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return "", err
	}
	password, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GetAccountPlainPassword: expected string, got %T", res)
	}
	return password, nil
}

// GetAccountRights возвращает список прав администратора сервера для аккаунта.
func (cli *Cli) GetAccountRights(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETACCOUNTRIGHTS " + account)
	if err != nil {
		return nil, err
	}
	list, ok := res.([]any)
	if !ok {
		if res == "" {
			return []string{}, nil
		}
		return nil, fmt.Errorf("GETACCOUNTRIGHTS: expected array, got %T", res)
	}
	rights := make([]string, len(list))
	for i, v := range list {
		rights[i] = fmt.Sprint(v)
	}
	return rights, nil
}

// GetAccountSettings возвращает только явно установленные настройки аккаунта.
func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETACCOUNTSETTINGS " + account)
	if err != nil {
		return nil, err
	}
	settings, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTSETTINGS: expected map, got %T", res)
	}
	return settings, nil
}

// GetAccountTelnums возвращает список телефонных номеров, назначенных аккаунту.
func (cli *Cli) GetAccountTelnums(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETACCOUNTTELNUMS " + account)
	if err != nil {
		return nil, err
	}
	list, ok := res.([]any)
	if !ok {
		if res == "" {
			return []string{}, nil
		}
		return nil, fmt.Errorf("GETACCOUNTTELNUMS: expected array, got %T", res)
	}
	nums := make([]string, len(list))
	for i, v := range list {
		nums[i] = fmt.Sprint(v)
	}
	return nums, nil
}

// KillAccountSessions прерывает все активные сессии аккаунта.
func (cli *Cli) KillAccountSessions(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.query("KILLACCOUNTSESSIONS " + account)
	return err
}

// RenameAccount переименовывает аккаунт или переносит его в другой домен.
func (cli *Cli) RenameAccount(oldName, newName string, storagePath string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new account names are required")
	}
	parts := []string{"RENAMEACCOUNT", oldName, "into", newName}
	if storagePath != "" {
		parts = append(parts, "PATH", storagePath)
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetAccountAliases устанавливает список алиасов для аккаунта (перезаписывает).
func (cli *Cli) SetAccountAliases(account string, aliases []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	anyAliases := make([]any, len(aliases))
	for i, a := range aliases {
		anyAliases[i] = a
	}
	parts := []string{"SETACCOUNTALIASES", account, Marshal(anyAliases)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetAccountPassword устанавливает новый пароль. Поддерживает METHOD (SIP/RADIUS) и NAME (tag).
func (cli *Cli) SetAccountPassword(account string, password string, check bool, method string, tag string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"SETACCOUNTPASSWORD", account, "PASSWORD", Marshal(password)}
	if method != "" {
		parts = append(parts, "METHOD", method)
	}
	if tag != "" {
		parts = append(parts, "NAME", Marshal(tag))
	}
	if check {
		parts = append(parts, "CHECK")
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetAccountSettings полностью заменяет словарь настроек аккаунта.
func (cli *Cli) SetAccountSettings(account string, settings map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"SETACCOUNTSETTINGS", account, Marshal(settings)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetAccountTelnums назначает список телефонных номеров (перезаписывает).
func (cli *Cli) SetAccountTelnums(account string, telnums []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	anyNums := make([]any, len(telnums))
	for i, n := range telnums {
		anyNums[i] = n
	}
	parts := []string{"SETACCOUNTTELNUMS", account, Marshal(anyNums)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// UpdateAccountInfo изменяет данные в информационном словаре аккаунта (недокументировано).
func (cli *Cli) UpdateAccountInfo(account string, info map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"UpdateAccountInfo", account, Marshal(info)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// UpdateAccountSettings обновляет настройки аккаунта (точечные изменения).
func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"UPDATEACCOUNTSETTINGS", account, Marshal(settings)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// VerifyAccountPassword проверяет корректность пароля аккаунта.
func (cli *Cli) VerifyAccountPassword(account string, password string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"VERIFYACCOUNTPASSWORD", account, "PASSWORD", Marshal(password)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// -----------------------------------------------------------------------------
// Helpers / Extended API
// -----------------------------------------------------------------------------

// CheckAccountPassword проверяет пароль, сравнивая его с паролем в открытом виде.
// Возвращает true, если пароли совпадают.
// Внимание: работает только если сервер хранит пароли в Plain/Reversible режиме.
func (cli *Cli) CheckAccountPassword(account string, password string) bool {
	plain, err := cli.GetAccountPlainPassword(account, "")
	if err != nil {
		return false
	}
	return plain == password
}

// IsAccountExists проверяет существование аккаунта.
func (cli *Cli) IsAccountExists(account string) (bool, error) {
	loc, err := cli.GetAccountLocation(account)
	if err != nil {
		return false, err
	}
	return loc != "", nil
}
