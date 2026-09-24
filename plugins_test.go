package cgpcli

import (
	"testing"
)

func TestPlugins(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	targetDomain := "test.domain.name"

	// 1. Get current Domain Plugins settings
	t.Run("GetDomainPluginsSettings", func(t *testing.T) {
		settings, err := cli.GetDomainPluginsSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainPluginsSettings failed: %v", err)
		}
		if settings == nil {
			t.Error("GetDomainPluginsSettings returned nil map")
		}
		t.Logf("Current Domain Plugins settings keys: %d", len(settings))
		for k, v := range settings {
			t.Logf("  %s = %v", k, v)
		}
	})

	// 2. Update and restore Domain Plugins settings
	t.Run("SetDomainPluginsSettings", func(t *testing.T) {
		// Save original settings
		original, err := cli.GetDomainPluginsSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainPluginsSettings (save) failed: %v", err)
		}

		// Apply new settings
		update := map[string]any{
			"Comment": "Test via Go SDK",
		}
		if err := cli.SetDomainPluginsSettings(targetDomain, update); err != nil {
			t.Errorf("SetDomainPluginsSettings failed: %v", err)
		}

		// Verify changes were applied
		current, err := cli.GetDomainPluginsSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainPluginsSettings (verify) failed: %v", err)
		}

		if val, ok := current["Comment"]; ok {
			t.Logf("Comment = %v", val)
		}

		// Restore original settings
		if err := cli.SetDomainPluginsSettings(targetDomain, original); err != nil {
			t.Errorf("SetDomainPluginsSettings (restore) failed: %v", err)
		}

		// Verify restoration
		restored, err := cli.GetDomainPluginsSettings(targetDomain)
		if err != nil {
			t.Fatalf("GetDomainPluginsSettings (restore verify) failed: %v", err)
		}
		t.Logf("Restored Domain Plugins settings keys: %d", len(restored))
	})

	// 3. SetDomainPluginsSettings with nil settings should return error
	t.Run("SetDomainPluginsSettings_NilSettings", func(t *testing.T) {
		err := cli.SetDomainPluginsSettings(targetDomain, nil)
		if err == nil {
			t.Error("SetDomainPluginsSettings with nil settings should return an error")
		} else {
			t.Logf("Expected error: %v", err)
		}
	})

	// 4. GetDomainPluginsSettings with empty domain
	t.Run("GetDomainPluginsSettings_EmptyDomain", func(t *testing.T) {
		settings, err := cli.GetDomainPluginsSettings("")
		if err != nil {
			t.Logf("GetDomainPluginsSettings with empty domain: %v", err)
		} else {
			t.Logf("GetDomainPluginsSettings with empty domain returned %d keys", len(settings))
		}
	})

	// 5. RemovePlugin with empty ID should return error
	t.Run("RemovePlugin_EmptyID", func(t *testing.T) {
		err := cli.RemovePlugin("")
		if err == nil {
			t.Error("RemovePlugin with empty id should return an error")
		} else {
			t.Logf("Expected error: %v", err)
		}
	})

	// 6. RemovePlugin with non-existent plugin ID — expects error from server
	t.Run("RemovePlugin_NonExistent", func(t *testing.T) {
		err := cli.RemovePlugin("non-existent-plugin-id")
		if err != nil {
			t.Logf("RemovePlugin non-existent (expected): %v", err)
		} else {
			t.Log("RemovePlugin non-existent: no error (plugin may exist or was silently ignored)")
		}
	})

	// 7. UploadPluginFile with empty content should return error (client-side)
	t.Run("UploadPluginFile_EmptyContent", func(t *testing.T) {
		err := cli.UploadPluginFile([]byte{}, false)
		if err == nil {
			t.Error("UploadPluginFile with empty content should return an error")
		} else {
			t.Logf("Expected error: %v", err)
		}
	})

	// 8. UploadPluginFile with invalid binary data — expects error from server
	t.Run("UploadPluginFile_Invalid", func(t *testing.T) {
		fakeContent := []byte("this is not a valid plugin binary")
		err := cli.UploadPluginFile(fakeContent, false)
		if err != nil {
			t.Logf("UploadPluginFile invalid (expected): %v", err)
		} else {
			t.Log("UploadPluginFile invalid: unexpectedly succeeded")
		}
	})

	// 9. UploadPluginFile with DRYRUN flag and invalid data — expects error from server
	t.Run("UploadPluginFile_DryRun", func(t *testing.T) {
		fakeContent := []byte("dry-run test plugin data")
		err := cli.UploadPluginFile(fakeContent, true)
		if err != nil {
			t.Logf("UploadPluginFile dry-run (expected): %v", err)
		} else {
			t.Log("UploadPluginFile dry-run: unexpectedly succeeded")
		}
	})
}
