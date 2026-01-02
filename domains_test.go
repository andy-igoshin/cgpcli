package cgpcli

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// helper для подключения, чтобы не дублировать код
func testGetClient(t *testing.T) *Cli {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Fatal("CGPUSER and CGPPASS env vars are required")
	}
	// Используем APOP для безопасности, как в ваших примерах
	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	return cli
}

// --- ТЕСТЫ ЧТЕНИЯ (Безопасны для рабочего сервера) ---

func TestDomainReadOperations(t *testing.T) {
	cli := testGetClient(t)
	defer cli.Logout()

	t.Run("ListDomains", func(t *testing.T) {
		domains, err := cli.ListDomains()
		if err != nil {
			t.Errorf("ListDomains failed: %v", err)
		}
		if len(domains) == 0 {
			t.Log("Warning: no domains found on server")
		}
		t.Logf("Found %d domains", len(domains))
	})

	t.Run("GetDomainSettingsMain", func(t *testing.T) {
		settings, err := cli.GetDomainSettings("") // Пустая строка - главный домен
		if err != nil {
			t.Fatalf("GetDomainSettings(main) failed: %v", err)
		}
		if _, ok := settings["Name"]; !ok && len(settings) > 0 {
			t.Log("Successfully fetched settings for main domain")
		}
	})

	t.Run("GetDomainLocation", func(t *testing.T) {
		loc, err := cli.GetDomainLocation("")
		if err != nil {
			t.Errorf("GetDomainLocation failed: %v", err)
		}
		t.Logf("Main domain location: %s", loc)
	})
}

// --- ТЕСТЫ ЗАПИСИ (Только для выделенного сервера) ---

func TestDomainWriteOperations(t *testing.T) {
	cli := testGetClient(t)
	defer cli.Logout()

	// Генерируем уникальное имя домена для теста
	tempDomain := fmt.Sprintf("test-%d.local", time.Now().Unix())

	// 1. Создание
	t.Run("CreateDomain", func(t *testing.T) {
		settings := map[string]any{
			"MaxAccountSize": "100M",
		}
		err := cli.CreateDomain(tempDomain, false, "", settings)
		if err != nil {
			t.Fatalf("Failed to create temp domain %s: %v", tempDomain, err)
		}
		
		exists, _ := cli.IsDomainExists(tempDomain)
		if !exists {
			t.Error("Domain should exist after CreateDomain")
		}
	})

	// 2. Настройки и Алиасы
	t.Run("UpdateAndGetSettings", func(t *testing.T) {
		newSettings := map[string]any{"Comment": "Test domain for unit tests"}
		if err := cli.UpdateDomainSettings(tempDomain, newSettings); err != nil {
			t.Errorf("UpdateDomainSettings failed: %v", err)
		}

		// Проверка алиасов домена
		testAliases := []string{"alias1." + tempDomain, "alias2." + tempDomain}
		if err := cli.SetDomainAliases(tempDomain, testAliases); err != nil {
			t.Errorf("SetDomainAliases failed: %v", err)
		}

		gotAliases, err := cli.GetDomainAliases(tempDomain)
		if err != nil || len(gotAliases) != 2 {
			t.Errorf("GetDomainAliases failed or count mismatch: %v", err)
		}
	})

	// 3. Состояние (Suspend/Resume)
	t.Run("SuspendResume", func(t *testing.T) {
		if err := cli.SuspendDomain(tempDomain); err != nil {
			t.Errorf("SuspendDomain failed: %v", err)
		}
		if err := cli.ResumeDomain(tempDomain); err != nil {
			t.Errorf("ResumeDomain failed: %v", err)
		}
	})

	// 4. Переименование
	newName := "renamed-" + tempDomain
	t.Run("RenameDomain", func(t *testing.T) {
		err := cli.RenameDomain(tempDomain, newName, "")
		if err != nil {
			t.Fatalf("RenameDomain failed: %v", err)
		}
		
		exists, _ := cli.IsDomainExists(newName)
		if !exists {
			t.Error("Renamed domain should exist")
		}
	})

	// 5. Удаление (Cleanup)
	t.Run("DeleteDomain", func(t *testing.T) {
		err := cli.DeleteDomain(newName, true)
		if err != nil {
			t.Fatalf("DeleteDomain failed: %v", err)
		}
		
		exists, _ := cli.IsDomainExists(newName)
		if exists {
			t.Error("Domain should be gone after DeleteDomain")
		}
	})
}

// --- ТЕСТЫ TELNUMS (Если поддерживаются лицензией) ---

func TestTelnums(t *testing.T) {
	cli := testGetClient(t)
	defer cli.Logout()

	t.Run("ListServerTelnums", func(t *testing.T) {
		res, err := cli.ListServerTelnums("", 10)
		if err != nil {
			// Если лицензия не поддерживает телефонию, CGP вернет ошибку, это нормально
			t.Skipf("Skipping Telnums: %v", err)
			return
		}
		t.Logf("Found %d server telnums (Total: %d)", len(res.Items), res.Total)
	})
}
