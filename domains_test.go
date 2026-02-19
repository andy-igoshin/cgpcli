package cgpcli

import (
	"testing"
)

func TestDomainAdministration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

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
		if err := cli.DeleteDirectoryRecords(targetDomain); err != nil {
			t.Errorf("DeleteDirectoryRecords failed: %v", err)
		}
		if err := cli.InsertDirectoryRecords(""); err != nil {
			t.Errorf("InsertDirectoryRecords failed: %v", err)
		}
		if err := cli.DeleteDirectoryRecords(""); err != nil {
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
	cli := getTestCli(t)
	defer cli.Close()

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

func TestDomainAdministration_Extra(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	targetDomain := "test.domain.name"

	t.Run("EffectiveAndExists", func(t *testing.T) {
		// IsDomainExists
		exists, err := cli.IsDomainExists(targetDomain)
		if err != nil || !exists {
			t.Errorf("IsDomainExists failed: exists=%v, err=%v", exists, err)
		}

		// GetDomainEffectiveSettings
		eff, err := cli.GetDomainEffectiveSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainEffectiveSettings failed: %v", err)
		}
		if len(eff) == 0 {
			t.Error("Effective settings map is empty")
		}
	})

	t.Run("Rules", func(t *testing.T) {
		// Mail Rules
		origMail, err := cli.GetDomainMailRules(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainMailRules failed: %v", err)
		}
		// Пытаемся установить (можно тот же список)
		if err := cli.SetDomainMailRules(targetDomain, origMail); err != nil {
			t.Errorf("SetDomainMailRules failed: %v", err)
		}

		// Signal Rules
		origSignal, err := cli.GetDomainSignalRules(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainSignalRules failed: %v", err)
		}
		if err := cli.SetDomainSignalRules(targetDomain, origSignal); err != nil {
			t.Errorf("SetDomainSignalRules failed: %v", err)
		}
	})

	t.Run("DefaultPrefsAndTemplates", func(t *testing.T) {
		// Работа с Account Default Prefs
		prefs, err := cli.GetAccountDefaultPrefs(targetDomain)
		if err != nil {
			t.Fatalf("GetAccountDefaultPrefs failed: %v", err)
		}

		// Update (добавляем/меняем язык по умолчанию)
		update := map[string]any{"Language": "English"}
		if err := cli.UpdateAccountDefaultPrefs(targetDomain, update); err != nil {
			t.Errorf("UpdateAccountDefaultPrefs failed: %v", err)
		}

		// Set (возвращаем старые)
		if err := cli.SetAccountDefaultPrefs(targetDomain, prefs); err != nil {
			t.Errorf("SetAccountDefaultPrefs failed: %v", err)
		}

		// Аналогично для AccountDefaults (Settings)
		defaults, _ := cli.GetAccountDefaults(targetDomain)
		if err := cli.UpdateAccountDefaults(targetDomain, map[string]any{"MaxStorage": "100M"}); err != nil {
			t.Logf("UpdateAccountDefaults info: %v", err)
		}
		_ = cli.SetAccountDefaults(targetDomain, defaults)

		// И для Template
		tpl, _ := cli.GetAccountTemplate(targetDomain)
		if err := cli.UpdateAccountTemplate(targetDomain, map[string]any{"RealName": "Template"}); err != nil {
			t.Logf("UpdateAccountTemplate info: %v", err)
		}
		_ = cli.SetAccountTemplate(targetDomain, tpl)
	})

	t.Run("Telnums", func(t *testing.T) {
		// ListDomainTelnums
		nums, err := cli.ListDomainTelnums(targetDomain, "", 10)
		if err != nil {
			// Метод может вернуть ошибку, если в домене вообще нет telnums/PBX
			t.Logf("ListDomainTelnums info: %v", err)
		} else {
			t.Logf("Domain telnums found: %d", len(nums))
		}
	})

	t.Run("DomainSettings_DirectSet", func(t *testing.T) {
		// Мы уже тестировали Update, теперь проверим полный Set
		settings, err := cli.GetDomainSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainSettings failed: %v", err)
		}
		// Пробуем перезаписать
		if err := cli.SetDomainSettings(targetDomain, settings); err != nil {
			t.Errorf("SetDomainSettings failed: %v", err)
		}
	})
}
