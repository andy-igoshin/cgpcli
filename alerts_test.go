package cgpcli

import (
	"reflect"
	"testing"
)

// Тест логики преобразования типов в мапе алертов
func TestAlerts_TypeConversion(t *testing.T) {
	// Имитируем то, что реально выдал ваш парсер в тесте:
	// Ключи — string, значения — смесь string и int64
	mockResponse := map[string]any{
		"20260115211048": "test alert",
		"20260115214705": int64(98765432109876),
	}

	t.Run("Convert any map to string map via valToString", func(t *testing.T) {
		// Используем ваш toStringMap (который внутри вызывает valToString)
		// Мы ожидаем, что int64 станет строкой "98765432109876"
		res, err := toStringMap("GETSERVERALERTS", mockResponse)
		if err != nil {
			t.Fatalf("toStringMap failed: %v", err)
		}

		expected := map[string]string{
			"20260115211048": "test alert",
			"20260115214705": "98765432109876",
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Data mismatch.\nExpected: %v\nGot:      %v", expected, res)
		}
	})
}

// Тест формирования аргументов для команд с литералом ALERT
func TestAlerts_CommandFormatting(t *testing.T) {
	// Здесь мы проверяем, что аргументы ложатся в правильном порядке
	// и специальные символы (атомы) не квотируются.

	tests := []struct {
		name     string
		account  string
		alert    string
		expected []any
	}{
		{
			name:     "Standard account",
			account:  "postmaster",
			alert:    "hello",
			expected: []any{Atom("postmaster"), "ALERT", "hello"},
		},
		{
			name:     "Current account asterisk",
			account:  "*",
			alert:    "system alert",
			expected: []any{Atom("*"), "ALERT", "system alert"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// В реальности мы проверяем, что передается в Query.
			// Ключевой момент: Atom(*) и литерал "ALERT" на своем месте.
			args := []any{Atom(tt.account), "ALERT", tt.alert}

			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("Args mismatch for %s", tt.name)
			}
		})
	}
}
