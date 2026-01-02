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
	data []byte
	len  int
	span int
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.len = len(data)
	d.span = 0
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

// Unmarshal преобразует байты ответа CGP в структуру Go.
func Unmarshal(data []byte, v any) error {
	var d decodeState
	d.init(data)
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("cgpcli: Unmarshal expects non-nil pointer")
	}

	// Начинаем парсинг. Корневой элемент рассматривается как значение (asString = false)
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

		// КЛЮЧИ словаря всегда читаем как строки
		keyRaw := d.wordInterface(true)
		key, ok := keyRaw.(string)
		if !ok {
			break
		}

		d.skipSpaces()
		if d.span < d.len && d.data[d.span] == '=' {
			d.span++
			// ЗНАЧЕНИЯ могут быть типизированы (числа, даты)
			val := d.valueInterface(false)

			// Специальная обработка для дат и квот по имени ключа
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
		// Элементы массива могут быть именами или числами (ID), используем авто-детект
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
	startSpan := d.span
	var result bytes.Buffer
	isTagged := d.data[d.span] == '#'
	isQuoted := d.data[d.span] == '"'
	isUnkData := false

	if isQuoted {
		d.span++
	} else if isTagged {
		if bytes.HasPrefix(d.data[d.span:], []byte("#(")) {
			isUnkData = true
			result.WriteString("#(")
			d.span += 2
		} else if bytes.HasPrefix(d.data[d.span:], []byte("#I[")) || (d.span+1 < d.len && d.data[d.span+1] == '[') {
			prefix := "#["
			if d.span+1 < d.len && d.data[d.span+1] == 'I' {
				prefix = "#I["
				d.span++
			}
			result.WriteString(prefix)
			d.span += 2
			return d.readCgpBlock(&result, '[', ']', startSpan)
		} else {
			d.span++
		}
	}

	for d.span < d.len {
		ch := d.data[d.span]
		if isQuoted {
			if ch == '\\' && d.span+3 < d.len {
				octStr := string(d.data[d.span+1 : d.span+4])
				if val, err := strconv.ParseInt(octStr, 8, 32); err == nil {
					result.WriteByte(byte(val))
					d.span += 4
					continue
				}
				result.WriteByte(d.data[d.span+1])
				d.span += 2
				continue
			} else if ch == '"' {
				d.span++
				d.skipSpaces()
				if d.span < d.len && d.data[d.span] == '"' {
					d.span++
					continue
				}
				break
			}
		} else if isUnkData {
			if ch == ')' {
				result.WriteByte(')')
				d.span++
				break
			}
		} else {
			if isSpace(ch) || ch == '=' || ch == ';' || ch == ',' || ch == ')' || ch == '}' {
				break
			}
		}
		result.WriteByte(ch)
		d.span++
	}

	raw := result.String()

	// 1. Если это КЛЮЧ или строка в КАВЫЧКАХ - всегда возвращаем string
	if asString || isQuoted {
		return raw
	}

	// 2. Если есть тег # (явные типы данных в CGP)
	if isTagged && !isUnkData {
		if strings.HasPrefix(raw, "T") {
			if t, err := time.Parse(cgpInternalTimeLayout, raw[1:]); err == nil {
				return t
			}
		}
		if val, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return val
		}
		return "#" + raw
	}

	// 3. Умное распознавание для незакавыченных литералов
	if val, err := strconv.ParseInt(raw, 10, 64); err == nil {
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

func (d *decodeState) readCgpBlock(buf *bytes.Buffer, open, close byte, start int) string {
	for d.span < d.len {
		if d.data[d.span] == close {
			d.span++
			d.skipSpaces()
			if d.span < d.len && d.data[d.span] == open {
				d.span++
				continue
			}
			buf.WriteByte(close)
			break
		}
		buf.WriteByte(d.data[d.span])
		d.span++
	}
	return buf.String()
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
					// Пропуск комментариев/деклараций до '>'
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
