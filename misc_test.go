package cgpcli

import (
	"net"
	"testing"
	"time"
)

func TestMiscellaneousLive(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

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
	t.Run("TempIPOperations_FullCycle", func(t *testing.T) {
		testIP1 := "192.168.250.250"
		testIP2 := "192.168.250.251"
		ip1 := net.ParseIP(testIP1)
		ip2 := net.ParseIP(testIP2)

		// 1. BACKUP
		originalList, _ := cli.GetTempBlacklistedIPs()

		// 2. ATOMIC ADD
		_ = cli.TempBlacklistIP(ip1, 10*time.Minute)

		// 3. BULK MODIFY
		currentList, _ := cli.GetTempBlacklistedIPs()
		_ = currentList.Add(ip2, 2*time.Minute, false)
		if err := cli.SetTempBlacklistedIPs(currentList); err != nil {
			t.Fatalf("Bulk Set failed: %v", err)
		}

		// 4. VERIFY (Проверяем только своих)
		verifiedList, _ := cli.GetTempBlacklistedIPs()
		for _, ip := range []string{testIP1, testIP2} {
			if found, _ := verifiedList.Exists(ip); !found {
				t.Errorf("IP %s missing after bulk operations", ip)
			}
		}

		// 5. CLEANUP & RESTORE
		// Чтобы реально удалить 251-й, нужно либо DELETE, либо знать,
		// что SET не удаляет. Удаляем явно то, что создали.
		_ = cli.TempBlacklistIP(ip1, 0)
		_ = cli.TempBlacklistIP(ip2, 0)

		// Возвращаем старые таймеры для тех, кто там был
		_ = cli.SetTempBlacklistedIPs(originalList)

		// 6. FINAL CHECK (Проверяем отсутствие нашего мусора)
		finalList, _ := cli.GetTempBlacklistedIPs()
		if found, _ := finalList.Exists(testIP2); found {
			t.Errorf("Cleanup failed: IP %s still exists", testIP2)
		}
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

		// GetMessageQueueInfo для несуществующей очереди
		info, err := cli.GetMessageQueueInfo("SMTP", "nonexistent.vsu.ru")
		if err != nil {
			t.Errorf("GetMessageQueueInfo failed: %v", err)
		}
		if info == nil {
			t.Error("Expected empty map (or initialized map), got nil")
		}
	})

	// --- Тест производительности ---
	t.Run("PerformanceTest", func(t *testing.T) {
		perf, err := cli.TestLoop(1)
		if err != nil {
			t.Errorf("TestLoop failed: %v", err)
		}
		t.Logf("CLI Performance (1s loop): %d units", perf)
	})
}

func TestMiscellaneous_Extended(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// 1. Команды эха и списков
	t.Run("EchoAndCommands", func(t *testing.T) {
		// Echo
		val := "Hello CGP"
		res, err := cli.Echo(val)
		if err != nil || res != val {
			t.Errorf("Echo failed: got %v, want %v", res, val)
		}

		// ListCLICommands
		cmds, err := cli.ListCLICommands()
		if err != nil || len(cmds) == 0 {
			t.Errorf("ListCLICommands failed or empty: %v", err)
		}
	})

	// 2. Управление временными Client/Unblockable IPs
	t.Run("TempIPsSpecial", func(t *testing.T) {
		testIP := net.ParseIP("127.0.0.1")

		// TempClientIPs
		_, _ = cli.GetTempClientIPs()

		// TempUnblockableIPs (Get/Set/Add/Delete)
		oldUnblock, _ := cli.GetTempUnblockableIPs()

		_ = cli.TempUnblockIP(testIP, 1*time.Minute)
		_ = cli.TempUnblockIP(testIP, 0) // Delete

		if oldUnblock != nil {
			_ = cli.SetTempUnblockableIPs(oldUnblock)
		}

		// ReportFailedLoginAddress
		_ = cli.ReportFailedLoginAddress(net.ParseIP("10.0.0.1"))
	})

	// 3. Диагностика и Трассировка
	t.Run("Diagnostics", func(t *testing.T) {
		//_ = cli.DumpAllObjects(false) // Пишет в лог сервера

		_ = cli.SetLogAll(true)
		_ = cli.SetLogAll(false)

		_ = cli.SetTrace(TraceFileIO, true)
		_ = cli.SetTrace(TraceFileIO, false)
	})

	// 4. Очереди и Кластер (методы с потенциальными ошибками логики)
	t.Run("QueuesAndCluster", func(t *testing.T) {
		// Эти вызовы скорее всего вернут ошибки (нет таких ID или не кластер),
		// но код методов в misc.go будет исполнен.

		_ = cli.RejectQueueMessage(999999, "Testing reject")
		_ = cli.RejectQueueMessages("nonexistent@domain.com", "Spam")
		_ = cli.ReleaseSMTPQueue("nonexistent.queue")

		// Кластерные команды
		_, _ = cli.GetCurrentController()
		_ = cli.ReconnectClusterAdmin()
	})
}
