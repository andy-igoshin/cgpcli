package cgpcli

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

// helper для получения строки из WriteCGP в тестах
func marshalIPList(t *testing.T, l IPList) string {
	var b bytes.Buffer
	if err := l.WriteCGP(&b); err != nil {
		t.Fatalf("WriteCGP failed: %v", err)
	}
	// Убираем обрамляющие кавычки для удобства сравнения в тестах
	res := b.String()
	if len(res) >= 2 && res[0] == '"' && res[len(res)-1] == '"' {
		return res[1 : len(res)-1]
	}
	return res
}

func TestIPList_AddAndMarshal(t *testing.T) {
	list := IPList{}

	_ = list.Add("1.1.1.1/32", "Single IP")
	_ = list.Add("192.168.1.1-192.168.1.10", "Internal range")
	_ = list.Add("1.1.1.1/32", "Updated comment")

	res := marshalIPList(t, list)

	// Теперь используем \e как разделитель (Honest)
	// 1.1.1.1 нормализуется (убирается /32 если маска полная)
	expected := `1.1.1.1 ;Updated comment\e192.168.1.1-192.168.1.10 ;Internal range`
	if res != expected {
		t.Errorf("Marshal mismatch.\nGot:      %q\nExpected: %q", res, expected)
	}
}

func TestIPList_LayoutPreservation(t *testing.T) {
	// Входящие данные с \e (как приходят от сервера)
	raw := `10.0.0.1          ; Long distance comment\e172.16.0.1-172.16.0.255 ; Short`

	// Unmarshal превратит \e в \n для парсера
	list := parseIPList(strings.ReplaceAll(raw, `\e`, "\n"))

	res := marshalIPList(t, list)
	if res != raw {
		t.Errorf("Layout not preserved.\nGot:      %q\nExpected: %q", res, raw)
	}

	// Проверка сохранения позиции при обновлении комментария
	_ = list.Add("10.0.0.1", "New comment")
	res = marshalIPList(t, list)

	if !strings.Contains(res, "10.0.0.1          ;New comment") {
		t.Errorf("Comment position shift failed.\nGot: %q", res)
	}
}

func TestIPList_Exclusion(t *testing.T) {
	// Эмулируем ответ сервера с инверсией
	raw := `10.0.0.0/24      ; Office Network\e!10.0.0.5          ; Blocked printer`
	list := parseIPList(strings.ReplaceAll(raw, `\e`, "\n"))

	res := marshalIPList(t, list)
	if res != raw {
		t.Errorf("Exclusion layout failed.\nGot:      %q\nExpected: %q", res, raw)
	}

	tests := []struct {
		ip   string
		want bool
	}{
		{"10.0.0.1", true},
		{"10.0.0.5", false}, // Инверсия сработала
		{"10.0.1.1", false}, // Вне сети
	}

	for _, tt := range tests {
		got := list.Contains(net.ParseIP(tt.ip))
		if got != tt.want {
			t.Errorf("Exclusion logic failed for %s: got %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIPList_SpecialCharsInComment(t *testing.T) {
	list := IPList{}
	// Проверяем, что наши новые правила эскейпинга работают и здесь
	comment := "Comment with \n and \t and \\"
	_ = list.Add("8.8.8.8", comment)

	res := marshalIPList(t, list)

	// Ожидаем: \n -> \e, \t -> \t, \ -> \\
	expectedComment := `8.8.8.8 ;Comment with \e and \t and \\`
	if !strings.Contains(res, expectedComment) {
		t.Errorf("Special chars handling failed.\nGot: %q", res)
	}
}

func TestIPList_Positioning(t *testing.T) {
	list := IPList{}
	_ = list.Add("1.1.1.1", "First")
	_ = list.Add("3.3.3.3", "Last")

	// Вставка в середину
	_ = list.Insert(1, "2.2.2.2", "Middle", 0)

	res := marshalIPList(t, list)
	expected := `1.1.1.1 ;First\e2.2.2.2 ;Middle\e3.3.3.3 ;Last`
	if res != expected {
		t.Errorf("Insert failed.\nGot: %q", res)
	}
}

func TestIPList_Delete(t *testing.T) {
	list := IPList{}
	_ = list.Add("1.1.1.1", "A")
	_ = list.Add("2.2.2.2", "B")

	if err := list.Delete("1.1.1.1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	res := marshalIPList(t, list)
	if strings.Contains(res, "1.1.1.1") {
		t.Error("Entry 1.1.1.1 still exists after delete")
	}
}
