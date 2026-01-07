package cgpcli

import (
	"fmt"
	"strings"
)

type AccountType int

const (
	TypeDefault AccountType = iota
	MultiMailbox
	TextMailbox
	MailDirMailbox
	SlicedMailbox
	AGrade
	BGrade
	CGrade
	DGrade
)

func (t AccountType) String() string {
	return [...]string{
		"", "MultiMailbox", "TextMailbox", "MailDirMailbox",
		"SlicedMailbox", "AGrade", "BGrade", "CGrade", "DGrade",
	}[t]
}

// CheckAccountPassword проверяет пароль, сравнивая его с паролем в открытом виде.
func (cli *Cli) CheckAccountPassword(account string, password string) bool {
	plain, err := cli.GetAccountPlainPassword(account, "")
	if err != nil {
		return false
	}
	return plain == password
}

// CreateAccount создает новый аккаунт с указанными настройками.
func (cli *Cli) CreateAccount(account string, settings map[string]any, accType AccountType, path string, legacy bool) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	args := make([]any, 0, 7)
	args = append(args, "CREATEACCOUNT", account)

	if accType != TypeDefault {
		args = append(args, accType.String())
	}
	if path != "" {
		args = append(args, "PATH", path)
	}
	if legacy {
		args = append(args, "LEGACY")
	}
	if settings != nil {
		args = append(args, settings)
	}

	_, err := cli.Query(args...)
	return err
}

// DeleteAccount удаляет указанный аккаунт.
func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("DELETEACCOUNT", account)
	return err
}

// DeleteAccountMailRule удаляет почтовое правило по его имени.
func (cli *Cli) DeleteAccountMailRule(account, ruleName string) error {
	if account == "" || ruleName == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTMAILRULE", account, "DELETE", ruleName)
	return err
}

// DeleteAccountSignalRule удаляет сигнальное правило по его имени.
func (cli *Cli) DeleteAccountSignalRule(account, ruleName string) error {
	if account == "" || ruleName == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTSIGNALRULE", account, "DELETE", ruleName)
	return err
}

// GetAccountACL возвращает словарь прав доступа (ACL) для аккаунта.
func (cli *Cli) GetAccountACL(account, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	if authAccount != "" {
		return cli.getMapAny("GETACCOUNTACL", fmt.Sprintf("%s AUTH %s", account, authAccount))
	}
	return cli.getMapAny("GETACCOUNTACL", account)
}

// GetAccountACLRights возвращает строку эффективных прав доступа.
func (cli *Cli) GetAccountACLRights(account, authAccount string) (string, error) {
	if account == "" || authAccount == "" {
		return "", fmt.Errorf("target account and auth account names are required")
	}
	res, err := cli.Query("GETACCOUNTACLRIGHTS", Atom(account), "AUTH", authAccount)
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETACCOUNTACLRIGHTS: expected string, got %T", res)
	}
	return s, nil
}

// GetAccountAccessModes возвращает список включенных сервисов.
func (cli *Cli) GetAccountAccessModes(account string) ([]string, error) {
	settings, err := cli.GetAccountSettings(account)
	if err != nil {
		return nil, err
	}

	raw, ok := settings["AccessModes"]
	if !ok {
		return nil, nil
	}

	switch val := raw.(type) {
	case string:
		if strings.EqualFold(val, "All") {
			return []string{"All"}, nil
		}
		return []string{}, nil
	case []any:
		if len(val) <= 1 {
			return []string{}, nil
		}
		return cli.toStringSlice(val[1:]), nil
	default:
		return []string{}, nil
	}
}

// GetAccountAliases возвращает список алиасов (псевдонимов) аккаунта.
func (cli *Cli) GetAccountAliases(account string) ([]string, error) {
	return cli.getSliceString("GETACCOUNTALIASES", account)
}

// GetAccountEffectiveSettings возвращает итоговые настройки (явные + унаследованные).
func (cli *Cli) GetAccountEffectiveSettings(account string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTEFFECTIVESETTINGS", account)
}

// GetAccountInfo возвращает словарь указанных элементов информации об аккаунте.
func (cli *Cli) GetAccountInfo(account string, keys []string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.Query("GETACCOUNTINFO", Atom(account), keys)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTINFO: expected map, got %T", res)
	}
	return m, nil
}

// GetAccountLocation возвращает путь к данным аккаунта в файловой системе.
func (cli *Cli) GetAccountLocation(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	res, err := cli.Query("GETACCOUNTLOCATION", account)
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETACCOUNTLOCATION: expected string, got %T", res)
	}
	return s, nil
}

// GetAccountMailRules возвращает список почтовых правил.
func (cli *Cli) GetAccountMailRules(account string) ([]MailRule, error) {
	return cli.getMailRules("GETACCOUNTMAILRULES", account)
}

// GetAccountOneInfo возвращает один конкретный элемент информации.
func (cli *Cli) GetAccountOneInfo(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTINFO", account, "Key", key)
}

// GetAccountOneSetting возвращает значение одной эффективной настройки.
func (cli *Cli) GetAccountOneSetting(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTONESETTING", Atom(account), key)
}

// GetAccountPlainPassword возвращает пароль в открытом виде.
func (cli *Cli) GetAccountPlainPassword(account string, method string) (string, error) {
	args := make([]any, 0, 4)
	args = append(args, "GETACCOUNTPLAINPASSWORD", account)
	if method == "SIP" || method == "RADIUS" {
		args = append(args, "METHOD", method)
	}
	res, err := cli.Query(args...)
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETACCOUNTPLAINPASSWORD: expected string, got %T", res)
	}
	return s, nil
}

// GetAccountPresence возвращает статус присутствия (Presence) аккаунта.
func (cli *Cli) GetAccountPresence(account string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.Query("GETACCOUNTPRESENCE", account)
}

// GetAccountRIMAPs возвращает словарь RIMAP записей.
func (cli *Cli) GetAccountRIMAPs(account string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTRIMAPS", account)
}

// GetAccountRPOPs возвращает словарь RPOP записей.
func (cli *Cli) GetAccountRPOPs(account string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTRPOPS", account)
}

// GetAccountRSIPs возвращает словарь RSIP записей.
func (cli *Cli) GetAccountRSIPs(account string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTRSIPS", account)
}

// GetAccountRights возвращает список административных прав сервера.
func (cli *Cli) GetAccountRights(account string) ([]string, error) {
	return cli.getSliceString("GETACCOUNTRIGHTS", account)
}

// GetAccountSettings возвращает только явно установленные настройки.
func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {
	return cli.getMapAny("GETACCOUNTSETTINGS", account)
}

// GetAccountSignalRules возвращает список сигнальных правил.
func (cli *Cli) GetAccountSignalRules(account string) ([]SignalRule, error) {
	return cli.getSignalRules("GETACCOUNTSIGNALRULES", account)
}

// GetAccountTelnums возвращает список телефонных номеров.
func (cli *Cli) GetAccountTelnums(account string) ([]string, error) {
	return cli.getSliceString("GETACCOUNTTELNUMS", account)
}

// IsAccountExists проверяет существование аккаунта.
func (cli *Cli) IsAccountExists(account string) (bool, error) {
	loc, err := cli.GetAccountLocation(account)
	if err != nil {
		return false, err
	}
	return loc != "", nil
}

// KillAccountSessions обрывает все активные сессии пользователя.
func (cli *Cli) KillAccountSessions(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("KILLACCOUNTSESSIONS", account)
	return err
}

// ModifyAccountTelnums выполняет атомарные операции с телефонами.
func (cli *Cli) ModifyAccountTelnums(account string, params map[string]any) (map[string]any, error) {
	if account == "" || params == nil {
		return nil, fmt.Errorf("account name and parameters dictionary are required")
	}
	res, err := cli.Query("MODIFYACCOUNTTELNUMS", account, params)
	if err != nil {
		return nil, err
	}
	if m, ok := res.(map[string]any); ok {
		return m, nil
	}
	return make(map[string]any), nil
}

// RenameAccount переименовывает или переносит аккаунт.
func (cli *Cli) RenameAccount(oldName, newName, path string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new names are required")
	}
	args := make([]any, 0, 6)
	args = append(args, "RENAMEACCOUNT", oldName, "INTO", newName)
	if path != "" {
		args = append(args, "PATH", path)
	}
	_, err := cli.Query(args...)
	return err
}

// SetAccountACL устанавливает или модифицирует права доступа.
func (cli *Cli) SetAccountACL(account, authAccount string, newACL map[string]any) error {
	if account == "" || newACL == nil {
		return fmt.Errorf("account name and ACL dictionary are required")
	}
	args := make([]any, 0, 5)
	args = append(args, "SETACCOUNTACL", Atom(account))
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	args = append(args, newACL)
	_, err := cli.Query(args...)
	return err
}

// SetAccountAccessModes устанавливает список доступных сервисов.
func (cli *Cli) SetAccountAccessModes(account string, modes []string) error {
	var value any
	if modes == nil {
		value = nil
	} else if len(modes) == 0 {
		value = "None"
	} else {
		raw := make([]any, 0, len(modes)+1)
		raw = append(raw, 30) // Код формата для CGP
		for _, m := range modes {
			raw = append(raw, m)
		}
		value = raw
	}
	return cli.UpdateAccountSettings(account, map[string]any{"AccessModes": value})
}

// SetAccountAliases полностью заменяет список алиасов.
func (cli *Cli) SetAccountAliases(account string, aliases []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETACCOUNTALIASES", account, aliases)
	return err
}

// SetAccountMailRules полностью заменяет список почтовых правил.
func (cli *Cli) SetAccountMailRules(account string, rules []MailRule) error {
	return setRulesGeneric(cli, "SETACCOUNTMAILRULES", account, rules)
}

// SetAccountPassword обновляет пароль аккаунта.
func (cli *Cli) SetAccountPassword(account, password, method, tag string, check bool) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	args := make([]any, 0, 8)
	args = append(args, "SETACCOUNTPASSWORD", account, "PASSWORD", password)
	if method != "" {
		args = append(args, "METHOD", method)
	}
	if tag != "" {
		args = append(args, "NAME", tag)
	}
	if check {
		args = append(args, "CHECK")
	}
	_, err := cli.Query(args...)
	return err
}

// SetAccountRIMAPs полностью заменяет словарь RIMAP записей.
func (cli *Cli) SetAccountRIMAPs(account string, records map[string]any) error {
	return cli.setMapAny("SETACCOUNTRIMAPS", account, records)
}

// SetAccountRPOPs полностью заменяет словарь RPOP записей.
func (cli *Cli) SetAccountRPOPs(account string, records map[string]any) error {
	return cli.setMapAny("SETACCOUNTRPOPS", account, records)
}

// SetAccountRSIPs полностью заменяет словарь RSIP записей.
func (cli *Cli) SetAccountRSIPs(account string, records map[string]any) error {
	return cli.setMapAny("SETACCOUNTRSIPS", account, records)
}

// SetAccountSettings полностью заменяет словарь настроек аккаунта.
func (cli *Cli) SetAccountSettings(account string, settings map[string]any) error {
	return cli.setMapAny("SETACCOUNTSETTINGS", account, settings)
}

// SetAccountSignalRules полностью заменяет список сигнальных правил.
func (cli *Cli) SetAccountSignalRules(account string, rules []SignalRule) error {
	return setRulesGeneric(cli, "SETACCOUNTSIGNALRULES", account, rules)
}

// SetAccountTelnums полностью заменяет список телефонных номеров.
func (cli *Cli) SetAccountTelnums(account string, telnums []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETACCOUNTTELNUMS", account, telnums)
	return err
}

// SetAccountType изменяет тип существующего аккаунта.
func (cli *Cli) SetAccountType(account string, accType AccountType) error {
	if account == "" || accType == TypeDefault {
		return fmt.Errorf("account name and account type are required")
	}
	_, err := cli.Query("SETACCOUNTTYPE", account, accType.String())
	return err
}

// UpdateAccountInfo изменяет данные в информационном словаре.
func (cli *Cli) UpdateAccountInfo(account string, info map[string]any) error {
	return cli.setMapAny("UPDATEACCOUNTINFO", account, info)
}

// UpdateAccountMailRule обновляет или создает почтовое правило.
func (cli *Cli) UpdateAccountMailRule(account string, rule MailRule) error {
	return updateRuleGeneric(cli, "UPDATEACCOUNTMAILRULE", account, rule)
}

// UpdateAccountSettings точечно обновляет настройки аккаунта.
func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {
	return cli.setMapAny("UPDATEACCOUNTSETTINGS", account, settings)
}

// UpdateAccountSignalRule обновляет или создает сигнальное правило.
func (cli *Cli) UpdateAccountSignalRule(account string, rule SignalRule) error {
	return updateRuleGeneric(cli, "UPDATEACCOUNTSIGNALRULE", account, rule)
}

// UpdateScheduledTask создает, обновляет или удаляет задачу.
func (cli *Cli) UpdateScheduledTask(account string, taskData map[string]any) error {
	if _, ok := taskData["id"]; !ok {
		return fmt.Errorf("task id is required in taskData")
	}
	return cli.setMapAny("UPDATESCHEDULEDTASK", account, taskData)
}

// VerifyAccountIdentity проверяет право использовать From-адрес.
func (cli *Cli) VerifyAccountIdentity(account, identity string) (bool, error) {
	if account == "" || identity == "" {
		return false, fmt.Errorf("account name and identity are required")
	}
	_, err := cli.Query("VERIFYACCOUNTIDENTITY", account, "FOR", identity)
	return err == nil, err
}

// VerifyAccountPassword проверяет корректность пароля.
func (cli *Cli) VerifyAccountPassword(account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("account name and password are required")
	}
	_, err := cli.Query("VERIFYACCOUNTPASSWORD", account, "PASSWORD", password)
	return err
}
