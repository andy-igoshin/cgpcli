package main

import (
	"fmt"
	"log"
	"os"

	"git.vsu.ru/ai/cgpcli"
)

func main() {
	// 1. Определяем сервер (local по умолчанию)
	mode := "local"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	suffix := ""
	defaultAccount := "testuser1@test.domain.name"

	if mode == "ro" {
		suffix = "_RO"
		defaultAccount = "testadmin@vsu.ru" // или любой другой боевой
		fmt.Println(">>> MODE: READ-ONLY (PRODUCTION) <<<")
	} else {
		fmt.Println(">>> MODE: LOCAL TEST <<<")
	}

	// 2. Определяем целевой аккаунт
	target := defaultAccount
	if len(os.Args) > 2 {
		target = os.Args[2]
	}

	// 3. Загружаем переменные
	host := os.Getenv("CGPHOST" + suffix)
	user := os.Getenv("CGPUSER" + suffix)
	pass := os.Getenv("CGPPASS" + suffix)

	if host == "" || user == "" || pass == "" {
		log.Fatalf("Error: env vars for %s are not set", "CGP..."+suffix)
	}

	// 4. Подключаемся
	cli, err := cgpcli.New(host, user, pass, cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", host, err)
	}
	defer cli.Close()

	fmt.Printf("--- Inspecting: %s ---\n", target)

	// Читаем настройки (Settings)
	fmt.Println("\n[Account Settings]")
	settings, err := cli.GetAccountSettings(target)
	if err != nil {
		log.Printf("Settings error: %v", err)
	} else {
		cgpcli.PrettyPrint(settings)
	}

	// Читаем настройки (Effective Settings)
	fmt.Println("\n[Account Effective Settings]")
	settings, err = cli.GetAccountSettings(target)
	if err != nil {
		log.Printf("Settings error: %v", err)
	} else {
		cgpcli.PrettyPrint(settings)
	}

	// Читаем состояние (Info)
	fmt.Println("\n[Account Info]")
	info, err := cli.GetAccountInfo(target, nil)
	if err != nil {
		log.Printf("Info error: %v", err)
	} else {
		cgpcli.PrettyPrint(info)
	}
}
