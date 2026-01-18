package cgpcli

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func setupSandbox(t *testing.T) *Cli {
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

	return cli
}

func TestIntegration_Directory_AccessRights(t *testing.T) {
	cli := setupSandbox(t)
	defer cli.Logout()

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
	cli := setupSandbox(t)
	defer cli.Logout()

	cli.SetDebug(true)

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
