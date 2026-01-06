package cgpcli

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// getTestCli — вспомогательная функция для получения клиента из переменных окружения
func getTestCli(t *testing.T) *Cli {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Skip("CGPUSER or CGPPASS not set, skipping integration test")
	}
	// Используем APOP для безопасности в тестах
	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	return cli
}

// 1. Тестирование жизненного цикла аккаунта (создание, атрибуты, удаление)
func TestAccountLifecycle(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

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

	t.Run("Cleanup", func(t *testing.T) {
		if err := cli.DeleteAccount(testAcc); err != nil {
			t.Errorf("DeleteAccount failed: %v", err)
		}
	})
}

// 2. Тестирование алиасов и тел. номеров
func TestAccountCollections(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

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
	defer cli.Logout()

	testAcc := fmt.Sprintf("pw-%d@test.domain.name", time.Now().Unix())
	testPass := "Secret123!"
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("SetAndVerify", func(t *testing.T) {
		// want (account, password, method, tag, check)
		if err := cli.SetAccountPassword(testAcc, testPass, "", "", true); err != nil {
			t.Fatalf("SetAccountPassword failed: %v", err)
		}
		if err := cli.VerifyAccountPassword(testAcc, testPass); err != nil {
			t.Errorf("VerifyAccountPassword failed: %v", err)
		}
	})

	t.Run("CheckPlain", func(t *testing.T) {
		// Работает только если в домене разрешено хранение паролей в Plain/Reversible
		match := cli.CheckAccountPassword(testAcc, testPass)
		t.Logf("CheckAccountPassword result: %v", match)
	})
}

// 4. Тестирование правил (Mail Rules)
func TestAccountMailRules(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	testAcc := fmt.Sprintf("rules-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	rule := AccountRule{
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
	defer cli.Logout()
	user := os.Getenv("CGPUSER")

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
	defer cli.Logout()

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
	defer cli.Logout()

	testAcc := fmt.Sprintf("acl-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil, TypeDefault, "", false)
	defer cli.DeleteAccount(testAcc)

	t.Run("SetAndGetACL", func(t *testing.T) {
		rights := map[string]any{"anyone": "i"}
		if err := cli.SetAccountACL(testAcc, "", rights); err != nil {
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
	defer cli.Logout()

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
	defer cli.Logout()

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
