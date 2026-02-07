package cgpcli

import (
	"net"
	"time"
)

// cgpInternalTimeLayout соответствует формату поля Created
const cgpInternalTimeLayout = "02-01-2006_15:04:05"

type FieldType int

const (
	TypeAuto       FieldType = iota
	TypeDateLegacy           // Формат Created (01-09-2024_17:55:37)
	TypeDuration             // Интервалы (10s, 1m, never, immediately)
	TypeIPList               // Список IP адресов
	TypeIPSB                 // IP in Square Brackets [1.2.3.4] <-> net.IP
	TypeIPSBSA               // Строка вида "[ip1], [ip2]" <-> []net.IP
	TypeSize                 // Квоты и размеры (10K, 5M, unlimited)
	TypeString               // Принудительное сохранение значения как строки
)

// FieldExceptions — глобальный реестр исключений типов.
var FieldExceptions = map[string]FieldType{
	// dateFields
	"Created": TypeDateLegacy,

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
	"DNSServers":  TypeIPSBSA,
	"DummyIPs":    TypeIPList,
	"LANAddress":  TypeIPSB,
	"WANAddress":  TypeIPSB,
	"WAN6Address": TypeIPSB,
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
	case TypeDateLegacy:
		// В CGP всегда UTC.
		if tm, err := time.Parse(cgpInternalTimeLayout, s); err == nil {
			return tm.UTC()
		}
	case TypeDuration:
		return parseCGPDuration(s)
	case TypeSize:
		return parseCGPSize(s)
	case TypeIPSB:
		return parseIPSB(s)
	case TypeIPSBSA:
		return parseIPSBSA(s)
	case TypeIPList:
		return parseIPList(s)
	}
	return val
}

// transformWrite подготавливает объект Go к маршалингу в синтаксис CGP.
func transformWrite(key string, v any) any {
	fType, exists := FieldExceptions[key]
	if !exists {
		return v
	}

	switch fType {
	case TypeSize:
		return cgpSize(toInt64(v))

	case TypeDuration:
		if d, ok := v.(time.Duration); ok {
			return cgpDuration(d)
		}

	case TypeIPSB:
		if ip, ok := v.(net.IP); ok {
			return cgpIPSB(ip)
		}

	case TypeIPSBSA:
		if ips, ok := v.([]net.IP); ok {
			return cgpIPSBSA(ips)
		}

	case TypeString:
		s, _ := toString("registry", v)
		return cgpForceString(s)

	case TypeDateLegacy:
		if tm, ok := v.(time.Time); ok {
			return cgpDateLegacy(tm)
		}
	}

	return v
}
