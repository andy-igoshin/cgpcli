package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.vsu.ru/ai/cgpcli"
)

func main() {
	// Используем переменные окружения из ваших тестов
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")

	if cgpuser == "" || cgppass == "" {
		log.Fatal("Error: please set CGPUSER and CGPPASS environment variables")
	}

	// Используем ваш конструктор New (APOP, без TLS для теста на 127.0.0.1)
	cli, err := cgpcli.New("127.0.0.1", cgpuser, cgppass, cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Logout()

	fmt.Println("Connected to 127.0.0.1. Enumerating OIDs...")
	fmt.Println(strings.Repeat("-", 40))

	// Запрашиваем список
	stats, err := cli.ListAllStats()
	if err != nil {
		log.Fatalf("Error during enumeration: %v", err)
	}

	// Выводим с порядковыми номерами
	for i, name := range stats {
		fmt.Printf("%04d | %s\n", i+1, name)
	}

	count := len(stats)
	recommended := ((count / 512) + 1) * 512

	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("SUMMARY:\n")
	fmt.Printf("  Total OIDs found: %d\n", count)
}
