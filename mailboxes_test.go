package cgpcli

import (
	"os"
	"strings"
	"testing"
)

func TestMailboxes(t *testing.T) {
        cli := getTestCli(t)
        defer cli.Logout()

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
		if err := cli.SetMailboxACL(testAcc, mbox, "", map[string]string{"anyone": "l"}); err != nil {
			t.Errorf("Initial ACL set failed: %v", err)
		}
		// Инкрементальное добавление
		if err := cli.SetMailboxACL(testAcc, mbox, "", map[string]string{"anyone": "+r"}); err != nil {
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
}

func TestMailboxesRO(t *testing.T) {
	// Проверка чтения из Production (RO)
	host := os.Getenv("CGPHOST_RO")
	if host == "" {
		t.Skip("RO environment (CGPHOST_RO) not set")
	}

	cli, err := New(host, os.Getenv("CGPUSER_RO"), os.Getenv("CGPPASS_RO"), APOP, false)
	if err != nil {
		t.Fatalf("RO Connect failed: %v", err)
	}
	defer cli.Logout()

	target := os.Getenv("USER_RO")

	t.Run("ListAndInfo", func(t *testing.T) {
		mboxes, err := cli.ListMailboxes(target, "*", "")
		if err != nil {
			t.Errorf("ListMailboxes failed: %v", err)
		}
		if len(mboxes) == 0 {
			t.Log("RO account has no mailboxes to list")
		}
	})
}
