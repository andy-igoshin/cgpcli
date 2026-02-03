package cgpcli

import (
	"reflect"
	"testing"
)

// 1. Тест логики преобразования типов в мапе алертов.
// Проверяем, что toStringMap корректно нормализует числовые алерты в строки.
func TestAlerts_TypeConversion(t *testing.T) {
	mockResponse := map[string]any{
		"20260115211048": "test alert",
		"20260115214705": int64(98765432109876),
	}

	t.Run("Convert any map to string map", func(t *testing.T) {
		res, err := toStringMap("GETSERVERALERTS", mockResponse)
		if err != nil {
			t.Fatalf("toStringMap failed: %v", err)
		}

		expected := map[string]string{
			"20260115211048": "test alert",
			"20260115214705": "98765432109876",
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Data mismatch.\nExpected: %v\nGot:      %v", expected, res)
		}
	})
}

// 2. Интеграционный тест жизненного цикла алертов на всех уровнях.
func TestAlerts_FullLifecycle(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testMsg := "Integration Test Alert"
	testDomain := "test.domain.name" // Должен существовать в тестовом окружении
	testAcc := "testuser1@" + testDomain

	// Вспомогательная функция для поиска ключа (таймстампа) по тексту алерта
	findAlertTS := func(m map[string]string, val string) string {
		for k, v := range m {
			if v == val {
				return k
			}
		}
		return ""
	}

	// SERVER ALERTS
	t.Run("ServerAlerts", func(t *testing.T) {
		if err := cli.PostServerAlert(testMsg); err != nil {
			t.Fatalf("PostServerAlert failed: %v", err)
		}
		list, err := cli.GetServerAlerts()
		if err != nil {
			t.Fatalf("GetServerAlerts failed: %v", err)
		}
		ts := findAlertTS(list, testMsg)
		if ts == "" {
			t.Errorf("Server alert not found in %v", list)
		} else {
			if err := cli.RemoveServerAlert(ts); err != nil {
				t.Errorf("RemoveServerAlert failed: %v", err)
			}
		}
	})

	// DOMAIN ALERTS
	t.Run("DomainAlerts", func(t *testing.T) {
		if err := cli.PostDomainAlert(testDomain, testMsg); err != nil {
			t.Fatalf("PostDomainAlert skipped (check domain exists): %v", err)
			return
		}
		list, err := cli.GetDomainAlerts(testDomain)
		if err != nil {
			t.Fatalf("GetDomainAlerts failed: %v", err)
		}
		ts := findAlertTS(list, testMsg)
		if ts != "" {
			if err := cli.RemoveDomainAlert(testDomain, ts); err != nil {
				t.Errorf("RemoveDomainAlert failed: %v", err)
			}
		}
	})

	// ACCOUNT ALERTS
	t.Run("AccountAlerts", func(t *testing.T) {
		// Используем "*" для текущего аккаунта или конкретный postmaster
		if err := cli.PostAccountAlert(testAcc, testMsg); err != nil {
			t.Logf("PostAccountAlert skipped: %v", err)
			return
		}
		list, err := cli.GetAccountAlerts(testAcc)
		if err != nil {
			t.Fatalf("GetAccountAlerts failed: %v", err)
		}
		ts := findAlertTS(list, testMsg)
		if ts != "" {
			if err := cli.RemoveAccountAlert(testAcc, ts); err != nil {
				t.Errorf("RemoveAccountAlert failed: %v", err)
			}
		}
	})

	// CLUSTER ALERTS (с обработкой отсутствия лицензии)
	t.Run("ClusterAlerts", func(t *testing.T) {
		err := cli.PostClusterAlert(testMsg)
		if err != nil {
			if !checkClusterEnabled(err) {
				t.Skip("Cluster license not found, skipping cluster alerts test")
			}
			t.Fatalf("PostClusterAlert failed: %v", err)
		}

		list, err := cli.GetClusterAlerts()
		if err != nil {
			t.Fatalf("GetClusterAlerts failed: %v", err)
		}

		ts := findAlertTS(list, testMsg)
		if ts != "" {
			if err := cli.RemoveClusterAlert(ts); err != nil {
				t.Errorf("RemoveClusterAlert failed: %v", err)
			}
		}
	})
}

// 3. Тест массовой установки (Set) алертов.
func TestAlerts_SetMethods(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Используем пустую мапу для проверки сброса/установки
	emptyAlerts := map[string]any{}

	t.Run("SetServerAlerts", func(t *testing.T) {
		if err := cli.SetServerAlerts(emptyAlerts); err != nil {
			t.Errorf("SetServerAlerts failed: %v", err)
		}
	})

	t.Run("SetAccountAlerts", func(t *testing.T) {
		// "*" — системный указатель на текущий аккаунт в CLI
		if err := cli.SetAccountAlerts("*", emptyAlerts); err != nil {
			t.Errorf("SetAccountAlerts failed: %v", err)
		}
	})

	t.Run("SetDomainAlerts", func(t *testing.T) {
		err := cli.SetDomainAlerts("test.domain.name", emptyAlerts)
		if err != nil {
			t.Logf("SetDomainAlerts info: %v", err)
		}
	})

	t.Run("SetClusterAlerts", func(t *testing.T) {
		err := cli.SetClusterAlerts(emptyAlerts)
		if err != nil {
			if !checkClusterEnabled(err) {
				t.Skip("Cluster license not found, skipping SetClusterAlerts")
			}
			t.Errorf("SetClusterAlerts failed: %v", err)
		}
	})
}
