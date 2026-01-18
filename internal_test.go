package cgpcli

import (
	"os"
	"reflect"
	"testing"
)

// getTestCli перенесен сюда как в базовый файл тестов internal.
// Убедитесь, что вы удалили это определение из accounts_test.go, чтобы избежать redeclared error.
func getTestCli(t *testing.T) *Cli {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Skip("CGPUSER or CGPPASS not set, skipping integration test")
	}
	// Используем APOP для безопасности в тестах
	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Fatalf("Connection failed: %v", err)
	}
	return cli
}

// --- Unit Tests (Логика конвертации без сети) ---

func TestToString(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{"string", "hello", "hello", false},
		{"int", 123, "123", false},
		{"int64", int64(456), "456", false},
		{"float64", 789.0, "789", false},
		{"invalid type", []int{1}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toString("TEST", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("toString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	cli := &Cli{}
	tests := []struct {
		input any
		want  int
	}{
		{10, 10},
		{int64(20), 20},
		{30.0, 30},
		{"not a number", 0},
	}

	for _, tt := range tests {
		if got := cli.toInt(tt.input); got != tt.want {
			t.Errorf("toInt(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestToStringMap(t *testing.T) {
	cli := &Cli{}
	input := map[string]any{
		"Key1": "StringValue",
		"Key2": 2026,
	}
	want := map[string]string{
		"Key1": "StringValue",
		"Key2": "2026",
	}

	got, err := cli.toStringMap("TEST", input)
	if err != nil {
		t.Fatalf("toStringMap failed: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("toStringMap() got = %v, want %v", got, want)
	}
}

// --- Integration Tests (Работа с живым сервером) ---

func TestInternalIntegration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.close()

	// 1. Сначала получаем список доменов, чтобы знать, к кому обращаться
	var targetDomain string
	t.Run("getSliceString_ListDomains", func(t *testing.T) {
		domains, err := cli.getSliceString("LISTDOMAINS", "")
		if err != nil {
			t.Fatalf("getSliceString failed: %v", err)
		}
		if len(domains) == 0 {
			t.Fatal("expected at least one domain on server")
		}
		targetDomain = domains[0] // Берем первый доступный домен
	})

	// 2. Тестируем получение мапы, используя живое имя домена
	t.Run("getMapAny_DomainSettings", func(t *testing.T) {
		if targetDomain == "" {
			t.Skip("No target domain available")
		}
		res, err := cli.getMapAny("GETDOMAINSETTINGS", targetDomain)
		if err != nil {
			t.Errorf("getMapAny failed for domain %s: %v", targetDomain, err)
		}
		if res == nil {
			t.Error("expected map, got nil")
		}
	})

	// 3. Тестируем отправку данных (UpdateInfo)
	t.Run("setMapAny_UpdateAccountInfo", func(t *testing.T) {
		testData := map[string]any{
			"IntegrationTest": "Passed",
			"Year":            2026,
		}
		err := cli.setMapAny("UPDATEACCOUNTINFO", "*", testData)
		if err != nil {
			t.Errorf("setMapAny failed: %v", err)
		}
	})

	// 4. Тестируем комплексную конвертацию ListAccounts на живом домене
	t.Run("LISTACCOUNTS_Mapping", func(t *testing.T) {
		if targetDomain == "" {
			t.Skip("No target domain available")
		}
		raw, err := cli.getMapAny("LISTACCOUNTS", targetDomain)
		if err != nil {
			t.Fatalf("LISTACCOUNTS call failed for %s: %v", targetDomain, err)
		}

		accounts, err := cli.toStringMap("LISTACCOUNTS", raw)
		if err != nil {
			t.Errorf("toStringMap conversion failed: %v", err)
		}

		for name, accType := range accounts {
			if name == "" || accType == "" {
				t.Errorf("empty key or value for account %s", name)
			}
		}
	})
}
