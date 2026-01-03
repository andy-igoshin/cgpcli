package cgpcli

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDomainSetAdministration(t *testing.T) {
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

	t.Run("ListDomains", func(t *testing.T) {
		domains, err := cli.ListDomains()
		if err != nil {
			t.Errorf("ListDomains failed: %v", err)
		}
		t.Logf("Found %d domains", len(domains))
	})

	t.Run("DomainLifecycle", func(t *testing.T) {
		tempName := fmt.Sprintf("set-test-%d.local", time.Now().Unix())
		newName := "renamed-" + tempName

		// Создание
		if err := cli.CreateDomain(tempName, false, "", nil); err != nil {
			t.Fatalf("CreateDomain failed: %v", err)
		}

		// Переименование
		if err := cli.RenameDomain(tempName, newName, ""); err != nil {
			t.Errorf("RenameDomain failed: %v", err)
		}

		// Удаление
		if err := cli.DeleteDomain(newName, true); err != nil {
			t.Errorf("DeleteDomain failed: %v", err)
		}
	})

	t.Run("ServerTelnums", func(t *testing.T) {
		_, err := cli.ListServerTelnums("", 10)
		if err != nil {
			t.Logf("ListServerTelnums (license check): %v", err)
		}
	})
}
