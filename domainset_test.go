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
	defer cli.Logout()

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

		if err := cli.CreateDomain(tempName, false, "", nil); err != nil {
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
