package cgpcli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type decodeState struct {
	data     []byte
	len      int
	span     int
	scratch  bytes.Buffer
	keyCache map[string]string
}

var (
	binYesUpper = []byte("YES")
	binYesLower = []byte("yes")
	binNoUpper  = []byte("NO")
	binNoLower  = []byte("no")
)

const (
	isCgpSpace      uint8 = 1 << 0
	isCgpDelimiter  uint8 = 1 << 1
	isCgpOpen       uint8 = 1 << 2 // {
	isCgpClose      uint8 = 1 << 3 // }
	isCgpSemi       uint8 = 1 << 4 // ;
	isCgpParenOpen  uint8 = 1 << 5 // (
	isCgpParenClose uint8 = 1 << 6 // )
	isCgpComma      uint8 = 1 << 7 // ,
)

var cgpCharTable [256]uint8

func init() {
	for _, c := range " \t\r\n" {
		cgpCharTable[c] |= isCgpSpace
	}
	for _, c := range " \t\r\n=;,)}" {
		cgpCharTable[c] |= isCgpDelimiter
	}
	// Специфические флаги для структуры словаря
	cgpCharTable['{'] |= isCgpOpen
	cgpCharTable['}'] |= isCgpClose
	cgpCharTable[';'] |= isCgpSemi
	cgpCharTable['('] |= isCgpParenOpen
	cgpCharTable[')'] |= isCgpParenClose
	cgpCharTable[','] |= isCgpComma
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.len = len(data)
	d.span = 0
	d.scratch.Grow(4096)
	if d.keyCache == nil {
		d.keyCache = make(map[string]string, 256)
	}
	return d
}

func isSpace(c byte) bool {
	return cgpCharTable[c]&isCgpSpace != 0
}

func (d *decodeState) skipSpaces() {
	for d.span < d.len && (cgpCharTable[d.data[d.span]]&isCgpSpace != 0) {
		d.span++
	}
}

func Unmarshal(data []byte, v any) error {
	var d decodeState
	d.init(data)

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("cgpcli: Unmarshal expects non-nil pointer")
	}

	result := d.valueInterface(false)
	target := reflect.ValueOf(v).Elem()

	if result == nil {
		if target.CanSet() {
			target.Set(reflect.Zero(target.Type()))
		}
		return nil
	}

	val := reflect.ValueOf(result)
	if val.Type().AssignableTo(target.Type()) {
		target.Set(val)
	} else if val.Type().ConvertibleTo(target.Type()) {
		target.Set(val.Convert(target.Type()))
	} else {
		return fmt.Errorf("cgpcli: cannot assign %T to %v", result, target.Type())
	}
	return nil
}

func (d *decodeState) valueInterface(asString bool) any {
	d.skipSpaces()
	if d.span >= d.len {
		return nil
	}
	switch d.data[d.span] {
	case '{':
		d.span++
		return d.dictionaryInterface()
	case '(':
		d.span++
		return d.arrayInterface()
	case '<':
		d.span++
		return d.xmlInterface()
	default:
		return d.wordInterface(asString)
	}
}

func (d *decodeState) dictionaryInterface() any {
	count := 0
	depth := 1
	ts := d.span
	data := d.data

	for ts < len(data) {
		f := cgpCharTable[data[ts]]
		if f != 0 {
			if f&isCgpOpen != 0 {
				depth++
			} else if f&isCgpClose != 0 {
				depth--
				if depth == 0 {
					break
				}
			} else if f&isCgpSemi != 0 && depth == 1 {
				count++
			}
		}
		ts++
	}

	res := make(map[string]any, count)

	for d.span < d.len {
		d.skipSpaces()
		if d.span >= d.len || d.data[d.span] == '}' {
			d.span++
			break
		}

		var key string
		start := d.span
		if d.data[d.span] != '"' {
			for d.span < d.len && cgpCharTable[d.data[d.span]]&isCgpDelimiter == 0 {
				d.span++
			}

			keyBytes := d.data[start:d.span]

			var ok bool
			if key, ok = d.keyCache[string(keyBytes)]; !ok {
				key = string(keyBytes)
				d.keyCache[key] = key
			}
		} else {
			kRaw := d.wordInterface(true)
			key, _ = kRaw.(string)
			if cached, ok := d.keyCache[key]; ok {
				key = cached
			} else {
				d.keyCache[key] = key
			}
		}

		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == '=' {
			d.span++
			val := d.valueInterface(false)
			res[key] = transformRead(key, val)
		}

		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == ';' {
			d.span++
		}
	}
	return res
}

func (d *decodeState) arrayInterface() any {
	count := 1
	depth := 1
	ts := d.span
	data := d.data

	for ts < len(data) {
		f := cgpCharTable[data[ts]]
		if f != 0 {
			if f&isCgpParenOpen != 0 {
				depth++
			} else if f&isCgpParenClose != 0 {
				depth--
				if depth == 0 {
					break
				}
			} else if f&isCgpComma != 0 && depth == 1 {
				count++
			}
		}
		ts++
	}

	// Если мы нашли элементы (или хотя бы один), count должен быть корректным.
	// Для пустого массива () count останется 0. Для (1) - тоже 0, надо +1 если не пусто.
	// Но проще: если мы встретили хотя бы один символ до закрытия, это 1 элемент или больше.
	if ts > d.span { // Если между ( и ) что-то было
		count++
	}

	res := make([]any, 0, count)

	for d.span < d.len {
		d.skipSpaces()
		if d.span >= d.len || d.data[d.span] == ')' {
			d.span++
			break
		}

		res = append(res, d.valueInterface(false))

		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == ',' {
			d.span++
		}
	}
	return res
}

func (d *decodeState) wordInterface(asString bool) any {
	d.skipSpaces()
	if d.span >= d.len {
		return ""
	}

	curr := d.data[d.span]

	// 1. Data Blocks [base64]
	if curr == '[' {
		d.span++
		return d.readCgpBlock('[', ']')
	}

	// 2. IP Addresses #I[...]
	if curr == '#' && d.span+2 < d.len && d.data[d.span+1] == 'I' && d.data[d.span+2] == '[' {
		d.span += 3
		ipData := d.readCgpBlock('[', ']')
		res := CGPIP{IP: net.ParseIP(string(ipData))}
		if d.span < d.len && d.data[d.span] == ':' {
			d.span++
			start := d.span
			for d.span < d.len && d.data[d.span] >= '0' && d.data[d.span] <= '9' {
				d.span++
			}
			res.Port, _ = strconv.Atoi(string(d.data[start:d.span]))
		}
		return res
	}

	// 3. Fast Path (Unquoted tokens: keys, numbers, constants)
	if curr != '"' && curr != '#' {
		start := d.span
		for d.span < d.len {
			if cgpCharTable[d.data[d.span]]&isCgpDelimiter != 0 {
				break
			}
			d.span++
		}

		subSlice := d.data[start:d.span]

		// Если нам не принудительно нужна строка (как для ключей мап)
		if !asString {
			// Проверка на YES/NO через bytes.Equal
			if val, ok := tryCgpBool(subSlice); ok {
				return val
			}
			// Проверка на числа через unsafe string
			if val, ok := tryParseInt(subSlice); ok {
				return val
			}
		}

		// Поиск в кеше строк (интернирование)
		if s, ok := d.keyCache[string(subSlice)]; ok {
			return s
		}

		// Единственная точка аллокации для новых уникальных строк/ключей
		raw := string(subSlice)
		d.keyCache[raw] = raw
		return raw
	}

	// 4. Slow Path (Quoted, Objects #(..), Tags)
	d.scratch.Reset()
	isQuoted := curr == '"'
	isTagged := curr == '#'
	isObject := false

	if isQuoted {
		for d.span < d.len && d.data[d.span] == '"' {
			d.span++
			for d.span < d.len {
				ch := d.data[d.span]
				if ch == '\\' {
					if d.span+1 < d.len {
						next := d.data[d.span+1]
						if next >= '0' && next <= '7' && d.span+3 < d.len {
							octStr := string(d.data[d.span+1 : d.span+4])
							if val, err := strconv.ParseInt(octStr, 8, 32); err == nil {
								d.scratch.WriteByte(byte(val))
								d.span += 4
								continue
							}
						}
						switch next {
						case 'e':
							d.scratch.WriteByte('\n')
						case 't':
							d.scratch.WriteByte('\t')
						case '\\':
							d.scratch.WriteByte('\\')
						case '"':
							d.scratch.WriteByte('"')
						default:
							d.scratch.WriteByte(next)
						}
						d.span += 2
						continue
					}
				}
				if ch == '"' {
					d.span++
					if d.span < d.len && d.data[d.span] == '"' {
						d.scratch.WriteByte('"')
						d.span++
						continue
					}
					d.skipSpaces()
					break
				}
				d.scratch.WriteByte(ch)
				d.span++
			}
		}
	} else if isTagged {
		if bytes.HasPrefix(d.data[d.span:], []byte("#(")) {
			isObject = true
		} else {
			d.scratch.WriteByte('#')
			d.span++
		}
		for d.span < d.len {
			ch := d.data[d.span]
			if isObject {
				d.scratch.WriteByte(ch)
				d.span++
				if ch == ')' {
					break
				}
				continue
			}
			if isSpace(ch) || ch == '=' || ch == ';' || ch == ',' || ch == ')' || ch == '}' {
				break
			}
			d.scratch.WriteByte(ch)
			d.span++
		}
	}

	raw := d.scratch.String()
	if asString || isQuoted {
		return raw
	}

	if isTagged && !isObject {
		if raw == "#NULL#" {
			return nil
		}
		if strings.HasPrefix(raw, "#T") {
			switch raw {
			case "#TPAST":
				return TimePast
			case "#TFUTURE":
				return TimeFuture
			default:
				if t, err := time.Parse(cgpInternalTimeLayout, raw[2:]); err == nil {
					return t
				}
			}
		}
		return parseCGPNumber(raw)
	}

	if val, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return val
	}
	return raw
}

func (d *decodeState) xmlInterface() any {
	startSpan := d.span - 1
	depth := 1
	inString := false
	var quoteChar byte

	for d.span < d.len {
		ch := d.data[d.span]
		if inString {
			if ch == quoteChar {
				inString = false
			}
			d.span++
			continue
		}
		if ch == '"' || ch == '\'' {
			inString = true
			quoteChar = ch
			d.span++
			continue
		}
		if ch == '<' && d.span+1 < d.len {
			if d.data[d.span+1] == '/' {
				depth--
			} else if d.data[d.span+1] != '!' && d.data[d.span+1] != '?' {
				depth++
			}
		} else if ch == '>' {
			if d.span > 0 && d.data[d.span-1] == '/' {
				depth--
			}
			if depth == 0 {
				d.span++
				break
			}
		}
		d.span++
	}
	return string(d.data[startSpan:d.span])
}

func parseCGPDuration(s string) time.Duration {
	if s == "" {
		return 0
	}
	switch strings.ToLower(s) {
	case "never":
		return -1
	case "immediately":
		return 0
	}

	lastIdx := len(s) - 1
	multiplier := time.Second
	numPart := s
	hasSuffix := true

	switch s[lastIdx] {
	case 's', 'S':
		multiplier = time.Second
	case 'm', 'M':
		multiplier = time.Minute
	case 'h', 'H':
		multiplier = time.Hour
	case 'd', 'D':
		multiplier = time.Hour * 24
	default:
		hasSuffix = false
	}

	if hasSuffix {
		numPart = s[:lastIdx]
	}

	val, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0
	}

	return time.Duration(val) * multiplier
}

func parseCGPNumber(raw string) any {
	if len(raw) < 2 {
		return raw
	}
	s := raw[1:] // Отрезаем '#'

	base := 10
	start := 0
	if s[0] == '-' {
		start = 1
	}

	if len(s[start:]) >= 2 && s[start] == '0' {
		switch s[start+1] {
		case 'x':
			base = 16
		case 'o':
			base = 8
		case 'b':
			base = 2
		}
		if base != 10 {
			// Убираем индикатор базы (0x, 0o, 0b) для strconv
			toParse := s[:start] + s[start+2:]
			if val, err := strconv.ParseInt(toParse, base, 64); err == nil {
				return val
			}
		}
	}

	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	return raw
}

func parseCGPSize(s string) any {
	if s == "unlimited" {
		return int64(-1)
	}
	s = strings.Trim(s, "\"")
	if len(s) == 0 {
		return s
	}
	multiplier := int64(1)
	last := s[len(s)-1]
	hasSuffix := true
	switch last {
	case 'K', 'k':
		multiplier = 1024
	case 'M', 'm':
		multiplier = 1024 * 1024
	case 'G', 'g':
		multiplier = 1024 * 1024 * 1024
	case 'T', 't':
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		hasSuffix = false
	}
	numPart := s
	if hasSuffix {
		numPart = s[:len(s)-1]
	}
	if val, err := strconv.ParseInt(numPart, 10, 64); err == nil {
		return val * multiplier
	}
	return s
}

func parseIPSB(v any) net.IP {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}
	s = strings.Trim(s, "[] ")
	return net.ParseIP(s)
}

func parseIPSBSA(s string) []net.IP {
	if s == "" {
		return []net.IP{}
	}
	s = strings.NewReplacer("[", "", "]", "").Replace(s)
	parts := strings.Split(s, ",")
	res := make([]net.IP, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		if ip := net.ParseIP(trimmed); ip != nil {
			res = append(res, ip)
		}
	}
	return res
}

func (d *decodeState) readCgpBlock(open, close byte) []byte {
	tempSpan := d.span
	totalB64Len := 0
	for {
		for tempSpan < d.len && d.data[tempSpan] != close {
			totalB64Len++
			tempSpan++
		}
		if tempSpan < d.len {
			tempSpan++
		}
		for tempSpan < d.len && isSpace(d.data[tempSpan]) {
			tempSpan++
		}
		if tempSpan < d.len && d.data[tempSpan] == open {
			tempSpan++
			continue
		}
		break
	}

	result := make([]byte, 0, base64.StdEncoding.DecodedLen(totalB64Len))
	for d.span < d.len {
		d.scratch.Reset()
		for d.span < d.len && d.data[d.span] != close {
			d.scratch.WriteByte(d.data[d.span])
			d.span++
		}
		if d.span < d.len {
			d.span++
		}

		if d.scratch.Len() > 0 {
			src := d.scratch.Bytes()
			start := len(result)
			n := base64.StdEncoding.DecodedLen(len(src))
			result = result[:start+n]
			actualLen, err := base64.StdEncoding.Decode(result[start:], src)
			if err == nil {
				result = result[:start+actualLen]
			} else {
				result = result[:start]
				result = append(result, src...)
			}
		}
		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == open {
			d.span++
			continue
		}
		break
	}
	return result
}

func tryCgpBool(b []byte) (any, bool) {
	switch len(b) {
	case 2:
		if bytes.Equal(b, binNoUpper) || bytes.Equal(b, binNoLower) {
			return false, true
		}
	case 3:
		if bytes.Equal(b, binYesUpper) || bytes.Equal(b, binYesLower) {
			return true, true
		}
	}
	return nil, false
}

func tryParseInt(b []byte) (any, bool) {
	if len(b) > 0 && ((b[0] >= '0' && b[0] <= '9') || b[0] == '-') {
		s := unsafe.String(unsafe.SliceData(b), len(b))
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			return val, true
		}
	}
	return nil, false
}
