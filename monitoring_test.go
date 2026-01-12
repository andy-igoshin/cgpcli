package cgpcli

import (
	"fmt"
	"testing"
)

func TestMonitoring(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	t.Run("GetSystemStat", func(t *testing.T) {
		// Проверяем базовый элемент, который есть всегда
		val, err := cli.GetStatElement("1.3.6.1.2.1.1.1")
		if err != nil {
			t.Errorf("Failed to get system stat: %v", err)
		}
		t.Logf("Version: %v", val)
	})

	t.Run("CustomStats", func(t *testing.T) {
		// Custom элементы должны начинаться с "Custom."
		obj := "Custom.MyTestCounter"

		// Пытаемся установить значение
		err := cli.SetStatElement(obj, "SET", 100)
		if err != nil {
			// Может не быть прав или Custom элементы не поддерживаются в данной версии/лицензии
			t.Skipf("SetStatElement skipped: %v", err)
			return
		}

		// Прибавляем
		_ = cli.SetStatElement(obj, "INC", 5)

		val, _ := cli.GetStatElement(obj)
		if fmt.Sprint(val) != "105" {
			t.Errorf("Expected 105, got %v", val)
		}
	})

	t.Run("EnumerateStats", func(t *testing.T) {
		first, err := cli.GetNextStatName("")
		if err != nil {
			t.Fatalf("GetNextStatName(first) failed: %v", err)
		}
		if first == "" {
			t.Error("Returned empty first stat name")
		}

		second, _ := cli.GetNextStatName(first)
		if second == "" || second == first {
			t.Error("Enumeration logic failed")
		}
	})
}

func TestIdentifyEOF(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Logout()

	// 1. Сначала найдем последний элемент
	current := ""
	var lastValid string
	for {
		next, err := cli.GetNextStatName(current)
		if err != nil {
			t.Logf("Protocol error at the end: %v", err)
			break
		}
		lastValid = next
		current = next
	}

	// 2. Теперь вызываем GETNEXTSTATNAME от последнего валидного элемента
	// чтобы увидеть точный текст ошибки от сервера
	res, err := cli.Query("GETNEXTSTATNAME", lastValid)
	if err != nil {
		t.Logf("EOF Error Details: '%v'", err)
	} else {
		t.Errorf("Expected error at EOF, but got: %v", res)
	}
}
