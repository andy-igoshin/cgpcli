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
			"Flags":       []byte("ABC"), // Проверка []byte
		},
	}

	t.Run("ConsoleOutput", func(t *testing.T) {
		// Мы просто визуально проверяем, что не паникует
		// В реальных тестах можно перехватить stdout, но для дебаг-функции достаточно прогона
		PrettyPrint(data)
	})
}
