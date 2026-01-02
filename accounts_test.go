package cgpcli

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// Вспомогательная функция для получения клиента
func getTestCli(t *testing.T) *Cli {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	return cli
}

// 1. Тестирование расширенных атрибутов и прав
func TestAccountAttributes(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()
	user := os.Getenv("CGPUSER")

	t.Run("EffectiveSettings", func(t *testing.T) {
		res, err := cli.GetAccountEffectiveSettings(user)
		if err != nil {
			t.Errorf("GetAccountEffectiveSettings failed: %v", err)
		}
		if len(res) == 0 {
			t.Error("Expected some effective settings, got empty map")
		}
	})

	t.Run("AccountRights", func(t *testing.T) {
		rights, err := cli.GetAccountRights(user)
		if err != nil {
			t.Errorf("GetAccountRights failed: %v", err)
		}
		t.Logf("Admin rights for %s: %v", user, rights)
	})

	t.Run("OneSetting", func(t *testing.T) {
		val, err := cli.GetAccountOneSetting(user, "ServiceLevel")
		if err != nil {
			t.Errorf("GetAccountOneSetting failed: %v", err)
		}
		t.Logf("ServiceLevel: %v", val)
	})
}

// 2. Тестирование алиасов и тел. номеров (через временный аккаунт)
func TestAccountCollections(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	testAcc := fmt.Sprintf("test-coll-%d@test.domain.name", time.Now().Unix())
	_ = cli.CreateAccount(testAcc, nil)
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

	t.Run("Telnums", func(t *testing.T) {
		nums := []string{"701", "702"}
		if err := cli.SetAccountTelnums(testAcc, nums); err != nil {
			t.Fatalf("SetAccountTelnums failed: %v", err)
		}
		got, err := cli.GetAccountTelnums(testAcc)
		if err != nil || len(got) != 2 {
			t.Errorf("GetAccountTelnums mismatch: %v", got)
		}
	})
}

// 3. Тестирование паролей (Verify и PlainPassword)
func TestAccountPasswords(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	testAcc := fmt.Sprintf("test-pw-%d@test.domain.name", time.Now().Unix())
	testPass := "Secret123!"
	_ = cli.CreateAccount(testAcc, nil)
	defer cli.DeleteAccount(testAcc)

	t.Run("SetAndVerify", func(t *testing.T) {
		if err := cli.SetAccountPassword(testAcc, testPass, true, "", ""); err != nil {
			t.Fatalf("SetAccountPassword failed: %v", err)
		}
		if err := cli.VerifyAccountPassword(testAcc, testPass); err != nil {
			t.Errorf("VerifyAccountPassword failed: %v", err)
		}
	})

	t.Run("CheckPlain", func(t *testing.T) {
		// Работает только если в домене/классе разрешено хранение паролей в Plain/Reversible
		match := cli.CheckAccountPassword(testAcc, testPass)
		t.Logf("CheckAccountPassword (Plain) result: %v", match)
	})
}

// 4. Тестирование управления сессиями и локации
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

		missing, _ := cli.IsAccountExists("non-existent-user@invalid")
		if missing {
			t.Error("IsAccountExists should be false for non-existent user")
		}
	})

	t.Run("KillSessions", func(t *testing.T) {
		err := cli.KillAccountSessions(user)
		if err != nil {
			// Если произошел EOF, вероятно, мы убили свою собственную сессию
			if strings.Contains(err.Error(), "EOF") {
				t.Log("KillAccountSessions: connection closed by server (likely killed own session)")
				// После EOF нужно переподключиться, если планируются дальнейшие тесты,
				// либо просто завершить этот тест.
				return
			}
			t.Logf("KillAccountSessions info: %v (this is normal if no sessions exist)", err)
		}
	})
}

// 5. Тестирование переименования
func TestAccountRename(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	oldName := fmt.Sprintf("old-%d@test.domain.name", time.Now().Unix())
	newName := fmt.Sprintf("new-%d@test.domain.name", time.Now().Unix())

	_ = cli.CreateAccount(oldName, nil)

	t.Run("RenameProcess", func(t *testing.T) {
		if err := cli.RenameAccount(oldName, newName, ""); err != nil {
			t.Fatalf("RenameAccount failed: %v", err)
		}

		// Проверяем, что старого нет, а новый есть
		existsOld, _ := cli.IsAccountExists(oldName)
		existsNew, _ := cli.IsAccountExists(newName)

		if existsOld || !existsNew {
			t.Errorf("Rename logic failed: OldExists=%v, NewExists=%v", existsOld, existsNew)
		}

		cli.DeleteAccount(newName)
	})
}
