package cgpcli

import (
	"fmt"
	"testing"
	"time"
)

func TestIntegration_Directory_AccessRights(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	t.Run("Verify Structure and Expert Parsing", func(t *testing.T) {
		rights, err := cli.GetDirectoryAccessRights(false)
		if err != nil {
			t.Fatalf("Expert parsing failed or server returned unexpected format: %v", err)
		}

		if len(rights) == 0 {
			t.Log("Directory Access Rights list is empty on this server")
			return
		}

		// Проверяем наличие эталонного правила ReadOptions
		var found bool
		for _, r := range rights {
			if r.Name == "ReadOptions" {
				found = true
				if r.Type == "" {
					t.Error("Rule 'ReadOptions' has empty Type")
				}
				if len(r.ReadAttrs) == 0 {
					t.Error("Rule 'ReadOptions' has empty ReadAttrs (expected at least '*')")
				}
				break
			}
		}

		if !found {
			t.Log("Rule 'ReadOptions' not found, verifying general data integrity")
			r := rights[0]
			if r.Name == "" {
				t.Error("First rule has empty name")
			}
		}
	})
}

func TestIntegration_Directory_Units(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Используем уникальное имя для каждого запуска теста,
	// чтобы избежать "хвостов" от предыдущих падений.
	testUnit := fmt.Sprintf("GoUnit_%d", time.Now().UnixNano()%10000)
	mount := fmt.Sprintf("tmp/go_test_%s", testUnit)

	t.Run("Lifecycle", func(t *testing.T) {
		// Очистка перед тестом (на всякий случай)
		_ = cli.DeleteDirectoryUnit(testUnit, false)

		// Создание
		err := cli.CreateDirectoryUnit(testUnit, mount, false, false)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// Проверка
		units, err := cli.ListDirectoryUnits(false)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if _, ok := units[testUnit]; !ok {
			t.Errorf("Unit %s not found in list. Got: %v", testUnit, units)
		}

		// Удаление
		if err := cli.DeleteDirectoryUnit(testUnit, false); err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})
}

func TestIntegration_Directory_FullMethods(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	testUnit := fmt.Sprintf("GoUnitFull_%d", time.Now().UnixNano()%1000)
	mount := fmt.Sprintf("tmp/full_%s", testUnit)
	newMount := fmt.Sprintf("tmp/relocated_%s", testUnit)

	t.Run("UnitSettings_and_Relocate", func(t *testing.T) {
		_ = cli.CreateDirectoryUnit(testUnit, mount, false, false)
		defer cli.DeleteDirectoryUnit(testUnit, false)

		// 1. Читаем настройки. Ожидаем, что парсер уже превратил YES/NO в bool.
		settings, err := cli.GetDirectoryUnit(testUnit, false)
		if err != nil {
			t.Fatalf("GetDirectoryUnit failed: %v", err)
		}

		// 2. Инвертируем флаг. Если здесь не bool — тест упадет с паникой,
		// и это правильно, значит парсер работает не так, как мы ожидаем.
		key := "ComposeDisplayName"
		newVal := !settings[key].(bool)

		settings[key] = newVal

		// 3. Сохраняем. Маршалер должен превратить bool обратно в YES/NO.
		if err := cli.SetDirectoryUnit(testUnit, false, settings); err != nil {
			t.Errorf("SetDirectoryUnit failed: %v", err)
		}

		// 4. Проверяем результат.
		updated, _ := cli.GetDirectoryUnit(testUnit, false)
		if updated[key] != newVal {
			t.Errorf("Value for %s mismatch: expected %v (type %T), got %v (type %T)",
				key, newVal, newVal, updated[key], updated[key])
		}

		// 5. Relocate
		if err := cli.RelocateDirectoryUnit(testUnit, newMount, false); err != nil {
			t.Fatalf("RelocateDirectoryUnit failed: %v", err)
		}
	})

	t.Run("AccessRights_ReadWrite", func(t *testing.T) {
		// Читаем текущие права
		original, err := cli.GetDirectoryAccessRights(false)
		if err != nil {
			t.Fatalf("Initial GetDirectoryAccessRights failed: %v", err)
		}

		// Создаем временное правило (копируем структуру ваших DirAccessRight)
		testRule := DirAccessRight{
			Name:      "GoTmpRule",
			Target:    "*",
			BindDN:    "*",
			Type:      "allow",
			ReadAttrs: []string{"*"},
		}

		// 1. Тестируем SetDirectoryAccessRights
		newRights := append(original, testRule)
		if err := cli.SetDirectoryAccessRights(false, newRights); err != nil {
			t.Fatalf("SetDirectoryAccessRights failed: %v", err)
		}

		// 2. Проверяем наличие
		current, _ := cli.GetDirectoryAccessRights(false)
		found := false
		for _, r := range current {
			if r.Name == "GoTmpRule" {
				found = true
				break
			}
		}

		// 3. Восстанавливаем оригинальные права (важно!)
		if err := cli.SetDirectoryAccessRights(false, original); err != nil {
			t.Errorf("Failed to restore original ACL: %v", err)
		}

		if !found {
			t.Error("SetDirectoryAccessRights succeeded but rule not found in subsequent Get")
		}
	})
}
