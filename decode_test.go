package cgpcli

import (
	"bytes"
	"math"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	// Фиксируем зону для тестов дат
	testTime := time.Date(2025, 12, 31, 13, 0, 15, 0, time.UTC)
	// Строка времени для Created в формате CGP: 31-12-2025_13:00:15

	tests := []struct {
		name    string
		input   string
		want    any
		wantErr bool
	}{
		{
			name:  "ListDomains",
			input: `(vsu.ru, ai.vsu.ru, "test.domain.name")`,
			want:  []any{"vsu.ru", "ai.vsu.ru", "test.domain.name"},
		},
		{
			name:  "GetAccountProperties",
			input: `{ RealName = "Ivan Ivanov"; ServiceMask = *; StorageUsed = 1024; }`,
			want: map[string]any{
				"RealName":    "Ivan Ivanov",
				"ServiceMask": "*",
				"StorageUsed": int64(1024), // Теперь это int64
			},
		},
		{
			name:  "Created_Date_Exception",
			input: `{ Created = 31-12-2025_13:00:15; }`,
			want: map[string]any{
				"Created": testTime,
			},
		},
		{
			name:  "Quota_Unlimited",
			input: `{ MaxAccountSize = unlimited; }`,
			want: map[string]any{
				"MaxAccountSize": int64(-1),
			},
		},
		{
			name:  "Quota_Suffix_G",
			input: `{ MaxWebSize = 30G; }`,
			want: map[string]any{
				"MaxWebSize": int64(30 * 1024 * 1024 * 1024),
			},
		},
		{
			name:  "Octal_Cyrillic",
			input: `"\317\360\350\342\345\362"`, // Привет в Win1251
			want:  "\317\360\350\342\345\362",
		},
		{
			name:  "Quoted_Concatenation",
			input: `"Part1" "Part2"`,
			want:  "Part1Part2",
		},
		{
			name:  "Complex_XML",
			input: `<root><item id="1" desc="a > b" /><nested><![CDATA[ < > ]]></nested></root>`,
			want:  `<root><item id="1" desc="a > b" /><nested><![CDATA[ < > ]]></nested></root>`,
		},
		{
			name:  "Deep_Recursion",
			input: `{ a = { b = (1, 2, { c = 3 }); }; }`,
			want: map[string]any{
				"a": map[string]any{
					"b": []any{int64(1), int64(2), map[string]any{"c": int64(3)}},
				},
			},
		},
		{
			name:  "Tagged_Long_and_Time",
			input: `{ Seq = #1767165; Last = #T31-12-2025_13:00:15; }`,
			want: map[string]any{
				"Seq":  int64(1767165),
				"Last": testTime,
			},
		},
		{
			name:  "Empty_Structures",
			input: `{ dict = {}; list = (); }`,
			want: map[string]any{
				"dict": map[string]any{},
				"list": []any{},
			},
		},
		{
			name:  "CGP_IP_Address",
			input: "#I[127.0.0.1]:8080",
			want:  CGPIP{IP: net.ParseIP("127.0.0.1"), Port: 8080},
		},
		{
			name:  "CGP_IPv6_Address",
			input: "#I[::1]",
			want:  CGPIP{IP: net.ParseIP("::1"), Port: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			err := Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: Unmarshal() got = %v (%T), want %v (%T)", tt.name, got, got, tt.want, tt.want)
			}
		})
	}
}

func TestDecodeDataBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		// "ZGF0YQ==" это "data" в Base64
		{"SimpleBlock", "[ZGF0YQ==]", []byte("data")},

		// "cGFydDE=" (part1) + "cGFydDI=" (part2)
		{"Concatenated", "[cGFydDE=][cGFydDI=]", []byte("part1part2")},

		// С пробелами между блоками
		{"WithSpaces", "[cGFydDE=]  [cGFydDI=]", []byte("part1part2")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res any
			// Используем Unmarshal, который вызывает readCgpBlock
			err := Unmarshal([]byte(tt.input), &res)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			b, ok := res.([]byte)
			if !ok {
				t.Fatalf("expected []byte, got %T", res)
			}

			if !bytes.Equal(b, tt.expected) {
				t.Errorf("expected %q, got %q (hex: %x)", tt.expected, b, b)
			}
		})
	}
}

// TestArrayCounting проверяет корректность подсчёта элементов в arrayInterface.
func TestArrayCounting(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"Empty", "()", 0},
		{"One", "(a)", 1},
		{"Two", "(a,b)", 2},
		{"Three", "(a,b,c)", 3},
		{"Nested", "(a,(b,c),d)", 3},
		{"NestedDict", "(a,{k=v;},c)", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			if err := Unmarshal([]byte(tt.input), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			arr, ok := got.([]any)
			if !ok {
				if tt.want == 0 {
					return
				}
				t.Fatalf("expected []any, got %T: %v", got, got)
			}
			if len(arr) != tt.want {
				t.Errorf("len=%d, want %d, got %v", len(arr), tt.want, arr)
			}
		})
	}
}

// TestParseCGPSizeQuotes проверяет нужен ли strings.Trim в parseCGPSize.
func TestParseCGPSizeQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  any
	}{
		{"NoQuotes", "10K", int64(10 * 1024)},
		{"WithQuotes", `"10K"`, int64(10 * 1024)},
		{"Unlimited", "unlimited", int64(-1)},
		{"Plain", "1024", int64(1024)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCGPSize(tt.input)
			if got != tt.want {
				t.Errorf("parseCGPSize(%q) = %v (%T), want %v (%T)", tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

// TestParseIPSBSA проверяет parseIPSBSA.
func TestParseIPSBSA(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"Empty", "", 0},
		{"Single", "[192.168.1.1]", 1},
		{"Multiple", "[192.168.1.1], [10.0.0.1]", 2},
		{"IPv6", "[::1], [fe80::1]", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseIPSBSA(tt.input)
			if len(got) != tt.want {
				t.Errorf("parseIPSBSA(%q) = %d IPs, want %d: %v", tt.input, len(got), tt.want, got)
			}
		})
	}
}

// TestParseInt64Boundaries проверяет граничные значения int64 при разборе чисел.
// Парсер должен корректно принимать MinInt64/MaxInt64 и откатываться к строке
// при переполнении как для обычных, так и для тегированных (#…) чисел.
func TestParseInt64Boundaries(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  any
	}{
		{"MaxInt64", "9223372036854775807", int64(math.MaxInt64)},
		{"MinInt64", "-9223372036854775808", int64(math.MinInt64)},
		{"OverflowPositive", "9223372036854775808", "9223372036854775808"},
		{"OverflowNegative", "-9223372036854775809", "-9223372036854775809"},
		{"TaggedMaxInt64", "#9223372036854775807", int64(math.MaxInt64)},
		{"TaggedMinInt64", "#-9223372036854775808", int64(math.MinInt64)},
		{"TaggedOverflowPositive", "#9223372036854775808", "#9223372036854775808"},
		{"TaggedOverflowNegative", "#-9223372036854775809", "#-9223372036854775809"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			if err := Unmarshal([]byte(tt.input), &got); err != nil {
				t.Fatalf("Unmarshal(%q) error = %v", tt.input, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal(%q) = %v (%T), want %v (%T)", tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

// FuzzUnmarshal проверяет устойчивость парсера к произвольному (в т.ч. malformed) входу.
// Инвариант: Unmarshal не должен паниковать ни на каких входных данных.
func FuzzUnmarshal(f *testing.F) {
	seeds := []string{
		`(vsu.ru, ai.vsu.ru)`,
		`{ RealName = "Ivan"; StorageUsed = 1024; }`,
		`{ Created = 31-12-2025_13:00:15; MaxAccountSize = unlimited; }`,
		`[ZGF0YQ==]`,
		`[cGFydDE=][cGFydDI=]`,
		`#I[127.0.0.1]:8080`,
		`#TPAST`,
		`<root><item id="1"/></root>`,
		`{ a = { b = (1, 2, { c = 3 }); }; }`,
		`"Part1" "Part2"`,
		`#-9223372036854775809`,
		`"\317\360\350\342\345\362"`,
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var v any
		_ = Unmarshal(data, &v)
	})
}
