package cgpcli

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDomainSetAdministration(t *testing.T) {
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

	// --- Базовые списки и инфо ---

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

		storages, err := cli.ListDomainStorage(false)
		if err != nil {
			t.Errorf("ListDomainStorage failed: %v", err)
		}
		t.Logf("Found %d storage paths", len(storages))
	})

	// --- Жизненный цикл домена ---

	t.Run("DomainLifecycle", func(t *testing.T) {
		tempName := fmt.Sprintf("test-dom-%d.local", time.Now().Unix())
		newName := "renamed-" + tempName

		// Создание
		if err := cli.CreateDomain(tempName, false, "", nil); err != nil {
			t.Fatalf("CreateDomain failed: %v", err)
		}

		// Переименование
		if err := cli.RenameDomain(tempName, newName, ""); err != nil {
			t.Errorf("RenameDomain failed: %v", err)
		}

		// Удаление
		if err := cli.DeleteDomain(newName, true); err != nil {
			t.Errorf("DeleteDomain failed: %v", err)
		}
	})

	// --- Настройки по умолчанию (Server) ---

	t.Run("ServerDefaults", func(t *testing.T) {
		// Domain Defaults
		if dd, err := cli.GetDomainDefaults(); err == nil {
			if err := cli.UpdateDomainDefaults(dd); err != nil {
				t.Errorf("UpdateDomainDefaults failed: %v", err)
			}
		} else {
			t.Errorf("GetDomainDefaults failed: %v", err)
		}

		// Account Defaults
		if ad, err := cli.GetServerAccountDefaults(); err == nil {
			if err := cli.UpdateServerAccountDefaults(ad); err != nil {
				t.Errorf("UpdateServerAccountDefaults failed: %v", err)
			}
		} else {
			t.Errorf("GetServerAccountDefaults failed: %v", err)
		}

		// Account Prefs
		if ap, err := cli.GetServerAccountPrefs(); err == nil {
			if err := cli.UpdateServerAccountPrefs(ap); err != nil {
				t.Errorf("UpdateServerAccountPrefs failed: %v", err)
			}
		} else {
			t.Errorf("GetServerAccountPrefs failed: %v", err)
		}
	})

	// --- Сертификаты и Интеграция ---

	t.Run("CertsAndIntegration", func(t *testing.T) {
		// Trusted Certs
		certs, err := cli.GetServerTrustedCerts()
		if err != nil {
			t.Errorf("GetServerTrustedCerts failed: %v", err)
		} else if len(certs) > 0 {
			// Идемпотентная проверка записи
			if err := cli.SetServerTrustedCerts(certs); err != nil {
				t.Errorf("SetServerTrustedCerts failed: %v", err)
			}
		}

		// Directory Integration
		if di, err := cli.GetDirectoryIntegration(); err == nil {
			if err := cli.SetDirectoryIntegration(di); err != nil {
				t.Errorf("SetDirectoryIntegration failed: %v", err)
			}
		} else {
			t.Errorf("GetDirectoryIntegration failed: %v", err)
		}
	})

	// --- Телефонные номера ---

	t.Run("TelNums", func(t *testing.T) {
		// Обычно требует лицензии на PBX, поэтому не фаталим при ошибке
		nums, err := cli.ListServerTelnums("", 10)
		if err != nil {
			t.Logf("ListServerTelnums info (check license): %v", err)
		} else {
			t.Logf("Found %d server telnums", len(nums))
		}
	})

	// --- Кластерные команды (с проверкой лицензии) ---

	t.Run("ClusterChecks", func(t *testing.T) {
		// Пробуем получить настройки кластера.
		// Если это не кластер или нет лицензии, CGP вернет ошибку.
		cdd, err := cli.GetClusterDomainDefaults()
		if err != nil {
			t.Logf("Cluster license not active or not a dynamic cluster: %v", err)
			return
		}

		t.Log("Cluster detected, running extended checks...")

		// Проверка Update Cluster
		if err := cli.UpdateClusterDomainDefaults(cdd); err != nil {
			t.Errorf("UpdateClusterDomainDefaults failed: %v", err)
		}

		// Кластерные сертификаты
		if cc, err := cli.GetClusterTrustedCerts(); err == nil {
			t.Logf("Found %d cluster certs", len(cc))
		}

		// Кластерное хранилище
		if cs, err := cli.ListDomainStorage(true); err != nil {
			t.Errorf("ListDomainStorage (SHARED) failed: %v", err)
		} else {
			t.Logf("Found %d shared storage paths", len(cs))
		}
	})

	// --- Прочее ---

	t.Run("Maintenance", func(t *testing.T) {
		if err := cli.ReloadDirectoryDomains(); err != nil {
			t.Errorf("ReloadDirectoryDomains failed: %v", err)
		}
	})
}

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

func TestTrustedCertsLive(t *testing.T) {
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

		// 2. Добавляем наш тестовый
		newCerts := append(original, testCert)
		if err := cli.SetServerTrustedCerts(newCerts); err != nil {
			t.Fatalf("Failed to set certs: %v", err)
		}

		// 3. Проверяем наличие
		updated, err := cli.GetServerTrustedCerts()
		if err != nil {
			t.Fatalf("Failed to verify certs: %v", err)
		}

		found := false
		for _, c := range updated {
			if c == clean {
				found = true
				break
			}
		}

		// 4. Возвращаем как было (важно!)
		defer cli.SetServerTrustedCerts(original)

		if !found {
			t.Error("Test certificate not found in the list after update")
		}
	})
}
