// # IP List Administration
//
// The IP List section provides specialized types and methods for handling CommuniGate Pro
// Network Lists (Client IP addresses, Blacklists, etc.). It implements the unique
// CommuniGate Pro format that allows for mixed IP addresses, CIDR masks, ranges, and sequences
// within a single list, including support for negative (exclusion) rules and comments.
//
// Key capabilities include:
//   - Flexible Management: Adding, inserting, and deleting entries with [Cli.Add],
//     [Cli.InsertBefore], and [Cli.InsertAfter] to maintain strict rule ordering.
//   - Rule Logic: Full support for exclusion rules (prefixed with "!") and complex
//     entry types like [IPRange] and [IPSeq].
//   - IP Matching: The [IPList.Contains] method allows for high-performance
//     checks to determine if a specific net.IP is matched by any rule in the list.
//   - Format Preservation: Methods for parsing and writing lists while maintaining
//     comment alignment and special character encoding required by the CommuniGate Pro CLI.
//
// This section is essential for managing security settings, relay permissions,
// and access control lists based on network topology.
package cgpcli

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type IPCIDR struct{ net.IPNet }

// IPEntry represents a single line in a CommuniGate Pro IP List, including the value,
// optional comment, and whether it is an exclusion rule.
type IPEntry struct {
	Comment    string // Human-readable comment after the semicolon.
	CommentPos int    // Original position of the semicolon for alignment preservation.
	IsNegative bool   // True if the entry starts with "!" (exclusion).
	Value      any    // Can be IPCIDR, IPRange, IPSeq, or IPUnknown.
}

type (
	IPList    []IPEntry
	IPRange   struct{ Start, End net.IP } // IPRange represents a contiguous range of IP addresses from Start to End.
	IPSeq     []any                       // IPSeq represents a comma-separated sequence of multiple IP rules.
	IPUnknown string
)

func (l *IPList) Add(input string, comment string) error {
	isNegative := false
	trimmedInput := strings.TrimSpace(input)

	// Распознаем "!" в начале входной строки
	if strings.HasPrefix(trimmedInput, "!") {
		isNegative = true
		trimmedInput = strings.TrimSpace(trimmedInput[1:])
	}

	val := parseDataPart(trimmedInput)
	if _, ok := val.(IPUnknown); ok {
		return fmt.Errorf("invalid IP format: %s", input)
	}

	inputNormalized := dataToString(val)

	for i, entry := range *l {
		if dataToString(entry.Value) == inputNormalized {
			(*l)[i].Comment = comment
			(*l)[i].IsNegative = isNegative
			return nil
		}
	}

	*l = append(*l, IPEntry{
		Value:      val,
		Comment:    comment,
		IsNegative: isNegative,
	})
	return nil
}

func (l IPList) Contains(ip net.IP) bool {
	if ip == nil {
		return false
	}
	matched := false
	for _, entry := range l {
		if matchAny(entry.Value, ip) {
			if entry.IsNegative {
				matched = false
			} else {
				matched = true
			}
		}
	}
	return matched
}

func (l *IPList) Delete(input string) error {
	cleanInput := strings.TrimPrefix(strings.TrimSpace(input), "!")
	val := parseDataPart(cleanInput)

	if _, ok := val.(IPUnknown); ok {
		return fmt.Errorf("invalid IP format: %s", input)
	}

	inputNormalized := dataToString(val)

	for i, entry := range *l {
		if dataToString(entry.Value) == inputNormalized {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("entry not found: %s", input)
}

func (l IPList) GetEntry(input string) *IPEntry {
	idx := l.indexOf(input)
	if idx == -1 {
		return nil
	}
	return &l[idx]
}

func (l *IPList) Insert(index int, input string, comment string, alignPos int) error {
	isNegative := false
	trimmedInput := strings.TrimSpace(input)

	if strings.HasPrefix(trimmedInput, "!") {
		isNegative = true
		trimmedInput = strings.TrimSpace(trimmedInput[1:])
	}

	val := parseDataPart(trimmedInput)
	if _, ok := val.(IPUnknown); ok {
		return fmt.Errorf("invalid IP format: %s", input)
	}

	_ = l.Delete(trimmedInput)

	newEntry := IPEntry{
		Value:      val,
		Comment:    comment,
		CommentPos: alignPos,
		IsNegative: isNegative,
	}

	if index < 0 {
		index = 0
	}
	if index > len(*l) {
		index = len(*l)
	}

	*l = append((*l)[:index], append([]IPEntry{newEntry}, (*l)[index:]...)...)
	return nil
}

func (l *IPList) InsertAfter(marker string, input string, comment string, alignPos int) error {
	idx := l.indexOf(marker)
	if idx == -1 {
		return fmt.Errorf("marker not found: %s", marker)
	}
	return l.Insert(idx+1, input, comment, alignPos)
}

func (l *IPList) InsertBefore(marker string, input string, comment string, alignPos int) error {
	idx := l.indexOf(marker)
	if idx == -1 {
		return fmt.Errorf("marker not found: %s", marker)
	}
	return l.Insert(idx, input, comment, alignPos)
}

func dataToString(val any) string {
	switch v := val.(type) {
	case IPCIDR:
		ones, bits := v.Mask.Size()
		if ones == bits {
			return v.IP.String()
		}
		return v.IP.String() + "/" + strconv.Itoa(ones)
	case IPRange:
		return v.Start.String() + "-" + v.End.String()
	case IPSeq:
		parts := make([]string, len(v))
		for i, sub := range v {
			parts[i] = dataToString(sub)
		}
		return strings.Join(parts, ",")
	case IPUnknown:
		return string(v)
	default:
		return ""
	}
}

func (l IPList) indexOf(input string) int {
	cleanInput := strings.TrimPrefix(strings.TrimSpace(input), "!")
	val := parseDataPart(cleanInput)

	if targetCIDR, ok := val.(IPCIDR); ok {
		for i, entry := range l {
			if matchAny(entry.Value, targetCIDR.IP) {
				return i
			}
		}
	}

	clean := dataToString(val)
	for i, entry := range l {
		if dataToString(entry.Value) == clean {
			return i
		}
	}
	return -1
}

func matchAny(val any, ip net.IP) bool {
	switch v := val.(type) {
	case IPCIDR:
		return v.Contains(ip)
	case IPRange:
		return bytes.Compare(ip, v.Start) >= 0 && bytes.Compare(ip, v.End) <= 0
	case IPSeq:
		for _, sub := range v {
			if matchAny(sub, ip) {
				return true
			}
		}
	case IPUnknown:
		if target := net.ParseIP(string(v)); target != nil {
			return target.Equal(ip)
		}
	}
	return false
}

func parseDataPart(s string) any {
	s = strings.TrimSpace(s)
	if strings.Contains(s, ",") {
		parts := strings.Split(s, ",")
		seq := make(IPSeq, 0, len(parts))
		for _, p := range parts {
			seq = append(seq, parseDataPart(p))
		}
		return seq
	}

	if idx := strings.Index(s, "-"); idx != -1 {
		start := net.ParseIP(strings.TrimSpace(s[:idx]))
		end := net.ParseIP(strings.TrimSpace(s[idx+1:]))
		if start != nil && end != nil {
			return IPRange{Start: start, End: end}
		}
	}

	if !strings.Contains(s, "/") {
		if ip := net.ParseIP(s); ip != nil {
			if ip4 := ip.To4(); ip4 != nil {
				ip = ip4
			}
			mask := net.CIDRMask(len(ip)*8, len(ip)*8)
			return IPCIDR{net.IPNet{IP: ip, Mask: mask}}
		}
	} else if _, ipnet, err := net.ParseCIDR(s); err == nil {
		return IPCIDR{*ipnet}
	}

	return IPUnknown(s)
}

func parseIPList(raw string) IPList {
	lines := strings.Split(raw, "\n")
	list := make(IPList, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		entry := IPEntry{}
		idx := strings.IndexByte(line, ';')
		dataPart := line

		if idx != -1 {
			entry.CommentPos = idx
			entry.Comment = line[idx+1:]
			dataPart = line[:idx]
		}

		trimmedData := strings.TrimSpace(dataPart)
		if strings.HasPrefix(trimmedData, "!") {
			entry.IsNegative = true
			trimmedData = strings.TrimSpace(trimmedData[1:])
		}

		if trimmedData != "" {
			entry.Value = parseDataPart(trimmedData)
			list = append(list, entry)
		} else if entry.Comment != "" {
			list = append(list, entry)
		}
	}
	return list
}

func (l IPList) writeCGP(w encodeWriter) (err error) {
	if err = w.WriteByte('"'); err != nil {
		return
	}

	for i, entry := range l {
		if i > 0 {
			if _, err = w.WriteString(`\e`); err != nil {
				return
			}
		}

		if entry.IsNegative {
			if err = w.WriteByte('!'); err != nil {
				return
			}
		}

		if err = writeDataToString(w, entry.Value); err != nil {
			return
		}

		if entry.Comment != "" {
			dataStrLen := len(dataToString(entry.Value))
			if entry.IsNegative {
				dataStrLen++
			}

			if entry.CommentPos > dataStrLen {
				w.WriteString(strings.Repeat(" ", entry.CommentPos-dataStrLen))
			} else {
				w.WriteByte(' ')
			}

			if err = w.WriteByte(';'); err != nil {
				return
			}

			if _, err = w.WriteString(replaceSpecChars(entry.Comment)); err != nil {
				return
			}
		}
	}

	return w.WriteByte('"')
}

// writeDataToString пишет строковое представление IP-данных напрямую в writer
func writeDataToString(w encodeWriter, val any) error {
	switch v := val.(type) {
	case IPCIDR:
		ones, bits := v.Mask.Size()
		w.WriteString(v.IP.String())
		if ones != bits {
			w.WriteByte('/')
			w.WriteString(strconv.Itoa(ones))
		}
	case IPRange:
		w.WriteString(v.Start.String())
		w.WriteByte('-')
		w.WriteString(v.End.String())
	case IPSeq:
		for i, sub := range v {
			if i > 0 {
				w.WriteByte(',')
			}
			if err := writeDataToString(w, sub); err != nil {
				return err
			}
		}
	case IPUnknown:
		w.WriteString(string(v))
	}
	return nil
}
