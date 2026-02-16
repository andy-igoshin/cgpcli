// # Data Transformation and Field Registry
//
// The Registry section provides a high-level abstraction for data transformation
// between CommuniGate Pro string-based protocol and Go native types. Since CommuniGate Pro
// CLI often returns complex types (like "10M" for sizes or "10-02-2026_12:00:00"
// for dates) as plain strings, this registry ensures consistent and
// type-safe conversion.
//
// Key capabilities include:
//   - Field Exceptions: A global map [FieldExceptions] that defines how specific
//     keys (e.g., "Created", "MaxAccountSize") should be interpreted.
//   - Specialized Parsers: Built-in logic for handling legacy date formats,
//     storage quotas with suffixes (K, M, G), and time durations (s, m, h, never).
//   - Networking Support: Automatic marshaling of IP addresses in square
//     brackets (TypeIPSB) and comma-separated IP lists (TypeIPSBSA).
//   - Extensibility: The [RegisterFieldType] function allows developers to
//     add support for custom fields or override existing transformation rules.
//   - Symmetric Transformation: Automatic handling of both [transformRead]
//     (from CGP to Go) and [transformWrite] (from Go to CGP syntax).
//
// This registry is applied automatically during Query operations, making the
// interaction with CommuniGate Pro feel like working with a strongly-typed Go application.
package cgpcli

import (
	"net"
	"time"
)

// cgpInternalTimeLayout defines the standard CommuniGate Pro date format.
const cgpInternalTimeLayout = "02-01-2006_15:04:05"

// FieldType represents a data transformation category.
// It instructs the library how to map a specific field to a native Go type.
type FieldType int

const (
	// TypeAuto relies on the parser's basic detection (Atoms, Numbers, Arrays).
	TypeAuto FieldType = iota

	// TypeDateLegacy represents the "DD-MM-YYYY_HH:MM:SS" format.
	// Maps to: time.Time (UTC).
	TypeDateLegacy

	// TypeDuration represents time intervals (e.g., "10s", "1m", "never").
	// Maps to: time.Duration.
	TypeDuration

	// TypeIPList represents IP address lists as formatted in the CommuniGate Pro WebAdmin UI
	// (e.g., space-separated or multi-line addresses lists).
	// Maps to: IPList.
	TypeIPList

	// TypeIPSB represents an IP address in square brackets [1.2.3.4].
	// Maps to: net.IP.
	TypeIPSB

	// TypeIPSBSA represents a string of comma-separated square bracketed IPs: "[ip1],[ip2]".
	// Maps to: []net.IP.
	TypeIPSBSA

	// TypeSize represents storage quotas with suffixes (e.g., "10K", "5M", "unlimited").
	// Maps to: int64 (total bytes).
	TypeSize

	// TypeString forces the value to be treated as a raw string,
	// bypassing any automatic conversion.
	TypeString
)

// FieldExceptions is the global registry of field-to-type mappings.
// It contains predefined rules for core CommuniGate Pro fields.
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

// RegisterFieldType adds or overrides a field mapping in the global registry [FieldExceptions].
// Use this to support new CommuniGate Pro fields or change how existing fields
// are interpreted by the library's transformation engine.
//
// Parameters:
//   - name: The exact key name as it appears in the CommuniGate Pro protocol (e.g., "MyCustomTimeout").
//   - t: The target [FieldType] which defines the transformation logic.
//
// Note: This function must be called before initializing any new CLI sessions
// via [New]. The recommended approach is to call [RegisterFieldType] within
// an init() function in your application. You can call it as many times
// as necessary to register all required field exceptions.
//
// Example:
//
//	func init() {
//	    cgpcli.RegisterFieldType("CustomQuota", cgpcli.TypeSize)
//	}
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

// transformWrite подготавливает объект Go к маршалингу в синтаксис CommuniGate Pro.
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
