package cgpcli

import (
	"net"
	"time"
)

// cgpInternalTimeLayout соответствует формату поля Created
const cgpInternalTimeLayout = "02-01-2006_15:04:05"

type FieldType int

const (
	TypeAuto     FieldType = iota
	TypeDate               // Формат Created (02-Jan-2006_15:04:05)
	TypeDuration           // Интервалы (10s, 1m, never, immediately)
	TypeIP                 // Одиночный IP адрес в скобках [1.2.3.4]
	TypeIPList             // Список IP адресов (DNSServers и др.)
	TypeIPSA               // Строка вида "[ip1], [ip2]" <-> []net.IP
	TypeSize               // Квоты и размеры (10K, 5M, unlimited)
	TypeString             // Принудительное сохранение значения как строки
)

// FieldExceptions — глобальный реестр исключений типов.
var FieldExceptions = map[string]FieldType{
	// dateFields
	"Created": TypeDate,

	// durationFields
	"BlockedLoginTime":      TypeDuration,
	"BlockedLoginsPeriod":   TypeDuration,
	"CacheNegative":         TypeDuration,
	"DefaultRegister":       TypeDuration,
	"DeleteInterval":        TypeDuration,
	"DequeuerWaitTime":      TypeDuration,
	"DialogMediaIdle":       TypeDuration,
	"DialogSignalIdle":      TypeDuration,
	"DialogTimeLimit":       TypeDuration,
	"InitialWait":           TypeDuration,
	"LDAPPagedCacheTimeout": TypeDuration,
	"MaxRegister":           TypeDuration,
	"MinRegister":           TypeDuration,
	"NATTCPPingPeriod":      TypeDuration,
	"NATUDPPingPeriod":      TypeDuration,
	"ProtocolErrorsPeriod":  TypeDuration,
	"RestartPause":          TypeDuration,
	"SessionTimeLimit":      TypeDuration,
	"SwapInterval":          TypeDuration,
	"TCPTimeWait":           TypeDuration,
	"TempClientTime":        TypeDuration,
	"TempUnBlockableTime":   TypeDuration,
	"TimeLimit":             TypeDuration,
	"Timeout":               TypeDuration,

	// quotaFields (TypeSize)
	"ArchiveSizeLimit": TypeSize,
	"DigestSizeLimit":  TypeSize,
	"MaxAccountSize":   TypeSize,
	"MaxFileSize":      TypeSize,
	"MaxMailboxSize":   TypeSize,
	"MaxMailOutSize":   TypeSize,
	"MaxMessageSize":   TypeSize,
	"MaxWebSize":       TypeSize,
	"MBOXSizeToIndex":  TypeSize,
	"SizeLimit":        TypeSize,

	// Сетевые исключения
	"DNSServers":  TypeIPSA,
	"DummyIPs":    TypeIPList,
	"LANAddress":  TypeIP,
	"WANAddress":  TypeIP,
	"WAN6Address": TypeIP,
}

func RegisterFieldType(name string, t FieldType) {
	FieldExceptions[name] = t
}

// transformRead преобразует сырое значение от парсера в целевой тип Go.
func transformRead(key string, val any) any {
	fType, exists := FieldExceptions[key]
	if !exists {
		return val
	}

	// Обработка TypeString: блокируем авто-типизацию
	if fType == TypeString {
		s, _ := toString("registry", val)
		return s
	}

	// Для всех остальных трансформаций CGP нам нужна непустая строка
	s, ok := val.(string)
	if !ok || s == "" {
		return val
	}

	switch fType {
	case TypeDate:
		// В CGP всегда UTC.
		if tm, err := time.Parse(cgpInternalTimeLayout, s); err == nil {
			return tm.UTC()
		}
	case TypeDuration:
		return parseCGPDuration(s)
	case TypeSize:
		return parseCGPSize(s)
	case TypeIP:
		return parseIPSingle(s)
	case TypeIPSA:
		return parseIPSA(s)
	case TypeIPList:
		return parseIPList(s)
	}
	return val
}

// transformWrite подготавливает объект Go к маршалингу в синтаксис CGP.
func transformWrite(key string, v any) (string, bool) {
	fType, exists := FieldExceptions[key]
	if !exists {
		return "", false
	}

	switch fType {
	case TypeSize:
		return formatCGPSize(toInt64(v)), true

	case TypeDuration:
		if d, ok := v.(time.Duration); ok {
			return formatCGPDuration(d), true
		}

	case TypeIP:
		if ip, ok := v.(net.IP); ok {
			return "[" + ip.String() + "]", true
		}

	case TypeIPSA:
		if ips, ok := v.([]net.IP); ok {
			return toIPsString(ips), true
		}

	case TypeString:
		s, _ := toString("registry", v)
		return s, true

	case TypeDate:
		if tm, ok := v.(time.Time); ok {
			if tm.IsZero() {
				return "#NULL#", true
			}
			// Возвращаем строку без префикса #T,
			// так как маршалер в зависимости от контекста может добавить его сам
			return tm.UTC().Format(cgpInternalTimeLayout), true
		}
	}
	return "", false
}
