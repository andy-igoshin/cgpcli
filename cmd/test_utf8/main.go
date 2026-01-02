package main

import (
	"fmt"
	"log"
	"os"

	"git.vsu.ru/ai/cgpcli"
)

func main() {
	// 1. Настройка окружения
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")

	if cgpuser == "" || cgppass == "" {
		log.Fatal("Error: please set CGPUSER and CGPPASS environment variables")
	}

	// 2. Инициализация клиента
	// Используем ваш конструктор New
	cli, err := cgpcli.New("127.0.0.1", cgpuser, cgppass, cgpcli.Plain, false)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Logout()

	targetAccount := "testuser1@test.domain.name" // замените на ваш боевой/тестовый
	russianName := "Алексей Самылов"              // Используем имя автора из ваших любимых книг для теста

	fmt.Printf("Updating RealName for %s to: %s\n", targetAccount, russianName)

	// 1. Устанавливаем русское имя
	// UpdateAccountSettings отправляет данные через Marshal
	err = cli.UpdateAccountSettings(targetAccount, map[string]any{
		"RealName": russianName,
	})
	if err != nil {
		log.Fatalf("Failed to update settings: %v", err)
	}
	fmt.Println("Update successful.")

	// 2. Читаем обратно
	settings, err := cli.GetAccountSettings(targetAccount)
	if err != nil {
		log.Fatalf("Failed to get settings: %v", err)
	}

	if val, ok := settings["RealName"]; ok {
		fmt.Printf("Read from CGP RealName: %v (Type: %T)\n", val, val)

		if val == russianName {
			fmt.Println("MATCH: UTF-8 integrity preserved by the library!")
		} else {
			fmt.Printf("MISMATCH! Expected %s, got %v\n", russianName, val)
		}
	} else {
		fmt.Println("RealName not found in settings.")
	}
}
