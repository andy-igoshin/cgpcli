package cgpcli

import (
	"os"
	"testing"
)

func TestRunScript_Integration(t *testing.T) {
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

	const (
		account = "testuser1@test.domain.name"
		prog    = "sdk_script"
	)

	t.Run("RunScript", func(t *testing.T) {
		if _, err := cli.RunScript(account, prog, "", nil); err != nil {
			t.Fatalf("Failed to run script: %v", err)
		}
	})
}
