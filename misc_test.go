package cgpcli

import (
	"net"
	"os"
	"testing"
	"time"
)

func TestMiscellaneousLive(t *testing.T) {
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

	// --- Информационные команды ---
	t.Run("SystemInfo", func(t *testing.T) {
		version, err := cli.GetVersion()
		if err != nil {
			t.Errorf("GetVersion failed: %v", err)
		}
		t.Logf("Server Version: %s", version)

		sysTime, err := cli.GetCurrentTime()
		if err != nil {
			t.Errorf("GetCurrentTime failed: %v", err)
		}
		t.Logf("Server Time: %v", sysTime)

		// Проверка Enum для GetSystemInfo
		info, err := cli.GetSystemInfo(SysOS)
		if err != nil {
			t.Errorf("GetSystemInfo(SysOS) failed: %v", err)
		}
		t.Logf("Server OS: %v", info)
	})

	// --- Работа с временными IP (Blacklist/Unblock) ---
	t.Run("TempIPOperations", func(t *testing.T) {
		testIP := net.ParseIP("192.168.250.250")

		// 1. Сохраняем текущий блэклист для восстановления
		// Так как GetTempBlacklistedIPs теперь возвращает []TempIP,
		// нам нужно извлечь только IP для последующего восстановления через Set...
		oldData, err := cli.GetTempBlacklistedIPs()
		if err != nil {
			t.Fatalf("Failed to backup blacklist: %v", err)
		}

		oldIPs := make([]net.IP, len(oldData))
		for i, v := range oldData {
			oldIPs[i] = v.IP
		}

		// 2. Тестируем точечную блокировку
		duration := 5 * time.Minute
		if err := cli.TempBlacklistIP(testIP, duration); err != nil {
			t.Errorf("TempBlacklistIP failed: %v", err)
		}

		// 3. Проверяем наличие в списке и работу парсера
		list, err := cli.GetTempBlacklistedIPs()
		if err != nil {
			t.Errorf("GetTempBlacklistedIPs failed: %v", err)
		}

		found := false
		for _, item := range list {
			if item.IP.Equal(testIP) {
				found = true
				if item.TimeLeft == 0 && !item.Permanent {
					t.Error("TimeLeft should be greater than 0")
				}
				break
			}
		}
		if !found {
			t.Error("Test IP not found in blacklist")
		}

		// 4. Тестируем удаление через DELETE (timeout <= 0)
		if err := cli.TempBlacklistIP(testIP, 0); err != nil {
			t.Errorf("TempBlacklistIP DELETE failed: %v", err)
		}

		// 5. Восстанавливаем исходное состояние
		defer func() {
			if err := cli.SetTempBlacklistedIPs(oldIPs); err != nil {
				t.Errorf("Cleanup: failed to restore blacklist: %v", err)
			}
		}()
	})

	// --- Очереди и Логи ---
	t.Run("QueuesAndLogs", func(t *testing.T) {
		// NOOP для проверки связи
		if err := cli.Noop(); err != nil {
			t.Errorf("Noop failed: %v", err)
		}

		// WriteLog (уровень 4 - Low Level)
		if err := cli.WriteLog(4, "CGP-API-TEST: Testing WriteLog from Go SDK"); err != nil {
			t.Errorf("WriteLog failed: %v", err)
		}

		// GetMessageQueueInfo для несуществующей очереди (проверка на возврат пустой мапы)
		info, err := cli.GetMessageQueueInfo("SMTP", "nonexistent.vsu.ru")
		if err != nil {
			t.Errorf("GetMessageQueueInfo failed: %v", err)
		}
		if info == nil {
			t.Error("Expected empty map, got nil")
		}
	})

	// --- Тест производительности (опционально) ---
	t.Run("PerformanceTest", func(t *testing.T) {
		perf, err := cli.TestLoop(1)
		if err != nil {
			t.Errorf("TestLoop failed: %v", err)
		}
		t.Logf("CLI Performance (1s loop): %d units", perf)
	})
}
