// # Router Table Administration
//
// The Router section provides specialized structures and methods for managing the
// CommuniGate Pro Routing Table. The Router is the central engine that determines
// how every message, signal, or call is directed within the server or to
// external destinations.
//
// Key capabilities include:
//   - Ordered Management: [RouterList] ensures that the strict top-to-bottom
//     evaluation order of CommuniGate Pro routing rules is preserved during manipulation.
//   - Pattern-Based Insertion: Methods like [RouterList.InsertBefore] and
//     [RouterList.InsertAfter] allow for precise rule placement without
//     overwriting the entire table manually.
//   - Metadata Preservation: Full support for comments and alignment positions
//     (Pos), ensuring that the table remains readable in the WebAdmin interface.
//   - Format Compatibility: Automatic handling of the internal multi-line
//     representation and special character escaping required by the CLI.
//
// Use these structures to automate complex routing tasks, such as dynamic
// gateway switching, anti-spam redirection, or multi-domain mail flow
// management.
package cgpcli

import (
	"fmt"
	"strings"
)

// RouterEntry represents a single line in the CommuniGate Pro Routing Table.
type RouterEntry struct {
	Rule    string // The routing rule (e.g., "<*@domain.com> = user@other.com").
	Comment string // Optional text after the semicolon.
	Pos     int    // The original character position of the semicolon for alignment.
}

// RouterList is an ordered slice of [RouterEntry] providing high-level
// manipulation methods.
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

// parseRouterTable преобразует строку маршрутов в RouterList.
func parseRouterTable(raw string) RouterList {
	lines := strings.Split(raw, "\n")
	res := make(RouterList, 0, len(lines))
	for _, l := range lines {
		res = append(res, parseRouterLine(l))
	}
	return res
}

func (l RouterList) writeCGP(w encodeWriter) error {
	w.WriteByte('"')
	for i, entry := range l {
		if i > 0 {
			w.WriteString(`\e`)
		}

		// Обрабатываем Rule и Comment через нашу прозрачную функцию
		w.WriteString(replaceSpecChars(entry.Rule))

		if entry.Comment != "" {
			// Если правило есть, нужно разделение. Если правила нет (пустая строка/заголовок),
			// пишем ';' сразу, без ведущего пробела.
			if entry.Rule != "" {
				pad := entry.Pos - len(entry.Rule)
				if pad > 0 {
					w.WriteString(strings.Repeat(" ", pad))
				} else {
					w.WriteByte(' ')
				}
			}
			w.WriteByte(';')
			w.WriteString(replaceSpecChars(entry.Comment))
		}
	}
	w.WriteByte('"')
	return nil
}
