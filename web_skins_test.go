package cgpcli

import (
	"bytes"
	"testing"
)

func TestWebSkins_Integration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	const (
		testDomain  = "test.domain.name"
		namedSkinID = "sdk_named_skin"
		testFile    = "test.txt"
	)
	testContent := []byte("Hello CommuniGate!")

	// --- SECTION: DOMAIN SKINS ---
	t.Run("DomainSkins", func(t *testing.T) {
		// 1. Unnamed Skin
		t.Run("Unnamed", func(t *testing.T) {
			_ = cli.DeleteDomainSkin(testDomain, "") // Clean
			if err := cli.CreateDomainSkin(testDomain, ""); err != nil {
				t.Fatalf("Failed to create unnamed domain skin: %v", err)
			}

			if err := cli.StoreDomainSkinFile(testDomain, "", testFile, testContent); err != nil {
				t.Errorf("Store to unnamed failed: %v", err)
			}

			if err := cli.StoreDomainSkinFile(testDomain, "", testFile, nil); err != nil {
				t.Errorf("Store to unnamed failed: %v", err)
			}

			_ = cli.DeleteDomainSkin(testDomain, "")
		})

		// 2. Named Skin (requires unnamed to exist)
		t.Run("Named", func(t *testing.T) {
			_ = cli.CreateDomainSkin(testDomain, "") // Ensure base exists
			defer cli.DeleteDomainSkin(testDomain, "")

			if err := cli.CreateDomainSkin(testDomain, namedSkinID); err != nil {
				t.Fatalf("Failed to create named domain skin: %v", err)
			}

			// Операция Store/Read
			if err := cli.StoreDomainSkinFile(testDomain, namedSkinID, testFile, testContent); err != nil {
				t.Fatalf("Store failed: %v", err)
			}

			data, _, err := cli.ReadDomainSkinFile(testDomain, namedSkinID, testFile)
			if err != nil || !bytes.Equal(data, testContent) {
				t.Errorf("Read mismatch or error: %v", err)
			}

			if err := cli.StoreDomainSkinFile(testDomain, namedSkinID, testFile, nil); err != nil {
				t.Fatalf("Store failed: %v", err)
			}

			_ = cli.DeleteDomainSkin(testDomain, namedSkinID)
		})
	})

	// --- SECTION: SERVER SKINS ---
	t.Run("ServerSkins", func(t *testing.T) {
		t.Run("Named", func(t *testing.T) {
			_ = cli.DeleteServerSkin(namedSkinID)
			if err := cli.CreateServerSkin(namedSkinID); err != nil {
				t.Fatalf("Failed to create named server skin: %v", err)
			}

			if err := cli.StoreServerSkinFile(namedSkinID, testFile, testContent); err != nil {
				t.Fatalf("Store failed: %v", err)
			}

			if err := cli.StoreServerSkinFile(namedSkinID, testFile, nil); err != nil {
				t.Fatalf("Store failed: %v", err)
			}

			// Проверка Rename (специфично для Server)
			const newName = namedSkinID + "_renamed"
			_ = cli.DeleteServerSkin(newName)
			if err := cli.RenameServerSkin(namedSkinID, newName); err != nil {
				t.Errorf("Rename failed: %v", err)
			} else {
				_ = cli.DeleteServerSkin(newName)
			}
		})
	})
}
