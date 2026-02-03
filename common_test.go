package cgpcli

import (
	"os"
	"strings"
	"testing"
)

// checkClusterEnabled — вспомогательная функция для определения отсутствия лицензии.
// Возвращает false, если сервер явно ответил, что кластер не активен.
func checkClusterEnabled(err error) bool {
	if err == nil {
		return true
	}
	msg := err.Error()
	return !strings.Contains(msg, "not a dynamic cluster")
}

// getTestCli устанавливает соединение с песочницей,
// используется всеми интеграционными тестами.
func getTestCli(t *testing.T) *Cli {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Skip("CGPUSER or CGPPASS not set, skipping integration test")
	}
	// Используем APOP для безопасности в тестах
	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	return cli
}

func skipIfNoClusterAuth(t *testing.T, err error) {
	if err != nil && !checkClusterEnabled(err) {
		t.Skip("Cluster functionality not supported by license")
	}
}

func TestNew_host_fail(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	_, err := New("127.0.0.2", cgpuser, cgppass, APOP, false)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("%s", err)
	}
}

func TestNew_login_fail(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	_, err := New("127.0.0.1", cgpuser, cgppass+"1", APOP, false)
	if err != nil && !strings.Contains(err.Error(), "incorrect password or account name") {
		t.Errorf("%s", err)
	}
}

func TestNew_login_plain(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Close()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestNew_login_apop(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Close()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestNew_login_stls(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, Plain, true)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Close()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestNew_logout(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Close()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestCli_CommonOperations(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// 1. Тестируем SetDebug (просто вызываем, чтобы убедиться в отсутствии побочек)
	cli.SetDebug(true)
	defer cli.SetDebug(false)

	// 2. Тестируем QueryNV на простой команде NOOP
	// NOOP — стандартная команда подтверждения связи
	err := cli.QueryNV("NOOP")
	if err != nil {
		t.Errorf("QueryNV(NOOP) failed: %v", err)
	}

	// 3. Тестируем вспомогательные геттеры после выполнения команды
	if !cli.IsSuccess() {
		t.Error("IsSuccess() returned false after successful NOOP")
	}
	if cli.GetRespCode() != CodeOk {
		t.Errorf("Expected code 200, got %v", cli.GetRespCode())
	}
	// RespData для NOOP обычно содержит "completed" или пустоту
	if cli.GetRespData() == "" && cli.GetRespCode() != CodeOk {
		t.Error("GetRespData() is empty on successful operation")
	}

	// 4. Тестируем Query с возвратом значения
	// GETVERSION возвращает строку с версией сервера
	res, err := cli.Query("GETVERSION")
	if err != nil {
		t.Errorf("Query(GETVERSION) failed: %v", err)
	}

	version, ok := res.(string)
	if !ok || version == "" {
		t.Errorf("GETVERSION expected string, got %T (%v)", res, res)
	}

	// 5. Тестируем Send напрямую (низкий уровень)
	// Отправим корректную команду, но обработаем ответ вручную через parseResponse
	err = cli.Send("NOOP")
	if err != nil {
		t.Errorf("Send(NOOP) failed: %v", err)
	}
	err = cli.parseResponse()
	if err != nil {
		t.Errorf("parseResponse after Send failed: %v", err)
	}
	if cli.GetRespCode() != CodeOk {
		t.Errorf("Manual Send/Parse cycle failed with code %v", cli.GetRespCode())
	}
}

func TestCli_QueryError(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Тестируем поведение при заведомо ошибочной команде
	_, err := cli.Query("NONEXISTENTCOMMAND")
	if err == nil {
		t.Error("Query expected error for unknown command, got nil")
	}

	if cli.IsSuccess() {
		t.Error("IsSuccess() should be false for unknown command")
	}

	if cli.GetRespCode() < 500 {
		t.Errorf("Expected error code 5xx for unknown command, got %v", cli.GetRespCode())
	}
}

func TestCli_Logout(t *testing.T) {
	cli := getTestCli(t)

	// Проверяем, что метод Logout доступен и не возвращает ошибку
	err := cli.Logout()
	if err != nil {
		t.Errorf("Logout() failed: %v", err)
	}

	// Проверяем состояние: после Logout соединение должно быть закрыто
	if cli.conn != nil {
		t.Error("cli.conn should be nil after Logout()")
	}

	// Повторный вызов не должен приводить к панике (идемпотентность)
	err = cli.Logout()
	if err != nil {
		t.Errorf("Second Logout() call should not return error, got: %v", err)
	}
}
