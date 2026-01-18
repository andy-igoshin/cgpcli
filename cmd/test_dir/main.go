package main

import (
	"fmt"
	"log"
	"os"

	"git.vsu.ru/ai/cgpcli"
)

func main() {
	cli, err := cgpcli.New("127.0.0.1", os.Getenv("CGPUSER"), os.Getenv("CGPPASS"), cgpcli.APOP, false)
	if err != nil {
		log.Fatalf("Connect error: %v", err)
	}
	defer cli.Logout()

	fmt.Println("--- Directory Units ---")
	units, err := cli.ListDirectoryUnits(false)
	if err != nil {
		fmt.Printf("ListUnits error: %v\n", err)
	} else {
		cgpcli.PrettyPrint(units)
	}

	fmt.Println("\n--- Directory Access Rights ---")
	rights, err := cli.GetDirectoryAccessRights(false)
	if err != nil {
		fmt.Printf("GetRights error: %v\n", err)
	} else {
		cgpcli.PrettyPrint(rights)
	}

	fmt.Println("\n--- Directory Integration ---")
	// Настройки интеграции (LDAP, Active Directory и т.д.)
	integration, err := cli.Query("GETDIRECTORYINTEGRATION")
	if err != nil {
		fmt.Printf("GetIntegration error: %v\n", err)
	} else {
		cgpcli.PrettyPrint(integration)
	}
}
