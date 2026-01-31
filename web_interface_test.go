package cgpcli

import (
	"testing"
)

func TestWebInterface_Integration(t *testing.T) {
        cli := getTestCli(t)
        defer cli.Logout()

	const (
		testAcc  = "webtestuser@test.domain.name"
		testPass = "webPass123"
		testAddr = "127.0.0.1"
	)

	// Подготовка: создаем аккаунт для тестов сессий
	_ = cli.DeleteAccount(testAcc) // На случай, если остался от прошлого прогона
	if err := cli.CreateAccount(testAcc, map[string]any{"Password": testPass}, TypeDefault, "", false); err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}
	defer cli.DeleteAccount(testAcc)

	var sessionID string

	// 1. Тест создания XIMSS сессии
	t.Run("CreateXimssSession", func(t *testing.T) {
		id, err := cli.CreateXimssSession(testAcc, testAddr, "")
		if err != nil {
			t.Fatalf("CreateXimssSession failed: %v", err)
		}
		if id == "" {
			t.Fatal("Received empty sessionID")
		}
		sessionID = id
	})

	// 2. Тест получения данных сессии (GetSession)
	t.Run("GetSession_Data", func(t *testing.T) {
		data, err := cli.GetSession(sessionID, "")
		if err != nil {
			t.Fatalf("GetSession failed: %v", err)
		}

		// В CLI PWD возвращается ключ 'accountName', а не 'Account'
		if acc, ok := data["accountName"].(string); !ok || acc != "webtestuser" {
			t.Errorf("Unexpected account in session: %v (data: %+v)", data["accountName"], data)
		}
	})

	// 3. Тест обновления параметров сессии (UpdateSession)
	t.Run("UpdateSession_CustomData", func(t *testing.T) {
		custom := map[string]any{"MyParam": "MyValue"}
		err := cli.UpdateSession(sessionID, "", custom)
		if err != nil {
			t.Fatalf("UpdateSession failed: %v", err)
		}

		// Проверяем, что данные применились
		data, _ := cli.GetSession(sessionID, "")
		if val, ok := data["MyParam"].(string); !ok || val != "MyValue" {
			t.Errorf("Custom param not found or incorrect: %v", data["MyParam"])
		}
	})

	// 4. Тест поиска сессии (FindAccountSession)
	t.Run("FindAccountSession", func(t *testing.T) {
		id, err := cli.FindAccountSession(testAcc, testAddr, "", "", "", "")
		if err != nil {
			t.Fatalf("FindAccountSession failed: %v", err)
		}
		if id != sessionID {
			t.Errorf("FindAccountSession returned %q, expected %q", id, sessionID)
		}
	})

	// 5. Тест списка сессий (ListAccountSessions)
	t.Run("ListAccountSessions", func(t *testing.T) {
		ids, err := cli.ListAccountSessions(testAcc, "", "", "", "", "")
		if err != nil {
			t.Fatalf("ListAccountSessions failed: %v", err)
		}
		found := false
		for _, id := range ids {
			if id == sessionID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("sessionID %q not found in account sessions list %v", sessionID, ids)
		}
	})

	// 6. Тест удаления сессии (KillSession)
	t.Run("KillSession", func(t *testing.T) {
		err := cli.KillSession(sessionID, "")
		if err != nil {
			t.Fatalf("KillSession failed: %v", err)
		}

		// Проверяем, что сессия пропала
		_, err = cli.GetSession(sessionID, "")
		if err == nil {
			t.Error("GetSession should have failed for killed session")
		}
	})
}
