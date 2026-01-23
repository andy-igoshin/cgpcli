package cgpcli

import (
	"strings"
	"testing"
)

func TestRouterRoundTrip(t *testing.T) {
	input := `localhost = ; main domain\e\e; --- SIP Section ---\eSignal:<911@*> = emergency ; emergency call`

	parts := strings.Split(input, `\e`)
	anyParts := make([]any, len(parts))
	for i, v := range parts {
		anyParts[i] = v
	}

	table := parseRouterTable(anyParts)

	if len(table) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(table))
	}

	if table[0].Rule != "localhost =" {
		t.Errorf("expected rule 'localhost =', got '%s'", table[0].Rule)
	}
	// Исправлено: ожидаем строку с ведущим пробелом, так как он есть в input после ';'
	if table[2].Comment != " --- SIP Section ---" {
		t.Errorf("expected header comment, got '%s'", table[2].Comment)
	}
	if table[2].Rule != "" {
		t.Errorf("expected empty rule for header, got '%s'", table[2].Rule)
	}

	output := table.MarshalCGP()
	if output != input {
		t.Errorf("Round-trip failed.\nExpected: %s\nGot     : %s", input, output)
	}
}

func TestRouterAlignmentProtection(t *testing.T) {
	entry := RouterEntry{
		Rule:    "NewLongRuleName = value",
		Comment: "added comment",
		Pos:     5,
	}

	list := RouterList{entry}
	output := list.MarshalCGP()

	expected := "NewLongRuleName = value ;added comment"
	if output != expected {
		t.Errorf("Alignment protection failed.\nExpected: %s\nGot     : %s", expected, output)
	}
}

func TestRouterInsertAndPos(t *testing.T) {
	table := RouterList{
		{Rule: "rule1", Comment: "comm1", Pos: 30},
		{Rule: "rule3", Comment: "comm3", Pos: 30},
	}

	err := table.InsertBefore("rule3", "rule2", "comm2", 30)
	if err != nil {
		t.Fatal(err)
	}

	output := table.MarshalCGP()
	// Проверка на наличие 25 пробелов (30 - len("rule2"))
	expectedPart := "rule2                         ;comm2"
	if !strings.Contains(output, expectedPart) {
		t.Errorf("Alignment generation failed.\nExpected to contain: %s\nGot: %s", expectedPart, output)
	}
}
