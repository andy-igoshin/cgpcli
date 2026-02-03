package cgpcli

import (
	"testing"
)

func TestNamedTasks_Integration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	const (
		domain      = "test.domain.name"
		baseTask    = "sdk_test_task"
		renamedBase = "sdk_renamed_task"
		testTask    = baseTask + "@" + domain
		renamedTask = renamedBase + "@" + domain
		testAcc     = "testuser1@" + domain
	)

	// Предварительная очистка
	_ = cli.DeleteNamedTask(testTask)
	_ = cli.DeleteNamedTask(renamedTask)

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

	t.Run("Lifecycle_With_Rename_And_AccountList", func(t *testing.T) {
		// 1. Создание
		if err := cli.CreateNamedTask(testTask, testAcc); err != nil {
			t.Fatalf("CreateNamedTask failed: %v", err)
		}

		// 2. [!] ЗАКРЫВАЕМ: ListAccountNamedTasks (line 37)
		// Проверяем, видит ли аккаунт свою задачу
		accTasks, err := cli.ListAccountNamedTasks(testAcc)
		if err != nil {
			t.Fatalf("ListAccountNamedTasks failed: %v", err)
		}
		// В списке задач аккаунта обычно возвращается короткое имя
		if _, ok := accTasks[baseTask]; !ok {
			t.Errorf("Task %s not found in account tasks (got: %+v)", baseTask, accTasks)
		}

		// 3. [!] ЗАКРЫВАЕМ: RenameNamedTask (line 50)
		if err := cli.RenameNamedTask(testTask, renamedTask); err != nil {
			t.Fatalf("RenameNamedTask failed: %v", err)
		}

		// Проверяем, что под старым именем задачи нет, а под новым — есть
		domainTasks, _ := cli.ListDomainNamedTasks(domain)
		if _, ok := domainTasks[baseTask]; ok {
			t.Error("Old task name still exists after rename")
		}
		if _, ok := domainTasks[renamedBase]; !ok {
			t.Errorf("Renamed task %s not found in domain list", renamedBase)
		}

		// 4. Проверка метаданных после переименования
		settings, err := cli.GetNamedTask(renamedTask)
		if err != nil {
			t.Fatalf("GetNamedTask failed for renamed task: %v", err)
		}
		if settings["ownerName"] != "testuser1" {
			t.Errorf("Owner mismatch after rename: got %v", settings["ownerName"])
		}

		// 5. Cleanup
		if err := cli.DeleteNamedTask(renamedTask); err != nil {
			t.Errorf("DeleteNamedTask failed: %v", err)
		}
	})
}
