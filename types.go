package cgpcli

import (
	"net"
	"strconv"
	"strings"
	"time"
)

const cgpInternalTimeLayout = "02-01-2006_15:04:05"

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

// CGPIP представляет специфический для CGP тип данных IP Address (#I).
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

// MarshalCGP формирует строку для отправки на сервер CGP.
func (c CGPIP) MarshalCGP() string {
	if c.IsZero() {
		return "#NULL#"
	}
	var sb strings.Builder
	sb.WriteString("#I[")
	sb.WriteString(c.IP.String())
	sb.WriteByte(']')
	if c.Port > 0 {
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(c.Port))
	}
	return sb.String()
}
