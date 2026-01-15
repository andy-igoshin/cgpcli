package cgpcli

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type IPKey [16]byte

func ToKey(ip net.IP) IPKey {
	var key IPKey
	if ip != nil {
		copy(key[:], ip.To16())
	}
	return key
}

type TempIP struct {
	IP        net.IP
	TimeLeft  time.Duration
	Permanent bool
}

type TempIPList struct {
	items map[IPKey]TempIP
}

// Add обновляет локальный черновик.
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

// Delete удаляет из локального черновика.
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

// Exists — короткая проверка наличия.
func (l *TempIPList) Exists(input any) (bool, error) {
	_, ok, err := l.Get(input)
	return ok, err
}

// Get возвращает данные, признак наличия и ошибку валидации входа.
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

// IPs экспортирует слайс для Set-команд.
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

// Len возвращает количество записей в списке.
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

	b := []byte(raw)
	count := bytes.Count(b, []byte{','}) + 1
	res := make(map[IPKey]TempIP, count)

	for i := 0; i < len(b); {
		start := bytes.IndexByte(b[i:], '[')
		if start == -1 {
			break
		}
		start += i

		end := bytes.IndexByte(b[start:], ']')
		if end == -1 {
			break
		}
		end += start

		ipStr := string(b[start+1 : end])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			i = end + 1
			continue
		}

		item := TempIP{IP: ip}
		curr := end + 1

		if curr < len(b) && b[curr] == '-' {
			curr++
			if curr < len(b) && b[curr] == '*' {
				item.Permanent = true
				curr++
			} else {
				val := 0
				foundDigit := false
				for curr < len(b) && b[curr] >= '0' && b[curr] <= '9' {
					val = val*10 + int(b[curr]-'0')
					curr++
					foundDigit = true
				}
				if foundDigit {
					item.TimeLeft = time.Duration(val) * time.Second
				}
			}
		}

		res[ToKey(ip)] = item

		nextComma := bytes.IndexByte(b[curr:], ',')
		if nextComma == -1 {
			break
		}
		i = curr + nextComma + 1
	}
	return res
}
