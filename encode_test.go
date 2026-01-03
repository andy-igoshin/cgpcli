package cgpcli

import (
	"strings"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	// Фиксируем зону для тестов времени, чтобы результат был воспроизводим
	testZone := time.FixedZone("UTC", 0)

	tests := []struct {
		name     string
		input    any
		contains []string // подстроки, которые должны быть в результате
		excludes []string // подстроки, которых быть не должно (например, лишние кавычки)
	}{
		{
			name: "Basic Types",
			input: map[string]any{
				"A_String": "Regular",
				"A_Number": 12345,
				"A_Bool":   true,
			},
			contains: []string{`A_String = "Regular"`, `A_Number = 12345`, `A_Bool = YES`},
		},
		{
			name: "Quota Unlimited",
			input: map[string]any{
				"MaxAccountSize": int64(-1),
			},
			contains: []string{`MaxAccountSize = unlimited`},
			excludes: []string{`"unlimited"`},
		},
		{
			name: "Quota Suffixes",
			input: map[string]any{
				"Giga": int64(30 * 1024 * 1024 * 1024),
				"Mega": int64(100 * 1024 * 1024),
				"Kilo": int64(100 * 1024),
				"Zero": int64(0),
			},
			contains: []string{`Giga = 30G`, `Mega = 100M`, `Kilo = 100K`, `Zero = 0`},
		},
		{
			name: "Quota Exact Bytes",
			input: map[string]any{
				"Exact": int64(31457291), // 30MB + 11 байт
			},
			contains: []string{`Exact = 31457291`},
		},
		{
			name: "Time Formatting",
			input: map[string]any{
				"Created": time.Date(2025, 12, 31, 13, 0, 0, 0, testZone),
			},
			contains: []string{`Created = "31-Dec-2025 13:00:00 +0000"`},
		},
		{
			name: "Binary Blocks",
			input: map[string]any{
				"Data": []byte{0x01, 0x02, 0x0A, 0xFF},
			},
			contains: []string{`Data = #[01020AFF]`},
		},
		{
			name: "DataBlock vs Binary",
			input: map[string]any{
				"AsHex":   []byte("Hello"),          // Обычный слайс байт
				"AsBlock": DataBlock("MIIDlTCC..."), // Наш новый тип
			},
			contains: []string{
				"AsHex = #[48656C6C6F]",   // HEX представление
				"AsBlock = [MIIDlTCC...]", // Datablock (без кавычек!)
			},
			excludes: []string{
				`"[MIIDlTCC...]"`, // Убеждаемся, что кавычек нет
			},
		},
		{
			name: "System Literals",
			input: map[string]any{
				"Mode": "default",
				"Rule": "none",
			},
			contains: []string{`Mode = default`, `Rule = none`},
			excludes: []string{`"default"`, `"none"`},
		},
		{
			name: "Nested and Slices",
			input: map[string]any{
				"List": []any{int64(1024), "unlimited", "text"},
				"Sub": map[string]any{
					"Size": int64(2048),
				},
			},
			contains: []string{`( 1K, unlimited, "text" )`, `Size = 2K`},
		},
		{
			name: "Strict String Map (ACL/Aliases)",
			// Проверяем то, что реально уходит в SETMAILBOXACL
			input: map[string]string{
				"anyone":  "l",
				"admin":   "lirs",
				"trusted": "+rs",
			},
			contains: []string{
				"anyone = \"l\"",
				"admin = \"lirs\"",
				"trusted = \"+rs\"",
			},
		},
		{
			name: "Strict String Slice (Subscriptions)",
			// Проверяем то, что уходит в SETMAILBOXSUBSCRIPTION
			input: []string{"INBOX", "Sent Messages", "Trash"},
			contains: []string{
				"( \"INBOX\", \"Sent Messages\", \"Trash\" )",
			},
		},
		{
			name: "Mixed Case System Literals",
			// CGP понимает YES/yes, но наш маршал должен быть консистентен
			input: map[string]any{
				"Enabled":  true,
				"Disabled": false,
				"Value":    "unlimited",
			},
			contains: []string{
				"Enabled = YES",
				"Disabled = NO",
				"Value = unlimited",
			},
			excludes: []string{"\"YES\"", "\"unlimited\""},
		},
		{
			name: "Nested Map in map[string]any",
			input: map[string]any{
				"Rights": map[string]string{
					"user1": "r",
				},
			},
			contains: []string{
				"Rights = { user1 = \"r\"; }",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Marshal(tt.input)

			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("%s: Marshal() output missing [%s]\nGot: %s", tt.name, s, got)
				}
			}

			for _, s := range tt.excludes {
				if strings.Contains(got, s) {
					t.Errorf("%s: Marshal() output contains forbidden [%s]\nGot: %s", tt.name, s, got)
				}
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
		{1073741824, "1G"},
		{1099511627776, "1T"},
		{1024 * 1024 * 500, "500M"},
		{31457291, "31457291"}, // не делится нацело
	}

	for _, c := range cases {
		got := formatCGPSize(c.in)
		if got != c.out {
			t.Errorf("formatCGPSize(%d) = %s, want %s", c.in, got, c.out)
		}
	}
}

func TestLiterals(t *testing.T) {
	// Тестируем логику КЛЮЧЕЙ
	t.Run("Identifiers", func(t *testing.T) {
		cases := []struct {
			in   string
			want bool
		}{
			{"RealName", true},
			{"Storage_Used", true},
			{"vsu.ru", false},    // точка требует кавычек
			{"my-domain", false}, // тире требует кавычек
		}
		for _, c := range cases {
			if got := isCGPIdentifier(c.in); got != c.want {
				t.Errorf("isCGPIdentifier(%q) = %v, want %v", c.in, got, c.want)
			}
		}
	})

	// Тестируем логику ЗНАЧЕНИЙ
	t.Run("SystemLiterals", func(t *testing.T) {
		cases := []struct {
			in   string
			want bool
		}{
			{"unlimited", true},
			{"YES", true},
			{"10M", true},
			{"text", false}, // обычная строка-значение должна быть в кавычках
		}
		for _, c := range cases {
			if got := isCGPSystemLiteral(c.in); got != c.want {
				t.Errorf("isCGPSystemLiteral(%q) = %v, want %v", c.in, got, c.want)
			}
		}
	})
}
