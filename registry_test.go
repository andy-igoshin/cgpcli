package cgpcli

import (
	"bufio"
	"bytes"
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
		name     string
		key      string
		input    any
		validate func(any) bool
	}{
		{
			name:  "Size exception",
			key:   "MaxMessageSize",
			input: int64(1024 * 1024),
			validate: func(v any) bool {
				val, ok := v.(cgpSize)
				return ok && int64(val) == 1024*1024
			},
		},
		{
			name:  "Duration exception",
			key:   "Timeout",
			input: 10 * time.Minute,
			validate: func(v any) bool {
				val, ok := v.(cgpDuration)
				return ok && time.Duration(val) == 10*time.Minute
			},
		},
		{
			name:  "IPSB exception",
			key:   "WANAddress",
			input: net.ParseIP("8.8.8.8"),
			validate: func(v any) bool {
				val, ok := v.(cgpIPSB)
				return ok && net.IP(val).Equal(net.ParseIP("8.8.8.8"))
			},
		},
		{
			name:  "ForceString exception",
			key:   "SomeStringField", // Добавьте его в реестр как TypeString для теста
			input: 99.9,
			validate: func(v any) bool {
				val, ok := v.(cgpForceString)
				return ok && string(val) == "99.9"
			},
		},
		{
			name:  "No exception (pass-through)",
			key:   "RegularField",
			input: "hello",
			validate: func(v any) bool {
				val, ok := v.(string)
				return ok && val == "hello"
			},
		},
	}

	// Регистрируем тестовое поле для TypeString
	RegisterFieldType("SomeStringField", TypeString)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformWrite(tt.key, tt.input)
			if !tt.validate(got) {
				t.Errorf("%s: transformWrite returned unexpected type or value: %T(%v)", tt.name, got, got)
			}
		})
	}
}

func TestMarshal_EdgeCases(t *testing.T) {
	// Регистрируем поля для теста
	RegisterFieldType("MaxMessageSize", TypeSize)
	RegisterFieldType("Timeout", TypeDuration)
	RegisterFieldType("TestStringField", TypeString)
	RegisterFieldType("WANAddress", TypeIPSB)

	tests := []struct {
		name  string
		key   string
		input any
		want  string
	}{
		{
			name:  "Size zero",
			key:   "MaxMessageSize",
			input: 0,
			want:  "0",
		},
		{
			name:  "Size unlimited",
			key:   "MaxMessageSize",
			input: -1,
			want:  "unlimited",
		},
		{
			name:  "Duration immediately",
			key:   "Timeout",
			input: time.Duration(0),
			want:  "immediately",
		},
		{
			name:  "TypeString nil",
			key:   "TestStringField",
			input: nil,
			want:  `""`, // marshal nil через toString даст "", а Quote даст ""
		},
		{
			name:  "IP nil (pass-through to marshal)",
			key:   "WANAddress",
			input: nil,
			want:  "#NULL#", // transformWrite вернет nil, marshal обработает как nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)

			err := marshal(w, tt.input, tt.key)
			w.Flush()

			if err != nil {
				t.Errorf("%s: marshal error: %v", tt.name, err)
			}
			if got := b.String(); got != tt.want {
				t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
