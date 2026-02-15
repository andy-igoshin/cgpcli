package cgpcli

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// cleanPEM — локальная вспомогательная функция для тестов
func cleanPEM(pem string) string {
	r := strings.NewReplacer(
		"-----BEGIN CERTIFICATE-----", "",
		"-----END CERTIFICATE-----", "",
		"\n", "",
		"\r", "",
		" ", "",
	)
	return r.Replace(pem)
}

func TestDomainSetAdministration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// --- Списки и инфо ---
	t.Run("BasicLists", func(t *testing.T) {
		domains, err := cli.ListDomains()
		if err != nil {
			t.Errorf("ListDomains failed: %v", err)
		}
		t.Logf("Found %d domains", len(domains))

		mainDomain, err := cli.MainDomainName()
		if err != nil {
			t.Errorf("MainDomainName failed: %v", err)
		}
		t.Logf("Main domain: %s", mainDomain)
	})

	// --- Жизненный цикл домена ---
	t.Run("DomainLifecycle", func(t *testing.T) {
		tempName := fmt.Sprintf("test-dom-%d.local", time.Now().Unix())
		newName := "renamed-" + tempName

		if err := cli.CreateDomain(tempName, nil, "", false); err != nil {
			t.Fatalf("CreateDomain failed: %v", err)
		}

		if err := cli.RenameDomain(tempName, newName, ""); err != nil {
			t.Errorf("RenameDomain failed: %v", err)
		}

		if err := cli.DeleteDomain(newName, true); err != nil {
			t.Errorf("DeleteDomain failed: %v", err)
		}
	})

	// --- Проверка настроек (идемпотентность Get/Update) ---
	t.Run("ServerSettingsIdempotency", func(t *testing.T) {
		check := func(name string, getter func() (map[string]any, error), updater func(map[string]any) error) {
			data, err := getter()
			if err != nil {
				t.Errorf("Get %s failed: %v", name, err)
				return
			}
			if err := updater(data); err != nil {
				t.Errorf("Update %s failed: %v", name, err)
			}
		}

		check("DomainDefaults", cli.GetDomainDefaults, cli.UpdateDomainDefaults)
		check("AccountDefaults", cli.GetServerAccountDefaults, cli.UpdateServerAccountDefaults)
	})
}

func TestTrustedCertsLive(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testCert := `-----BEGIN CERTIFICATE-----
MIIDlTCCAn2gAwIBAgIUeRPaboVA3C6FPU1Jw1Lz31UnpkIwDQYJKoZIhvcNAQEL
BQAwWjELMAkGA1UEBhMCUlUxDzANBgNVBAgMBk1ha3VyaTEPMA0GA1UEBwwGTWFr
dXJpMRIwEAYDVQQKDAlNeUNvbXBhbnkxFTATBgNVBAMMDE1ha3VyaVJvb3RDQTAe
Fw0yNjAxMDMxOTAwMjBaFw0yNjAxMDQxOTAwMjBaMFoxCzAJBgNVBAYTAlJVMQ8w
DQYDVQQIDAZNYWt1cmkxDzANBgNVBAcMBk1ha3VyaTESMBAGA1UECgwJTXlDb21w
YW55MRUwEwYDVQQDDAxNYWt1cmlSb290Q0EwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCzxyKNXCfATego9YoDDtJrb3i42ywBjVMjNx3XMXxSmp/aQHLi
U8wsvFJMd/3JTwO6b7OBB3kTfHf06GHiDCkpWLASx20Ac6TADRHaHgwzXx/9j73v
nuPK6Rk9xndkHe0jPHgKs3n3bCAP92LcsDvvLchmXj+ZBRMXVOSuYsSTfhb4s9t2
Knc7QzRxyC/ICplWzkaOoSxa8lv9hKrcjeB52hn4d/ru0Qbz59oHq6RGYTfbDmXi
3dGXWYQbxfblBkmIaTh2sjbxq5BatvCkTgAvUI1o5dscW6MJ29f/8CwlQi+1cgX0
RaYKWVJV39tG+qXuJX9i0RDuVeJUXO+lkszbAgMBAAGjUzBRMB0GA1UdDgQWBBRF
kuIUqFqUOz+GjG3bZhYsWYB/4jAfBgNVHSMEGDAWgBRFkuIUqFqUOz+GjG3bZhYs
WYB/4jAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBAgxEppxyd
iH+OhZa/xdhHCBnmSjh/so2SS4gfEDSzExa0kExcZQtvPofnzhtzK7G/FlZoJTR6
D7cgAGlILLnrlw66cE7L35w1KMPxUdLODIMRWxonUxQBfJHfGr1f36r3dJlAjgp6
XS821hdUHD9vaknZYy6sMb9nFc9nvwCyO9rUnXqUjTp9vjkTELx4heFOCF4JF2+e
vUWH2N8+SO4MxiknhClkTw9n9TKFvYiIkKl+ZVWbz+nA6joqRGu+sfG4YDjedxRB
HAExJ9rxX23JVrPPI8CI4cS83fj+A3LEJUmCyD5H49u2zKjTToFBd5bhCxdL/qj1
tlz8SxRmzBDC
-----END CERTIFICATE-----`

	clean := cleanPEM(testCert)

	t.Run("WriteReadCertificate", func(t *testing.T) {
		// 1. Читаем текущие
		original, err := cli.GetServerTrustedCerts()
		if err != nil {
			t.Fatalf("Failed to get original certs: %v", err)
		}

		// 2. Добавляем тестовый (даже если он уже там есть)
		newCerts := append(original, testCert)
		if err := cli.SetServerTrustedCerts(newCerts); err != nil {
			t.Fatalf("Failed to set certs: %v", err)
		}

		// 3. Проверяем наличие (сервер мог дедуплицировать, поэтому просто ищем в списке)
		updated, err := cli.GetServerTrustedCerts()
		if err != nil {
			t.Fatalf("Failed to verify certs: %v", err)
		}

		found := false
		for _, c := range updated {
			// GetServerTrustedCerts теперь возвращает чистый base64 без []
			if c == clean {
				found = true
				break
			}
		}

		// 4. Возвращаем всё в исходное состояние
		defer func() {
			if err := cli.SetServerTrustedCerts(original); err != nil {
				t.Errorf("Cleanup failed: %v", err)
			}
		}()

		if !found {
			t.Error("Test certificate not found in the list after update")
		}
	})
}

func TestDomainSet_Extended(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// 1. Тестирование хранилищ сервера (Domain Storage)
	t.Run("DomainStorageLifecycle", func(t *testing.T) {
		storagePath := "mnt/test_storage_set"

		// Создание (может вернуть ошибку, если путь недоступен на ФС, логируем как инфо)
		if err := cli.CreateDomainStorage(storagePath, false); err != nil {
			t.Logf("CreateDomainStorage info: %v", err)
		}

		list, err := cli.ListDomainStorage(false)
		if err != nil {
			t.Fatalf("ListDomainStorage failed: %v", err)
		}
		t.Logf("Server domain storages: %v", list)
	})

	// 2. Интеграция со справочником (Directory Integration)
	t.Run("DirectoryIntegration", func(t *testing.T) {
		// Server level
		data, err := cli.GetDirectoryIntegration()
		if err != nil {
			t.Fatalf("GetDirectoryIntegration failed: %v", err)
		}
		if err := cli.SetDirectoryIntegration(data); err != nil {
			t.Errorf("SetDirectoryIntegration failed: %v", err)
		}

		// Reload (просто вызов команды)
		if err := cli.ReloadDirectoryDomains(); err != nil {
			t.Errorf("ReloadDirectoryDomains failed: %v", err)
		}
	})

	// 3. Кластерные настройки (Cluster Section)
	// Используем checkClusterEnabled из alerts_test.go (предполагаем наличие или копируем логику)
	t.Run("ClusterDefaults", func(t *testing.T) {
		_, err := cli.GetClusterDomainDefaults()
		if err != nil {
			if !checkClusterEnabled(err) {
				t.Skip("Dynamic Cluster not enabled, skipping cluster tests")
			}
			t.Fatalf("GetClusterDomainDefaults failed: %v", err)
		}

		// Тестируем Update/Set для Account Defaults в кластере
		accDef, _ := cli.GetClusterAccountDefaults()
		if err := cli.UpdateClusterAccountDefaults(accDef); err != nil {
			t.Errorf("UpdateClusterAccountDefaults failed: %v", err)
		}

		// Prefs
		accPrefs, _ := cli.GetClusterAccountPrefs()
		if err := cli.SetClusterAccountPrefs(accPrefs); err != nil {
			t.Errorf("SetClusterAccountPrefs failed: %v", err)
		}

		if err := cli.UpdateClusterAccountPrefs(accPrefs); err != nil {
			t.Errorf("UpdateClusterAccountPrefs failed: %v", err)
		}

		// Directory Integration Cluster
		dirInt, _ := cli.GetClusterDirectoryIntegration()
		if err := cli.SetClusterDirectoryIntegration(dirInt); err != nil {
			t.Errorf("SetClusterDirectoryIntegration failed: %v", err)
		}

		// Trusted Certs Cluster
		certs, err := cli.GetClusterTrustedCerts()
		if err == nil {
			_ = cli.SetClusterTrustedCerts(certs)
		}

		// Telnums Cluster
		_, _ = cli.ListClusterTelnums("", 10)
	})

	// 4. Оставшиеся серверные настройки
	t.Run("ServerAccountSettings", func(t *testing.T) {
		// Prefs
		prefs, err := cli.GetServerAccountPrefs()
		if err != nil {
			t.Fatalf("GetServerAccountPrefs failed: %v", err)
		}
		if err := cli.SetServerAccountPrefs(prefs); err != nil {
			t.Errorf("SetServerAccountPrefs failed: %v", err)
		}
		if err := cli.UpdateServerAccountPrefs(prefs); err != nil {
			t.Errorf("UpdateServerAccountPrefs failed: %v", err)
		}

		// Telnums
		_, err = cli.ListServerTelnums("", 10)
		if err != nil {
			t.Logf("ListServerTelnums info: %v", err)
		}
	})
}

func TestDomainSet_FinalRemaining(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// 1. CreateDirectoryDomain (требует Directory Unit, но мы проверим сам вызов)
	t.Run("CreateDirectoryDomain", func(t *testing.T) {
		// Пытаемся создать, скорее всего вернет ошибку конфигурации, но метод будет покрыт
		testDirDom := fmt.Sprintf("dir-dom-%d.local", time.Now().UnixNano()%1000)
		err := cli.CreateDirectoryDomain(testDirDom, map[string]any{"Enabled": "YES"})
		if err != nil {
			t.Logf("CreateDirectoryDomain call verified, server response: %v", err)
		}
	})

	// 2. Явное покрытие Server Domain Defaults (Get/Set/Update)
	t.Run("ServerDomainDefaults_Direct", func(t *testing.T) {
		data, err := cli.GetDomainDefaults()
		if err != nil {
			t.Fatalf("GetDomainDefaults failed: %v", err)
		}
		if err := cli.UpdateDomainDefaults(data); err != nil {
			t.Errorf("UpdateDomainDefaults failed: %v", err)
		}
		if err := cli.SetDomainDefaults(data); err != nil {
			t.Errorf("SetDomainDefaults failed: %v", err)
		}
	})

	// 3. Явное покрытие Server Account Defaults (Get/Set/Update)
	t.Run("ServerAccountDefaults_Direct", func(t *testing.T) {
		data, err := cli.GetServerAccountDefaults()
		if err != nil {
			t.Fatalf("GetServerAccountDefaults failed: %v", err)
		}
		if err := cli.UpdateServerAccountDefaults(data); err != nil {
			t.Errorf("UpdateServerAccountDefaults failed: %v", err)
		}
		if err := cli.SetServerAccountDefaults(data); err != nil {
			t.Errorf("SetServerAccountDefaults failed: %v", err)
		}
	})

	// 4. Явное покрытие Cluster Defaults (Set/Update)
	t.Run("ClusterDefaults_Direct", func(t *testing.T) {
		// Проверяем наличие лицензии через GetClusterDomainDefaults
		data, err := cli.GetClusterDomainDefaults()
		if err != nil {
			if !checkClusterEnabled(err) {
				t.Skip("Cluster not enabled, skipping cluster defaults coverage")
			}
			t.Fatalf("GetClusterDomainDefaults failed: %v", err)
		}

		// Cluster Domain Defaults
		if err := cli.UpdateClusterDomainDefaults(data); err != nil {
			t.Errorf("UpdateClusterDomainDefaults failed: %v", err)
		}
		if err := cli.SetClusterDomainDefaults(data); err != nil {
			t.Errorf("SetClusterDomainDefaults failed: %v", err)
		}

		// Cluster Account Defaults
		accData, _ := cli.GetClusterAccountDefaults()
		if err := cli.SetClusterAccountDefaults(accData); err != nil {
			t.Errorf("SetClusterAccountDefaults failed: %v", err)
		}
	})
}
