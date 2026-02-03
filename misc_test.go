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
	t.Run("TempIPOperations", func(t *testing.T) {
		testIP := "192.168.250.250"

		// 1. Backup
		oldList, err := cli.GetTempBlacklistedIPs()
		if err != nil {
			t.Fatalf("Step 1 (Backup) failed: %v", err)
		}
		if oldList == nil {
			t.Fatal("Step 1: GetTempBlacklistedIPs returned nil object")
		}
		oldIPs := oldList.IPs()
		t.Logf("Backup done. Current blacklist size: %d", oldList.Len())

		// 2. Add
		duration := 5 * time.Minute
		ipObj := net.ParseIP(testIP)
		if err := cli.TempBlacklistIP(ipObj, duration); err != nil {
			t.Fatalf("Step 2 (Add) failed: %v", err)
		}

		// 3. Verify
		newList, err := cli.GetTempBlacklistedIPs()
		if err != nil {
			t.Fatalf("Step 3 (Verify) failed: %v", err)
		}

		found, err := newList.Exists(testIP)
		if err != nil {
			t.Fatalf("Step 3 (Exists) errored: %v", err)
		}
		if !found {
			t.Error("Step 3: Test IP not found in blacklist after addition")
		}

		// 4. Delete
		if err := cli.TempBlacklistIP(ipObj, 0); err != nil {
			t.Errorf("Step 4 (Delete) failed: %v", err)
		}

		// 5. Restore
		if err := cli.SetTempBlacklistedIPs(oldIPs); err != nil {
			t.Errorf("Step 5 (Restore) failed: %v", err)
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
			_ = cli.SetTempUnblockableIPs(oldUnblock.IPs())
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
