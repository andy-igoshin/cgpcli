package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	defer cli.Close()

	fmt.Println("Connected to 127.0.0.1. Enumerating OIDs...")
	fmt.Println(strings.Repeat("-", 40))

	stats, err := cli.ListAllStats()
	if err != nil {
		log.Fatalf("Error during enumeration: %v", err)
	}

	for i, name := range stats {
		fmt.Printf("%04d | %s\n", i+1, name)
	}

	count := len(stats)

	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("SUMMARY:\n")
	fmt.Printf("  Total OIDs found: %d\n", count)
}
