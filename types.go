// # Specialized Protocol Types
//
// The Types section defines the fundamental building blocks of the CommuniGate Pro
// CLI protocol syntax. It maps CGP-specific data representations to native
// Go structures, ensuring that complex data types are marshaled and
// unmarshaled with bit-perfect accuracy.
//
// Key components include:
//   - Atoms: the [Atom] type represents unquoted strings used for commands,
//     keywords, and object identifiers where standard quoting would cause
//     syntax errors.
//   - Network Addresses: [CGPIP] handles the server's unique #I[address]:port
//     format, supporting both IPv4 and IPv6 with optional port numbers.
//   - Temporal Constants: predefined [TimePast] (#TPAST) and [TimeFuture]
//     (#TFUTURE) constants for handling expiration dates and schedule limits.
//   - Binary Data: [DataBlock] provides a container for raw bytes (like SSL
//     certificates) that require specific encapsulation during transmission.
//   - Session Offsets: enumerated constants for navigating data streams
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

// Atom represents a string that is always transmitted without quotes.
// It is used for commands, keywords, and object identifiers where
// standard CLI quoting would result in a syntax error.
type Atom string

// DataBlock represents binary data, such as SSL certificates or keys.
// When marshaled, it is enclosed in square brackets [].
//
// Unlike the standard []byte type, the data inside a DataBlock MUST already
// be Base64-encoded. This type serves as a transparent container and does
// not perform any automatic encoding.
type DataBlock []byte

// FileInfo represents metadata for a file in the CommuniGate Pro storage
// (Skins, PBX, etc.).
type FileInfo struct {
	Name     string    // File name
	Size     int64     // File size in bytes (STFileSize)
	Created  time.Time // File creation timestamp (STCreated)
	Modified time.Time // Last modification timestamp (STModified)
}

// CGPIP represents the specific CommuniGate Pro IP Address data type.
// It follows the protocol format: #I[address] or #I[address]:port.
type CGPIP struct {
	IP   net.IP // IP is the network address (IPv4 or IPv6).
	Port int    // Port is the optional port number (0 means no port).
}

// IsZero returns true if the IP address is not set.
func (c CGPIP) IsZero() bool {
	return len(c.IP) == 0
}

// String returns the standard human-readable representation of the IP.
// For example: "1.2.3.4:25" or "2001::1".
func (c CGPIP) String() string {
	if c.IsZero() {
		return "<nil>"
	}
	if c.Port > 0 {
		return net.JoinHostPort(c.IP.String(), strconv.Itoa(c.Port))
	}
	return c.IP.String()
}

// writeCGP формирует данные для отправки на сервер CGP напрямую в буфер.
func (c CGPIP) writeCGP(w encodeWriter) error {
	if c.IsZero() {
		_, err := w.WriteString("#NULL#")
		return err
	}

	w.WriteString("#I[")
	w.WriteString(c.IP.String())
	w.WriteByte(']')

	if c.Port > 0 {
		w.WriteByte(':')
		w.WriteString(strconv.Itoa(c.Port))
	}

	return nil
}
