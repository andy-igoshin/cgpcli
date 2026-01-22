package cgpcli

import (
	"os"
	"testing"
)

func TestPBXControl_Integration(t *testing.T) {
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

	const (
		testAccount = "postmaster"
		testDomain  = "customer.domain.name"
		testProgram = "clock"
	)

	t.Run("TaskLifecycle", func(t *testing.T) {
		// 1. Запуск задачи
		startParam := map[string]any{"mode": "test"}
		taskID, err := cli.StartPBXTask(testAccount, testProgram, "Main", startParam)
		if err != nil {
			t.Fatalf("StartPBXTask failed: %v (make sure pbx.sppr exists)", err)
		}
		t.Logf("Task started with ID: %s", taskID)

		// 2. Чтение статуса
		status, err := cli.ReadNodeStatus(taskID)
		if err != nil {
			t.Errorf("ReadNodeStatus failed: %v", err)
		}
		t.Logf("Initial task status: %v", status)

		// 3. Отправка события
		// Проверяем SendTaskEvent без параметра
		if err := cli.SendTaskEvent(taskID, "ping", nil); err != nil {
			t.Errorf("SendTaskEvent (ping) failed: %v", err)
		}

		// Проверяем SendTaskEvent с параметром
		eventData := map[string]any{"code": 200}
		if err := cli.SendTaskEvent(taskID, "update", eventData); err != nil {
			t.Errorf("SendTaskEvent (update) failed: %v", err)
		}

		// 4. Завершение задачи
		if err := cli.KillNode(taskID); err != nil {
			t.Errorf("KillNode failed: %v", err)
		}
	})
}
