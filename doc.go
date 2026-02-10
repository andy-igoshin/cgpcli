/*
Package cgpcli implements a high-performance parser and marshaler for CommuniGate Pro CLI.

# Type Mapping

The library automatically maps CommuniGate Pro data structures to native Go types.
The mapping logic is divided into standard protocol types and specialized field
transformations based on the CGP documentation.

Standard Protocol Types:
  - Atoms & Strings: Mapped to Go string.
  - Numbers: All CGP numeric values are parsed as int64.
  - Booleans: CGP "YES" and "NO" are mapped to Go bool.
  - Arrays: Parenthesized lists (val1, val2, ...) are mapped to []any or []string.
  - Dictionaries: Key-value structures {key=val; ...} are mapped to map[string]any or map[string]string.
  - Data Blocks: Block data [base64] is mapped to Go []byte.
  - IP addresses: Block data #I[ip] is converted to net.IP.

Specialized Transformations (Registry-based):
To provide a more ergonomic API, cgpcli recognizes specific field names and
converts their string values into rich Go types:

  - Dates: Fields like "Created" or "LastLogin" are parsed using the
    "DD-MM-YYYY_HH:MM:SS" layout and returned as time.Time in UTC.
  - Time Intervals (Durations): Fields representing time-outs or periods
    (e.g., "Timeout") support CGP-specific units like "10s", "1m", "2h",
    as well as "never" and "immediately". These are mapped to time.Duration.
  - Network Addresses: IP addresses in square brackets [1.2.3.4] are
    automatically converted to net.IP.
  - Storage Sizes: Quota and file size fields (e.g., "MaxAccountSize")
    supporting suffixes like "K", "M", "G" or "unlimited" are
    converted to int64 bytes.

# Specialized Values

The library handles unique CGP constants:
  - #NULL# is mapped to nil.
  - #TPAST and #TFUTURE are mapped to specific time.Time constants
    provided by the library.
*/
package cgpcli
