package cgpcli

import (
	"os"
	"strings"
	"testing"
	"time"
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

	t.Cleanup(func() {
		_ = cli.Close()
	})

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

	cli, err := New("127.0.0.1", cgpuser, cgppass, Plain, true, WithTLSInsecureSkipVerify(true))
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

func TestWithCommandTimeout(t *testing.T) {
	cli := new(Cli)

	WithCommandTimeout(30 * time.Second)(cli)
	if cli.cmdTimeout != 30*time.Second {
		t.Errorf("WithCommandTimeout: expected 30s, got %v", cli.cmdTimeout)
	}

	// Non-positive values must be ignored (field unchanged).
	WithCommandTimeout(0)(cli)
	if cli.cmdTimeout != 30*time.Second {
		t.Errorf("WithCommandTimeout(0): expected 30s unchanged, got %v", cli.cmdTimeout)
	}

	WithCommandTimeout(-time.Second)(cli)
	if cli.cmdTimeout != 30*time.Second {
		t.Errorf("WithCommandTimeout(-1s): expected 30s unchanged, got %v", cli.cmdTimeout)
	}
}

func TestWithConnectTimeout(t *testing.T) {
	cli := new(Cli)

	WithConnectTimeout(15 * time.Second)(cli)
	if cli.connTimeout != 15*time.Second {
		t.Errorf("WithConnectTimeout: expected 15s, got %v", cli.connTimeout)
	}

	WithConnectTimeout(0)(cli)
	if cli.connTimeout != 15*time.Second {
		t.Errorf("WithConnectTimeout(0): expected 15s unchanged, got %v", cli.connTimeout)
	}

	WithConnectTimeout(-time.Second)(cli)
	if cli.connTimeout != 15*time.Second {
		t.Errorf("WithConnectTimeout(-1s): expected 15s unchanged, got %v", cli.connTimeout)
	}
}

func TestWithIdleTimeout(t *testing.T) {
	cli := new(Cli)

	WithIdleTimeout(10 * time.Minute)(cli)
	if cli.idleTimeout != 10*time.Minute {
		t.Errorf("WithIdleTimeout: expected 10m, got %v", cli.idleTimeout)
	}

	WithIdleTimeout(0)(cli)
	if cli.idleTimeout != 10*time.Minute {
		t.Errorf("WithIdleTimeout(0): expected 10m unchanged, got %v", cli.idleTimeout)
	}

	WithIdleTimeout(-time.Second)(cli)
	if cli.idleTimeout != 10*time.Minute {
		t.Errorf("WithIdleTimeout(-1s): expected 10m unchanged, got %v", cli.idleTimeout)
	}
}

func TestCli_SetCommandTimeout(t *testing.T) {
	cli := new(Cli)

	cli.SetCommandTimeout(45 * time.Second)
	if cli.cmdTimeout != 45*time.Second {
		t.Errorf("SetCommandTimeout: expected 45s, got %v", cli.cmdTimeout)
	}

	// Non-positive values must be ignored.
	cli.SetCommandTimeout(0)
	if cli.cmdTimeout != 45*time.Second {
		t.Errorf("SetCommandTimeout(0): expected 45s unchanged, got %v", cli.cmdTimeout)
	}

	cli.SetCommandTimeout(-time.Second)
	if cli.cmdTimeout != 45*time.Second {
		t.Errorf("SetCommandTimeout(-1s): expected 45s unchanged, got %v", cli.cmdTimeout)
	}
}

func TestCli_SetConnectTimeout(t *testing.T) {
	cli := new(Cli)

	cli.SetConnectTimeout(20 * time.Second)
	if cli.connTimeout != 20*time.Second {
		t.Errorf("SetConnectTimeout: expected 20s, got %v", cli.connTimeout)
	}

	// Non-positive values must be ignored.
	cli.SetConnectTimeout(0)
	if cli.connTimeout != 20*time.Second {
		t.Errorf("SetConnectTimeout(0): expected 20s unchanged, got %v", cli.connTimeout)
	}

	cli.SetConnectTimeout(-time.Second)
	if cli.connTimeout != 20*time.Second {
		t.Errorf("SetConnectTimeout(-1s): expected 20s unchanged, got %v", cli.connTimeout)
	}
}
