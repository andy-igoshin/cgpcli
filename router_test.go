package cgpcli

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestRouterRoundTrip(t *testing.T) {
	// Исходные данные (внутри кавычек)
	rawTable := "localhost = ; main domain\n\n; --- SIP Section ---\nSignal:<911@*> = emergency ; emergency call"
	expTable := `localhost = ; main domain\e\e; --- SIP Section ---\eSignal:<911@*> = emergency ; emergency call`
	// Ожидаем, что маршалер обернет всё в кавычки
	expectedOutput := `"` + expTable + `"`

	table := parseRouterTable(rawTable)

	if len(table) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(table))
	}

	if table[0].Rule != "localhost =" {
		t.Errorf("expected rule 'localhost =', got '%s'", table[0].Rule)
	}

	// Проверка заголовка (пустое правило, только комментарий)
	if table[2].Comment != " --- SIP Section ---" {
		t.Errorf("expected header comment, got '%s'", table[2].Comment)
	}

	// Проверка через WriteCGP
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := table.WriteCGP(w); err != nil {
		t.Fatalf("WriteCGP failed: %v", err)
	}
	w.Flush()

	output := b.String()
	if output != expectedOutput {
		t.Errorf("Round-trip failed.\nExpected: %s\nGot     : %s", expectedOutput, output)
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

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := table.WriteCGP(w); err != nil {
		t.Fatal(err)
	}
	w.Flush()

	output := b.String()

	// "rule2" имеет длину 5 символов. При Pos: 30 ожидаем 25 пробелов перед ';'
	spaces := strings.Repeat(" ", 30-len("rule2"))
	expectedPart := "rule2" + spaces + ";comm2"

	if !strings.Contains(output, expectedPart) {
		t.Errorf("Alignment generation failed.\nExpected to contain: %s\nGot: %s", expectedPart, output)
	}
}
