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
