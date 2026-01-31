package cgpcli

import (
	"testing"
)

func TestMailingLists_Integration(t *testing.T) {
        cli := getTestCli(t)
        defer cli.Logout()

	const (
		testDomain = "test.domain.name"
		testOwner  = "testuser1@test.domain.name"
		testList   = "test-sdk-list@test.domain.name"
		testSub    = "testuser2@test.domain.name"
	)

	// 1. Создание (проверка обрезки домена)
	t.Run("CreateList_AutoTruncate", func(t *testing.T) {
		err := cli.CreateList(testList, testOwner)
		if err != nil {
			t.Fatalf("CreateList failed with full email: %v", err)
		}
	})

	// 2. Тест на корректность маппинга структур
	t.Run("GetSubscriberInfo_Types", func(t *testing.T) {
		// Сначала подпишем
		if err := cli.List(testList, "feed", testSub, true, false); err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		info, err := cli.GetSubscriberInfo(testList, testSub)
		if err != nil {
			t.Fatalf("GetSubscriberInfo failed: %v", err)
		}
		if info == nil {
			t.Fatal("subscriber not found")
		}

		// Проверка ConfirmationID (int64)
		if info.ConfirmationID == 0 {
			t.Log("Warning: ConfirmationID is 0, check if this is expected for direct subscription")
		}

		// Проверка дат
		if info.SubscribeTime.IsZero() {
			t.Error("SubscribeTime was not parsed correctly")
		}
	})

	// 3. Тест GetList / UpdateList (трансляция текста)
	t.Run("TextTranslation", func(t *testing.T) {
		text := "First Line\nSecond Line"
		err := cli.UpdateList(testList, map[string]any{"PolicyText": text})
		if err != nil {
			t.Fatalf("UpdateList failed: %v", err)
		}

		res, err := cli.GetList(testList)
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}

		if res["PolicyText"] != text {
			t.Errorf("Translation failed: expected %q, got %q", text, res["PolicyText"])
		}
	})

	// 4. Очистка
	t.Run("Cleanup", func(t *testing.T) {
		if err := cli.DeleteList(testList); err != nil {
			t.Errorf("DeleteList failed: %v", err)
		}
	})
}

func TestMapToSubscriber_Unit(t *testing.T) {
	data := map[string]any{
		"Sub":            "postmaster@test.domain.name",
		"RealName":       `System\eEngineer`,
		"ConfirmationID": int64(987654321),
		"timeSubscribed": int64(20260118180000),
	}

	res := mapToSubscriber(data)

	// Теперь ожидаем "как есть"
	if res.RealName != `System\eEngineer` {
		t.Errorf("RealName changed: expected %q, got %q", `System\eEngineer`, res.RealName)
	}

	if res.ConfirmationID != 987654321 {
		t.Errorf("ConfirmationID failed: %v", res.ConfirmationID)
	}

	if res.TimeSubscribed.Year() != 2026 {
		t.Errorf("TimeSubscribed parsing failed: %v", res.TimeSubscribed)
	}
}
