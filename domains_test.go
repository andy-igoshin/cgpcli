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

	targetDomain := "test.domain.name"

	t.Run("DomainSettingsAndLocation", func(t *testing.T) {
		// Проверка получения настроек
		if _, err := cli.GetDomainSettings(targetDomain); err != nil {
			t.Fatalf("GetDomainSettings failed: %v", err)
		}

		// Обновление комментария
		update := map[string]any{"Comment": "Updated via Go SDK"}
		if err := cli.UpdateDomainSettings(targetDomain, update); err != nil {
			t.Errorf("UpdateDomainSettings failed: %v", err)
		}

		// Проверка локации
		if _, err := cli.GetDomainLocation(targetDomain); err != nil {
			t.Errorf("GetDomainLocation failed: %v", err)
		}
	})

	t.Run("DomainAliases", func(t *testing.T) {
		original, _ := cli.GetDomainAliases(targetDomain)
		
		newAliases := append(original, "extra-alias."+targetDomain)
		if err := cli.SetDomainAliases(targetDomain, newAliases); err != nil {
			t.Errorf("SetDomainAliases failed: %v", err)
		}

		// Восстановление оригинальных алиасов
		_ = cli.SetDomainAliases(targetDomain, original)
	})

	t.Run("ListDomainObjects", func(t *testing.T) {
		// Простой список аккаунтов
		accs, err := cli.ListAccounts(targetDomain)
		if err != nil {
			t.Errorf("ListAccounts failed: %v", err)
		}
		t.Logf("Accounts count: %d", len(accs))

		// Детальные объекты
		res, err := cli.ListDomainObjects(targetDomain, 20, "", TypeAll, "")
		if err != nil {
			t.Errorf("ListDomainObjects failed: %v", err)
		} else {
			t.Logf("Total objects in domain: %d", res.Total)
		}
	})

	t.Run("SuspendResume", func(t *testing.T) {
		if err := cli.SuspendDomain(targetDomain); err != nil {
			t.Errorf("Suspend failed: %v", err)
		}
		if err := cli.ResumeDomain(targetDomain); err != nil {
			t.Errorf("Resume failed: %v", err)
		}
	})
}
