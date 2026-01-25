package cgpcli

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestTransformRead(t *testing.T) {
	// Предварительно регистрируем тестовое поле
	RegisterFieldType("TestStringField", TypeString)

	tests := []struct {
		name     string
		key      string
		input    any
		expected any
	}{
		{
			name:     "Date transformation",
			key:      "Created",
			input:    "24-01-2026_21:00:00",
			expected: time.Date(2026, time.January, 24, 21, 0, 0, 0, time.UTC),
		},
		{
			name:     "Duration transformation (never)",
			key:      "Timeout",
			input:    "never",
			expected: time.Duration(-1),
		},
		{
			name:     "Size transformation (10M)",
			key:      "MaxAccountSize",
			input:    "10M",
			expected: int64(10 * 1024 * 1024),
		},
		{
			name:     "IP transformation",
			key:      "LANAddress",
			input:    "[192.168.1.1]",
			expected: net.ParseIP("192.168.1.1"),
		},
		{
			name:     "TypeString override (integer to string)",
			key:      "TestStringField",
			input:    12345,
			expected: "12345",
		},
		{
			name:     "Auto type (no exception)",
			key:      "UnknownField",
			input:    int64(100),
			expected: int64(100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformRead(tt.key, tt.input)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.expected) {
				t.Errorf("transformRead() type = %T, want %T", got, tt.expected)
			}
			if fmt.Sprint(got) != fmt.Sprint(tt.expected) {
				t.Errorf("transformRead() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTransformWrite(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		input   any
		wantRaw string
		wantOk  bool
	}{
		{
			name:    "Size to CGP string",
			key:     "MaxMessageSize",
			input:   int64(1024 * 1024),
			wantRaw: "1M",
			wantOk:  true,
		},
		{
			name:    "Duration to CGP string",
			key:     "Timeout",
			input:   10 * time.Minute,
			wantRaw: "10m",
			wantOk:  true,
		},
		{
			name:    "IP to CGP string",
			key:     "WANAddress",
			input:   net.ParseIP("8.8.8.8"),
			wantRaw: "[8.8.8.8]",
			wantOk:  true,
		},
		{
			name:    "TypeString to string",
			key:     "TestStringField",
			input:   99.9,
			wantRaw: "99.9",
			wantOk:  true,
		},
		{
			name:    "Non-exception field",
			key:     "SomeRegularField",
			input:   "hello",
			wantRaw: "",
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRaw, gotOk := transformWrite(tt.key, tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("transformWrite() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if gotRaw != tt.wantRaw {
				t.Errorf("transformWrite() raw = %v, want %v", gotRaw, tt.wantRaw)
			}
		})
	}
}

func TestTransformRead_EdgeCases(t *testing.T) {
	RegisterFieldType("TestStringField", TypeString)

	tests := []struct {
		name     string
		key      string
		input    any
		expected any
	}{
		{
			name:     "Nil input for TypeString",
			key:      "TestStringField",
			input:    nil,
			expected: "",
		},
		{
			name:     "Invalid date format (fallback to string)",
			key:      "Created",
			input:    "invalid-date",
			expected: "invalid-date",
		},
		{
			name:     "Empty string for Duration",
			key:      "Timeout",
			input:    "",
			expected: "", // Должен вернуть как есть, если не распарсилось
		},
		{
			name:     "Size with 'unlimited' keyword",
			key:      "MaxAccountSize",
			input:    "unlimited",
			expected: int64(-1),
		},
		{
			name:     "IP with missing brackets (should still work if parseSingleIP handles it)",
			key:      "LANAddress",
			input:    "127.0.0.1",
			expected: net.ParseIP("127.0.0.1"),
		},
		{
			name:     "TypeString with boolean input",
			key:      "TestStringField",
			input:    true,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformRead(tt.key, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("%s: got %v (%T), want %v (%T)", tt.name, got, got, tt.expected, tt.expected)
			}
		})
	}
}

func TestTransformWrite_EdgeCases(t *testing.T) {
	RegisterFieldType("TestStringField", TypeString)

	tests := []struct {
		name    string
		key     string
		input   any
		wantRaw string
		wantOk  bool
	}{
		{
			name:    "Size with zero value",
			key:     "MaxMessageSize",
			input:   0,
			wantRaw: "0",
			wantOk:  true,
		},
		{
			name:    "Size with negative value (unlimited)",
			key:     "MaxMessageSize",
			input:   -1,
			wantRaw: "unlimited",
			wantOk:  true,
		},
		{
			name:    "Nil input for TypeString in write",
			key:     "TestStringField",
			input:   nil,
			wantRaw: "",
			wantOk:  true,
		},
		{
			name:    "Duration zero (immediately)",
			key:     "Timeout",
			input:   time.Duration(0),
			wantRaw: "immediately",
			wantOk:  true,
		},
		{
			name:    "IP nil input",
			key:     "WANAddress",
			input:   nil,
			wantRaw: "",
			wantOk:  false, // Реестр не должен обрабатывать nil как валидный IP
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRaw, gotOk := transformWrite(tt.key, tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("%s: gotOk %v, wantOk %v", tt.name, gotOk, tt.wantOk)
			}
			if gotRaw != tt.wantRaw {
				t.Errorf("%s: gotRaw %q, wantRaw %q", tt.name, gotRaw, tt.wantRaw)
			}
		})
	}
}
