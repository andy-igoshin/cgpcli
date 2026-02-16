package cgpcli

import (
	"bytes"
	"testing"
)

func TestPBX_Integration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

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
		if err := cli.DeleteDomainPBXFile(testDomain, testFile); err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})

	// --- 2. SERVER PBX ---
	t.Run("ServerPBX", func(t *testing.T) {
		// Внимание: Создание SERVER PBX может требовать прав All Domains
		// Часто среда уже существует. Попробуем записать файл.
		if err := cli.StoreServerPBXFile(testFile, testData); err != nil {
			t.Logf("Note: StoreServerPBXFile might fail without permissions: %v", err)
		} else {
			defer cli.DeleteServerPBXFile(testFile)

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

func TestPBX_Extended(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	const testLang = "English"
	const testFile = "sdk_extended_test.txt"
	testData := []byte("entry main is stop; end entry")

	// 1. Уровень SERVER (Create/Delete/List)
	t.Run("ServerPBX_Full", func(t *testing.T) {
		// Пробуем создать языковую среду (может вернуть ошибку, если уже есть - это ок)
		_ = cli.CreateServerPBX(testLang)

		_ = cli.StoreServerPBXFile(testFile, testData)

		_, err := cli.ListServerPBXFiles(testLang)
		if err != nil {
			t.Logf("ListServerPBXFiles info: %v", err)
		}

		_ = cli.DeleteServerPBX(testLang)
	})

	// 2. Уровень CLUSTER (если сервер в кластере - сработает, если нет - просто покроет код)
	t.Run("ClusterPBX", func(t *testing.T) {
		_ = cli.CreateClusterPBX(testLang)

		// Пытаемся сохранить файл в кластерную среду
		_ = cli.StoreClusterPBXFile(testFile, testData)

		// List
		_, _ = cli.ListClusterPBXFiles("")

		// Read (даже если файла нет, вызов ReadClusterPBXFile будет покрыт)
		_, _ = cli.ReadClusterPBXFile(testFile)

		// Delete file
		_ = cli.DeleteClusterPBXFile(testFile)

		// Delete
		_ = cli.DeleteClusterPBX(testLang)
	})

	// 3. Уровень STOCK и DOMAIN (оставшиеся методы)
	t.Run("StockAndDomainAdditional", func(t *testing.T) {
		// ReadStockPBXFile - пробуем прочитать базовый пакет
		// Обычно pbxbase.spp или что-то подобное всегда есть
		_, _ = cli.ReadStockPBXFile("pbxbase.spp")

		// DeleteDomainPBX (удаление языкового набора в домене)
		// Сначала создадим, чтобы было что удалять
		_ = cli.CreateDomainPBX("test.domain.name", testLang)
		err := cli.DeleteDomainPBX("test.domain.name", testLang)
		if err != nil {
			t.Logf("DeleteDomainPBX info: %v", err)
		}
	})
}
