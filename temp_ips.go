// # Temporary IP List Management
//
// The TempIP section provides advanced structures for handling transient IP
// address lists in CommuniGate Pro. These lists are primarily used for
// dynamic security measures, such as temporary blacklisting of failed login
// attempts or short-term "unblockable" exceptions.
//
// Key capabilities include:
//   - Structured Representation: The [TempIP] struct associates an IP address
//     with its remaining time-to-live (TTL) or a permanent status flag.
//   - Efficient Lookups: [TempIPList] uses a fixed-size [IPKey] ([16]byte)
//     based map for O(1) IP address lookups, handling both IPv4 and IPv6 transparently.
//   - Robust Parsing: Optimized parsing logic in [parseTempIPsToMap] that
//     extracts IPs and their durations from raw protocol strings with
//     minimal allocations.
//   - Local Drafting: Methods like [TempIPList.Add] and [TempIPList.Delete]
//     allow for building or modifying lists locally before committing
//     them to the server.
//
// These structures are vital for building security orchestrators,
// fail2ban-style integrators, and automated traffic controllers.
package cgpcli

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// IPKey is a fixed-size byte array representing an IP address.
// It is used as a map key to provide efficient O(1) lookups and
// minimal memory footprint compared to string keys.
type IPKey [16]byte

// ToKey converts a net.IP address into an IPKey.
// It transparently handles both IPv4 and IPv6 by converting
// addresses to their 16-byte representation.
//
// Parameters:
//   - ip: the net.IP address to convert.
//
// Returns:
//   - IPKey: a 16-byte array suitable for use as a map key.
func ToKey(ip net.IP) IPKey {
	var key IPKey
	if ip != nil {
		copy(key[:], ip.To16())
	}
	return key
}

// TempIP represents an entry in a temporary IP list.
// It combines the network address with its expiration logic.
type TempIP struct {
	IP        net.IP        // IP is the network address (IPv4 or IPv6).
	TimeLeft  time.Duration // TimeLeft specifies the remaining duration before expiration.
	Permanent bool          // Permanent indicates if the IP should remain in the list indefinitely.
}

// TempIPList manages a collection of temporary IP addresses.
// It provides methods for local manipulation of the list before
// it is serialized and sent to the CommuniGate Pro server.
type TempIPList struct {
	items map[IPKey]TempIP
}

// Add updates the local draft by adding or updating an IP address.
//
// Parameters:
//   - input: the IP address to add (supported types: string or net.IP).
//   - duration: the time-to-live for the entry.
//   - permanent: if true, the entry is marked as permanent (TimeLeft is ignored).
//
// Returns:
//   - error: an error if the input type is unsupported or the IP is invalid.
func (l *TempIPList) Add(input any, duration time.Duration, permanent bool) error {
	ip, key, err := l.parseInput(input)
	if err != nil {
		return err
	}
	if l.items == nil {
		l.items = make(map[IPKey]TempIP)
	}
	l.items[key] = TempIP{
		IP:        ip,
		TimeLeft:  duration,
		Permanent: permanent,
	}
	return nil
}

// Delete removes an IP address from the local draft.
//
// Parameters:
//   - input: the IP address to remove (supported types: string or net.IP).
//
// Returns:
//   - error: an error if the input type is unsupported or the IP is invalid.
func (l *TempIPList) Delete(input any) error {
	_, key, err := l.parseInput(input)
	if err != nil {
		return err
	}
	if l.items == nil {
		return nil
	}
	delete(l.items, key)
	return nil
}

// Exists provides a quick check for the presence of an IP in the list.
//
// Parameters:
//   - input: the IP address to check (supported types: string or net.IP).
//
// Returns:
//   - bool: true if the IP exists in the current list.
//   - error: an error if the input IP is invalid.
func (l *TempIPList) Exists(input any) (bool, error) {
	_, ok, err := l.Get(input)
	return ok, err
}

// Get retrieves the details of a specific IP entry.
//
// Parameters:
//   - input: the IP address to look up (supported types: string or net.IP).
//
// Returns:
//   - [TempIP]: the structure containing IP, TimeLeft, and Permanent status.
//   - bool: true if the entry was found.
//   - error: an error if the input IP is invalid.
func (l *TempIPList) Get(input any) (TempIP, bool, error) {
	_, key, err := l.parseInput(input)
	if err != nil {
		return TempIP{}, false, err
	}
	if l.items == nil {
		return TempIP{}, false, nil
	}
	item, ok := l.items[key]
	return item, ok, nil
}

// IPs exports all IP addresses currently present in the list as a slice.
//
// Returns:
//   - []net.IP: a slice of IP addresses. Returns nil if the list is empty.
func (l *TempIPList) IPs() []net.IP {
	if l == nil || l.items == nil {
		return nil
	}
	res := make([]net.IP, 0, len(l.items))
	for _, item := range l.items {
		res = append(res, item.IP)
	}
	return res
}

// Len returns the number of entries in the current list.
//
// Returns:
//   - int: the total count of IP entries.
func (l *TempIPList) Len() int {
	if l == nil || l.items == nil {
		return 0
	}
	return len(l.items)
}

func (l *TempIPList) parseInput(input any) (net.IP, IPKey, error) {
	var ip net.IP
	switch v := input.(type) {
	case net.IP:
		if v == nil {
			return nil, IPKey{}, fmt.Errorf("IP address is nil")
		}
		ip = v
	case string:
		ip = net.ParseIP(v)
		if ip == nil {
			return nil, IPKey{}, fmt.Errorf("invalid IP address: %s", v)
		}
	default:
		return nil, IPKey{}, fmt.Errorf("unsupported type %T, need string or net.IP", v)
	}
	return ip, ToKey(ip), nil
}

func parseTempIPsToMap(raw string) map[IPKey]TempIP {
	if raw == "" {
		return nil
	}

	// Мы не знаем точного количества, но можем оценить по количеству запятых
	commas := 0
	for i := 0; i < len(raw); i++ {
		if raw[i] == ',' {
			commas++
		}
	}
	res := make(map[IPKey]TempIP, commas+1)

	b := []byte(raw)
	n := len(b)
	for i := 0; i < n; {
		// Ищем начало IP '['
		for i < n && b[i] != '[' {
			i++
		}
		if i >= n {
			break
		}
		i++ // Пропускаем '['

		start := i
		for i < n && b[i] != ']' {
			i++
		}
		if i >= n {
			break
		}

		// Извлекаем IP.
		ipStr := string(b[start:i])
		ip := net.ParseIP(ipStr)
		i++ // Пропускаем ']'

		item := TempIP{IP: ip}

		// Проверяем наличие суффикса '-'
		if i < n && b[i] == '-' {
			i++
			if i < n && b[i] == '*' {
				item.Permanent = true
				i++
			} else {
				val := 0
				foundDigit := false
				for i < n && b[i] >= '0' && b[i] <= '9' {
					val = val*10 + int(b[i]-'0')
					i++
					foundDigit = true
				}
				if foundDigit {
					item.TimeLeft = time.Duration(val) * time.Second
				}
			}
		}

		if ip != nil {
			res[ToKey(ip)] = item
		}

		// Ищем следующую запятую или конец
		for i < n && b[i] != ',' {
			i++
		}
		i++ // Пропускаем запятую для следующей итерации
	}
	return res
}

func (l *TempIPList) writeCGP(w encodeWriter) (err error) {
	if err = w.WriteByte('"'); err != nil {
		return
	}

	first := true
	for _, item := range l.items {
		if !first {
			if err = w.WriteByte(','); err != nil {
				return
			}
		}
		first = false

		if err = w.WriteByte('['); err != nil {
			return
		}
		if _, err = w.WriteString(item.IP.String()); err != nil {
			return
		}
		if err = w.WriteByte(']'); err != nil {
			return
		}

		if item.Permanent {
			if _, err = w.WriteString("-*"); err != nil {
				return
			}
		} else if item.TimeLeft > 0 {
			if err = w.WriteByte('-'); err != nil {
				return
			}
			seconds := int64(item.TimeLeft.Seconds())
			if _, err = w.WriteString(strconv.FormatInt(seconds, 10)); err != nil {
				return
			}
		}
	}

	return w.WriteByte('"')
}
