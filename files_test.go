package cgpcli

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestIntegration_StorageFiles(t *testing.T) {
	cli := setupSandbox(t)
	defer cli.Logout()

	// Используем технический аккаунт песочницы (например, postmaster)
	account := os.Getenv("CGPUSER")
	if account == "" {
		account = "postmaster"
	}

	// Уникальное имя файла для изоляции
	fileName := fmt.Sprintf("test_file_%d.bin", time.Now().UnixNano()%10000)
	content := []byte("CommuniGate Pro Binary Data: \x00\x01\x02\x03\xff")

	t.Run("Write and Read Binary File", func(t *testing.T) {
		// 1. Запись файла
		err := cli.WriteStorageFile(account, fileName, nil, content, "")
		if err != nil {
			t.Fatalf("WriteStorageFile failed: %v", err)
		}

		// 2. Чтение файла
		res, err := cli.ReadStorageFile(account, fileName, 0, 0, "")
		if err != nil {
			t.Fatalf("ReadStorageFile failed: %v", err)
		}

		// Сверяем содержимое
		if !bytes.Equal(res.Data, content) {
			t.Errorf("Content mismatch.\nGot:  %x\nWant: %x", res.Data, content)
		}

		t.Logf("File size on server: %d", res.CurrentSize)
	})

	t.Run("List and Info", func(t *testing.T) {
		// Проверяем наличие файла в списке
		files, err := cli.ListStorageFiles(account, "", "")
		if err != nil {
			t.Fatalf("ListStorageFiles failed: %v", err)
		}

		if _, ok := files[fileName]; !ok {
			t.Errorf("File %s not found in list", fileName)
		}

		// Проверяем статистику хранилища
		info, err := cli.GetStorageFileInfo(account, "", "")
		if err != nil {
			t.Fatalf("GetStorageFileInfo failed: %v", err)
		}
		if info.FilesCount == 0 {
			t.Error("Expected at least one file in storage info")
		}
	})

	t.Run("Rename and Delete", func(t *testing.T) {
		newName := fileName + ".bak"

		// Переименование
		err := cli.RenameStorageFile(account, fileName, newName, "")
		if err != nil {
			t.Fatalf("Rename failed: %v", err)
		}

		// Удаление
		err = cli.DeleteStorageFile(account, newName, "")
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}

		// Проверка, что файла больше нет
		files, _ := cli.ListStorageFiles(account, "", "")
		if _, ok := files[newName]; ok {
			t.Error("File still exists after deletion")
		}
	})
}
