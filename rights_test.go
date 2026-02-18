package cgpcli

import (
	"fmt"
	"testing"
	"time"
)

func TestAccountRights(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Создаем временного "админа" для теста
	domain := "test.domain.name"
	testAcc := fmt.Sprintf("test-admin-%d@%s", time.Now().Unix(), domain)

	err := cli.CreateAccount(testAcc, nil, MultiMailbox, "", false)
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}
	defer cli.DeleteAccount(testAcc)

	t.Run("SetAndVerifyRights", func(t *testing.T) {
		// Список прав для установки
		// CanModifySettings, CanCreateAccounts
		targetRights := []string{"CanModifySettings", "CanCreateDomains"}

		err := cli.SetAccountRights(testAcc, targetRights)
		if err != nil {
			t.Fatalf("SetAccountRights failed: %v", err)
		}

		// Проверяем через GetAccountRights (метод из accounts.go)
		currentRights, err := cli.GetAccountRights(testAcc)
		if err != nil {
			t.Fatalf("GetAccountRights failed: %v", err)
		}

		if len(currentRights) != len(targetRights) {
			t.Errorf("Expected %d rights, got %d", len(targetRights), len(currentRights))
		}

		// Простая проверка наличия
		found := 0
		for _, r := range currentRights {
			for _, tR := range targetRights {
				if r == tR {
					found++
				}
			}
		}
		if found != len(targetRights) {
			t.Errorf("Rights mismatch. Expected %v, got %v", targetRights, currentRights)
		}
	})

	t.Run("ClearRights", func(t *testing.T) {
		// Очистка прав (пустой список)
		err := cli.SetAccountRights(testAcc, []string{})
		if err != nil {
			t.Errorf("Failed to clear rights: %v", err)
		}

		current, _ := cli.GetAccountRights(testAcc)
		if len(current) > 0 {
			t.Errorf("Rights should be empty, got %v", current)
		}
	})
}
