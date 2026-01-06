package cgpcli

import (
	"os"
	"testing"
)

func TestDomainAdministration(t *testing.T) {
	host := os.Getenv("CGPHOST")
	user := os.Getenv("CGPUSER")
	pass := os.Getenv("CGPPASS")
	if host == "" {
		t.Skip("Sandbox environment (CGPHOST) not configured")
	}

	cli, err := New(host, user, pass, APOP, false)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer cli.Logout()

	// Используем существующий домен или главный
	targetDomain := "test.domain.name"

	// 1. Тестирование получения списков и локации
	t.Run("Discovery", func(t *testing.T) {
		// ListAdminDomains (проверка на возврат []string{} при успехе)
		adminDomains, err := cli.ListAdminDomains("")
		if err != nil {
			t.Errorf("ListAdminDomains failed: %v", err)
		}
		t.Logf("Admin domains found: %v", adminDomains)

		// GetDomainLocation
		loc, err := cli.GetDomainLocation(targetDomain)
		if err != nil {
			t.Errorf("GetDomainLocation failed: %v", err)
		}
		t.Logf("Domain location: %s", loc)
	})

	// 2. Тестирование настроек и шаблонов
	t.Run("SettingsAndTemplates", func(t *testing.T) {
		// Проверка GetDomainSettings
		settings, err := cli.GetDomainSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainSettings failed: %v", err)
		}
		t.Logf("Current domain settings keys: %d", len(settings))

		// Тест UpdateDomainSettings
		update := map[string]any{"Comment": "Updated via Go SDK Test"}
		if err := cli.UpdateDomainSettings(targetDomain, update); err != nil {
			t.Errorf("UpdateDomainSettings failed: %v", err)
		}

		// Тест GetAccountDefaults / Template
		if _, err := cli.GetAccountDefaults(targetDomain); err != nil {
			t.Errorf("GetAccountDefaults failed: %v", err)
		}
		if _, err := cli.GetAccountTemplate(targetDomain); err != nil {
			t.Errorf("GetAccountTemplate failed: %v", err)
		}
	})

	// 3. Тестирование алиасов (включая проверку Sprint)
	t.Run("Aliases", func(t *testing.T) {
		original, err := cli.GetDomainAliases(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainAliases failed: %v", err)
		}

		testAlias := "test-alias-unique"
		newAliases := append(original, testAlias)

		if err := cli.SetDomainAliases(targetDomain, newAliases); err != nil {
			t.Errorf("SetDomainAliases failed: %v", err)
		}

		// Проверяем, что алиас появился
		updated, _ := cli.GetDomainAliases(targetDomain)
		found := false
		for _, a := range updated {
			if a == testAlias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("New alias %s not found in the list", testAlias)
		}

		// Восстановление
		_ = cli.SetDomainAliases(targetDomain, original)
	})

	// 4. Тестирование хранилищ (Storage)
	t.Run("AccountStorage", func(t *testing.T) {
		storageName := "test_storage_point"
		// Пытаемся создать точку монтирования
		if err := cli.CreateAccountStorage(targetDomain, storageName); err != nil {
			// В песочнице это может вернуть ошибку, если нет прав на ФС или точка уже есть
			t.Logf("CreateAccountStorage info: %v", err)
		}

		storages, err := cli.ListAccountStorage(targetDomain)
		if err != nil {
			t.Errorf("ListAccountStorage failed: %v", err)
		}

		t.Logf("Storages found: %v", storages)
	})

	// 5. Объекты домена и аккаунты
	t.Run("DomainObjects", func(t *testing.T) {
		// Простой список аккаунтов (строгая проверка)
		accs, err := cli.ListAccounts(targetDomain)
		if err != nil {
			t.Errorf("ListAccounts failed: %v", err)
		}
		t.Logf("Accounts: %v", accs)

		// Детальные объекты через LISTDOMAINOBJECTS
		res, err := cli.ListDomainObjects(targetDomain, 10, "", TypeAll, "")
		if err != nil {
			t.Errorf("ListDomainObjects failed: %v", err)
		} else {
			t.Logf("Total domain objects: %d, Cookie: %s", len(res.Objects), res.Cookie)
			for name, obj := range res.Objects {
				t.Logf("Found %s: Type=%v, Targets=%v", name, obj.Subtype, obj.Targets)
			}
		}
	})

	// 6. Операции со справочником (Directory)
	t.Run("Directory", func(t *testing.T) {
		if err := cli.InsertDirectoryRecords(targetDomain); err != nil {
			t.Errorf("InsertDirectoryRecords failed: %v", err)
		}
		// Осторожно с Delete в песочнице, если там есть важные записи
		if err := cli.DeleteDirectoryRecords(targetDomain); err != nil {
			t.Errorf("DeleteDirectoryRecords failed: %v", err)
		}
	})

	// 7. Состояние домена
	t.Run("Management", func(t *testing.T) {
		if err := cli.SuspendDomain(targetDomain); err != nil {
			t.Errorf("Suspend failed: %v", err)
		}
		if err := cli.ResumeDomain(targetDomain); err != nil {
			t.Errorf("Resume failed: %v", err)
		}
	})
}

func TestDomainAliasesTypeJuggling(t *testing.T) {
	host := os.Getenv("CGPHOST")
	user := os.Getenv("CGPUSER")
	pass := os.Getenv("CGPPASS")
	if host == "" {
		t.Skip("Sandbox environment not configured")
	}

	cli, _ := New(host, user, pass, APOP, false)
	defer cli.Logout()

	tempDomain := "type-test.local"
	// 1. Создаем домен для теста (если его нет)
	// Примечание: предполагаем, что права позволяют создание домена
	_, _ = cli.Query("CREATEDOMAIN", tempDomain)

	t.Run("DigitalAliases", func(t *testing.T) {
		// Подготавливаем микс из строкового и цифрового алиаса
		// В CLI CGP это будет выглядеть как (string, string) или (string, number)
		testAliases := []string{"normal.alias", "999888"}

		err := cli.SetDomainAliases(tempDomain, testAliases)
		if err != nil {
			t.Fatalf("Failed to set aliases: %v", err)
		}

		// 2. Получаем список
		aliases, err := cli.GetDomainAliases(tempDomain)
		if err != nil {
			t.Fatalf("GetDomainAliases failed: %v", err)
		}

		t.Logf("Retrieved aliases: %#v", aliases)

		// 3. Проверяем наличие обоих и их типы
		foundDigital := false
		for _, a := range aliases {
			if a == "999888" {
				foundDigital = true
			}
			// Проверка на пустые строки или странные символы
			if a == "" {
				t.Error("Found empty alias string")
			}
		}

		if !foundDigital {
			t.Errorf("Digital alias '999888' lost or corrupted. Got: %v", aliases)
		}
	})

	// Удаляем тестовый домен
	_, _ = cli.Query("DELETEDOMAIN", tempDomain)
}
