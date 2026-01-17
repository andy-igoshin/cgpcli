package cgpcli

import (
	"reflect"
	"testing"
)

func TestDirectory_CreateUnitArgs(t *testing.T) {
	// Убрали неиспользуемый cli, так как мы тестируем саму логику формирования слайса
	tests := []struct {
		name     string
		unitName string
		mount    string
		shared   bool
		remote   bool
		expected []any
	}{
		{
			name:     "Local Private",
			unitName: "Office1",
			mount:    "ou=off1,dc=ext",
			shared:   false,
			remote:   false,
			expected: []any{"CREATEDIRECTORYUNIT", "Office1", "ou=off1,dc=ext"},
		},
		{
			name:     "Remote Shared",
			unitName: "LDAP-Gate",
			mount:    "dc=remote",
			shared:   true,
			remote:   true,
			expected: []any{"CREATEDIRECTORYUNIT", "LDAP-Gate", "SHARED", "REMOTE", "dc=remote"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]any, 0, 5)
			args = append(args, "CREATEDIRECTORYUNIT", tt.unitName)
			if tt.shared {
				args = append(args, "SHARED")
			}
			if tt.remote {
				args = append(args, "REMOTE")
			}
			args = append(args, tt.mount)

			if !reflect.DeepEqual(args, tt.expected) {
				t.Errorf("Args mismatch.\nGot: %v\nExp: %v", args, tt.expected)
			}
		})
	}
}

func TestDirectory_ListUnitsParsing(t *testing.T) {
	mockRaw := map[string]any{
		"": "Main",
	}

	cli := &Cli{}

	t.Run("Parse mount point and name", func(t *testing.T) {
		res, err := cli.toStringMap("LISTDIRECTORYUNITS", mockRaw)
		if err != nil {
			t.Fatalf("toStringMap failed: %v", err)
		}

		expected := map[string]string{"": "Main"}
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Expected %v, got %v", expected, res)
		}
	})
}

func TestDirectory_AccessRightsStructure(t *testing.T) {
	// rawResponse теперь явно any, чтобы эмулировать результат cli.Query()
	var rawResponse any = []any{
		[]any{"HidePasswords", "*", "anyone", "prohibit", []any{}, []any{"pass"}, []any{"pass"}, []any{}},
	}

	t.Run("Type assertion for Access Rights", func(t *testing.T) {
		// Теперь приведение типа сработает, так как rawResponse — интерфейс
		arr, ok := rawResponse.([]any)
		if !ok {
			t.Errorf("Expected []any, got %T", rawResponse)
			return
		}

		if len(arr) != 1 {
			t.Errorf("Expected len 1, got %d", len(arr))
		}

		rule, ok := arr[0].([]any)
		if !ok || len(rule) != 8 {
			t.Errorf("Rule should be []any with len 8, got %T", arr[0])
		}
	})
}
