package cgpcli

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// 1. Тестирование жизненного цикла аккаунта (создание, атрибуты, удаление)
func TestAccountLifecycle(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("life-%d@test.domain.name", time.Now().Unix())

	t.Run("CreateWithType", func(t *testing.T) {
		settings := map[string]any{"RealName": "Test Lifecycle"}
		// want (string, map[string]any, AccountType, string, bool)
		err := cli.CreateAccount(testAcc, settings, MultiMailbox, "", false)
		if err != nil {
			t.Fatalf("CreateAccount failed: %v", err)
		}
	})

	t.Run("EffectiveSettings", func(t *testing.T) {
		res, err := cli.GetAccountEffectiveSettings(testAcc)
		if err != nil {
			t.Errorf("GetAccountEffectiveSettings failed: %v", err)
		}
		if len(res) == 0 {
			t.Error("Expected some effective settings, got empty map")
		}
	})

	t.Run("OneSetting", func(t *testing.T) {
		val, err := cli.GetAccountOneSetting(testAcc, "ServiceLevel")
		if err != nil {
			t.Errorf("GetAccountOneSetting failed: %v", err)
		}
		t.Logf("ServiceLevel: %v", val)
	})

	t.Run("OneInfo", func(t *testing.T) {
		val, err := cli.GetAccountOneInfo(testAcc, "Created")
		if err != nil {
			t.Errorf("GetAccountOneInfo failed: %v", err)
		}
		t.Logf("Created: %v", val)
	})

	t.Run("Cleanup", func(t *testing.T) {
		if err := cli.DeleteAccount(testAcc); err != nil {
			t.Errorf("DeleteAccount failed: %v", err)
		}
	})
}

// 2. Тестирование алиасов и тел. номеров
func TestAccountCollections(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("coll-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("Aliases", func(t *testing.T) {
		list := []string{"alias1", "alias2"}
		if err := cli.SetAccountAliases(testAcc, list); err != nil {
			t.Fatalf("SetAccountAliases failed: %v", err)
		}
		got, err := cli.GetAccountAliases(testAcc)
		if err != nil || len(got) != 2 {
			t.Errorf("GetAccountAliases mismatch: %v (err: %v)", got, err)
		}
	})

	t.Run("TelnumsAtomic", func(t *testing.T) {
		// Атомарное добавление
		_, err := cli.ModifyAccountTelnums(testAcc, map[string]any{"add": true, "telnum": "701"})
		if err != nil {
			// Если сервер вернул "not implemented", это ожидаемое поведение для некоторых конфигураций
			// Мы не считаем это ошибкой теста, а выводим как предупреждение.
			if strings.Contains(err.Error(), "not implemented") {
				t.Log("ModifyAccountTelnums is not implemented on this server/account")
				return
			}
			t.Fatalf("ModifyAccountTelnums failed: %v", err)
		}

		// Если команда прошла успешно, проверяем результат
		got, err := cli.GetAccountTelnums(testAcc)
		if err != nil {
			t.Fatalf("GetAccountTelnums failed: %v", err)
		}

		found := false
		for _, n := range got {
			if n == "701" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Telnum modification succeeded but number 701 not found in: %v", got)
		}
	})
}

// 3. Тестирование паролей (Verify и PlainPassword)
func TestAccountPasswords(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("pw-%d@test.domain.name", time.Now().Unix())
	testPass := "Secret123!"
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("SetAndVerify", func(t *testing.T) {
		// want (account, password, method, tag, check)
		if err := cli.SetAccountPassword(testAcc, testPass, "", "", true); err != nil {
			t.Fatalf("SetAccountPassword failed: %v", err)
		}
		if _, err := cli.VerifyAccountPassword(testAcc, testPass); err != nil {
			t.Errorf("VerifyAccountPassword failed: %v", err)
		}
	})

	t.Run("CheckPlain", func(t *testing.T) {
		// Работает только если в домене разрешено хранение паролей в Plain/Reversible
		match, _ := cli.CheckAccountPassword(testAcc, testPass)
		t.Logf("CheckAccountPassword result: %v", match)
	})
}

// 4. Тестирование правил (Mail Rules)
func TestAccountMailRules(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("rules-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	rule := MailRule{
		Priority: 7,
		Name:     "TestRule",
		Conditions: []any{
			[]any{"Subject", "is", "Test"},
		},
		Actions: []any{
			[]any{"Discard"},
		},
	}

	t.Run("UpdateAndGet", func(t *testing.T) {
		if err := cli.UpdateAccountMailRule(testAcc, rule); err != nil {
			t.Fatalf("UpdateAccountMailRule failed: %v", err)
		}
		rules, err := cli.GetAccountMailRules(testAcc)
		if err != nil || len(rules) == 0 {
			t.Fatalf("GetAccountMailRules failed: %v", err)
		}
		if rules[0].Name != "TestRule" {
			t.Errorf("Rule name mismatch: %s", rules[0].Name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := cli.DeleteAccountMailRule(testAcc, "TestRule"); err != nil {
			t.Errorf("DeleteAccountMailRule failed: %v", err)
		}
	})
}

// 5. Тестирование статуса, локации и сессий
func TestAccountStatus(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	user := "postmaster"

	t.Run("Location", func(t *testing.T) {
		loc, err := cli.GetAccountLocation(user)
		if err != nil {
			t.Errorf("GetAccountLocation failed: %v", err)
		}
		t.Logf("Account location: %s", loc)
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := cli.IsAccountExists(user)
		if err != nil || !exists {
			t.Errorf("IsAccountExists failed: %v", err)
		}
	})

	t.Run("KillSessions", func(t *testing.T) {
		err := cli.KillAccountSessions(user)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				t.Log("KillAccountSessions: own session likely killed")
				return
			}
			t.Logf("KillAccountSessions info: %v", err)
		}
	})
}

// 6. Тестирование переименования
func TestAccountRename(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	oldName := fmt.Sprintf("old-%d@test.domain.name", time.Now().Unix())
	newName := fmt.Sprintf("new-%d@test.domain.name", time.Now().Unix())

	_ = cli.CreateAccount(oldName, nil, TypeDefault, "", false)

	t.Run("RenameProcess", func(t *testing.T) {
		if err := cli.RenameAccount(oldName, newName, ""); err != nil {
			t.Fatalf("RenameAccount failed: %v", err)
		}

		existsOld, _ := cli.IsAccountExists(oldName)
		existsNew, _ := cli.IsAccountExists(newName)

		if existsOld || !existsNew {
			t.Errorf("Rename logic failed: OldExists=%v, NewExists=%v", existsOld, existsNew)
		}
		_ = cli.DeleteAccount(newName)
	})
}

// 7. Тестирование ACL
func TestAccountACL(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("acl-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("SetAndGetACL", func(t *testing.T) {
		rights := map[string]any{"anyone": "i"}
		if err := cli.SetAccountACL(testAcc, rights, ""); err != nil {
			t.Fatalf("SetAccountACL failed: %v", err)
		}

		acl, err := cli.GetAccountACL(testAcc, "")
		if err != nil {
			t.Fatalf("GetAccountACL failed: %v", err)
		}
		if r, ok := acl["anyone"]; !ok || r != "i" {
			t.Errorf("Unexpected ACL value: %v", acl["anyone"])
		}
	})
}

func TestAccountAccessModesInheritance(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("modes-inh-%d@test.domain.name", time.Now().Unix())
	// Создаем чистый аккаунт
	err := cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
	defer cli.DeleteAccount(testAcc)

	t.Run("1. Set Specific Modes", func(t *testing.T) {
		target := []string{"Mail", "IMAP"}
		err := cli.SetAccountAccessModes(testAcc, target)
		if err != nil {
			t.Fatalf("SetAccountAccessModes failed: %v", err)
		}

		got, err := cli.GetAccountAccessModes(testAcc)
		if err != nil {
			t.Fatalf("GetAccountAccessModes failed: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("Expected 2 modes, got %d: %v", len(got), got)
		}
	})

	t.Run("2. Set None (Explicit Empty)", func(t *testing.T) {
		// Передаем инициализированный, но пустой слайс
		err := cli.SetAccountAccessModes(testAcc, []string{})
		if err != nil {
			t.Fatalf("SetAccountAccessModes(None) failed: %v", err)
		}

		// В GetAccountSettings мы должны увидеть "None"
		settings, _ := cli.GetAccountSettings(testAcc)
		if settings["AccessModes"] != "None" {
			t.Errorf("Expected 'None' in settings, got: %v", settings["AccessModes"])
		}

		got, _ := cli.GetAccountAccessModes(testAcc)
		if len(got) != 0 {
			t.Error("Expected empty slice for None")
		}
	})

	t.Run("3. Reset to Default (Inheritance)", func(t *testing.T) {
		// Передаем nil — это должно сгенерировать #NULL#
		err := cli.SetAccountAccessModes(testAcc, nil)
		if err != nil {
			t.Fatalf("SetAccountAccessModes(nil) failed: %v", err)
		}

		// Проверяем через GetAccountSettings: ключа AccessModes вообще не должно быть в словаре
		settings, err := cli.GetAccountSettings(testAcc)
		if err != nil {
			t.Fatalf("GetAccountSettings failed: %v", err)
		}

		if _, exists := settings["AccessModes"]; exists {
			t.Error("AccessModes key still exists in settings after #NULL# reset (inheritance failed)")
		}

		// GetAccountAccessModes должен вернуть nil (согласно нашей логике прозрачности)
		got, err := cli.GetAccountAccessModes(testAcc)
		if err != nil {
			t.Errorf("GetAccountAccessModes failed: %v", err)
		}
		if got != nil {
			t.Error("Expected nil result for inherited settings")
		}
	})
}

func TestUpdateAccountInfo(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("info-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("Update and Verify", func(t *testing.T) {
		// Записываем кастомные данные
		update := map[string]any{
			"CustomField": "TestData",
			"Note":        "Updated via Go",
		}
		if err := cli.UpdateAccountInfo(testAcc, update); err != nil {
			t.Fatalf("UpdateAccountInfo failed: %v", err)
		}

		// Проверяем через GetAccountInfo
		info, err := cli.GetAccountInfo(testAcc, []string{"CustomField", "Note"})
		if err != nil {
			t.Fatalf("GetAccountInfo failed: %v", err)
		}

		if info["CustomField"] != "TestData" {
			t.Errorf("Expected TestData, got %v", info["CustomField"])
		}
	})

	t.Run("Delete Field via NULL", func(t *testing.T) {
		// Удаляем поле CustomField
		update := map[string]any{
			"CustomField": nil, // Превратится в #NULL#
		}
		if err := cli.UpdateAccountInfo(testAcc, update); err != nil {
			t.Fatalf("UpdateAccountInfo (delete) failed: %v", err)
		}

		// Проверяем, что поля больше нет в словаре
		info, _ := cli.GetAccountInfo(testAcc, []string{"CustomField"})
		if _, exists := info["CustomField"]; exists {
			t.Error("CustomField should have been deleted via #NULL#")
		}
	})
}

func TestAccountSignalRulesLogic(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := "*"

	// Тестовые сценарии: проверяем все типы упаковки
	testCases := []SignalRule{
		{
			Name:       "#ForkOnTimer",
			Priority:   1,
			Delay:      10, // 10 сек -> 99001
			Conditions: []any{[]any{"Method", "is", "INVITE"}},
			Actions:    []any{[]any{"Discard", ""}},
		},
		{
			Name:       "LongDelayRule",
			Priority:   50,
			Delay:      600, // 10 мин -> 40050
			Conditions: []any{[]any{"From", "is", "external@domain.tld"}},
			Actions:    []any{[]any{"Fork to", "manager@domain.tld"}},
		},
		{
			Name:       "EventBusyRule",
			Priority:   5,
			EventCode:  EventBusy, // 3 -> 305
			Conditions: []any{[]any{"From", "is", "8*@*"}},
			Actions:    []any{[]any{"Reject with", "486 Busy Here"}},
		},
		{
			Name:       "#Block",
			Priority:   99,
			Delay:      30, // Должно быть проигнорировано (не в белом списке)
			Conditions: []any{[]any{"From", "in", "#Blocked"}},
			Actions:    []any{[]any{"Reject with", "680 No registered devices"}},
		},
	}

	t.Run("WriteAndVerifyRules", func(t *testing.T) {
		// 1. Записываем набор правил
		if err := cli.SetAccountSignalRules(testAcc, testCases); err != nil {
			t.Fatalf("SetAccountSignalRules failed: %v", err)
		}

		// 2. Читаем актуальное состояние с сервера
		readRules, err := cli.GetAccountSignalRules(testAcc)
		if err != nil {
			t.Fatalf("GetAccountSignalRules failed: %v", err)
		}

		// 3. Сверяем каждое правило по имени
		for _, expected := range testCases {
			found := false
			for _, read := range readRules {
				if read.Name == expected.Name {
					found = true

					// Проверка приоритета (базовое поле)
					if read.Priority != expected.Priority {
						t.Errorf("[%s] Priority mismatch: want %d, got %d", expected.Name, expected.Priority, read.Priority)
					}

					// Отработка специфичной логики
					switch expected.Name {
					case "#ForkOnTimer", "LongDelayRule":
						if read.Delay != expected.Delay {
							t.Errorf("[%s] Delay mismatch: want %d, got %d", expected.Name, expected.Delay, read.Delay)
						}
					case "EventBusyRule":
						if read.EventCode != expected.EventCode {
							t.Errorf("[%s] EventCode mismatch: want %d, got %d", expected.Name, expected.EventCode, read.EventCode)
						}
					case "#Block":
						// Проверка защиты: задержка для этого правила должна быть 0
						if read.Delay != 0 {
							t.Errorf("[%s] Security breach: delay should be 0, got %d", expected.Name, read.Delay)
						}
					}
					break
				}
			}
			if !found {
				t.Errorf("Rule %s not found in server response", expected.Name)
			}
		}
	})
}

// 8. Тестирование парсинга сложных системных данных (имитация ответа сервера)
func TestAccountInfoComplexParsing(t *testing.T) {
	// Мы не делаем запрос к серверу, а тестируем Unmarshal напрямую на реальном дампе
	// Это гарантирует, что проблема в парсере будет поймана
	complexData := `{
		Created = "31-12-2025_07:14:13";
		LastLogin = #T11-01-2026_07:00:48;
		MailboxSeq = #1767981;
		Subscription = (INBOX, Sent);
		MailboxDeleted = {
			1767894 = GoSDKTestRenamed;
			1767897 = GoSDKTestRenamed;
		};
		StorageUsed = 0;
	}`

	var res any
	err := Unmarshal([]byte(complexData), &res)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	info, ok := res.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", res)
	}

	// Если парсер сломался на #T или на (, словарь будет неполным или пустым
	t.Run("CheckFieldsExistence", func(t *testing.T) {
		expectedFields := []string{"Created", "LastLogin", "MailboxSeq", "Subscription", "MailboxDeleted"}
		for _, f := range expectedFields {
			if _, exists := info[f]; !exists {
				t.Errorf("Field %s is missing in parsed map (parser probably stopped early)", f)
			}
		}
	})

	t.Run("CheckTypes", func(t *testing.T) {
		// Проверяем, что Subscription распарсился как слайс, а не как строка "("
		if sub, ok := info["Subscription"].([]any); !ok {
			t.Errorf("Subscription should be []any, got %T", info["Subscription"])
		} else if len(sub) != 2 {
			t.Errorf("Subscription length mismatch: want 2, got %d", len(sub))
		}

		// Проверяем вложенный словарь
		if del, ok := info["MailboxDeleted"].(map[string]any); !ok {
			t.Errorf("MailboxDeleted should be map, got %T", info["MailboxDeleted"])
		} else if len(del) != 2 {
			t.Errorf("MailboxDeleted length mismatch: want 2, got %d", len(del))
		}
	})
}

func TestAccountExtraSettings(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("extra-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	// 1. Тестируем RIMAP / RPOP / RSIP (они идентичны по структуре)
	t.Run("RemoteSyncs", func(t *testing.T) {
		data := map[string]any{
			"Record1": map[string]any{"Server": "imap.test.com", "User": "test"},
		}

		// RIMAP
		if err := cli.SetAccountRIMAPs(testAcc, data); err != nil {
			t.Errorf("SetAccountRIMAPs failed: %v", err)
		}
		_, err := cli.GetAccountRIMAPs(testAcc)
		if err != nil {
			t.Errorf("GetAccountRIMAPs failed: %v", err)
		}

		// RPOP
		if err := cli.SetAccountRPOPs(testAcc, data); err != nil {
			t.Errorf("SetAccountRPOPs failed: %v", err)
		}
		_, err = cli.GetAccountRPOPs(testAcc)
		if err != nil {
			t.Errorf("GetAccountRPOPs failed: %v", err)
		}

		// RSIP
		if err := cli.SetAccountRSIPs(testAcc, data); err != nil {
			t.Errorf("SetAccountRSIPs failed: %v", err)
		}
		_, err = cli.GetAccountRSIPs(testAcc)
		if err != nil {
			t.Errorf("GetAccountRSIPs failed: %v", err)
		}
	})

	// 2. Тестируем специфичные геттеры
	t.Run("SpecificGetters", func(t *testing.T) {
		// GetAccountACLRights
		_, err := cli.GetAccountACLRights(testAcc, "postmaster")
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("GetAccountACLRights failed: %v", err)
		}

		// GetAccountEffectivePrefs
		_, err = cli.GetAccountEffectivePrefs(testAcc)
		if err != nil {
			t.Errorf("GetAccountEffectivePrefs failed: %v", err)
		}

		// GetAccountPresence
		_, err = cli.GetAccountPresence(testAcc)
		if err != nil {
			t.Errorf("GetAccountPresence failed: %v", err)
		}

		// GetAccountRights (Системные права)
		_, err = cli.GetAccountRights(testAcc)
		if err != nil {
			t.Errorf("GetAccountRights failed: %v", err)
		}
	})

	// 3. Тестируем сигнальные правила (удаление)
	t.Run("DeleteSignalRule", func(t *testing.T) {
		// Сначала создадим правило, чтобы было что удалять
		rule := SignalRule{Name: "DeleteMe", Priority: 1}
		_ = cli.UpdateAccountSignalRule(testAcc, rule)

		if err := cli.DeleteAccountSignalRule(testAcc, "DeleteMe"); err != nil {
			t.Errorf("DeleteAccountSignalRule failed: %v", err)
		}
	})
}

func TestAccountFinalStub(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testAcc := fmt.Sprintf("final-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("PrefsAndSettings", func(t *testing.T) {
		prefs := map[string]any{"Language": "English"}
		// SetAccountPrefs
		if err := cli.SetAccountPrefs(testAcc, prefs); err != nil {
			t.Errorf("SetAccountPrefs failed: %v", err)
		}
		// UpdateAccountPrefs
		if err := cli.UpdateAccountPrefs(testAcc, map[string]any{"Timezone": "UTC"}); err != nil {
			t.Errorf("UpdateAccountPrefs failed: %v", err)
		}
		// GetAccountPrefs
		if _, err := cli.GetAccountPrefs(testAcc); err != nil {
			t.Errorf("GetAccountPrefs failed: %v", err)
		}
		// SetAccountSettings
		if err := cli.SetAccountSettings(testAcc, map[string]any{"ServiceLevel": 10}); err != nil {
			t.Errorf("SetAccountSettings failed: %v", err)
		}
		// UpdateAccountSettings
		if err := cli.UpdateAccountSettings(testAcc, map[string]any{"CanChangePassword": true}); err != nil {
			t.Errorf("UpdateAccountSettings failed: %v", err)
		}
	})

	t.Run("RulesAndTypes", func(t *testing.T) {
		// SetAccountMailRules (массовая установка)
		rules := []MailRule{{Name: "BulkRule", Priority: 5}}
		if err := cli.SetAccountMailRules(testAcc, rules); err != nil {
			t.Errorf("SetAccountMailRules failed: %v", err)
		}
		// SetAccountType
		if err := cli.SetAccountType(testAcc, MultiMailbox); err != nil {
			t.Errorf("SetAccountType failed: %v", err)
		}
	})

	t.Run("SecurityAndIdentity", func(t *testing.T) {
		// Сначала установим пароль, чтобы GetAccountPlainPassword имел шанс сработать
		tempPass := "TestPass123"
		_ = cli.SetAccountPassword(testAcc, tempPass, "", "", false)
		// GetAccountPlainPassword
		_, err := cli.GetAccountPlainPassword(testAcc, "")
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("GetAccountPlainPassword failed: %v", err)
		}
		// VerifyAccountIdentity (проверяем себя же)
		ok, err := cli.VerifyAccountIdentity(testAcc, testAcc)
		if err != nil || !ok {
			t.Errorf("VerifyAccountIdentity failed: ok=%v, err=%v", ok, err)
		}
	})

	t.Run("TelnumsAndTasks", func(t *testing.T) {
		// SetAccountTelnums
		if err := cli.SetAccountTelnums(testAcc, []string{"101", "102"}); err != nil {
			t.Errorf("SetAccountTelnums failed: %v", err)
		}

		// UpdateScheduledTask
		// Просто вызываем метод для покрытия кода.
		// Т.к. CGP может возвращать OK даже для несуществующих ID в задачах,
		// мы просто фиксируем факт отсутствия паники/ошибок маршалинга.
		task := map[string]any{"id": "nonexistent", "Enabled": false}
		err := cli.UpdateScheduledTask(testAcc, task)
		if err != nil {
			t.Logf("UpdateScheduledTask returned info: %v", err)
		}
	})
}
