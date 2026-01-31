package cgpcli

import (
	"testing"
)

func TestStatistics(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	// Используем текущий аккаунт для тестов
	testAcc := "*"
	// Используем текущий домен
	testDomain := "*"

	t.Run("AccountFullStat", func(t *testing.T) {
		res, err := cli.GetAccountStat(testAcc, "")
		if err != nil {
			t.Errorf("GetAccountStat (full) failed: %v", err)
			return
		}
		m, ok := res.(map[string]any)
		if !ok {
			t.Errorf("Expected map[string]any, got %T", res)
		}
		if len(m) == 0 {
			t.Log("Warning: Account stats map is empty (check if stats are enabled in Settings)")
		}
	})

	t.Run("AccountKeyStat", func(t *testing.T) {
		// Проверяем количество логинов — этот ключ есть всегда, если статистика включена
		val, err := cli.GetAccountStat(testAcc, "Logins")
		if err != nil {
			t.Errorf("GetAccountStat (Logins) failed: %v", err)
		}
		t.Logf("Logins count: %v", val)
	})

	t.Run("DomainFullStat", func(t *testing.T) {
		res, err := cli.GetDomainStat(testDomain, "")
		if err != nil {
			t.Errorf("GetDomainStat (full) failed: %v", err)
			return
		}
		if _, ok := res.(map[string]any); !ok {
			t.Errorf("Expected map, got %T", res)
		}
	})

	t.Run("ResetStats", func(t *testing.T) {
		// Сброс конкретного ключа (безопаснее, чем полный сброс в тестах)
		err := cli.ResetAccountStat(testAcc, "Logins")
		if err != nil {
			t.Errorf("ResetAccountStat failed: %v", err)
		}
	})
}
