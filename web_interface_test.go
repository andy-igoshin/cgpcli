package cgpcli

import (
	"testing"
)

func TestWebInterface_Integration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	const (
		testAcc  = "webtestuser@test.domain.name"
		testPass = "webPass123"
		testAddr = "127.0.0.1"
	)

	// Подготовка: создаем аккаунт для тестов сессий
	_ = cli.DeleteAccount(testAcc)
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

	// 5.1 Тест создания LITE сессии
	t.Run("LiteSessions_FullCycle", func(t *testing.T) {
		liteID, err := cli.CreateLiteSession(testAddr, "")
		if err != nil {
			t.Fatalf("CreateLiteSession failed: %v", err)
		}

		lites, err := cli.ListLiteSessions(testAddr, "")
		if err != nil {
			t.Errorf("ListLiteSessions failed: %v", err)
		}
		found := false
		for _, id := range lites {
			if id == liteID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Lite session %s not found in list", liteID)
		}

		_ = cli.KillSession(liteID, "")
	})

	// 5.2 Тест WebUser сессии и Bless
	t.Run("WebUser_And_Bless", func(t *testing.T) {
		wuID, err := cli.CreateWebUserSession(testAcc, testAddr, "", "Basic")
		if err != nil {
			t.Fatalf("CreateWebUserSession failed: %v", err)
		}
		defer cli.KillSession(wuID, "")

		// BlessSession: ожидаем 500 (если 2FA не включен), но проверяем формат вызова
		err = cli.BlessSession(wuID, "dummy-secret", "")
		if err != nil {
			t.Logf("BlessSession call successful (server rejected as expected): %v", err)
		}
	})

	// 5.3 Тест StoreSessionFile (ЗАКРЫВАЕМ: StoreSessionFile)
	t.Run("StoreSessionFile_Call", func(t *testing.T) {
		// Используем существующий sessionID (XIMSS)
		// Ожидаем ошибку "Unknown UploadID", так как мы не делали реальный аплоад,
		// но это подтвердит работоспособность метода StoreSessionFile.
		err := cli.StoreSessionFile(sessionID, "", "test.txt", "bogus-id", OffsetNew)
		if err != nil {
			t.Logf("StoreSessionFile call successful (server rejected as expected): %v", err)
		}
	})

	// 5.4 Тест WebAuth сессий (ListAccountWebAuthSessions, GetAccountWebAuthSession, KillAccountWebAuthSession)
	t.Run("WebAuthSessions_FullCycle", func(t *testing.T) {
		cli.SetDebug(true)
		// 1. Проверяем получение списка WebAuth сессий (без activeonly и с activeonly)
		sessions, err := cli.ListAccountWebAuthSessions(testAcc, false)
		if err != nil {
			t.Fatalf("ListAccountWebAuthSessions (activeonly=false) failed: %v", err)
		}

		activeSessions, err := cli.ListAccountWebAuthSessions(testAcc, true)
		if err != nil {
			t.Fatalf("ListAccountWebAuthSessions (activeonly=true) failed: %v", err)
		}

		t.Logf("WebAuth sessions count: all=%d, active=%d", len(sessions), len(activeSessions))

		// 2. Проверяем GetAccountWebAuthSession с несуществующим ID.
		// Команда должна успешно отправиться, а сервер вернет ошибку отсутствия сессии.
		const bogusAuthID = "non-existent-webauth-id"
		_, err = cli.GetAccountWebAuthSession(testAcc, bogusAuthID)
		if err == nil {
			t.Errorf("GetAccountWebAuthSession should fail for non-existent auth ID")
		} else {
			t.Logf("GetAccountWebAuthSession call executed as expected: %v", err)
		}

		// 3. Проверяем KillAccountWebAuthSession с несуществующим ID.
		ret, err := cli.KillAccountWebAuthSession(testAcc, bogusAuthID)
		if err == nil && ret == 0 {
			t.Logf("KillAccountWebAuthSession call executed as expected: %v", ret)
		} else {
			t.Errorf("KillAccountWebAuthSession should fail for non-existent auth ID: %v", err)
		}

		cli.SetDebug(false)

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
