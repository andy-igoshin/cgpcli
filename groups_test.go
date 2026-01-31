package cgpcli

import (
	"fmt"
	"testing"
	"time"
)

func TestGroups(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	domain := "test.domain.name"
	// Уникальное имя группы для каждого запуска
	ts := time.Now().Unix()
	groupName := fmt.Sprintf("test-grp-%d@%s", ts, domain)
	renamedGroupName := fmt.Sprintf("renamed-grp-%d@%s", ts, domain)

	t.Run("Lifecycle", func(t *testing.T) {
		// 1. Create
		err := cli.CreateGroup(groupName, map[string]any{"RealName": "GoTest"})
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// 2. Exists
		ex, err := cli.IsGroupExists(groupName)
		if err != nil {
			t.Errorf("IsGroupExists error: %v", err)
		}
		if !ex {
			t.Error("Group should exist")
		}

		// 3. Set/Get
		err = cli.SetGroup(groupName, map[string]any{"Comment": "Test"}, []string{"postmaster@" + domain})
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		res, err := cli.GetGroup(groupName)
		if err != nil {
			t.Fatalf("GetGroup failed: %v", err)
		}
		if len(res.Members) == 0 {
			t.Error("Members should not be empty")
		}

		// 4. Rename
		err = cli.RenameGroup(groupName, renamedGroupName)
		if err != nil {
			t.Errorf("Rename failed: %v", err)
		}

		// 5. Delete
		err = cli.DeleteGroup(renamedGroupName)
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}

		ex, _ = cli.IsGroupExists(renamedGroupName)
		if ex {
			t.Error("Group should be deleted")
		}
	})
}
