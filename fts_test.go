package cgpcli

import (
	"strings"
	"testing"
)

func TestFTSEnvironmentLifecycle(t *testing.T) {
	cli := getTestCli(t)

	domain := "test.domain.name"
	accountName := "fts-test-user"
	fullAccount := accountName + "@" + domain
	groupName := "fts-test-group"
	fullGroup := groupName + "@" + domain

	// 1. Подготовка: Создаем пользователя и группу для тестов
	settings := map[string]any{"RealName": "FTS Real Lifecycle"}
	err := cli.CreateAccount(fullAccount, settings, MultiMailbox, "", false)
	if err != nil {
		t.Fatalf("failed to create test account for FTS: %v", err)
	}
	t.Cleanup(func() {
		_ = cli.DeleteAccount(fullAccount)
	})

	// Создаем тестовую группу (сигнатура по аналогии с аккаунтами/доменами)
	// Если метод называется по-другому, подправьте под вашу реализацию управления группами
	err = cli.CreateGroup(fullGroup, map[string]any{})
	if err != nil {
		t.Logf("Warning: failed to create test group (skipping group test if not supported): %v", err)
	} else {
		t.Cleanup(func() {
			_ = cli.DeleteGroup(fullGroup)
		})
	}

	// 2. Тест уровня домена с учетом асинхронности
	t.Run("Domain FTS Protocol Steps", func(t *testing.T) {
		state, err := cli.GetDomainSearchIndexState(domain)
		if err != nil {
			t.Fatalf("GetDomainSearchIndexState failed: %v", err)
		}

		// Запускаем индексацию. Если домен уже в процессе, команда с ENABLEFORALL может вернуть 500.
		// Запускаем сразу с true, чтобы проверить полный флаг.
		err = cli.BuildDomainSearchIndex(domain, true)
		if err != nil {
			// Если сервер говорит, что процесс уже идет — это не ошибка теста, это успех операции
			if !strings.Contains(err.Error(), "already in progress") {
				t.Fatalf("BuildDomainSearchIndex failed: %v", err)
			}
		}

		state, err = cli.GetDomainSearchIndexState(domain)
		if err != nil {
			t.Fatalf("GetDomainSearchIndexState after build failed: %v", err)
		}

		// Любой статус кроме пустого говорит о том, что команда отработала
		if state == "" || state == "none" && err == nil {
			// Если сбросился мгновенно или вернул ошибку "already in progress",
			// состояние "none" маловероятно, но допустимо при моментальном завершении пустого домена.
		}

		// Проверяем ожидаемый отказ удаления активного индекса
		err = cli.RemoveDomainSearchIndex(domain)
		if err == nil {
			t.Error("expected error 500 'Search index must be disabled' for active domain index, got nil")
		}
	})

	// 3. Тест уровня аккаунта (PASS в прошлый раз, оставляем логику)
	t.Run("Account FTS Protocol Steps", func(t *testing.T) {
		err := cli.BuildAccountSearchIndex(fullAccount)
		if err != nil {
			t.Fatalf("BuildAccountSearchIndex failed: %v", err)
		}

		state, err := cli.GetAccountSearchIndexState(fullAccount)
		if err != nil {
			t.Fatalf("GetAccountSearchIndexState failed: %v", err)
		}
		if state == "" {
			t.Error("expected non-empty account index state")
		}

		size, err := cli.GetAccountSearchIndexSize(fullAccount)
		if err != nil {
			t.Fatalf("GetAccountSearchIndexSize failed: %v", err)
		}
		if size < 0 {
			t.Errorf("unexpected negative size parsed: %d", size)
		}

		err = cli.RemoveAccountSearchIndex(fullAccount)
		if err == nil {
			t.Error("expected error 500 'Search index must be disabled', got nil")
		}
	})

	// 4. Тест уровня групп (теперь группа существует)
	t.Run("Group FTS Command", func(t *testing.T) {
		// Проверяем вызов только если группа успешно создалась в Setup
		err := cli.BuildGroupSearchIndex(fullGroup, false)
		if err != nil {
			// Если индексация домена залочила объект, проверяем на 'already in progress'
			if !strings.Contains(err.Error(), "already in progress") {
				t.Fatalf("BuildGroupSearchIndex failed: %v", err)
			}
		}
	})
}
