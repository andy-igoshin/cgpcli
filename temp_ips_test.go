package cgpcli

import (
	"bufio"
	"bytes"
	"net"
	"testing"
	"time"
)

func TestTempIPList_All(t *testing.T) {
	// Имитируем реальный ответ CommuniGate Pro
	// 1. IPv4 с таймаутом в секундах
	// 2. IPv6 с флагом постоянства
	// 3. IPv4 без дополнительных флагов
	raw := "[192.168.1.1]-3600,[2001:db8::1]-*,[10.0.0.5]"

	list := &TempIPList{items: parseTempIPsToMap(raw)}

	t.Run("Parsing_Basics", func(t *testing.T) {
		if list.Len() != 3 {
			t.Errorf("Expected 3 items, got %d", list.Len())
		}
	})

	t.Run("Get_Existing_IPv4", func(t *testing.T) {
		info, ok, err := list.Get("192.168.1.1")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("IP should be found")
		}
		if info.TimeLeft != 3600*time.Second {
			t.Errorf("Expected 3600s, got %v", info.TimeLeft)
		}
		if info.Permanent {
			t.Error("Should not be permanent")
		}
	})

	t.Run("Get_Existing_IPv6_Permanent", func(t *testing.T) {
		// Проверяем поиск IPv6 по строке
		info, ok, err := list.Get("2001:db8::1")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !ok || !info.Permanent {
			t.Error("IPv6 permanent record not found or incorrect")
		}
	})

	t.Run("IPv6_Normalization_Consistency", func(t *testing.T) {
		// Проверяем, что разные формы записи IPv6 дают одно попадание в мапу
		// В мапе лежит [2001:db8::1], ищем полную форму
		_, ok, _ := list.Get("2001:0db8:0000:0000:0000:0000:0000:0001")
		if !ok {
			t.Error("Failed to find IPv6 using expanded string format")
		}
	})

	t.Run("Modification_Lifecycle", func(t *testing.T) {
		l := &TempIPList{}
		testIP := "172.16.0.1"

		// Test Add
		err := l.Add(testIP, 10*time.Minute, false)
		if err != nil {
			t.Errorf("Add failed: %v", err)
		}

		if exists, _ := l.Exists(testIP); !exists {
			t.Error("IP should exist after Add")
		}

		// Test Delete
		err = l.Delete(testIP)
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}

		if exists, _ := l.Exists(testIP); exists {
			t.Error("IP should not exist after Delete")
		}
	})

	t.Run("Invalid_Inputs", func(t *testing.T) {
		l := &TempIPList{}

		// 1. Некорректная строка
		_, _, err := l.Get("invalid-ip")
		if err == nil {
			t.Error("Expected error for malformed string")
		}

		// 2. Неподдерживаемый тип (float)
		_, _, err = l.Get(123.45)
		if err == nil {
			t.Error("Expected error for unsupported type")
		}

		// 3. nil в качестве net.IP
		var nilIP net.IP
		_, _, err = l.Get(nilIP)
		if err == nil {
			t.Error("Expected error for nil net.IP")
		}
	})

	t.Run("Export_IPs", func(t *testing.T) {
		l := &TempIPList{}
		l.Add("1.1.1.1", 0, false)
		l.Add("2.2.2.2", 0, false)

		exported := l.IPs()
		if len(exported) != 2 {
			t.Errorf("Expected 2 exported IPs, got %d", len(exported))
		}

		// Убеждаемся, что это именно net.IP, пригодные для UpdateServerSettings
		for _, ip := range exported {
			if ip.To4() == nil {
				t.Errorf("%v is not a valid IPv4", ip)
			}
		}
	})
}

type mockEncodeWriter struct {
	*bufio.Writer
}

func (m mockEncodeWriter) WriteByte(c byte) error {
	return m.Writer.WriteByte(c)
}

func (m mockEncodeWriter) WriteString(s string) (int, error) {
	return m.Writer.WriteString(s)
}

func TestTempIPList_RoundTrip(t *testing.T) {
	list := &TempIPList{items: make(map[IPKey]TempIP)}
	ip := net.ParseIP("192.168.250.250")
	// Добавляем запись с 300 секундами
	_ = list.Add(ip, 300*time.Second, false)

	var b bytes.Buffer
	mw := mockEncodeWriter{bufio.NewWriter(&b)}

	if err := list.WriteCGP(mw); err != nil {
		t.Fatalf("WriteCGP failed: %v", err)
	}
	mw.Flush()

	gotRaw := b.String()
	// Ожидаем: "[192.168.250.250]-300" в кавычках
	expected := `"[192.168.250.250]-300"`
	if gotRaw != expected {
		t.Errorf("WriteCGP output mismatch.\nGot:      %s\nExpected: %s", gotRaw, expected)
	}

	// Проверяем, что парсер понимает наш собственный вывод (без кавычек)
	trimmed := gotRaw[1 : len(gotRaw)-1]
	parsedMap := parseTempIPsToMap(trimmed)

	item, ok := parsedMap[ToKey(ip)]
	if !ok {
		t.Fatal("Parser failed to find IP in the resulting map")
	}
	if item.TimeLeft.Seconds() != 300 {
		t.Errorf("TimeLeft mismatch: got %.0f, want 300", item.TimeLeft.Seconds())
	}
}
