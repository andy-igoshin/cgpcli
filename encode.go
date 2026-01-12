package cgpcli

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Atom string
type DataBlock []byte

type encodeWriter interface {
	io.Writer
	WriteByte(byte) error
	WriteString(string) (int, error)
}

var cgpSafeTable = [256]byte{
	'0': 1, '1': 1, '2': 1, '3': 1, '4': 1, '5': 1, '6': 1, '7': 1, '8': 1, '9': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1,
	'K': 1, 'L': 1, 'M': 1, 'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1,
	'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1,
	'k': 1, 'l': 1, 'm': 1, 'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1,
	'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
	'_': 1,
}

// isCGPIdentifier проверяет, может ли строка быть КЛЮЧОМ без кавычек.
// Разрешены только буквы, цифры и подчеркивание. Точки и тире требуют кавычек.
func isCGPIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

// isCGPSystemLiteral проверяет, может ли строка быть ЗНАЧЕНИЕМ без кавычек.
// Разрешены системные константы и числовые квоты с суффиксами.
func isCGPSystemLiteral(s string) bool {
	if s == "" {
		return false
	}

	switch strings.ToLower(s) {
	case "default", "unlimited", "yes", "no", "none", "any":
		return true
	}

	// Проверка на число или квоту (например, 10G)
	isNumeric := true
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		// Суффикс допустим только в конце строки
		if i == len(s)-1 && strings.ContainsRune("KMGTkmgt", r) {
			continue
		}
		isNumeric = false
		break
	}
	return isNumeric
}

// isCGPSafeString проверяет, можно ли передать строку без кавычек.
// В CGP это буквенно-цифровые строки без спецсимволов.
func isCGPSafeString(s string) bool {
	if s == "" {
		return false
	}
	// Идем по байтам. len(s) возвращает количество байт, а не рун.
	for i := 0; i < len(s); i++ {
		// s[i] возвращает байт (byte/uint8)
		if cgpSafeTable[s[i]] == 0 {
			return false
		}
	}
	return true
}

// formatCGPSize преобразует байты в формат CGP (1K, 5M, 10G) или "unlimited".
func formatCGPSize(v int64) string {
	if v == -1 {
		return "unlimited"
	}
	if v == 0 {
		return "0"
	}
	units := []struct {
		suffix string
		value  int64
	}{
		{"T", 1024 * 1024 * 1024 * 1024},
		{"G", 1024 * 1024 * 1024},
		{"M", 1024 * 1024},
		{"K", 1024},
	}
	for _, u := range units {
		if v%u.value == 0 {
			return fmt.Sprintf("%d%s", v/u.value, u.suffix)
		}
	}
	return strconv.FormatInt(v, 10)
}

func MarshalTo(w *bufio.Writer, v any) error {
	return marshal(w, v, false)
}

func MarshalToString(v any) (string, error) {
	var b bytes.Buffer
	err := marshal(&b, v, false)
	return b.String(), err
}

func marshal(w encodeWriter, v any, asKey bool) error {
	switch val := v.(type) {
	case Atom:
		_, err := w.WriteString(string(val))
		return err

	case DataBlock:
		w.WriteByte('[')
		w.Write(val)
		return w.WriteByte(']')

	case []byte:
		w.WriteString("#[")
		enc := hex.NewEncoder(w)
		enc.Write(val)
		return w.WriteByte(']')

	case string:
		if isCGPSafeString(val) {
			_, err := w.WriteString(val)
			return err
		}
		_, err := w.WriteString(strconv.Quote(val))
		return err

	case int, int32, int64:
		var n int64
		switch i := val.(type) {
		case int:
			n = int64(i)
		case int32:
			n = int64(i)
		case int64:
			n = i
		}
		_, err := w.WriteString(formatCGPSize(n))
		return err

	case bool:
		var err error
		if val {
			_, err = w.WriteString("YES")
		} else {
			_, err = w.WriteString("NO")
		}
		return err

	case time.Time:
		if val.IsZero() {
			_, err := w.WriteString(`""`)
			return err
		}
		w.WriteByte('"')
		w.WriteString(val.Format("02-Jan-2006 15:04:05 -0700"))
		return w.WriteByte('"')

	case map[string]any:
		w.WriteByte('{')
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			marshal(w, k, true)
			w.WriteString("=")
			marshal(w, val[k], false)
			w.WriteByte(';')
		}
		return w.WriteByte('}')

	case []any:
		w.WriteByte('(')
		for i, item := range val {
			marshal(w, item, false)
			if i < len(val)-1 {
				w.WriteByte(',')
			}
		}
		return w.WriteByte(')')

	case map[string]string:
		w.WriteByte('{')
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			marshal(w, k, true)
			w.WriteString("=")
			marshal(w, val[k], false)
			w.WriteByte(';')
		}
		return w.WriteByte('}')

	case []string:
		w.WriteByte('(')
		for i, s := range val {
			marshal(w, s, false)
			if i < len(val)-1 {
				w.WriteByte(',')
			}
		}
		return w.WriteByte(')')

	default:
		if val == nil {
			_, err := w.WriteString("#NULL#")
			return err
		}
		_, err := w.WriteString(fmt.Sprint(val))
		return err
	}
}
