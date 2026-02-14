// # Data Serialization (Encoder)
//
// The Encoding section provides methods for converting Go data structures into
// the CommuniGate Pro CLI protocol format. It ensures that all commands,
// keywords, and settings are formatted correctly according to the server's
// unique grammar rules.
//
// Key capabilities include:
//   - Smart String Quoting: Automatically determines if a string is "safe"
//     to send as an unquoted token or if it requires escaping and double quotes.
//   - Unit Transformation: Converts numeric values into CommuniGate Pro CLI size units
//     (K, M, G, T) and durations into time units (s, m, h, d) automatically
//     based on the context (key name).
//   - Native Binary Support: Efficiently encodes byte slices into
//     Base64-wrapped [data blocks] required by the server for file transfers
//     and certificates.
//   - Temporal Handling: Maps Go's [time.Time] to CGP's #T, #TPAST, and
//     #TFUTURE syntax, specifically managing the UTC offset to match the
//     server's internal time representation.
//   - Recursive Structures: Deep serialization of maps ({...}) and
//     slices ((...)) with support for nested custom types.
//
// This encoder is the core communication bridge, responsible for
// constructing every command and data structure sent to the server.
package cgpcli

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type (
	cgpDateLegacy  time.Time
	cgpDuration    time.Duration
	cgpIPSB        net.IP
	cgpIPSBSA      []net.IP
	cgpSize        int64
	cgpForceString string
)

type encodeWriter interface {
	io.Writer
	WriteByte(byte) error
	WriteString(string) (int, error)
}

type cgpWriterTo interface {
	writeCGP(w encodeWriter) error
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

// MarshalTo serializes v directly into a buffered writer. This is the
// preferred function for high-performance operations or large data sets
// to avoid intermediate string allocations.
//
// Parameters:
//   - w: A pointer to a *bufio.Writer where the serialized data will be written.
//   - v: The data structure to be serialized.
//
// Returns:
//   - An error if writing to the buffer fails or if the data is invalid.
func MarshalTo(w *bufio.Writer, v any) error {
	return marshal(w, v, "")
}

// MarshalToString returns the CGP-formatted string representation of v.
//
// Parameters:
//   - v: The data to be serialized. Can be any basic Go type (string, int,
//     bool, time.Time), a map[string]any, or a slice.
//
// Returns:
//   - A string ready to be sent to the CommuniGate Pro CLI.
//   - An error if the data structure contains types that cannot be serialized.
func MarshalToString(v any) (string, error) {
	var b bytes.Buffer
	err := marshal(&b, v, "")
	return b.String(), err
}

func marshal(w encodeWriter, v any, keyName string) (err error) {
	// 1. Прямая запись (самый быстрый путь для кастомных типов)
	if m, ok := v.(cgpWriterTo); ok {
		return m.writeCGP(w)
	}

	// 2. Трансформация исключений (меняет тип на cgp*)
	if keyName != "" {
		v = transformWrite(keyName, v)
	}

	switch val := v.(type) {
	case cgpDateLegacy:
		tm := time.Time(val)
		if tm.IsZero() {
			_, err = w.WriteString("#NULL#")
		} else {
			w.WriteByte('"')
			w.WriteString(tm.UTC().Format(cgpInternalTimeLayout))
			err = w.WriteByte('"')
		}
		return err

	case cgpIPSB:
		w.WriteByte('"')
		w.WriteByte('[')
		w.WriteString(net.IP(val).String())
		w.WriteByte(']')
		err = w.WriteByte('"')
		return err

	case cgpIPSBSA:
		w.WriteByte('"')
		ips := []net.IP(val)
		for i, ip := range ips {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteByte('[')
			w.WriteString(ip.String())
			w.WriteByte(']')
		}
		err = w.WriteByte('"')
		return err

	case cgpSize:
		_, err = w.WriteString(formatCGPSize(int64(val)))
		return err

	case cgpDuration:
		_, err = w.WriteString(formatCGPDuration(time.Duration(val)))
		return err

	case Atom:
		_, err = w.WriteString(string(val))
		return err

	case DataBlock:
		w.WriteByte('[')
		w.Write(val)
		err = w.WriteByte(']')
		return err

	case []byte:
		w.WriteByte('[')
		enc := base64.NewEncoder(base64.StdEncoding, w)
		enc.Write(val)
		enc.Close()
		err = w.WriteByte(']')
		return err

	// --- Строки ---
	case string:
		if isCGPSafeString(val) {
			_, err = w.WriteString(val)
			return err
		}
		if err = w.WriteByte('"'); err != nil {
			return err
		}
		if _, err = w.WriteString(replaceSpecChars(val)); err != nil {
			return err
		}
		return w.WriteByte('"')

	// --- Числа ---
	case int, int32, int64:
		return writeInt64(w, toInt64(val))

	// --- Время ---
	case time.Time:
		if val.IsZero() {
			_, err = w.WriteString("#NULL#")
			return err
		}
		if val.Equal(TimePast) {
			_, err = w.WriteString("#TPAST")
			return err
		}
		if val.Equal(TimeFuture) {
			_, err = w.WriteString("#TFUTURE")
			return err
		}
		w.WriteString("#T")
		_, err = w.WriteString(val.UTC().Format(cgpInternalTimeLayout))
		return err

	case time.Duration:
		_, err = w.WriteString(formatCGPDuration(val))
		return err

	// --- Словари (Maps) ---
	case map[string]any:
		return marshalMap(w, val)

	case map[string]string:
		return marshalMap(w, val)

	// --- Списки (Slices) ---
	case []any:
		w.WriteByte('(')
		for i, item := range val {
			if i > 0 {
				w.WriteByte(',')
			}
			marshal(w, item, "")
		}
		err = w.WriteByte(')')
		return err

	case []string:
		w.WriteByte('(')
		for i, s := range val {
			if i > 0 {
				w.WriteByte(',')
			}
			marshal(w, s, "")
		}
		err = w.WriteByte(')')
		return err

	// --- Остальное ---
	case bool:
		if val {
			_, err = w.WriteString("YES")
			return err
		}
		_, err = w.WriteString("NO")
		return err

	case nil:
		_, err = w.WriteString("#NULL#")
		return err

	default:
		strVal := fmt.Sprint(val)
		if isCGPSafeString(strVal) {
			_, err = w.WriteString(strVal)
			return err
		}
		if err = w.WriteByte('"'); err != nil {
			return err
		}
		if _, err = w.WriteString(replaceSpecChars(strVal)); err != nil {
			return err
		}
		return w.WriteByte('"')
	}
}

func marshalMap[V any](w encodeWriter, m map[string]V) error {
	w.WriteByte('{')

	for k, v := range m {
		if isCGPSafeString(k) {
			w.WriteString(k)
		} else {
			w.WriteByte('"')
			w.WriteString(replaceSpecChars(k))
			w.WriteByte('"')
		}

		w.WriteByte('=')
		marshal(w, v, k)
		w.WriteByte(';')
	}

	return w.WriteByte('}')
}

func writeInt64(w encodeWriter, i int64) error {
	if i == 0 {
		return w.WriteByte('0')
	}

	if i < 0 {
		if err := w.WriteByte('-'); err != nil {
			return err
		}
	} else {
		i = -i
	}

	var buf [20]byte
	pos := len(buf)

	for i != 0 {
		pos--
		buf[pos] = byte('0' - (i % 10))
		i /= 10
	}

	for ; pos < len(buf); pos++ {
		if err := w.WriteByte(buf[pos]); err != nil {
			return err
		}
	}
	return nil
}
