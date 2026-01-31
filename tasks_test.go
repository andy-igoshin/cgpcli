package cgpcli

import (
	"testing"
)

func TestNamedTasks_Integration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	const (
		domain   = "test.domain.name"
		baseTask = "sdk_test_task"
		testTask = baseTask + "@" + domain
		testAcc  = "testuser1@" + domain
	)

	_ = cli.DeleteNamedTask(testTask)

	t.Run("Lifecycle", func(t *testing.T) {
		// 1. Создание (проверяем отсечение домена в FOR)
		if err := cli.CreateNamedTask(testTask, testAcc); err != nil {
			t.Fatalf("CreateNamedTask failed: %v", err)
		}

		// 2. Листинг домена
		domainTasks, err := cli.ListDomainNamedTasks(domain)
		if err != nil {
			t.Errorf("ListDomainNamedTasks failed: %v", err)
		}
		// Сервер возвращает имена без @domain в ключах словаря
		if _, ok := domainTasks[baseTask]; !ok {
			t.Errorf("Task %s not found in list (got: %v)", baseTask, domainTasks)
		}

		// 3. Get/Update/Delete (здесь используем полное имя с @domain)
		settings, err := cli.GetNamedTask(testTask)
		if err != nil || settings["ownerName"] != "testuser1" {
			t.Errorf("GetNamedTask error or owner mismatch: %v", err)
		}

		if err := cli.UpdateNamedTask(testTask, map[string]any{"RealName": "SDK Task"}); err != nil {
			t.Errorf("UpdateNamedTask failed: %v", err)
		}

		if err := cli.DeleteNamedTask(testTask); err != nil {
			t.Errorf("DeleteNamedTask failed: %v", err)
		}
	})
}
