package cgpcli

import (
	"fmt"
	"strings"
)

// AccountRule представляет структуру правила (Mail или Signal).
// Согласно CLI, это массив: [priority, name, conditions, actions, comment]
type AccountRule struct {
	Priority   int    `json:"priority"`
	Name       string `json:"name"`
	Conditions []any  `json:"conditions"`
	Actions    []any  `json:"actions"`
	Comment    string `json:"comment,omitempty"`
}

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

// CreateAccount создает новый аккаунт с указанными настройками.
func (cli *Cli) CreateAccount(account string, settings map[string]any, accType AccountType, path string, legacy bool) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	// Аллоцируем с запасом под все возможные опциональные части
	args := make([]any, 0, 7)
	args = append(args, "CREATEACCOUNT", account)

	// [ accountType ]
	if accType != TypeDefault {
		args = append(args, accType.String())
	}

	// [ PATH storage ]
	if path != "" {
		args = append(args, "PATH", path)
	}

	// [ LEGACY ]
	if legacy {
		args = append(args, "LEGACY")
	}

	// [ settings ] - всегда последний
	if settings != nil {
		args = append(args, settings)
	}

	_, err := cli.Query(args...)
	return err
}

// DeleteAccount удаляет указанный аккаунт.
// Синтаксис: DELETEACCOUNT accountName
func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("DELETEACCOUNT", account)
	return err
}

// DeleteAccountMailRule удаляет почтовое правило по его имени.
// Синтаксис: UPDATEACCOUNTMAILRULE accountName DELETE ruleName
func (cli *Cli) DeleteAccountMailRule(account, ruleName string) error {
	if account == "" || ruleName == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTMAILRULE", account, "DELETE", ruleName)
	return err
}

// DeleteAccountSignalRule удаляет сигнальное правило по его имени.
// Синтаксис: UPDATEACCOUNTSIGNALRULE accountName DELETE ruleName
func (cli *Cli) DeleteAccountSignalRule(account, ruleName string) error {
	if account == "" || ruleName == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTSIGNALRULE", account, "DELETE", ruleName)
	return err
}

// GetAccountAccessModes возвращает список явно включенных сервисов.
func (cli *Cli) GetAccountAccessModes(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	settings, err := cli.GetAccountSettings(account)
	if err != nil {
		return nil, err
	}

	raw, ok := settings["AccessModes"]
	if !ok {
		// Ключа нет — значит, всё по дефолту (наследование)
		return nil, nil
	}

	switch val := raw.(type) {
	case string:
		if strings.EqualFold(val, "All") {
			return []string{"All"}, nil
		}
		// "None" или любая другая строка трактуется как пустой список включенных сервисов
		return []string{}, nil

	case []any:
		// Формат: (число, Service1, Service2...)
		if len(val) <= 1 {
			return []string{}, nil
		}
		res := make([]string, 0, len(val)-1)
		for i := 1; i < len(val); i++ {
			res = append(res, fmt.Sprint(val[i]))
		}
		return res, nil

	default:
		return []string{}, nil
	}
}

// GetAccountACL возвращает словарь прав доступа (ACL) для аккаунта.
// Если указан authAccount, возвращаются права только для этого субъекта.
func (cli *Cli) GetAccountACL(account, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	args := make([]any, 0, 4)
	args = append(args, "GETACCOUNTACL", account)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	if m, ok := res.(map[string]any); ok {
		return m, nil
	}
	return make(map[string]any), nil
}

// GetAccountACLRights возвращает строку эффективных прав доступа для конкретного субъекта.
func (cli *Cli) GetAccountACLRights(account, authAccount string) (string, error) {
	if account == "" || authAccount == "" {
		return "", fmt.Errorf("target account and auth account names are required")
	}
	res, err := cli.Query("GETACCOUNTACLRIGHTS", account, "AUTH", authAccount)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(res), nil
}

// GetAccountAliases возвращает список алиасов (псевдонимов) указанного аккаунта.
func (cli *Cli) GetAccountAliases(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	res, err := cli.Query("GETACCOUNTALIASES", account)
	if err != nil {
		return nil, err
	}

	list, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	return cli.toStrictStringSlice(list), nil
}

// GetAccountEffectiveSettings возвращает итоговые настройки (явные + унаследованные).
func (cli *Cli) GetAccountEffectiveSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSettingsMapAlt("GETACCOUNTEFFECTIVESETTINGS", account)
}

// GetAccountInfo возвращает словарь указанных элементов информации об аккаунте.
// Пример: GetAccountInfo("user@dom.ru", []string{"LastLogin", "LastAddress"})
func (cli *Cli) GetAccountInfo(account string, keys []string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.Query("GETACCOUNTINFO", account, keys)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTINFO: expected map, got %T", res)
	}
	return m, nil
}

// GetAccountOneInfo возвращает один конкретный элемент информации.
// Синтаксис: GETACCOUNTINFO account Key keyName
func (cli *Cli) GetAccountOneInfo(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTINFO", account, "Key", key)
}

// GetAccountLocation возвращает путь к данным аккаунта в файловой системе.
// Синтаксис: GETACCOUNTLOCATION accountName
func (cli *Cli) GetAccountLocation(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	res, err := cli.Query("GETACCOUNTLOCATION", account)
	if err != nil {
		return "", err
	}
	location, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", res)
	}
	return location, nil
}

// GetAccountMailRules теперь может возвращать слайс структур вместо "сырого" слайса.
func (cli *Cli) GetAccountMailRules(account string) ([]AccountRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	res, err := cli.Query("GETACCOUNTMAILRULES", account)
	if err != nil {
		return nil, err
	}

	rawList, ok := res.([]any)
	if !ok {
		return []AccountRule{}, nil
	}

	rules := make([]AccountRule, 0, len(rawList))
	for _, item := range rawList {
		if r, ok := item.([]any); ok && len(r) >= 4 {
			rule := AccountRule{
				Priority:   cli.toInt(r[0]),
				Name:       fmt.Sprint(r[1]),
				Conditions: cli.toAnySlice(r[2]),
				Actions:    cli.toAnySlice(r[3]),
			}
			if len(r) > 4 {
				rule.Comment = fmt.Sprint(r[4])
			}
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

// GetAccountOneSetting возвращает значение одной конкретной эффективной настройки.
func (cli *Cli) GetAccountOneSetting(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}

	return cli.Query("GETACCOUNTONESETTING", account, key)
}

// GetAccountPlainPassword возвращает пароль аккаунта в открытом виде (недокументировано).
// Параметр method может быть "SIP" или "RADIUS" для получения альтернативных паролей.
func (cli *Cli) GetAccountPlainPassword(account string, method string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	args := make([]any, 0, 4)
	args = append(args, "GETACCOUNTPLAINPASSWORD", account)
	if method == "SIP" || method == "RADIUS" {
		args = append(args, "METHOD", method)
	}
	res, err := cli.Query(args...)
	if err != nil {
		return "", err
	}
	password, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETACCOUNTPLAINPASSWORD: expected string, got %T", res)
	}
	return password, nil
}

// GetAccountPresence возвращает статус присутствия (Presence) аккаунта.
func (cli *Cli) GetAccountPresence(account string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	return cli.Query("GETACCOUNTPRESENCE", account)
}

// GetAccountRights возвращает список административных прав сервера, назначенных аккаунту.
// Если у аккаунта нет прав администратора, возвращается пустой слайс.
func (cli *Cli) GetAccountRights(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	// Используем Query для автоматического маршалинга аргументов
	res, err := cli.Query("GETACCOUNTRIGHTS", account)
	if err != nil {
		return nil, err
	}

	// Обработка ответа
	list, ok := res.([]any)
	if !ok {
		// Если сервер вернул пустую строку или не массив — прав нет
		return []string{}, nil
	}

	// Используем хелпер для конвертации []any -> []string
	return cli.toStrictStringSlice(list), nil
}

// GetAccountRIMAPs возвращает словарь RIMAP записей (сбор почты через IMAP).
func (cli *Cli) GetAccountRIMAPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSettingsMapAlt("GETACCOUNTRIMAPS", account)
}

// GetAccountRPOPs возвращает словарь RPOP записей аккаунта.
func (cli *Cli) GetAccountRPOPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSettingsMapAlt("GETACCOUNTRPOPS", account)
}

// GetAccountRSIPs возвращает словарь RSIP записей (Remote SIP registration).
func (cli *Cli) GetAccountRSIPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSettingsMapAlt("GETACCOUNTRSIPS", account)
}

// GetAccountSettings возвращает только явно установленные настройки аккаунта.
// Поддерживает "*" для текущего пользователя.
func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSettingsMapAlt("GETACCOUNTSETTINGS", account)
}

// GetAccountTelnums возвращает список телефонных номеров, назначенных аккаунту.
func (cli *Cli) GetAccountTelnums(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.Query("GETACCOUNTTELNUMS", account)
	if err != nil {
		return nil, err
	}
	list, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	return cli.toStrictStringSlice(list), nil
}

// GetAccountSignalRules возвращает список сигнальных правил (Signal Rules) аккаунта.
func (cli *Cli) GetAccountSignalRules(account string) ([]AccountRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	res, err := cli.Query("GETACCOUNTSIGNALRULES", account)
	if err != nil {
		return nil, err
	}

	rawList, ok := res.([]any)
	if !ok {
		return []AccountRule{}, nil
	}

	rules := make([]AccountRule, 0, len(rawList))
	for _, item := range rawList {
		if r, ok := item.([]any); ok && len(r) >= 4 {
			rule := AccountRule{
				Priority:   cli.toInt(r[0]),
				Name:       fmt.Sprint(r[1]),
				Conditions: cli.toAnySlice(r[2]),
				Actions:    cli.toAnySlice(r[3]),
			}
			if len(r) > 4 {
				rule.Comment = fmt.Sprint(r[4])
			}
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

// KillAccountSessions немедленно обрывает все активные сессии пользователя.
func (cli *Cli) KillAccountSessions(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("KILLACCOUNTSESSIONS", account)
	return err
}

// ModifyAccountTelnums выполняет атомарные операции с телефонными номерами.
// Операции: "add", "del", "pop". Для "add" и "del" обязателен номер в ключе "telnum".
func (cli *Cli) ModifyAccountTelnums(account string, params map[string]any) (map[string]any, error) {
	if account == "" || params == nil {
		return nil, fmt.Errorf("account name and parameters dictionary are required")
	}

	res, err := cli.Query("MODIFYACCOUNTTELNUMS", account, params)
	if err != nil {
		return nil, err
	}

	// Для операции "pop" вернется словарь с удаленным номером.
	if m, ok := res.(map[string]any); ok {
		return m, nil
	}
	return make(map[string]any), nil
}

// RenameAccount переименовывает аккаунт или переносит его в другой домен.
// Синтаксис: RENAMEACCOUNT oldAccountName INTO newAccountName [ PATH storage ]
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

func (cli *Cli) SetAccountAccessModes(account string, modes []string) error {
	var value any

	if modes == nil {
		value = nil // MarshalTo отправит AccessModes = #NULL#;
	} else if len(modes) == 0 {
		value = "None" // Явное отключение всех сервисов
	} else {
		raw := make([]any, 0, len(modes)+1)
		raw = append(raw, 30)
		for _, m := range modes {
			raw = append(raw, m)
		}
		value = raw
	}

	return cli.UpdateAccountSettings(account, map[string]any{
		"AccessModes": value,
	})
}

// SetAccountACL устанавливает или модифицирует права доступа.
// В словаре newACL ключи — это идентификаторы (имена аккаунтов),
// а значения — строки прав (например, "i", "+w", "-d").
func (cli *Cli) SetAccountACL(account, authAccount string, newACL map[string]any) error {
	if account == "" || newACL == nil {
		return fmt.Errorf("account name and ACL dictionary are required")
	}
	args := make([]any, 0, 5)
	args = append(args, "SETACCOUNTACL", account)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	args = append(args, newACL)

	_, err := cli.Query(args...)
	return err
}

// SetAccountAliases полностью заменяет список алиасов аккаунта.
// Все старые алиасы удаляются.
func (cli *Cli) SetAccountAliases(account string, aliases []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	// Даже если aliases == nil, маршалер сформирует пустой массив (), что корректно для CGP.
	_, err := cli.Query("SETACCOUNTALIASES", account, aliases)
	return err
}

// SetAccountMailRules полностью заменяет список почтовых правил аккаунта.
func (cli *Cli) SetAccountMailRules(account string, rules []AccountRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	// Превращаем слайс структур в слайс слайсов для маршалера
	rawRules := make([]any, len(rules))
	for i, r := range rules {
		rawRules[i] = r.ToSlice()
	}

	_, err := cli.Query("SETACCOUNTMAILRULES", account, rawRules)
	return err
}

// SetAccountPassword обновляет пароль аккаунта.
// method: "SIP" или "RADIUS" для альтернативных паролей.
// tag: имя для пароля конкретного приложения (Application-Specific Password).
// check: если true, сервер проверит пароль на соответствие политикам безопасности.
func (cli *Cli) SetAccountPassword(account, password, method, tag string, check bool) error {
	if account == "" {
		// Несмотря на то, что пароль может быть пустым (для удаления ASP),
		// имя аккаунта обязательно всегда.
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
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	_, err := cli.Query("SETACCOUNTRIMAPS", account, records)
	return err
}

// SetAccountRPOPs полностью заменяет словарь RPOP записей.
func (cli *Cli) SetAccountRPOPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	_, err := cli.Query("SETACCOUNTRPOPS", account, records)
	return err
}

// SetAccountRSIPs полностью заменяет словарь RSIP записей.
func (cli *Cli) SetAccountRSIPs(account string, records map[string]any) error {
	if account == "" || records == nil {
		return fmt.Errorf("account name and records dictionary are required")
	}
	_, err := cli.Query("SETACCOUNTRSIPS", account, records)
	return err
}

// SetAccountSettings полностью заменяет словарь настроек аккаунта.
// Все настройки, не указанные в новом словаре, удаляются (сбрасываются к шаблону).
func (cli *Cli) SetAccountSettings(account string, settings map[string]any) error {
	if account == "" || settings == nil {
		return fmt.Errorf("account name and settings dictionary are required")
	}

	_, err := cli.Query("SETACCOUNTSETTINGS", account, settings)
	return err
}

// SetAccountSignalRules полностью заменяет список сигнальных правил аккаунта.
func (cli *Cli) SetAccountSignalRules(account string, rules []AccountRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}

	rawRules := make([]any, len(rules))
	for i, r := range rules {
		rawRules[i] = r.ToSlice()
	}

	_, err := cli.Query("SETACCOUNTSIGNALRULES", account, rawRules)
	return err
}

// SetAccountTelnums полностью заменяет список телефонных номеров аккаунта.
func (cli *Cli) SetAccountTelnums(account string, telnums []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETACCOUNTTELNUMS", account, telnums)
	return err
}

// SetAccountType изменяет тип существующего аккаунта.
// Синтаксис: SETACCOUNTTYPE accountName accountType
func (cli *Cli) SetAccountType(account string, accType AccountType) error {
	if account == "" || accType == TypeDefault {
		return fmt.Errorf("account name and account type are required")
	}

	_, err := cli.Query("SETACCOUNTTYPE", account, accType.String())
	return err
}

// UpdateAccountInfo изменяет данные в информационном словаре аккаунта.
// Эта функция недокументирована и полезна для обновления метаданных,
// таких как LastLogin, LastAddress или кастомных полей.
// Синтаксис: UpdateAccountInfo accountName { key = value; ... }
func (cli *Cli) UpdateAccountInfo(account string, info map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	if info == nil {
		return nil
	}
	_, err := cli.Query("UPDATEACCOUNTINFO", account, info)
	return err
}

// UpdateAccountMailRule обновляет или создает почтовое правило.
func (cli *Cli) UpdateAccountMailRule(account string, rule AccountRule) error {
	if account == "" || rule.Name == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTMAILRULE", account, rule.ToSlice())
	return err
}

// UpdateAccountSettings точечно обновляет настройки аккаунта.
// Пропущенные в словаре ключи сохраняют свои текущие значения.
func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {
	if account == "" || settings == nil {
		return fmt.Errorf("account name and settings dictionary are required")
	}

	_, err := cli.Query("UPDATEACCOUNTSETTINGS", account, settings)
	return err
}

// UpdateAccountSignalRule обновляет или создает сигнальное правило.
func (cli *Cli) UpdateAccountSignalRule(account string, rule AccountRule) error {
	if account == "" || rule.Name == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTSIGNALRULE", account, rule.ToSlice())
	return err
}

// UpdateScheduledTask создает, обновляет или удаляет запланированную задачу.
// Чтобы удалить задачу, не указывайте в taskData ключи "program" и "script".
func (cli *Cli) UpdateScheduledTask(account string, taskData map[string]any) error {
	if account == "" || taskData == nil {
		return fmt.Errorf("account name and task data dictionary are required")
	}
	// Обязательный ключ id должен быть в словаре согласно API
	if _, ok := taskData["id"]; !ok {
		return fmt.Errorf("task id is required in taskData")
	}

	_, err := cli.Query("UPDATESCHEDULEDTASK", account, taskData)
	return err
}

// VerifyAccountIdentity проверяет право аккаунта использовать указанный From-адрес.
// Синтаксис: VERIFYACCOUNTIDENTITY accountName FOR identity
func (cli *Cli) VerifyAccountIdentity(account, identity string) (bool, error) {
	if account == "" || identity == "" {
		return false, fmt.Errorf("account name and identity are required")
	}
	_, err := cli.Query("VERIFYACCOUNTIDENTITY", account, "FOR", identity)
	if err != nil {
		// Если сервер вернул ошибку, это может означать "не разрешено"
		// или техническую ошибку. Обычно CGP возвращает ошибку протокола, если identity запрещен.
		return false, err
	}
	return true, nil
}

// VerifyAccountPassword проверяет корректность пароля аккаунта.
// Синтаксис: VERIFYACCOUNTPASSWORD accountName PASSWORD password
func (cli *Cli) VerifyAccountPassword(account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("account name and password are required")
	}
	_, err := cli.Query("VERIFYACCOUNTPASSWORD", account, "PASSWORD", password)
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

func (cli *Cli) toAnySlice(v any) []any {
	if res, ok := v.([]any); ok {
		return res
	}
	return []any{}
}

// ToSlice преобразует структуру в формат, ожидаемый маршалером для CLI (массив).
func (r AccountRule) ToSlice() []any {
	res := []any{r.Priority, r.Name, r.Conditions, r.Actions}
	if r.Comment != "" {
		res = append(res, r.Comment)
	}
	return res
}

func (t AccountType) String() string {
	return [...]string{
		"", "MultiMailbox", "TextMailbox", "MailDirMailbox",
		"SlicedMailbox", "AGrade", "BGrade", "CGrade", "DGrade",
	}[t]
}
