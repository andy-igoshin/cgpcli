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
	cli, err := cgpcli.New("127.0.0.1", cgpuser, cgppass, cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Close()

	fmt.Println("--- Testing LISTDOMAINOBJECTS with Enum Flags ---")

	// 3. Вызов функции с использованием битовой маски
	mask := cgpcli.TypeAccount | cgpcli.TypeForwarder | cgpcli.TypeAlias

	res, err := cli.ListDomainObjects("test.domain.name", 100, "", mask, "")
	if err != nil {
		log.Fatalf("ListDomainObjects error: %v", err)
	}

	// 4. Красивый вывод
	fmt.Println("\nStructure via PrettyPrint:")
	cgpcli.PrettyPrint(res)

	// 5. Демонстрация работы с типизированными данными
	fmt.Printf("\n--- Summary (Count: %d, Cookie: '%s') ---\n", len(res.Objects), res.Cookie)

	keys := make([]string, 0, len(res.Objects))
	for k := range res.Objects {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		obj := res.Objects[name]

		if len(obj.Targets) > 0 {
			fmt.Printf("[%v] %-15s -> %v\n", obj.Type, obj.Name, obj.Targets)
		} else {
			fmt.Printf("[%v] %-15s\n", obj.Type, obj.Name)
		}
	}
}
