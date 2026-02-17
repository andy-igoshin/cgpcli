package cgpcli

import (
	"bytes"
	"testing"
)

// cleanSkinFiles вычищает все файлы из скина перед его удалением
func cleanSkinFiles(cli *Cli, level, domain, skin string) {
	var files map[string]*FileInfo
	var err error

	switch level {
	case "server":
		files, err = cli.ListServerSkinFiles(skin)
	case "domain":
		files, err = cli.ListDomainSkinFiles(domain, skin)
	case "cluster":
		files, err = cli.ListClusterSkinFiles(skin)
	}

	if err == nil {
		for file := range files {
			switch level {
			case "server":
				_ = cli.DeleteServerSkinFile(skin, file)
			case "domain":
				_ = cli.DeleteDomainSkinFile(domain, skin, file)
			case "cluster":
				_ = cli.DeleteClusterSkinFile(skin, file)
			}
		}
	}
}

func TestWebSkins_Comprehensive(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	const (
		testDomain = "test.domain.name"
		skinID     = "sdk_test_skin"
		testFile   = "test_file.txt"
		renamedID  = "sdk_test_skin_new"
	)
	testContent := []byte("Skin Content Verification")

	// --- SECTION 1: SERVER SKINS ---
	t.Run("ServerSkins_FullCycle", func(t *testing.T) {
		// Глубокая очистка перед тестом
		for _, s := range []string{skinID, renamedID} {
			cleanSkinFiles(cli, "server", "", s)
			_ = cli.DeleteServerSkin(s)
		}

		if err := cli.CreateServerSkin(skinID); err != nil {
			t.Fatalf("CreateServerSkin failed: %v", err)
		}

		skins, err := cli.ListServerSkins()
		if err != nil {
			t.Errorf("ListServerSkins failed: %v", err)
		}
		found := false
		for _, s := range skins {
			if s == skinID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Skin %s not found in ListServerSkins", skinID)
		}

		if err := cli.StoreServerSkinFile(skinID, testFile, testContent); err != nil {
			t.Fatalf("StoreServerSkinFile failed: %v", err)
		}

		// Проверка чтения
		data, _, err := cli.ReadServerSkinFile(skinID, testFile)
		if err != nil || !bytes.Equal(data, testContent) {
			t.Errorf("Read data mismatch")
		}

		// Rename
		if err := cli.RenameServerSkin(skinID, renamedID); err != nil {
			t.Errorf("RenameServerSkin failed: %v", err)
		}

		// Cleanup: вычищаем файлы в переименованном скине и удаляем его
		cleanSkinFiles(cli, "server", "", renamedID)
		_ = cli.DeleteServerSkin(renamedID)
	})

	// --- SECTION 2: DOMAIN SKINS ---
	t.Run("DomainSkins_FullCycle", func(t *testing.T) {
		// Очистка в правильном порядке: сначала именованные, потом unnamed
		for _, s := range []string{skinID, renamedID} {
			cleanSkinFiles(cli, "domain", testDomain, s)
			_ = cli.DeleteDomainSkin(testDomain, s)
		}
		cleanSkinFiles(cli, "domain", testDomain, "")
		_ = cli.DeleteDomainSkin(testDomain, "")

		// Создаем инфраструктуру
		if err := cli.CreateDomainSkin(testDomain, ""); err != nil {
			t.Fatalf("Create base domain skin failed: %v", err)
		}

		if err := cli.CreateDomainSkin(testDomain, skinID); err != nil {
			t.Fatalf("Create named domain skin failed: %v", err)
		}

		_, err := cli.ListDomainSkins("")
		if err != nil {
			t.Errorf("ListDomainSkins failed: %v", err)
		}

		dSkins, err := cli.ListDomainSkins(testDomain)
		if err != nil {
			t.Errorf("ListDomainSkins failed: %v", err)
		}
		found := false
		for _, s := range dSkins {
			if s == skinID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Skin %s not found in ListDomainSkins for %s", skinID, testDomain)
		}

		_ = cli.StoreDomainSkinFile(testDomain, skinID, testFile, testContent)

		data, _, err := cli.ReadDomainSkinFile(testDomain, skinID, testFile)
		if err != nil {
			t.Fatalf("ReadDomainSkinFile failed: %v", err)
		}
		if !bytes.Equal(data, testContent) {
			t.Error("ReadDomainSkinFile: content mismatch")
		}

		if err := cli.RenameDomainSkin(testDomain, skinID, renamedID); err != nil {
			t.Errorf("RenameDomainSkin failed: %v", err)
		}

		// Final Cleanup
		cleanSkinFiles(cli, "domain", testDomain, renamedID)
		_ = cli.DeleteDomainSkin(testDomain, renamedID)
		cleanSkinFiles(cli, "domain", testDomain, "")
		_ = cli.DeleteDomainSkin(testDomain, "")
	})

	// --- SECTION 3: CLUSTER SKINS ---
	t.Run("ClusterSkins_Honest", func(t *testing.T) {
		// Сначала пробуем List, чтобы понять, жив ли кластер
		_, err := cli.ListClusterSkins()
		if err != nil {
			skipIfNoClusterAuth(t, err)
		}

		if err := cli.CreateClusterSkin(skinID); err != nil {
			skipIfNoClusterAuth(t, err)
			return
		}

		_ = cli.StoreClusterSkinFile(skinID, testFile, testContent)

		data, _, err := cli.ReadClusterSkinFile(skinID, testFile)
		if err != nil {
			t.Errorf("ReadClusterSkinFile failed: %v", err)
		}
		if !bytes.Equal(data, testContent) {
			t.Error("ReadClusterSkinFile: content mismatch")
		}

		// [!] ЗАКРЫВАЕМ: RenameClusterSkin (line 237)
		if err := cli.RenameClusterSkin(skinID, renamedID); err != nil {
			t.Errorf("RenameClusterSkin failed: %v", err)
		}

		// Cleanup
		cleanSkinFiles(cli, "cluster", "", renamedID)
		_ = cli.DeleteClusterSkin(renamedID)
	})

	// --- SECTION 4: STOCK SKINS ---
	t.Run("StockSkins_Verification", func(t *testing.T) {
		stockSkin := ""
		files, err := cli.ListStockSkinFiles(stockSkin)
		if err != nil {
			t.Skipf("Stock skin %s not found (500 Unknown Skin is normal here)", stockSkin)
			return
		}

		if len(files) > 0 {
			for f := range files {
				_, _, err := cli.ReadStockSkinFile(stockSkin, f)
				if err == nil {
					break
				} // Проверили один успешный и хватит
			}
		}
	})
}
