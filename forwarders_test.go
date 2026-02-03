package cgpcli

import (
	"fmt"
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

		// 5. List
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

		// 6. Delete
		if err := cli.DeleteForwarder(renameFull); err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})
}
