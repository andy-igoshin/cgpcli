// # Specialized Protocol Types
//
// The Types section defines the fundamental building blocks of the CommuniGate Pro
// CLI protocol syntax. It maps CGP-specific data representations to native
// Go structures, ensuring that complex data types are marshaled and
// unmarshaled with bit-perfect accuracy.
//
// Key components include:
//   - Atoms: The [Atom] type represents unquoted strings used for commands,
//     keywords, and object identifiers where standard quoting would cause
//     syntax errors.
//   - Network Addresses: [CGPIP] handles the server's unique #I[address]:port
//     format, supporting both IPv4 and IPv6 with optional port numbers.
//   - Temporal Constants: Predefined [TimePast] (#TPAST) and [TimeFuture]
//     (#TFUTURE) constants for handling expiration dates and schedule limits.
//   - Binary Data: [DataBlock] provides a container for raw bytes (like SSL
//     certificates) that require specific encapsulation during transmission.
//   - Session Offsets: Enumerated constants for navigating data streams
//     ([OffsetBeg], [OffsetEnd], [OffsetNew]).
//
// These types are used throughout the library to maintain strict compliance
// with the CommuniGate Pro data model while providing a developer-friendly API.
package cgpcli

import (
	"net"
	"strconv"
	"time"
)

// Специальные метки времени CGP.
// Используются для представления #TPAST и #TFUTURE.
var (
	TimePast   = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	TimeFuture = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
)

// Atom представляет собой строку, которая всегда передается без кавычек.
type Atom string

// DataBlock представляет собой бинарные данные (сертификаты и т.д.).
// При маршалинге оборачивается в [] и считается, что данные ВНУТРИ
// уже закодированы в Base64.
type DataBlock []byte

// CGPIP представляет специфический для CGP тип данных IP Address.
// Формат в протоколе: #I[address] или #I[address]:port
type CGPIP struct {
	IP   net.IP
	Port int // 0 означает отсутствие порта
}

// IsZero возвращает true, если IP не задан.
func (c CGPIP) IsZero() bool {
	return len(c.IP) == 0
}

// String возвращает привычное представление: 1.2.3.4:25 или 2001::1
func (c CGPIP) String() string {
	if c.IsZero() {
		return "<nil>"
	}
	if c.Port > 0 {
		return net.JoinHostPort(c.IP.String(), strconv.Itoa(c.Port))
	}
	return c.IP.String()
}

// WriteCGP формирует данные для отправки на сервер CGP напрямую в буфер.
func (c CGPIP) WriteCGP(w encodeWriter) (err error) {
	if c.IsZero() {
		_, err = w.WriteString("#NULL#")
		return
	}

	if _, err = w.WriteString("#I["); err != nil {
		return
	}
	if _, err = w.WriteString(c.IP.String()); err != nil {
		return
	}
	if err = w.WriteByte(']'); err != nil {
		return
	}

	if c.Port > 0 {
		if err = w.WriteByte(':'); err != nil {
			return
		}
		if _, err = w.WriteString(strconv.Itoa(c.Port)); err != nil {
			return
		}
	}
	return
}

type SessionOffset int

const (
	OffsetBeg SessionOffset = -1
	OffsetEnd SessionOffset = -2
	OffsetNew SessionOffset = -3
)
