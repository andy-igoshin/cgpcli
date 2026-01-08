package cgpcli

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	testZone := time.FixedZone("UTC", 0)

	tests := []struct {
		name     string
		input    any
		contains []string
		excludes []string
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
			name: "DataBlock and Hex",
			input: map[string]any{
				"Raw":   []byte{0x01, 0x02},
				"Block": DataBlock("CERT"),
			},
			contains: []string{
				"Raw=#[0102]",
				"Block=[CERT]",
			},
		},
		{
			name: "Quotas and Sizes",
			input: map[string]any{
				"Limit": int64(-1),
				"SizeK": 1024,
				"Exact": 123,
			},
			contains: []string{
				"Limit=unlimited",
				"SizeK=1K",
				"Exact=123",
			},
		},
		{
			name:  "Time Formatting",
			input: time.Date(2026, 1, 8, 10, 0, 0, 0, testZone),
			contains: []string{
				`"08-Jan-2026 10:00:00 +0000"`,
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
				"List=(Default,text)", // Без пробелов внутри ()
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
				`ACL={user=l;}`,    // Кавычки ушли, так как строки "safe"
				`Sub=(INBOX,Sent)`, // То же самое здесь
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			w := bufio.NewWriter(&b)
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
			// Проверка на отсутствие пробелов вне кавычек (упрощенно)
			if strings.Contains(res, "= ") || strings.Contains(res, "; ") || strings.Contains(res, ", ") {
				t.Errorf("%s: Result contains unexpected spaces: %s", tt.name, res)
			}
		})
	}
}

func TestFormatCGPSize(t *testing.T) {
	cases := []struct {
		in  int64
		out string
	}{
		{-1, "unlimited"},
		{0, "0"},
		{1024, "1K"},
		{1048576, "1M"},
	}
	for _, c := range cases {
		if got := formatCGPSize(c.in); got != c.out {
			t.Errorf("formatCGPSize(%d) = %s, want %s", c.in, got, c.out)
		}
	}
}
