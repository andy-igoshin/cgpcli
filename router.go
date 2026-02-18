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
//   - Pattern-Based Insertion: methods like [RouterList.InsertBefore] and
//     [RouterList.InsertAfter] allow for precise rule placement without
//     overwriting the entire table manually.
//   - Metadata Preservation: full support for comments and alignment positions
//     (Pos), ensuring that the table remains readable in the WebAdmin interface.
//   - Format Compatibility: automatic handling of the internal multi-line
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
	Rule    string // routing rule (e.g., "<*@domain.com> = user@other.com").
	Comment string // optional text after the semicolon.
	Pos     int    // original character position of the semicolon for alignment.
}

// RouterList is an ordered slice of [RouterEntry] providing high-level
// manipulation methods. Since CommuniGate Pro evaluates rules from
// top to bottom, this structure ensures the order remains intact.
type RouterList []RouterEntry

// Add appends a new rule to the end of the routing table.
//
// Parameters:
//   - rule: the routing rule string (e.g., "<*@domain.com> = user@other.com").
//   - comment: optional text to be placed after the semicolon.
//   - alignPos: the character position where the semicolon should be placed.
//
// Returns:
//   - None.
func (l *RouterList) Add(rule, comment string, alignPos int) {
	*l = append(*l, RouterEntry{
		Rule:    rule,
		Comment: comment,
		Pos:     alignPos,
	})
}

// Find returns the index of the first entry where the rule contains the pattern.
//
// Parameters:
//   - pattern: the substring to search for within the routing rules.
//
// Returns:
//   - int: the zero-based index of the first match, or -1 if no match is found.
func (l RouterList) Find(pattern string) int {
	for i, entry := range l {
		if strings.Contains(entry.Rule, pattern) {
			return i
		}
	}
	return -1
}

// InsertAfter inserts a new rule immediately after the first entry that
// matches the specified pattern.
//
// Parameters:
//   - pattern: the substring to locate the anchor rule.
//   - rule: the new routing rule string to insert.
//   - comment: optional text for the new rule's comment.
//   - alignPos: the character position for semicolon alignment in the new rule.
//
// Returns:
//   - error: an error if the anchor pattern is not found in the list.
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

// InsertBefore inserts a new rule immediately before the first entry that
// matches the specified pattern.
//
// Parameters:
//   - pattern: the substring to locate the anchor rule.
//   - rule: the new routing rule string to insert.
//   - comment: optional text for the new rule's comment.
//   - alignPos: the character position for semicolon alignment in the new rule.
//
// Returns:
//   - error: an error if the anchor pattern is not found in the list.
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
