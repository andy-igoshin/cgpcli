package cgpcli

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestForwarders(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	domain := "test.domain.name"
	ts := time.Now().Unix()

	// Локальные имена для проверок
	fwLocal := fmt.Sprintf("fw-%d", ts)
	renameLocal := fmt.Sprintf("renamed-fw-%d", ts)

	// Полные имена для команд
	fwFull := fwLocal + "@" + domain
	renameFull := renameLocal + "@" + domain

	target := "backup@external.com"

	t.Run("FullLifecycle", func(t *testing.T) {
		// 1. Create
		if err := cli.CreateForwarder(fwFull, target); err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// 2. Get
		res, err := cli.GetForwarder(fwFull)
		if err != nil || res != target {
			t.Errorf("Get failed: got %q, error %v", res, err)
		}

		// 3. Find
		found, err := cli.FindForwarders(domain, target)
		if err != nil {
			t.Errorf("Find failed: %v", err)
		}

		isThere := false
		for _, f := range found {
			// Проверяем и полное совпадение, и локальное
			if f == fwFull || f == fwLocal {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("Forwarder %s (or %s) not found in %v by target %s", fwFull, fwLocal, found, target)
		}

		// 4. Rename
		if err := cli.RenameForwarder(fwFull, renameFull); err != nil {
			t.Errorf("Rename failed: %v", err)
		}

		// 5. Update
		newTarget := "newbackup@external.com"
		settings := map[string]any{
			"To":       newTarget,
			"RealName": "Backup Forwarder",
			"Comment":  "updated by test",
		}
		if err := cli.UpdateForwarder(renameFull, settings); err != nil {
			t.Errorf("Update failed: %v", err)
		}

		res, err = cli.GetForwarder(renameFull)
		if err != nil {
			t.Errorf("Get after Update failed: %v", err)
		}
		if !strings.Contains(res, newTarget) {
			t.Errorf("Get after Update: target %q not found in %q", newTarget, res)
		}

		// 6. List
		_, err = cli.ListForwarders("")
		if err != nil {
			t.Errorf("List failed: %v", err)
		}

		list, err := cli.ListForwarders(domain)
		if err != nil {
			t.Errorf("List failed: %v", err)
		}

		isRenamed := false
		for _, f := range list {
			if f == renameFull || f == renameLocal {
				isRenamed = true
				break
			}
		}
		if !isRenamed {
			t.Errorf("Renamed forwarder %s (or %s) not in list %v", renameFull, renameLocal, list)
		}

		// 7. Delete
		if err := cli.DeleteForwarder(renameFull); err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})
}
