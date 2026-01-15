package main

import (
	"fmt"
	"log"
	"os"

	"git.vsu.ru/ai/cgpcli"
)

func main() {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		log.Fatal("Error: please set CGPUSER and CGPPASS environment variables")
	}

	cli, err := cgpcli.New("127.0.0.1", cgpuser, cgppass, cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Logout()

	// 1. Создаем два алерта: обычный и полностью цифровой
	alertsToTest := []string{
		"Standard alert text with spaces",
		"98765432109876", // Полностью цифровая строка
	}

	for _, msg := range alertsToTest {
		fmt.Printf("--- Posting: %s ---\n", msg)
		if err := cli.PostServerAlert(msg); err != nil {
			log.Printf("Post error for '%s': %v", msg, err)
		}
	}

	// 2. Читаем алерты
	fmt.Println("\n--- GETSERVERALERTS (Checking types) ---")
	// Пока используем GetServerAlerts, который возвращает map[string]any (или уже map[string]string?)
	// Чтобы увидеть "сырые" типы от парсера, проверим через Any.
	res, err := cli.Query("GETSERVERALERTS")
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}

	alerts, ok := res.(map[string]any)
	if !ok {
		log.Fatalf("Expected map, got %T", res)
	}

	// 3. Анализируем типы ключей и значений
	for k, v := range alerts {
		if k == "" {
			continue
		} // Пропускаем сервисный ключ
		fmt.Printf("Key: [%s] (type %T) | Value: [%v] (type %T)\n", k, k, v, v)
	}

	// 4. Очистка
	fmt.Println("\n--- Cleanup ---")
	for k := range alerts {
		if k == "" {
			continue
		}
		if err := cli.RemoveServerAlert(k); err != nil {
			fmt.Printf("Remove error for %s: %v\n", k, err)
		} else {
			fmt.Println("Removed:", k)
		}
	}
}
