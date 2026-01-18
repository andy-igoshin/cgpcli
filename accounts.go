// # Section: Account Administration
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

var accountTypeNames = [...]string{
	"", "MultiMailbox", "TextMailbox", "MailDirMailbox",
	"SlicedMailbox", "AGrade", "BGrade", "CGrade", "DGrade",
}

func (t AccountType) String() string {
	if int(t) < 0 || int(t) >= len(accountTypeNames) {
		return ""
	}
	return accountTypeNames[t]
}

func (cli *Cli) CheckAccountPassword(account, password string) (bool, error) {
	if password == "" {
		return false, fmt.Errorf("password is required")
	}
	plain, err := cli.GetAccountPlainPassword(account, "")
	if err != nil {
		return false, err
	}
	return plain == password, nil
}

func (cli *Cli) CreateAccount(account string, settings map[string]any, accType AccountType, path string, legacy bool) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	args := make([]any, 0, 8)
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

func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("DELETEACCOUNT", account)
	return err
}

func (cli *Cli) DeleteAccountMailRule(account, ruleName string) error {
	if account == "" || ruleName == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTMAILRULE", account, "DELETE", ruleName)
	return err
}

func (cli *Cli) DeleteAccountSignalRule(account, ruleName string) error {
	if account == "" || ruleName == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	_, err := cli.Query("UPDATEACCOUNTSIGNALRULE", account, "DELETE", ruleName)
	return err
}

func (cli *Cli) GetAccountACL(account, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	args := make([]any, 0, 4)
	args = append(args, "GETACCOUNTACL", Atom(account))
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTACL: expected map, got %T", res)
	}
	return m, nil
}

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

func (cli *Cli) GetAccountAccessModes(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	settings, err := cli.GetAccountSettings(account)
	if err != nil {
		return nil, err
	}
	raw, ok := settings["AccessModes"]
	if !ok || raw == nil {
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
		return toStringSlice("GetAccountAccessModes", val[1:])
	default:
		return []string{}, nil
	}
}

func (cli *Cli) GetAccountAliases(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTALIASES", account)
}

func (cli *Cli) GetAccountEffectivePrefs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTEFFECTIVEPREFS", account)
}

func (cli *Cli) GetAccountEffectiveSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTEFFECTIVESETTINGS", account)
}

func (cli *Cli) GetAccountInfo(account string, keys []string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	var res any
	var err error
	if len(keys) == 0 {
		res, err = cli.Query("GETACCOUNTINFO", Atom(account))
	} else {
		res, err = cli.Query("GETACCOUNTINFO", Atom(account), keys)
	}
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTINFO: expected map, got %T", res)
	}
	return m, nil
}

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

func (cli *Cli) GetAccountMailRules(account string) ([]MailRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMailRules("GETACCOUNTMAILRULES", account)
}

func (cli *Cli) GetAccountOneInfo(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTINFO", Atom(account), "Key", key)
}

func (cli *Cli) GetAccountOneSetting(account, key string) (any, error) {
	if account == "" || key == "" {
		return nil, fmt.Errorf("account name and key name are required")
	}
	return cli.Query("GETACCOUNTONESETTING", Atom(account), key)
}

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
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETACCOUNTPLAINPASSWORD: expected string, got %T", res)
	}
	return s, nil
}

func (cli *Cli) GetAccountPrefs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTPREFS", account)
}

func (cli *Cli) GetAccountPresence(account string) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.Query("GETACCOUNTPRESENCE", account)
}

func (cli *Cli) GetAccountRIMAPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRIMAPS", account)
}

func (cli *Cli) GetAccountRPOPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRPOPS", account)
}

func (cli *Cli) GetAccountRSIPs(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTRSIPS", account)
}

func (cli *Cli) GetAccountRights(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTRIGHTS", account)
}

func (cli *Cli) GetAccountSettings(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("GETACCOUNTSETTINGS", account)
}

func (cli *Cli) GetAccountSignalRules(account string) ([]SignalRule, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSignalRules("GETACCOUNTSIGNALRULES", account)
}

func (cli *Cli) GetAccountTelnums(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETACCOUNTTELNUMS", account)
}

func (cli *Cli) IsAccountExists(account string) (bool, error) {
	if account == "" {
		return false, fmt.Errorf("account name is required")
	}
	location, err := cli.GetAccountLocation(account)
	if err != nil {
		return false, err
	}
	return location != "", nil
}

func (cli *Cli) KillAccountSessions(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("KILLACCOUNTSESSIONS", account)
	return err
}

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

func (cli *Cli) SetAccountAccessModes(account string, modes []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	var value any
	switch {
	case modes == nil:
		value = nil
	case len(modes) == 0:
		value = "None"
	default:
		raw := make([]any, 0, len(modes)+1)
		raw = append(raw, 30)
		for _, m := range modes {
			raw = append(raw, m)
		}
		value = raw
	}
	return cli.UpdateAccountSettings(account, map[string]any{"AccessModes": value})
}

func (cli *Cli) SetAccountAliases(account string, aliases []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETACCOUNTALIASES", account, aliases)
	return err
}

func (cli *Cli) SetAccountMailRules(account string, rules []MailRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return setRulesGeneric(cli, "SETACCOUNTMAILRULES", account, rules)
}

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

func (cli *Cli) SetAccountPrefs(account string, prefs map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("SETACCOUNTPREFS", account, prefs)
}

func (cli *Cli) SetAccountRIMAPs(account string, records map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("SETACCOUNTRIMAPS", account, records)
}

func (cli *Cli) SetAccountRPOPs(account string, records map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("SETACCOUNTRPOPS", account, records)
}

func (cli *Cli) SetAccountRSIPs(account string, records map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("SETACCOUNTRSIPS", account, records)
}

func (cli *Cli) SetAccountSettings(account string, settings map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("SETACCOUNTSETTINGS", account, settings)
}

func (cli *Cli) SetAccountSignalRules(account string, rules []SignalRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return setRulesGeneric(cli, "SETACCOUNTSIGNALRULES", account, rules)
}

func (cli *Cli) SetAccountTelnums(account string, telnums []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETACCOUNTTELNUMS", account, telnums)
	return err
}

func (cli *Cli) SetAccountType(account string, accType AccountType) error {
	if account == "" || accType == TypeDefault {
		return fmt.Errorf("account name and account type are required")
	}
	_, err := cli.Query("SETACCOUNTTYPE", account, accType.String())
	return err
}

func (cli *Cli) UpdateAccountInfo(account string, info map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTINFO", account, info)
}

func (cli *Cli) UpdateAccountMailRule(account string, rule MailRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return updateRuleGeneric(cli, "UPDATEACCOUNTMAILRULE", account, rule)
}

func (cli *Cli) UpdateAccountPrefs(account string, prefs map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTPREFS", account, prefs)
}

func (cli *Cli) UpdateAccountSettings(account string, settings map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.setMapAny("UPDATEACCOUNTSETTINGS", account, settings)
}

func (cli *Cli) UpdateAccountSignalRule(account string, rule SignalRule) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return updateRuleGeneric(cli, "UPDATEACCOUNTSIGNALRULE", account, rule)
}

func (cli *Cli) UpdateScheduledTask(account string, taskData map[string]any) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	if _, ok := taskData["id"]; !ok {
		return fmt.Errorf("task id is required in taskData")
	}
	return cli.setMapAny("UPDATESCHEDULEDTASK", account, taskData)
}

func (cli *Cli) VerifyAccountIdentity(account, identity string) (bool, error) {
	if account == "" || identity == "" {
		return false, fmt.Errorf("account name and identity are required")
	}
	_, err := cli.Query("VERIFYACCOUNTIDENTITY", account, "FOR", identity)
	return err == nil, err
}

func (cli *Cli) VerifyAccountPassword(account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("account name and password are required")
	}
	_, err := cli.Query("VERIFYACCOUNTPASSWORD", account, "PASSWORD", password)
	return err
}
