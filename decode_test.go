package cgpcli

import (
	"bytes"
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
		{"SimpleBlock", "[data]", []byte("data")},
		{"TaggedBlock", "#[more data]", []byte("more data")},
		{"TaggedIBlock", "#I[binary]", []byte("binary")},
		{"Concatenated", "[part1][part2]", []byte("part1part2")},
		{"WithSpaces", "[part1]  [part2]", []byte("part1part2")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res any
			err := Unmarshal([]byte(tt.input), &res)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			b, ok := res.([]byte)
			if !ok {
				t.Fatalf("expected []byte, got %T", res)
			}
			if !bytes.Equal(b, tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected, b)
			}
		})
	}
}
