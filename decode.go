package cgpcli

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const cgpInternalTimeLayout = "02-01-2006_15:04:05"

var (
	dateFields = map[string]bool{
		"Created": true,
	}
	quotaFields = map[string]bool{
		"MaxAccountSize": true,
		"MaxWebSize":     true,
		"MaxMailboxSize": true,
		"MaxMessageSize": true,
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
	// Предвыделяем память под атомы. 4КБ покрывает 99% нужд CLI.
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

			if dateFields[key] {
				if s, ok := val.(string); ok && s != "" {
					if t, err := time.Parse(cgpInternalTimeLayout, s); err == nil {
						val = t
					}
				}
			} else if quotaFields[key] {
				if s, ok := val.(string); ok {
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

	// 1. Проверка на блок данных: [..], #[..], #I[..]
	// Если это блок, мы СРАЗУ уходим в readCgpBlock и возвращаем []byte
	if curr == '[' {
		d.span++
		return d.readCgpBlock('[', ']')
	}

	if curr == '#' && d.span+1 < d.len {
		next := d.data[d.span+1]
		// Это блок данных с префиксом #[ или #I[
		if next == '[' {
			d.span += 2
			return d.readCgpBlock('[', ']')
		}
		if next == 'I' && d.span+2 < d.len && d.data[d.span+2] == '[' {
			d.span += 3
			return d.readCgpBlock('[', ']')
		}
		// Если это не блок, значит это #число, #Tдата или #(структура)
		// Оставляем это для обработки ниже (isTagged = true)
	}

	d.scratch.Reset()
	isTagged := curr == '#'
	isQuoted := curr == '"'
	isUnkData := false

	if isQuoted {
		// ... (Логика кавычек без изменений) ...
		for d.span < d.len && d.data[d.span] == '"' {
			d.span++
			for d.span < d.len {
				ch := d.data[d.span]
				if ch == '\\' && d.span+3 < d.len {
					octStr := string(d.data[d.span+1 : d.span+4])
					if val, err := strconv.ParseInt(octStr, 8, 32); err == nil {
						d.scratch.WriteByte(byte(val))
						d.span += 4
						continue
					}
					d.scratch.WriteByte(d.data[d.span+1])
					d.span += 2
					continue
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
	} else {
		// 2. Логика для обычных слов, чисел и тегов #(...)
		if isTagged {
			if bytes.HasPrefix(d.data[d.span:], []byte("#(")) {
				isUnkData = true
				d.scratch.WriteString("#(")
				d.span += 2
			} else {
				// НЕ пропускаем '#', он нам нужен в scratch для парсинга даты/числа
				// Но так как старый код делал d.span++, нам нужно решить:
				// либо писать '#' в scratch, либо обрезать его потом.
				// Чтобы не ломать вашу логику ниже (raw[1:]), записываем '#'
				d.scratch.WriteByte('#')
				d.span++
			}
		}

		for d.span < d.len {
			ch := d.data[d.span]
			if isUnkData {
				if ch == ')' {
					d.scratch.WriteByte(')')
					d.span++
					break
				}
			} else {
				if isSpace(ch) || ch == '=' || ch == ';' || ch == ',' || ch == ')' || ch == '}' {
					break
				}
			}
			d.scratch.WriteByte(ch)
			d.span++
		}
	}

	raw := d.scratch.String()
	if asString || isQuoted {
		return raw
	}

	if isTagged && !isUnkData {
		// raw сейчас выглядит как "#T02-..." или "#123"
		if strings.HasPrefix(raw, "#T") {
			if t, err := time.Parse(cgpInternalTimeLayout, raw[2:]); err == nil {
				return t
			}
		}
		// Пробуем парсить как число, отрезав '#'
		if val, err := strconv.ParseInt(raw[1:], 10, 64); err == nil {
			return val
		}
		return raw
	}

	if val, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return val
	}
	return raw
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
