package main

import (
	"fmt"
	"git.vsu.ru/ai/cgpcli"
	"log"
	"os"
)

func main() {
	cli, err := cgpcli.New("127.0.0.1", os.Getenv("CGPUSER"), os.Getenv("CGPPASS"), cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Connect error: %v", err)
	}
	defer cli.Close()

	targetAcc := "testuser1@test.domain.name"

	fmt.Printf("--- Testing File Storage for %s ---\n", targetAcc)

	// 1. Создаем директорию (в конце слеш, данных нет)
	fmt.Print("Creating directory 'uploads/'... ")
	err = cli.WriteStorageFile(targetAcc, "uploads/", nil, []byte{}, "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	// 2. Записываем файл в эту директорию
	fmt.Print("Writing file 'uploads/hello.txt'... ")
	content := []byte("Hello from Go SDK via CLI!")
	err = cli.WriteStorageFile(targetAcc, "uploads/hello.txt", "NEW", content, "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	// 3. Читаем информацию о файле
	fmt.Print("Getting file info... ")
	info, err := cli.ListStorageFiles(targetAcc, "uploads/", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found: %v\n", info)
	}

	// 4. Читаем содержимое файла
	fmt.Print("Reading file back... ")
	fileData, err := cli.ReadStorageFile(targetAcc, "uploads/hello.txt", 0, 0, "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Content: [%s], Size: %d, Modified: %s\n",
			string(fileData.Data), fileData.CurrentSize, fileData.ModifiedTime)
	}

	// 5. Проверяем подписки (обычно это массив строк)
	fmt.Print("Getting file subscriptions... ")
	subs, err := cli.GetFileSubscription(targetAcc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Subs: %v\n", subs)
	}
}
