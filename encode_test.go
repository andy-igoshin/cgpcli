package cgpcli

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	testZone := time.UTC

	tests := []struct {
		name     string
		input    any
		contains []string
	}{
		{
			name: "Atoms and Strings",
			input: map[string]any{
				"Command": Atom("SETSERVERTRUSTEDCERTS"),
				"Param":   "Simple",
				"Space":   "String with space",
			},
			contains: []string{
				"Command=SETSERVERTRUSTEDCERTS",
				"Param=Simple",
				`Space="String with space"`,
			},
		},
		{
			name: "DataBlock and IP",
			input: map[string]any{
				"Raw":   []byte("HI"),          // Маршалер сделает Base64 -> [SEk=]
				"Block": DataBlock("Q0VSVA=="), // Мы передаем уже готовый Base64 "CERT"
				"Addr":  CGPIP{IP: net.ParseIP("10.0.0.1"), Port: 25},
			},
			contains: []string{
				"Raw=[SEk=]",
				"Block=[Q0VSVA==]",
				"Addr=#I[10.0.0.1]:25",
			},
		},
		{
			name: "Certificate_DataBlock",
			input: map[string]any{
				// Здесь мы передаем уже валидный Base64 сертификата
				"Certificate": DataBlock("ZmFrZS1jZXJ0LWRhdGE="),
			},
			contains: []string{
				"Certificate=[ZmFrZS1jZXJ0LWRhdGE=]",
			},
		},
		{
			name: "Quotas and Sizes",
			input: map[string]any{
				"MaxAccountSize": int64(-1),
				"MaxWebSize":     int64(1024),
				"Exact":          int64(123),
			},
			contains: []string{
				"MaxAccountSize=unlimited",
				"MaxWebSize=1K",
				"Exact=123", // Обычное число, не квота
			},
		},
		{
			name: "Time and Tags",
			input: map[string]any{
				"Created": time.Date(2026, 1, 8, 10, 0, 0, 0, testZone),
				"LastLog": time.Date(2026, 1, 8, 10, 0, 0, 0, testZone),
				"Past":    TimePast,
			},
			contains: []string{
				`Created="08-01-2026_10:00:00"`, // Legacy формат
				`LastLog=#T08-01-2026_10:00:00`, // Modern формат
				`Past=#TPAST`,
			},
		},
		{
			name: "Nested Structures Compact",
			input: map[string]any{
				"Settings": map[string]any{
					"Enabled": true,
					"List":    []any{Atom("Default"), "text"},
				},
			},
			contains: []string{
				"Enabled=YES",
				"List=(Default,text)",
			},
		},
		{
			name: "Null Values",
			input: map[string]any{
				"Reset": nil,
			},
			contains: []string{
				"Reset=#NULL#",
			},
		},
		{
			name: "String Map and Slice",
			input: map[string]any{
				"ACL": map[string]string{"user": "l"},
				"Sub": []string{"INBOX", "Sent"},
			},
			contains: []string{
				`ACL={user=l;}`,
				`Sub=(INBOX,Sent)`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			w := bufio.NewWriter(&b)
			// Используем наш внутренний маршалер через MarshalTo или MarshalToString
			err := MarshalTo(w, tt.input)
			w.Flush()

			if err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			res := b.String()
			for _, s := range tt.contains {
				if !strings.Contains(res, s) {
					t.Errorf("%s: Result should contain %q, but got: %s", tt.name, s, res)
				}
			}
			// Проверка на отсутствие пробелов (согласно CLI.pm)
			if strings.Contains(res, "= ") || strings.Contains(res, "; ") || strings.Contains(res, ", ") {
				t.Errorf("%s: Result contains unexpected spaces: %s", tt.name, res)
			}
		})
	}
}

func TestWriteInt64(t *testing.T) {
	tests := []struct {
		name string
		val  int64
	}{
		{"Zero", 0},
		{"One digit positive", 5},
		{"One digit negative", -5},
		{"Multi digit positive", 123456789},
		{"Multi digit negative", -123456789},
		{"Max Int64", 9223372036854775807},
		{"Min Int64", -9223372036854775808},
		{"Power of 10", 1000},
		{"Negative power of 10", -1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			// Передаем нашу обертку или прямо Builder, если интерфейс совпадает
			err := writeInt64(&sb, tt.val)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := sb.String()
			want := strconv.FormatInt(tt.val, 10)

			if got != want {
				t.Errorf("writeInt64(%d) = %q, want %q", tt.val, got, want)
			}
		})
	}
}

// Дополнительный тест на отсутствие аллокаций
type discardWriter struct{}

func (d discardWriter) Write(p []byte) (int, error)     { return len(p), nil }
func (d discardWriter) WriteString(s string) (int, error) { return len(s), nil }
func (d discardWriter) WriteByte(c byte) error           { return nil }

func TestWriteInt64_Allocs(t *testing.T) {
	var dw discardWriter

	allocs := testing.AllocsPerRun(100, func() {
		_ = writeInt64(dw, -9223372036854775808)
	})

	if allocs > 0 {
		t.Errorf("writeInt64 allocated %f times, want 0", allocs)
	}
}
