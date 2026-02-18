package cgpcli

import (
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	// Имитируем сложный ответ от GetAccountSettings или GetDomainSettings
	data := map[string]any{
		"RealName": "Иван Иванов",
		"Storage":  1024 * 1024,
		"Enabled":  true,
		"Aliases":  []string{"ivan", "vanya"},
		"Limits": map[string]any{
			"MaxMessages": 500,
			"Flags":       []byte("ABC"),
		},
	}

	t.Run("ConsoleOutput", func(t *testing.T) {
		PrettyPrint(data)
	})
}
