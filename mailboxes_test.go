package cgpcli

import (
	"strings"
	"testing"
)

func TestMailboxes(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Используем существующий тестовый домен
	testAcc := "testuser1@test.domain.name"

	t.Run("FullLifecycle", func(t *testing.T) {
		mbox := "GoSDKTestFolder"
		newMbox := "GoSDKTestRenamed"

		// Предварительная очистка на случай прерванных ранее тестов
		_ = cli.DeleteMailbox(testAcc, mbox, true, "")
		_ = cli.DeleteMailbox(testAcc, newMbox, true, "")

		// 1. CreateMailbox
		if err := cli.CreateMailbox(testAcc, mbox, "", ""); err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// 2. SetMailboxClass
		if err := cli.SetMailboxClass(testAcc, mbox, "IPF.Note", ""); err != nil {
			t.Errorf("SetClass failed: %v", err)
		}

		// 3. SetMailboxACL & GetMailboxRights
		// Устанавливаем базовое право
		if err := cli.SetMailboxACL(testAcc, mbox, map[string]string{"anyone": "l"}, ""); err != nil {
			t.Errorf("Initial ACL set failed: %v", err)
		}
		// Инкрементальное добавление
		if err := cli.SetMailboxACL(testAcc, mbox, map[string]string{"anyone": "+r"}, ""); err != nil {
			t.Errorf("Incremental ACL update failed: %v", err)
		}

		rights, err := cli.GetMailboxRights(testAcc, mbox, "anyone")
		if err != nil {
			t.Errorf("GetRights failed: %v", err)
		}
		if !strings.Contains(rights, "r") {
			t.Errorf("ACL incremental update failed, expected 'r' in rights, got: %s", rights)
		}

		// 4. RenameMailbox
		if err := cli.RenameMailbox(testAcc, mbox, newMbox, false, ""); err != nil {
			t.Errorf("Rename failed: %v", err)
		}

		// 5. GetMailboxInfo
		info, err := cli.GetMailboxInfo(testAcc, newMbox, "")
		if err != nil {
			t.Errorf("GetInfo failed: %v", err)
		}
		if _, ok := info["Messages"]; !ok {
			t.Errorf("Mailbox info does not contain 'Messages' key")
		}

		// 6. DeleteMailbox
		if err := cli.DeleteMailbox(testAcc, newMbox, false, ""); err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})

	t.Run("SubscriptionsAndAliases", func(t *testing.T) {
		// Тест подписок
		targetSubs := []string{"INBOX", "Sent"}
		if err := cli.SetMailboxSubscription(testAcc, targetSubs); err != nil {
			t.Errorf("SetSubscription failed: %v", err)
		}

		subs, err := cli.GetMailboxSubscription(testAcc)
		if err != nil {
			t.Errorf("GetSubscription failed: %v", err)
		}
		if len(subs) < 1 {
			t.Errorf("Subscription list is empty, expected at least INBOX")
		}

		// Тест алиасов
		targetAliases := map[string]string{"OldInbox": "INBOX"}
		if err := cli.SetMailboxAliases(testAcc, targetAliases); err != nil {
			t.Errorf("SetAliases failed: %v", err)
		}

		aliases, err := cli.GetMailboxAliases(testAcc)
		if err != nil {
			t.Errorf("GetAliases failed: %v", err)
		}
		if aliases["OldInbox"] != "INBOX" {
			t.Errorf("Alias mismatch: expected INBOX, got %s", aliases["OldInbox"])
		}

		// Очистка алиасов (полная перезапись пустым словарем)
		_ = cli.SetMailboxAliases(testAcc, map[string]string{})
	})

	t.Run("ACLAndListing", func(t *testing.T) {
		mbox := "ACLTestFolder"
		_ = cli.CreateMailbox(testAcc, mbox, "", "")
		defer cli.DeleteMailbox(testAcc, mbox, false, "")

		// 1. Устанавливаем ACL
		targetACL := map[string]string{
			"anyone":     "l",
			"postmaster": "lrswipdka",
		}
		if err := cli.SetMailboxACL(testAcc, mbox, targetACL, ""); err != nil {
			t.Fatalf("SetMailboxACL failed: %v", err)
		}

		// 2. Проверяем GetMailboxACL
		acl, err := cli.GetMailboxACL(testAcc, mbox, "")
		if err != nil {
			t.Fatalf("GetMailboxACL failed: %v", err)
		}

		if acl["anyone"] != "l" {
			t.Errorf("ACL mismatch for 'anyone': expected 'l', got '%s'", acl["anyone"])
		}
		if !strings.Contains(acl["postmaster"], "s") {
			t.Errorf("ACL for postmaster missing 's' right: %s", acl["postmaster"])
		}

		// 3. Вызываем ListMailboxes для покрытия
		// Этот вызов закроет оставшиеся ветки в ListMailboxes (с фильтром)
		list, err := cli.ListMailboxes(testAcc, mbox, "")
		if err != nil {
			t.Errorf("ListMailboxes with filter failed: %v", err)
		}
		if _, ok := list[mbox]; !ok && len(list) > 0 {
			t.Logf("Note: %s not found in filtered list, but call was successful", mbox)
		}
	})
}
