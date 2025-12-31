package cgpcli

import (
	"os"
	"testing"
)

func TestListAccounts(t *testing.T) {
	// 1. Получаем учетные данные из окружения
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")

	if cgpuser == "" || cgppass == "" {
		t.Fatal("specify environment variables CGPUSER and CGPPASS")
	}

	// 2. Подключаемся
	cli, err := New("127.0.0.1", cgpuser, cgppass, Plain, false)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer cli.Logout()

	// 3. Тест 1: Общий список (без указания домена)
	t.Run("GlobalList", func(t *testing.T) {
		accounts, err := cli.ListAccounts("")
		if err != nil {
			t.Errorf("ListAccounts(global) failed: %v", err)
		} else {
			t.Logf("Found %d accounts in main domain", len(accounts))
		}
	})

	// 4. Тест 2: Конкретный домен
	t.Run("DomainList", func(t *testing.T) {
		domain := "test.domain.name"
		accounts, err := cli.ListAccounts(domain)
		if err != nil {
			// Используем Fatal, так как если домен существует, список должен быть доступен
			t.Fatalf("ListAccounts(%s) failed: %v", domain, err)
		}

		if len(accounts) == 0 {
			t.Logf("Notice: list of accounts for %s is empty", domain)
		} else {
			t.Logf("Found %d accounts in %s: %v", len(accounts), domain, accounts)
		}

		// Проверка на корректность данных
		for _, acc := range accounts {
			if acc == "" {
				t.Error("Found empty account name in result list")
			}
		}
	})
}

func TestGetAccountSettings(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	cli, _ := New("127.0.0.1", cgpuser, cgppass, Plain, false)
	defer cli.Logout()

	t.Run("GetPostmasterSettings", func(t *testing.T) {
		settings, err := cli.GetAccountSettings(cgpuser)
		if err != nil {
			t.Fatalf("Failed: %v", err)
		}

		// Проверим наличие какого-нибудь стандартного поля, например RealName или StorageQuota
		if _, ok := settings["RealName"]; ok {
			t.Logf("RealName: %v", settings["RealName"])
		}

		t.Logf("Total settings found: %d", len(settings))
	})
}

func TestGetAccountSettingsUser(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")

	cli, err := New("127.0.0.1", cgpuser, cgppass, Plain, false)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer cli.Logout()

	// Тестируем получение настроек обычного пользователя
	t.Run("GetTestUserSettings", func(t *testing.T) {
		targetUser := "testuser1@test.domain.name"
		settings, err := cli.GetAccountSettings(targetUser)
		if err != nil {
			t.Fatalf("Failed to get settings for %s: %v", targetUser, err)
		}

		t.Logf("Successfully fetched %d settings for %s", len(settings), targetUser)

		// Выведем ключевые поля для проверки парсинга
		importantFields := []string{"RealName", "StorageQuota", "ServiceLevels"}
		for _, field := range importantFields {
			if val, ok := settings[field]; ok {
				t.Logf("Field [%s]: %v (type: %T)", field, val, val)
			}
		}
	})
}

func TestUpdateAccountSettings(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	cli, _ := New("127.0.0.1", cgpuser, cgppass, Plain, false)
	defer cli.Logout()

	targetUser := "testuser1@test.domain.name"
	originalName := "Test User1"
	newName := "Go Robot User"

	// 1. Устанавливаем новое имя
	t.Run("SetNewName", func(t *testing.T) {
		settings := map[string]any{
			"RealName": newName,
		}
		err := cli.UpdateAccountSettings(targetUser, settings)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}
	})

	// 2. Проверяем, что оно применилось
	t.Run("VerifyChange", func(t *testing.T) {
		settings, _ := cli.GetAccountSettings(targetUser)
		if settings["RealName"] != newName {
			t.Errorf("Expected %s, got %v", newName, settings["RealName"])
		}
	})

	// 3. Возвращаем как было (Cleanup)
	t.Run("Cleanup", func(t *testing.T) {
		cli.UpdateAccountSettings(targetUser, map[string]any{"RealName": originalName})
	})
}

func TestAccountLifecycle(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	cli, _ := New("127.0.0.1", cgpuser, cgppass, Plain, false)
	defer cli.Logout()

	testAcc := "temporary-user@test.domain.name"

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		err := cli.CreateAccount(testAcc, map[string]any{"RealName": "Temp User"})
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	})

	// 2. Verify List
	t.Run("Verify", func(t *testing.T) {
		list, _ := cli.ListAccounts("test.domain.name")
		found := false
		for _, a := range list {
			if a == "temporary-user" { // CGP в списке часто отдает только имя без домена
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Account not found in list after creation")
		}
	})

	// 3. Delete
	t.Run("Delete", func(t *testing.T) {
		err := cli.DeleteAccount(testAcc)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
	})
}
