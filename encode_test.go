package cgpcli

import (
	"bufio"
	"net"
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
