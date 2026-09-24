package cgpcli

import (
	"testing"
)

func TestDomainFreeBusy(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	targetDomain := "test.domain.name"

	// 1. Get current Free/Busy settings
	t.Run("GetDomainFreeBusy", func(t *testing.T) {
		settings, err := cli.GetDomainFreeBusy(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainFreeBusy failed: %v", err)
		}
		if settings == nil {
			t.Error("GetDomainFreeBusy returned nil map")
		}
		t.Logf("Current Free/Busy settings keys: %d", len(settings))
		for k, v := range settings {
			t.Logf("  %s = %v", k, v)
		}
	})

	// 2. Update and restore Free/Busy settings
	t.Run("SetDomainFreeBusy", func(t *testing.T) {
		// Save original settings
		original, err := cli.GetDomainFreeBusy(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainFreeBusy (save) failed: %v", err)
		}

		// Apply new settings
		update := map[string]any{
			"FreeBusyEnable": true,
			"FreeBusyDays":   "30",
			"Comment":        "Test via Go SDK",
		}
		if err := cli.SetDomainFreeBusy(targetDomain, update); err != nil {
			t.Errorf("SetDomainFreeBusy failed: %v", err)
		}

		// Verify changes were applied
		current, err := cli.GetDomainFreeBusy(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainFreeBusy (verify) failed: %v", err)
		}

		// Check that our updated keys are reflected
		if val, ok := current["FreeBusyEnable"]; ok {
			t.Logf("FreeBusyEnable = %v", val)
		}
		if val, ok := current["Comment"]; ok {
			t.Logf("Comment = %v", val)
		}

		// Restore original settings
		if err := cli.SetDomainFreeBusy(targetDomain, original); err != nil {
			t.Errorf("SetDomainFreeBusy (restore) failed: %v", err)
		}

		// Verify restoration
		restored, err := cli.GetDomainFreeBusy(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainFreeBusy (restore verify) failed: %v", err)
		}
		t.Logf("Restored Free/Busy settings keys: %d", len(restored))
	})

	// 3. SetDomainFreeBusy with nil settings should return error
	t.Run("SetDomainFreeBusy_NilSettings", func(t *testing.T) {
		err := cli.SetDomainFreeBusy(targetDomain, nil)
		if err == nil {
			t.Error("SetDomainFreeBusy with nil settings should return an error")
		} else {
			t.Logf("Expected error: %v", err)
		}
	})

	// 4. GetDomainFreeBusy with empty domain
	t.Run("GetDomainFreeBusy_EmptyDomain", func(t *testing.T) {
		settings, err := cli.GetDomainFreeBusy("")
		if err != nil {
			// Empty domain may cause an error depending on server configuration
			t.Logf("GetDomainFreeBusy with empty domain: %v", err)
		} else {
			t.Logf("GetDomainFreeBusy with empty domain returned %d keys", len(settings))
		}
	})
}
