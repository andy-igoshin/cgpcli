package main

import (
	"fmt"
	"log"
	"os"
	"sort"

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
	cli, err := cgpcli.New("127.0.0.1", cgpuser, cgppass, cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Logout()

	fmt.Println("--- Testing LISTDOMAINOBJECTS with Enum Flags ---")

	// 3. Вызов функции с использованием битовой маски
	// Например: запрашиваем Аккаунты и Форвардеры, игнорируя Алиасы
	mask := cgpcli.TypeAccount | cgpcli.TypeForwarder | cgpcli.TypeAlias
	
	res, err := cli.ListDomainObjects("test.domain.name", 100, "", mask, "")
	if err != nil {
		log.Fatalf("ListDomainObjects error: %v", err)
	}

	// 4. Красивый вывод через нашу утилиту
	fmt.Println("\nStructure via PrettyPrint:")
	cgpcli.PrettyPrint(res)

	// 5. Демонстрация работы с типизированными данными
	fmt.Printf("\n--- Summary (Total: %d, Cookie: '%s') ---\n", res.Total, res.Cookie)

	// Сортируем ключи для красивого вывода в консоль
	keys := make([]string, 0, len(res.Objects))
	for k := range res.Objects {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		obj := res.Objects[name]
		
		// Благодаря Enum мы можем делать чистые проверки
		typeLabel := obj.Type.String()
		
		if len(obj.Targets) > 0 {
			fmt.Printf("[%s] %-15s -> %v\n", typeLabel, obj.Name, obj.Targets)
		} else {
			fmt.Printf("[%s] %-15s\n", typeLabel, obj.Name)
		}
	}
}
