package cgpcli

import (
	"net"
	"strings"
	"testing"
)

func TestIPList_AddAndMarshal(t *testing.T) {
	list := IPList{}

	err := list.Add("1.1.1.1/32", "Single IP")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	err = list.Add("192.168.1.1-192.168.1.10", "Internal range")
	if err != nil {
		t.Fatalf("Add range failed: %v", err)
	}

	err = list.Add("1.1.1.1/32", "Updated comment")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	err = list.Add("not-an-ip", "Bad entry")
	if err == nil {
		t.Error("Expected error for invalid IP format, got nil")
	}

	res := list.marshal()
	// Исправлено: 1.1.1.1 теперь без /32
	expected := "1.1.1.1 ;Updated comment\n192.168.1.1-192.168.1.10 ;Internal range"
	if res != expected {
		t.Errorf("Marshal mismatch.\nGot:\n%q\nExpected:\n%q", res, expected)
	}
}

func TestIPList_LayoutPreservation(t *testing.T) {
	// Исправлено: 10.0.0.1 в сырых данных без /32 для соответствия нормализации
	raw := "10.0.0.1         ; Long distance comment\n172.16.0.1-172.16.0.255 ; Short"
	list := parseIPList(raw)

	res := list.marshal()
	if res != raw {
		t.Errorf("Layout not preserved.\nGot:\n%q\nExpected:\n%q", res, raw)
	}

	_ = list.Add("10.0.0.1", "New comment")
	res = list.marshal()

	// Проверяем сохранение позиции ';'
	if !strings.Contains(res, "10.0.0.1         ;New comment") {
		t.Errorf("Comment position shift failed.\nGot: %q", res)
	}
}

func TestIPList_Contains(t *testing.T) {
	list := IPList{}
	_ = list.Add("10.0.0.0/24", "Network")
	_ = list.Add("192.168.1.50-192.168.1.60", "Range")
	_ = list.Add("8.8.8.8,9.9.9.9", "Sequence")

	tests := []struct {
		ip     string
		want   bool
		reason string
	}{
		{"10.0.0.5", true, "Inside CIDR"},
		{"10.0.1.5", false, "Outside CIDR"},
		{"192.168.1.55", true, "Inside Range"},
		{"192.168.1.49", false, "Below Range"},
		{"192.168.1.61", false, "Above Range"},
		{"8.8.8.8", true, "Inside Sequence (first)"},
		{"9.9.9.9", true, "Inside Sequence (second)"},
		{"8.8.4.4", false, "Outside all"},
	}

	for _, tt := range tests {
		got := list.Contains(net.ParseIP(tt.ip))
		if got != tt.want {
			t.Errorf("Contains(%s) = %v; want %v (%s)", tt.ip, got, tt.want, tt.reason)
		}
	}
}

func TestIPList_Delete(t *testing.T) {
	list := IPList{}
	_ = list.Add("1.1.1.1", "To be deleted")
	_ = list.Add("2.2.2.2", "To stay")

	// Delete найдет запись по нормализованному значению
	err := list.Delete("1.1.1.1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Delete failed, list size: %d", len(list))
	}
}

func TestIPList_Positioning(t *testing.T) {
	list := IPList{}

	_ = list.Add("1.1.1.1", "First")
	_ = list.Add("3.3.3.3", "Last")

	err := list.Insert(1, "2.2.2.2", "Middle", 0)
	// Исправлено: сравнение с 2.2.2.2 без маски
	if err != nil || len(list) != 3 || dataToString(list[1].Value) != "2.2.2.2" {
		t.Errorf("Insert by index failed. List: %v", list.marshal())
	}

	err = list.InsertBefore("1.1.1.1", "0.0.0.0", "Very First", 0)
	if err != nil || dataToString(list[0].Value) != "0.0.0.0" {
		t.Errorf("InsertBefore failed. Got first: %v", dataToString(list[0].Value))
	}

	err = list.InsertAfter("3.3.3.3", "4.4.4.4", "Extra Last", 0)
	if err != nil || dataToString(list[len(list)-1].Value) != "4.4.4.4" {
		t.Errorf("InsertAfter failed. Got last: %v", dataToString(list[len(list)-1].Value))
	}

	err = list.Insert(0, "4.4.4.4", "Moved to top", 0)
	if err != nil || dataToString(list[0].Value) != "4.4.4.4" || len(list) != 5 {
		t.Errorf("Re-insertion (move) failed. Size: %d, Top: %v", len(list), dataToString(list[0].Value))
	}
}

func TestIPList_SmartIndexOf(t *testing.T) {
	list := IPList{}
	_ = list.Add("10.0.0.1, 10.0.0.2", "Cluster nodes")
	_ = list.Add("192.168.1.1", "Gateway")

	idx := list.indexOf("10.0.0.2")
	if idx != 0 {
		t.Errorf("Smart indexOf failed to find IP inside sequence. Got index: %d", idx)
	}

	idx = list.indexOf("192.168.1.1")
	if idx != 1 {
		t.Errorf("Smart indexOf failed to match IP. Got index: %d", idx)
	}

	err := list.InsertBefore("10.0.0.2", "8.8.8.8", "Before cluster", 0)
	if err != nil || dataToString(list[0].Value) != "8.8.8.8" {
		t.Errorf("InsertBefore smart marker failed: %v", err)
	}
}

func TestIPList_Exclusion(t *testing.T) {
	raw := "10.0.0.0/24      ; Office Network\n!10.0.0.5         ; Blocked printer"
	list := parseIPList(raw)

	res := list.marshal()
	if res != raw {
		t.Errorf("Exclusion layout failed.\nGot:\n%q\nExpected:\n%q", res, raw)
	}

	tests := []struct {
		ip   string
		want bool
	}{
		{"10.0.0.1", true},
		{"10.0.0.5", false},
		{"10.0.1.1", false},
	}

	for _, tt := range tests {
		got := list.Contains(net.ParseIP(tt.ip))
		if got != tt.want {
			t.Errorf("Exclusion logic failed for %s: got %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIPList_AddNegative(t *testing.T) {
	list := IPList{}

	list = append(list, IPEntry{
		Value:      parseDataPart("192.168.1.0/24"),
		IsNegative: false,
		Comment:    "LAN",
		CommentPos: 20,
	})
	list = append(list, IPEntry{
		Value:      parseDataPart("192.168.1.1"),
		IsNegative: true,
		Comment:    "Gateway",
		CommentPos: 20,
	})

	res := list.marshal()
	// Исправлено: !192.168.1.1 (12 символов) + 8 пробелов = 20-я позиция
	expectedPart := "!192.168.1.1        ;Gateway"
	if !strings.Contains(res, expectedPart) {
		t.Errorf("Negative marshal alignment failed.\nGot: %q\nExpected to contain: %q", res, expectedPart)
	}
}

func TestIPList_AddNegativeVariantB(t *testing.T) {
	list := IPList{}

	_ = list.Add("10.0.0.0/24", "Office")
	_ = list.Add("!10.0.0.5", "Blocked")

	if !list[1].IsNegative {
		t.Error("Add failed to recognize '!' prefix")
	}

	if list.Contains(net.ParseIP("10.0.0.5")) {
		t.Error("Contains should be false for excluded IP")
	}
}
