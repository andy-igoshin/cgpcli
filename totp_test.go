package cgpcli

import (
	"testing"
)

func TestAccountTOTPManagement(t *testing.T) {
	cli := getTestCli(t)

	domain := "test.domain.name"
	accountName := "totp-test-user"
	fullAccount := accountName + "@" + domain

	// 1. Подготовка: Создаем тестовый аккаунт
	settings := map[string]any{"RealName": "TOTP Test User"}
	err := cli.CreateAccount(fullAccount, settings, MultiMailbox, "", false)
	if err != nil {
		t.Fatalf("failed to create test account for TOTP: %v", err)
	}
	t.Cleanup(func() {
		_ = cli.DeleteAccount(fullAccount)
	})

	// Тестовый секрет в Base32 (стандартный тестовый ключ RFC 6238 / Google Authenticator)
	testSecret := "JBSWY3DPEHPK3PXP"

	// 2. Сценарий 1: Установка секрета -> Сброс путем передачи пустой строки ""
	t.Run("Set secret and clear via empty string", func(t *testing.T) {
		// Устанавливаем секрет
		err := cli.SetTotpSecret(fullAccount, testSecret)
		if err != nil {
			t.Fatalf("SetTotpSecret with valid secret failed: %v", err)
		}

		// Сбрасываем передачей пустой строки (SETTOTPSECRET account)
		err = cli.SetTotpSecret(fullAccount, "")
		if err != nil {
			t.Fatalf("SetTotpSecret with empty secret failed: %v", err)
		}
	})

	// 3. Сценарий 2: Установка секрета -> Сброс через ResetTotpSecret
	t.Run("Set secret and reset via ResetTotpSecret", func(t *testing.T) {
		// Устанавливаем секрет повторно
		err := cli.SetTotpSecret(fullAccount, testSecret)
		if err != nil {
			t.Fatalf("SetTotpSecret with valid secret failed: %v", err)
		}

		// Вызываем явную команду сброса (RESETTOTPSECRET account)
		err = cli.ResetTotpSecret(fullAccount)
		if err != nil {
			t.Fatalf("ResetTotpSecret failed: %v", err)
		}
	})

	// 4. Валидация входных параметров (быстрые проверки без вызова CLI)
	t.Run("Input validation", func(t *testing.T) {
		if err := cli.SetTotpSecret("", testSecret); err == nil {
			t.Error("expected error for empty account in SetTotpSecret, got nil")
		}

		if err := cli.ResetTotpSecret(""); err == nil {
			t.Error("expected error for empty account in ResetTotpSecret, got nil")
		}
	})
}
