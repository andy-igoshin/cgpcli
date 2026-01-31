package cgpcli

import (
	"bytes"
	"testing"
)

func TestPBX_Integration(t *testing.T) {
        cli := getTestCli(t)
        defer cli.Logout()

	const (
		testDomain = "test.domain.name"
		testFile   = "sdk_test_script.txt"
	)
	testData := []byte("entry main is nop; end entry")

	// --- 1. DOMAIN PBX ---
	t.Run("DomainPBX", func(t *testing.T) {
		// Убеждаемся, что среда создана (игнорируем ошибку, если уже есть)
		_ = cli.CreateDomainPBX(testDomain, "")

		// Store
		if err := cli.StoreDomainPBXFile(testDomain, testFile, testData); err != nil {
			t.Fatalf("StoreDomainPBXFile failed: %v", err)
		}

		// List
		files, err := cli.ListDomainPBXFiles(testDomain, "")
		if err != nil {
			t.Errorf("ListDomainPBXFiles failed: %v", err)
		}
		if _, ok := files[testFile]; !ok {
			t.Errorf("File %s not found in list", testFile)
		}

		// Read
		readData, err := cli.ReadDomainPBXFile(testDomain, testFile)
		if err != nil {
			t.Fatalf("ReadDomainPBXFile failed: %v", err)
		}
		if !bytes.Equal(readData, testData) {
			t.Errorf("Data mismatch: expected %s, got %s", testData, readData)
		}

		// Delete
		if err := cli.StoreDomainPBXFile(testDomain, testFile, nil); err != nil {
			t.Errorf("Delete (Store with nil) failed: %v", err)
		}
	})

	// --- 2. SERVER PBX ---
	t.Run("ServerPBX", func(t *testing.T) {
		// Внимание: Создание SERVER PBX может требовать прав All Domains
		// Часто среда уже существует. Попробуем записать файл.
		if err := cli.StoreServerPBXFile(testFile, testData); err != nil {
			t.Logf("Note: StoreServerPBXFile might fail without permissions: %v", err)
		} else {
			defer cli.StoreServerPBXFile(testFile, nil)

			readData, err := cli.ReadServerPBXFile(testFile)
			if err == nil && !bytes.Equal(readData, testData) {
				t.Errorf("Server data mismatch")
			}
		}
	})

	// --- 3. STOCK PBX (Read-Only) ---
	t.Run("StockPBX", func(t *testing.T) {
		// Читаем стандартный файл, который точно есть в любой установке (например, default.sppkg или версионный файл)
		// Для теста просто попробуем List без языка
		files, err := cli.ListStockPBXFiles("")
		if err != nil {
			t.Errorf("ListStockPBXFiles failed: %v", err)
		}
		if len(files) == 0 {
			t.Log("Stock PBX is empty or not accessible")
		}
	})
}
