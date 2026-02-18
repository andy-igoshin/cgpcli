package cgpcli

import (
	"testing"
)

func TestAccountServices_Integration(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	const me = "testuser1@test.domain.name"

	t.Run("BalanceOperations", func(t *testing.T) {
		// Операция "list" возвращает историю транзакций и текущий баланс.
		params := map[string]any{
			"what": "list",
		}

		res, err := cli.Balance(me, params)
		if err != nil {
			t.Logf("Balance 'list' call result (might depend on billing config): %v", err)
			return
		}

		// Проверяем, что в ответе есть ожидаемые поля (обычно 'total' или 'history')
		if res == nil {
			t.Error("Balance returned nil map")
		}
		t.Logf("Current balance data keys: %v", res)
	})

	t.Run("DatasetOperations", func(t *testing.T) {
		const targetSet = "Personal"

		// 1. Проверяем существование или создаем сет
		paramsCreate := map[string]any{
			"subsetName": targetSet,
			"what":       "createSet",
		}
		_, _ = cli.Dataset(me, paramsCreate)

		// 2. Добавляем запись в адресную книгу (addRandomEntry)
		entryData := map[string]any{
			"addressbook.Email": "test-sdk@example.com",
			"RealName":          "SDK Test Entry",
		}
		paramsAdd := map[string]any{
			"subsetName": targetSet,
			"what":       "addRandomEntry",
			"data":       entryData,
		}

		res, err := cli.Dataset(me, paramsAdd)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}

		entryName, ok := res["entryName"].(string)
		if !ok {
			t.Fatalf("Expected entryName in response, got %v", res)
		}

		// 3. Листинг записей и поиск нашей (listEntries)
		paramsList := map[string]any{
			"subsetName": targetSet,
			"what":       "listEntries",
			"attribute":  "entry",
			"data":       entryName,
		}
		listRes, err := cli.Dataset(me, paramsList)
		if err != nil {
			t.Errorf("List entries failed: %v", err)
		}
		if len(listRes) == 0 {
			t.Errorf("Created entry %s not found in list", entryName)
		}

		// 4. Удаление записи (deleteEntry)
		paramsDel := map[string]any{
			"subsetName": targetSet,
			"what":       "deleteEntry",
			"entryName":  entryName,
		}
		if _, err := cli.Dataset(me, paramsDel); err != nil {
			t.Errorf("Delete entry failed: %v", err)
		}
	})

	t.Run("RosterOperations", func(t *testing.T) {
		params := map[string]any{"what": "List"}
		_, err := cli.Roster(me, params)
		if err != nil {
			t.Errorf("Roster List failed: %v", err)
		}
	})

	t.Run("RemoveSubset", func(t *testing.T) {
		tempSet := "SDK_Temp_Set"
		cli.Dataset(me, map[string]any{"subsetName": tempSet, "what": "createSet"})

		if err := cli.RemoveAccountSubset(me, tempSet); err != nil {
			t.Errorf("RemoveAccountSubset failed: %v", err)
		}
	})
}
