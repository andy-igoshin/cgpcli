package cgpcli

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	dateFields = map[string]bool{
		"Created": true,
	}
	durationFields = map[string]bool{
		"LDAPPagedCacheTimeout": true,
	}
	quotaFields = map[string]bool{
		"ArchiveSizeLimit": true,
		"DigestSizeLimit":  true,
		"MaxAccountSize":   true,
		"MaxFileSize":      true,
		"MaxMailboxSize":   true,
		"MaxMailOutSize":   true,
		"MaxMessageSize":   true,
		"MaxWebSize":       true,
		"SizeLimit":        true,
	}
)

type decodeState struct {
	data    []byte
	len     int
	span    int
	scratch bytes.Buffer
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.len = len(data)
	d.span = 0
	d.scratch.Grow(4096)
	return d
}

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
}

func (d *decodeState) skipSpaces() {
	for d.span < d.len && isSpace(d.data[d.span]) {
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
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	val := reflect.ValueOf(result)
	if val.Type().AssignableTo(target.Type()) {
		target.Set(val)
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
	res := make(map[string]any)
	for d.span < d.len {
		d.skipSpaces()
		if d.span >= d.len || d.data[d.span] == '}' {
			d.span++
			break
		}

		keyRaw := d.wordInterface(true)
		key, ok := keyRaw.(string)
		if !ok {
			break
		}

		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == '=' {
			d.span++
			val := d.valueInterface(false)

			if s, ok := val.(string); ok && s != "" {
				switch {
				case dateFields[key]:
					if t, err := time.Parse(cgpInternalTimeLayout, s); err == nil {
						val = t
					}
				case durationFields[key]:
					val = parseCGPDuration(s)
				case quotaFields[key]:
					val = parseCGPSize(s)
				}
			}
			res[key] = val
		}
		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == ';' {
			d.span++
		}
	}
	return res
}

func (d *decodeState) arrayInterface() any {
	res := make([]any, 0)
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

	// 1. Обработка блоков данных: [..], #[..], #I[..]
	if curr == '[' {
		d.span++
		return d.readCgpBlock('[', ']')
	}

	if curr == '#' && d.span+1 < d.len {
		next := d.data[d.span+1]
		// Специфический тип IP Address: #I[addr]:port
		if next == 'I' && d.span+2 < d.len && d.data[d.span+2] == '[' {
			d.span += 3
			ipData := d.readCgpBlock('[', ']')
			res := CGPIP{IP: net.ParseIP(string(ipData))}
			if d.span < d.len && d.data[d.span] == ':' {
				d.span++
				start := d.span
				for d.span < d.len && d.data[d.span] >= '0' && d.data[d.span] <= '9' {
					d.span++
				}
				portStr := string(d.data[start:d.span])
				res.Port, _ = strconv.Atoi(portStr)
			}
			return res
		}
	}

	d.scratch.Reset()
	isQuoted := curr == '"'
	isTagged := curr == '#'
	isUnkData := false

	if isQuoted {
		// Логика кавычек с поддержкой конкатенации "str1""str2"
		for d.span < d.len && d.data[d.span] == '"' {
			d.span++ // пропускаем открывающую кавычку
			for d.span < d.len {
				ch := d.data[d.span]
				if ch == '\\' && d.span+3 < d.len {
					// Восьмеричные последовательности \NNN
					octStr := string(d.data[d.span+1 : d.span+4])
					if val, err := strconv.ParseInt(octStr, 8, 32); err == nil {
						d.scratch.WriteByte(byte(val))
						d.span += 4
						continue
					}
					// Обычный escape
					d.scratch.WriteByte(d.data[d.span+1])
					d.span += 2
					continue
				}
				if ch == '"' {
					d.span++
					// Если за кавычкой сразу идет еще одна — это конкатенация
					if d.span < d.len && d.data[d.span] == '"' {
						// В некоторых версиях CGP "" внутри кавычек — это просто кавычка
						d.scratch.WriteByte('"')
						d.span++
						continue
					}
					// Конец квотированной строки
					d.skipSpaces()
					break
				}
				d.scratch.WriteByte(ch)
				d.span++
			}
		}
	} else {
		// Логика для слов, чисел и тегов #(...)
		if isTagged {
			if bytes.HasPrefix(d.data[d.span:], []byte("#(")) {
				isUnkData = true
				// Оставляем span на #, цикл ниже его заберет
			} else {
				// Забираем # и сдвигаем span, чтобы не зациклиться
				d.scratch.WriteByte('#')
				d.span++
			}
		}

		for d.span < d.len {
			ch := d.data[d.span]
			if isUnkData {
				d.scratch.WriteByte(ch)
				d.span++
				if ch == ')' {
					break
				}
				continue
			}
			// Разделители CGP
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

	switch raw {
	case "YES":
		return true
	case "NO":
		return false
	}

	// Пост-обработка накопленного "сырого" слова
	if isTagged && !isUnkData {
		if raw == "#NULL#" {
			return nil
		}
		// Временные метки
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
			return raw
		}
		// Числа с базами (0x, 0o, 0b)
		return parseCGPNumber(raw)
	}

	// Попытка распарсить как обычное число (без префикса #)
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
		if ch == '<' {
			if d.span+1 < d.len {
				next := d.data[d.span+1]
				if next == '/' {
					depth--
				} else if next == '!' || next == '?' {
					// skip comments
				} else {
					depth++
				}
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

func (d *decodeState) readCgpBlock(open, close byte) []byte {
	d.scratch.Reset()
	for d.span < d.len {
		if d.data[d.span] == close {
			d.span++
			d.skipSpaces()
			// Обработка продолжения блока: [data1][data2]
			if d.span < d.len && d.data[d.span] == open {
				d.span++
				continue
			}
			break
		}
		d.scratch.WriteByte(d.data[d.span])
		d.span++
	}
	res := make([]byte, d.scratch.Len())
	copy(res, d.scratch.Bytes())
	return res
}
