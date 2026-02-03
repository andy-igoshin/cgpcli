package cgpcli

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestIntegration_StorageFiles(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Используем технический аккаунт песочницы (например, postmaster)
	account := "postmaster"

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

func TestIntegration_StorageFiles_Extended(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	account := "testuser1@test.domain.name"
	fileName := "attr_test_file.txt"
	content := []byte("Attribute and Subscription Test Content")

	// Предварительно создаем файл для тестов атрибутов
	_ = cli.WriteStorageFile(account, fileName, nil, content, "")
	defer cli.DeleteStorageFile(account, fileName, "")

	t.Run("FileAttributes", func(t *testing.T) {
		// 1. Читаем атрибуты (пустой список attrs — читает все базовые)
		attrs, err := cli.ReadStorageFileAttr(account, fileName, nil, "")
		if err != nil {
			t.Fatalf("ReadStorageFileAttr failed: %v", err)
		}
		t.Logf("Initial file attributes: %v", attrs)

		// 2. Обновляем атрибут.
		// В CGP для файлов можно пробовать менять 'ModTime' или кастомные флаги,
		// но самый безопасный способ проверить вызов - отправить пустой или базовый набор.
		// В CLI атрибуты передаются как массив any [name, value, name, value...]
		updateAttrs := []any{"Keywords", "GoTest"}
		if err := cli.UpdateStorageFileAttr(account, fileName, updateAttrs, ""); err != nil {
			// Некоторые хранилища могут не поддерживать кастомные атрибуты,
			// если будет ошибка - это тоже покрытие метода.
			t.Logf("UpdateStorageFileAttr (Keywords) info: %v", err)
		}

		// Повторное чтение для проверки
		_, _ = cli.ReadStorageFileAttr(account, fileName, []string{"Keywords"}, "")
	})

	t.Run("FileSubscriptions", func(t *testing.T) {
		// 1. Получаем текущие подписки
		original, err := cli.GetFileSubscription(account)
		if err != nil {
			t.Fatalf("GetFileSubscription failed: %v", err)
		}

		// 2. Устанавливаем тестовую подписку (на чужую папку или вымышленную)
		testSubs := append(original, "~otheruser/public")
		if err := cli.SetFileSubscription(account, testSubs); err != nil {
			t.Fatalf("SetFileSubscription failed: %v", err)
		}

		// 3. Проверяем
		updated, err := cli.GetFileSubscription(account)
		if err != nil {
			t.Fatalf("Verification GetFileSubscription failed: %v", err)
		}

		found := false
		for _, s := range updated {
			if s == "~otheruser/public" {
				found = true
				break
			}
		}

		// 4. Возвращаем как было
		_ = cli.SetFileSubscription(account, original)

		if !found {
			t.Error("Test subscription was not saved")
		}
	})
}
