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
	args := make([]any, 0, 6)
	args = append(args, account)
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
	return cli.QueryNV("CREATEACCOUNT", args...)
}

func (cli *Cli) DeleteAccount(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("DELETEACCOUNT", account)
}

func (cli *Cli) DeleteAccountMailRule(account, rule string) error {
	if account == "" || rule == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	return cli.QueryNV("UPDATEACCOUNTMAILRULE", account, "DELETE", rule)
}

func (cli *Cli) DeleteAccountSignalRule(account, rule string) error {
	if account == "" || rule == "" {
		return fmt.Errorf("account name and rule name are required")
	}
	return cli.QueryNV("UPDATEACCOUNTSIGNALRULE", account, "DELETE", rule)
}

func (cli *Cli) GetAccountACL(account, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTACL"
	var res any
	var err error

	if authAccount != "" {
		res, err = cli.Query(cmd, Atom(account), "AUTH", authAccount)
	} else {
		res, err = cli.Query(cmd, Atom(account))
	}

	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

func (cli *Cli) GetAccountACLRights(account, authAccount string) (string, error) {
	if account == "" || authAccount == "" {
		return "", fmt.Errorf("target account and auth account names are required")
	}

	const cmd = "GETACCOUNTACLRIGHTS"

	res, err := cli.Query(cmd, Atom(account), "AUTH", authAccount)
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
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
	return cli.getMapAny("GETACCOUNTEFFECTIVESETTINGS", Atom(account))
}

func (cli *Cli) GetAccountInfo(account string, keys []string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETACCOUNTINFO"
	var res any
	var err error

	if len(keys) == 0 {
		res, err = cli.Query(cmd, Atom(account))
	} else {
		res, err = cli.Query(cmd, Atom(account), keys)
	}
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

func (cli *Cli) GetAccountLocation(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	return cli.getString("GETACCOUNTLOCATION", account)
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

	const cmd = "GETACCOUNTPLAINPASSWORD"
	var res any
	var err error

	if method == "SIP" || method == "RADIUS" {
		res, err = cli.Query(cmd, account, "METHOD", method)
	} else {
		res, err = cli.Query(cmd, account)
	}

	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
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
	return cli.getMapAny("GETACCOUNTSETTINGS", Atom(account))
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
	return cli.QueryNV("KILLACCOUNTSESSIONS", account)
}

func (cli *Cli) ModifyAccountTelnums(account string, params map[string]any) (map[string]any, error) {
	if account == "" || params == nil {
		return nil, fmt.Errorf("account name and parameters dictionary are required")
	}
	res, err := cli.Query("MODIFYACCOUNTTELNUMS", account, params)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return make(map[string]any), nil
	}
	return m, nil
}

func (cli *Cli) RenameAccount(oldName, newName, path string) (err error) {
	if oldName == "" || newName == "" {
		return fmt.Errorf("both old and new names are required")
	}
	const cmd = "RENAMEACCOUNT"
	if path != "" {
		err = cli.QueryNV(cmd, oldName, "INTO", newName, "PATH", path)
	} else {
		err = cli.QueryNV(cmd, oldName, "INTO", newName)
	}
	return
}

func (cli *Cli) SetAccountACL(account, authAccount string, newACL map[string]any) error {
	if account == "" || newACL == nil {
		return fmt.Errorf("account name and ACL dictionary are required")
	}
	args := make([]any, 0, 4)
	args = append(args, Atom(account))
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	args = append(args, newACL)
	return cli.QueryNV("SETACCOUNTACL", args...)
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
	return cli.QueryNV("SETACCOUNTALIASES", account, aliases)
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
	args = append(args, account, "PASSWORD", password)
	if method != "" {
		args = append(args, "METHOD", method)
	}
	if tag != "" {
		args = append(args, "NAME", tag)
	}
	if check {
		args = append(args, "CHECK")
	}
	return cli.QueryNV("SETACCOUNTPASSWORD", args...)
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
	return cli.QueryNV("SETACCOUNTTELNUMS", account, telnums)
}

func (cli *Cli) SetAccountType(account string, accType AccountType) error {
	if account == "" || accType == TypeDefault {
		return fmt.Errorf("account name and account type are required")
	}
	return cli.QueryNV("SETACCOUNTTYPE", account, accType.String())
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
	err := cli.QueryNV("VERIFYACCOUNTIDENTITY", account, "FOR", identity)
	return err == nil, err
}

func (cli *Cli) VerifyAccountPassword(account, password string) error {
	if account == "" || password == "" {
		return fmt.Errorf("account name and password are required")
	}
	return cli.QueryNV("VERIFYACCOUNTPASSWORD", account, "PASSWORD", password)
}
