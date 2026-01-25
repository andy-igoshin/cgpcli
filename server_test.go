package cgpcli

import (
	"fmt"
	"testing"
	"time"
)

// 1. Тестирование работы с таблицей маршрутизации (Router Table)
func TestServer_RouterTable(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	// 1. Получаем текущую таблицу
	originalTable, err := cli.GetRouterTable()
	if err != nil {
		t.Fatalf("Failed to get original router table: %v", err)
	}

	t.Run("UpdateAndVerify", func(t *testing.T) {
		testRule := fmt.Sprintf("test-rule-%d.local = trash", time.Now().Unix())

		// 2. Создаем рабочую копию
		workTable := make(RouterList, len(originalTable))
		copy(workTable, originalTable)

		// 3. Добавляем запись (согласно вашей реализации: rule, comment, alignPos)
		workTable.Add(testRule, "Go SDK Test", 40)

		// 4. Записываем на сервер
		if err := cli.SetRouterTable(workTable); err != nil {
			t.Fatalf("SetRouterTable failed: %v", err)
		}

		// 5. Читаем обратно и проверяем
		updated, err := cli.GetRouterTable()
		if err != nil {
			t.Fatalf("GetRouterTable after update failed: %v", err)
		}

		found := false
		for _, entry := range updated {
			// Используем поле Rule, как в вашем router.go
			if entry.Rule == testRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("New rule %q not found in table after reload", testRule)
		}
	})

	// 6. Восстанавливаем оригинал
	_ = cli.SetRouterTable(originalTable)
}

// 2. Тестирование списков IP (на примере LANIPs)
func TestServer_IPLists(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	originalLAN, err := cli.GetLANIPs()
	if err != nil {
		t.Fatalf("GetLANIPs failed: %v", err)
	}

	t.Run("ModifyLANIPs", func(t *testing.T) {
		testIP := "192.168.254.254"

		// Создаем копию
		newList := make(IPList, len(originalLAN))
		copy(newList, originalLAN)

		// Ваша реализация IPList.Add(input, comment)
		err := newList.Add(testIP, "SDK Temp Entry")
		if err != nil {
			t.Fatalf("IPList.Add failed: %v", err)
		}

		if err := cli.SetLANIPs(newList); err != nil {
			t.Fatalf("SetLANIPs failed: %v", err)
		}

		verified, err := cli.GetLANIPs()
		if err != nil {
			t.Fatalf("GetLANIPs verification failed: %v", err)
		}

		found := false
		for _, entry := range verified {
			// Используем вашу функцию dataToString для сравнения значений
			if dataToString(entry.Value) == testIP {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Test IP %s was not saved in LANIPs", testIP)
		}
	})

	// Возвращаем как было
	_ = cli.SetLANIPs(originalLAN)
}

// 3. Тестирование Banned Lines (Generic Map/Slice)
func TestServer_BannedLines(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	origBanned, err := cli.GetBanned()
	if err != nil {
		t.Fatalf("GetBanned failed: %v", err)
	}

	t.Run("UpdateBanned", func(t *testing.T) {
		testString := fmt.Sprintf("X-Spam-Test-%d", time.Now().Unix())

		newBanned := &BannedLines{
			HeaderLines: append(origBanned.HeaderLines, testString),
			BodyLines:   origBanned.BodyLines,
		}

		if err := cli.SetBanned(newBanned); err != nil {
			t.Fatalf("SetBanned failed: %v", err)
		}

		updated, err := cli.GetBanned()
		if err != nil {
			t.Fatalf("GetBanned verify failed: %v", err)
		}

		found := false
		for _, s := range updated.HeaderLines {
			if s == testString {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Banned header %s was not found after update", testString)
		}
	})

	_ = cli.SetBanned(origBanned)
}

// 4. Тестирование системных настроек
func TestServer_Settings(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	t.Run("GetNetworkSettings", func(t *testing.T) {
		settings, err := cli.GetNetwork()
		if err != nil {
			t.Fatalf("GetNetwork failed: %v", err)
		}

		if len(settings) == 0 {
			t.Error("Network settings map is empty")
		}
	})
}

func TestServer_Rules(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	t.Run("ServerMailRules", func(t *testing.T) {
		// 1. Получаем текущие
		original, err := cli.GetServerMailRules()
		if err != nil {
			t.Fatalf("GetServerMailRules failed: %v", err)
		}

		// 2. Создаем тестовое правило
		testRule := MailRule{
			Priority: 5,
			Name:     "SDK-Test-Rule",
			Conditions: []any{
				[]any{"Subject", "is", "test-alert"},
			},
			Actions: []any{
				[]any{"Discard"},
			},
			Comment: "Temporary test rule",
		}

		// 3. Записываем (добавляем к существующим)
		newList := append(original, testRule)
		if err := cli.SetServerMailRules(newList); err != nil {
			t.Fatalf("SetServerMailRules failed: %v", err)
		}

		// 4. Проверяем наличие
		updated, _ := cli.GetServerMailRules()
		found := false
		for _, r := range updated {
			if r.Name == testRule.Name {
				found = true
				break
			}
		}
		if !found {
			t.Error("Test Mail Rule was not saved")
		}

		// 5. Откат
		_ = cli.SetServerMailRules(original)
	})

	t.Run("ServerSignalRules", func(t *testing.T) {
		original, _ := cli.GetServerSignalRules()

		testRule := SignalRule{
			Priority:   7,
			Delay:      10, // Должно упаковаться через (1000-10)*100 + 7
			Name:       "SDK-Signal-Test",
			Conditions: []any{[]any{"Service", "is", "SIP"}},
			Actions:    []any{[]any{"Discard"}},
		}

		if err := cli.SetServerSignalRules(append(original, testRule)); err != nil {
			t.Fatalf("SetServerSignalRules failed: %v", err)
		}

		updated, _ := cli.GetServerSignalRules()
		found := false
		for _, r := range updated {
			if r.Name == testRule.Name {
				// Проверяем, что распаковка приоритета сработала верно
				if r.Priority != testRule.Priority || r.Delay != testRule.Delay {
					t.Errorf("Signal priority/delay mismatch: got p:%d d:%d, want p:%d d:%d",
						r.Priority, r.Delay, testRule.Priority, testRule.Delay)
				}
				found = true
				break
			}
		}
		if !found {
			t.Error("Test Signal Rule was not saved")
		}

		_ = cli.SetServerSignalRules(original)
	})
}

// 5. Тестирование системных настроек (DNR, Logs, Queue, и т.д.)
func TestServer_SystemSettings(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	// Используем таблицу для проверки нескольких MapAny функций
	targets := []struct {
		name string
		fn   func() (map[string]any, error)
	}{
		{"DNR", cli.GetDNRSettings},
		{"Log", cli.GetLogSettings},
		{"Queue", cli.GetQueueSettings},
		{"Session", cli.GetSessionSettings},
		{"Network", cli.GetNetwork},
	}

	for _, tt := range targets {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.fn()
			if err != nil {
				t.Logf("Skipping or info: %s failed (possibly cluster-only or permissions): %v", tt.name, err)
				return
			}
			if len(res) == 0 {
				t.Errorf("%s returned empty map", tt.name)
			}
		})
	}
}

// 6. Тестирование Модулей
func TestServer_Modules(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	t.Run("ListAndGet", func(t *testing.T) {
		modules, err := cli.ListModules()
		if err != nil {
			t.Fatalf("ListModules failed: %v", err)
		}
		if len(modules) == 0 {
			t.Fatal("No modules found")
		}

		// Берем первый попавшийся модуль и запрашиваем его настройки
		moduleName := modules[0]
		settings, err := cli.GetModule(moduleName)
		if err != nil {
			t.Fatalf("GetModule(%s) failed: %v", moduleName, err)
		}
		if len(settings) == 0 {
			t.Logf("Warning: Module %s returned empty settings", moduleName)
		}
	})
}

// 7. Тестирование Route и IPState
func TestServer_Utility(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	t.Run("Route", func(t *testing.T) {
		// Простой тест маршрутизации
		res, err := cli.Route("postmaster", "")
		if err != nil {
			t.Fatalf("Route failed: %v", err)
		}
		if len(res) == 0 {
			t.Error("Route returned empty slice")
		}
	})

	t.Run("GetIPState", func(t *testing.T) {
		// Проверяем состояние localhost
		res, err := cli.GetIPState("127.0.0.1", false)
		if err != nil {
			t.Fatalf("GetIPState failed: %v", err)
		}
		t.Logf("IP State of 127.0.0.1: %s", res)
	})
}

func TestServer_UpdateSettings(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	t.Run("UpdateLogSettings", func(t *testing.T) {
		// 1. Получаем текущие настройки логов
		orig, err := cli.GetLogSettings()
		if err != nil {
			t.Fatalf("Failed to get log settings: %v", err)
		}

		// 2. Меняем какой-нибудь безобидный параметр, например уровень детализации
		// В CGP это обычно ключи типа "STDOUT", "LogLevel" и т.д.
		// Чтобы не сломать продакшен, просто перепишем то же самое значение
		err = cli.UpdateLogSettings(orig)
		if err != nil {
			t.Fatalf("UpdateLogSettings failed: %v", err)
		}
		t.Log("Roundtrip UpdateLogSettings successful")
	})
}
