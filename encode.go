package cgpcli

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

type encodeWriter interface {
	io.Writer
	WriteByte(byte) error
	WriteString(string) (int, error)
}

// cgpMarshaler — контракт для типов с собственной логикой упаковки
type cgpMarshaler interface {
	MarshalCGP() string
}

type fieldMarshaler func(any) string

var customMarshalers = map[string]fieldMarshaler{
	"DNSServers": func(v any) string {
		ips, ok := v.([]net.IP)
		if !ok || len(ips) == 0 {
			return `""`
		}
		var sb strings.Builder
		for i, ip := range ips {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteByte('[')
			sb.WriteString(ip.String())
			sb.WriteByte(']')
		}
		return strconv.Quote(sb.String())
	},
	"LANAddress":  marshalSingleIP,
	"WANAddress":  marshalSingleIP,
	"WAN6Address": marshalSingleIP,
}

var cgpSafeTable = [256]byte{
	'0': 1, '1': 1, '2': 1, '3': 1, '4': 1, '5': 1, '6': 1, '7': 1, '8': 1, '9': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1,
	'K': 1, 'L': 1, 'M': 1, 'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1,
	'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1,
	'k': 1, 'l': 1, 'm': 1, 'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1,
	'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
	'_': 1, '@': 1, '.': 1, '-': 1,
}

func formatCGPDuration(d time.Duration) string {
	switch d {
	case -1:
		return "never"
	case 0:
		return "immediately"
	}

	totalSeconds := int64(d.Seconds())

	// Дни (24h)
	if totalSeconds%(24*3600) == 0 {
		return strconv.FormatInt(totalSeconds/(24*3600), 10) + "d"
	}
	// Часы
	if totalSeconds%3600 == 0 {
		return strconv.FormatInt(totalSeconds/3600, 10) + "h"
	}
	// Минуты
	if totalSeconds%60 == 0 {
		return strconv.FormatInt(totalSeconds/60, 10) + "m"
	}

	// Секунды (по умолчанию)
	return strconv.FormatInt(totalSeconds, 10) + "s"
}

// isCGPSafeString проверяет, можно ли передать строку без кавычек.
// В CGP это буквенно-цифровые строки без спецсимволов.
func isCGPSafeString(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
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
	return marshal(w, v, "")
}

func MarshalToString(v any) (string, error) {
	var b bytes.Buffer
	err := marshal(&b, v, "")
	return b.String(), err
}

func marshal(w encodeWriter, v any, keyName string) error {
	// Проверяем кастомные маршалеры по имени ключа
	if marshaler, ok := customMarshalers[keyName]; ok {
		_, err := w.WriteString(marshaler(v))
		return err
	}

	// Приоритет кастомного маршалинга (интерфейс cgpMarshaler)
	if m, ok := v.(cgpMarshaler); ok {
		_, err := w.WriteString(m.MarshalCGP())
		return err
	}

	switch val := v.(type) {
	case Atom:
		_, err := w.WriteString(string(val))
		return err

	case DataBlock:
		w.WriteByte('[')
		w.Write(val)
		return w.WriteByte(']')

	case []byte:
		w.WriteByte('[')
		enc := base64.NewEncoder(base64.StdEncoding, w)
		enc.Write(val)
		enc.Close()
		return w.WriteByte(']')

	// --- Строки ---
	case string:
		if val == "" {
			_, err := w.WriteString(`""`)
			return err
		}
		// Исключение для Created (legacy)
		if _, ok := dateFields[keyName]; ok {
			_, err := w.WriteString(strconv.Quote(val))
			return err
		}
		if isCGPSafeString(val) {
			_, err := w.WriteString(val)
			return err
		}
		_, err := w.WriteString(strconv.Quote(val))
		return err

	// --- Числа ---
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
		if _, ok := quotaFields[keyName]; ok {
			_, err := w.WriteString(formatCGPSize(n))
			return err
		}
		_, err := w.WriteString(strconv.FormatInt(n, 10))
		return err

	// --- Время ---
	case time.Time:
		if val.IsZero() {
			_, err := w.WriteString("#NULL#")
			return err
		}
		if _, ok := dateFields[keyName]; ok {
			fmt.Fprintf(w, "\"%s\"", val.Format(cgpInternalTimeLayout))
			return nil
		}
		if val.Equal(TimePast) {
			_, err := w.WriteString("#TPAST")
			return err
		}
		if val.Equal(TimeFuture) {
			_, err := w.WriteString("#TFUTURE")
			return err
		}
		w.WriteString("#T")
		w.WriteString(val.Format(cgpInternalTimeLayout))
		return nil

	case time.Duration:
		_, err := w.WriteString(formatCGPDuration(val))
		return err

	// --- Словари (Maps) ---
	case map[string]any:
		return marshalMapAny(w, val)

	case map[string]string:
		return marshalMapString(w, val)

	// --- Списки (Slices) ---
	case []any:
		w.WriteByte('(')
		for i, item := range val {
			if i > 0 {
				w.WriteByte(',')
			}
			marshal(w, item, "")
		}
		return w.WriteByte(')')

	case []string:
		w.WriteByte('(')
		for i, s := range val {
			if i > 0 {
				w.WriteByte(',')
			}
			marshal(w, s, "")
		}
		return w.WriteByte(')')

	// --- Остальное ---
	case bool:
		if val {
			_, err := w.WriteString("YES")
			return err
		}
		_, err := w.WriteString("NO")
		return err

	case nil:
		_, err := w.WriteString("#NULL#")
		return err

	default:
		_, err := w.WriteString(fmt.Sprint(val))
		return err
	}
}

func marshalMapAny(w encodeWriter, m map[string]any) error {
	w.WriteByte('{')
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		marshal(w, k, "")
		w.WriteByte('=')
		marshal(w, m[k], k)
		w.WriteByte(';')
	}
	return w.WriteByte('}')
}

func marshalMapString(w encodeWriter, m map[string]string) error {
	w.WriteByte('{')
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		marshal(w, k, "")
		w.WriteByte('=')
		marshal(w, m[k], k)
		w.WriteByte(';')
	}
	return w.WriteByte('}')
}

func marshalSingleIP(v any) string {
	if ip, ok := v.(net.IP); ok {
		return strconv.Quote("[" + ip.String() + "]")
	}
	return `""`
}
