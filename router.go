package cgpcli

import (
	"fmt"
	"strings"
)

// RouterEntry описывает одну строку в таблице маршрутизации.
type RouterEntry struct {
	Comment string
	Pos     int
	Rule    string
}

// RouterList представляет собой упорядоченный список правил маршрутизации.
type RouterList []RouterEntry

// Add добавляет правило в конец таблицы.
func (l *RouterList) Add(rule, comment string, alignPos int) {
	*l = append(*l, RouterEntry{
		Rule:    rule,
		Comment: comment,
		Pos:     alignPos,
	})
}

// Find возвращает индекс первого вхождения правила, содержащего pattern.
func (l RouterList) Find(pattern string) int {
	for i, entry := range l {
		if strings.Contains(entry.Rule, pattern) {
			return i
		}
	}
	return -1
}

// InsertAfter вставляет новое правило после найденного по паттерну.
func (l *RouterList) InsertAfter(pattern, rule, comment string, alignPos int) error {
	idx := l.Find(pattern)
	if idx == -1 {
		return fmt.Errorf("pattern '%s' not found", pattern)
	}

	entry := RouterEntry{
		Rule:    rule,
		Comment: comment,
		Pos:     alignPos,
	}
	*l = append((*l)[:idx+1], append([]RouterEntry{entry}, (*l)[idx+1:]...)...)
	return nil
}

// InsertBefore вставляет новое правило перед найденным по паттерну.
func (l *RouterList) InsertBefore(pattern, rule, comment string, alignPos int) error {
	idx := l.Find(pattern)
	if idx == -1 {
		return fmt.Errorf("pattern '%s' not found", pattern)
	}

	entry := RouterEntry{
		Rule:    rule,
		Comment: comment,
		Pos:     alignPos,
	}
	*l = append((*l)[:idx], append([]RouterEntry{entry}, (*l)[idx:]...)...)
	return nil
}

// MarshalCGP реализует интерфейс cgpMarshaler для генерации строки в формате CGP CLI.
func (l RouterList) MarshalCGP() string {
	var sb strings.Builder
	for i, entry := range l {
		if i > 0 {
			sb.WriteString(`\e`)
		}

		if entry.Comment == "" {
			sb.WriteString(entry.Rule)
		} else {
			sb.WriteString(entry.Rule)

			currentLen := len(entry.Rule)
			if entry.Rule != "" {
				if entry.Pos > currentLen {
					sb.WriteString(strings.Repeat(" ", entry.Pos-currentLen))
				} else {
					sb.WriteByte(' ')
				}
			}

			sb.WriteByte(';')
			sb.WriteString(entry.Comment)
		}
	}
	return sb.String()
}

// parseRouterLine разбирает одну строку, выделяя правило, комментарий и его позицию.
func parseRouterLine(line string) RouterEntry {
	idx := strings.IndexByte(line, ';')
	if idx == -1 {
		return RouterEntry{
			Rule: strings.TrimRight(line, " "),
			Pos:  0,
		}
	}

	rule := strings.TrimRight(line[:idx], " ")
	comment := line[idx+1:]

	return RouterEntry{
		Rule:    rule,
		Comment: comment,
		Pos:     idx,
	}
}

// parseRouterTable преобразует срез интерфейсов в RouterList.
func parseRouterTable(lines []any) RouterList {
	res := make(RouterList, 0, len(lines))
	for _, l := range lines {
		if s, ok := l.(string); ok {
			res = append(res, parseRouterLine(s))
		}
	}
	return res
}
