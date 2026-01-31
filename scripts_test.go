package cgpcli

import (
	"testing"
)

func TestRunScript_Integration(t *testing.T) {
        cli := getTestCli(t)
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
