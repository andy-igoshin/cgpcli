// # CLI Access
//
// The CLI Access section implements the core transport and session management for
// CommuniGate Pro. It handles the low-level communication protocol, ensuring
// reliable and secure data exchange between the Go application and the CommuniGate Pro CLI.
//
// Key capabilities include:
//   - Session Lifecycle: automatic connection establishment, authentication,
//     and graceful termination via [Cli.Close].
//   - Persistence: transparent re-connection logic if the session exceeds
//     the server's idle timeout (approx. 5 minutes).
//   - Security: support for standard [Plain] login, MD5-based [APOP], and [WebUser]
//     authentication, with optional TLS (STARTTLS) encryption.
//   - Efficiency: mandatory use of CommuniGate Pro CLI "INLINE" mode, which, combined with
//     the library's custom parser, ensures high-performance data processing.
//   - Concurrency: designed for single-goroutine execution per [Cli] session.
package cgpcli

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

// Secure defines the authentication method used for the CommuniGate Pro CLI session.
type Secure int

const (
	// Plain represents standard USER/PASS authentication.
	Plain Secure = iota

	// APOP represents Authenticated Post Office Protocol (MD5-based) security.
	APOP

	// WebUser represents authentication for web-specific users.
	WebUser
)

// Option is a functional option for configuring a [Cli] session.
type Option func(*Cli)

// WithCommandTimeout sets the read/write deadline for a single CLI command.
// Non-positive values are ignored (the default of 1 minute is kept).
//
// Parameters:
//   - d: the per-command read/write deadline.
//
// Returns:
//   - Option: a functional option that applies the configuration.
func WithCommandTimeout(d time.Duration) Option {
	return func(cli *Cli) {
		if d > 0 {
			cli.cmdTimeout = d
		}
	}
}

// WithConnectTimeout sets the TCP dial timeout.
// Non-positive values are ignored (the default of 7 seconds is kept).
//
// Parameters:
//   - d: the TCP dial timeout.
//
// Returns:
//   - Option: a functional option that applies the configuration.
func WithConnectTimeout(d time.Duration) Option {
	return func(cli *Cli) {
		if d > 0 {
			cli.connTimeout = d
		}
	}
}

// WithIdleTimeout sets the session idle timeout before automatic reconnect.
// Non-positive values are ignored (default ~5 minutes).
//
// Parameters:
//   - d: the session idle timeout before automatic reconnect.
//
// Returns:
//   - Option: a functional option that applies the configuration.
func WithIdleTimeout(d time.Duration) Option {
	return func(cli *Cli) {
		if d > 0 {
			cli.idleTimeout = d
		}
	}
}

// WithTLSConfig sets a custom TLS configuration for STARTTLS connections.
// If not provided, a secure default (certificate verification enabled) is used.
//
// Parameters:
//   - cfg: the TLS configuration to use for STARTTLS connections.
//
// Returns:
//   - Option: a functional option that applies the configuration.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(cli *Cli) {
		cli.tlsConfig = cfg
	}
}

// WithTLSInsecureSkipVerify controls whether TLS certificate verification is skipped.
// Use only for testing against servers with self-signed certificates.
//
// Parameters:
//   - skip: if true, TLS certificate verification is disabled.
//
// Returns:
//   - Option: a functional option that applies the configuration.
func WithTLSInsecureSkipVerify(skip bool) Option {
	return func(cli *Cli) {
		if cli.tlsConfig == nil {
			cli.tlsConfig = &tls.Config{}
		}
		cli.tlsConfig.InsecureSkipVerify = skip
	}
}

var (
	cmdSTLS   = []byte("STLS\n")
	cmdQUIT   = []byte("QUIT\n")
	cmdINLINE = []byte("INLINE\n")
)

const (
	cgpLineSize = 129 * 1024

	defaultCmdTimeout  = 1 * time.Minute
	defaultConnTimeout = 7 * time.Second
	defaultIdleTimeout = 5*time.Minute - 5*time.Second
)

// CliCode represents a numeric response status from the CommuniGate Pro server.
type CliCode int

const (
	CodeOk       CliCode = 200
	CodeOkInline CliCode = 201
	CodePassword CliCode = 300 // please send the PASS
	CodeErr500   CliCode = 500 // generic error
	CodeErr501   CliCode = 501 // generic error
	CodeErr510   CliCode = 510 // you do not have the required access right
	CodeErr512   CliCode = 512 // unknown secondary domain name
	CodeErr513   CliCode = 513 // unknown user account
	CodeErr515   CliCode = 515 // incorrect password or account name
	CodeErr516   CliCode = 516 // mailbox not found
	CodeErr520   CliCode = 520 // account with this name already exists
	CodeErr521   CliCode = 521 // account alias with this name already exists
	CodeErr522   CliCode = 522 // mailing list with this name already exists
	CodeErr523   CliCode = 523 // group with this name already exists
	CodeErr524   CliCode = 524 // forwarder with this name already exists
	CodeErr525   CliCode = 525
	CodeErr532   CliCode = 532 // mailbox with this name already exists
	CodeErr552   CliCode = 552 // group is not found
	CodeErr553   CliCode = 553 // forwarder is not found
	CodeErr599   CliCode = 599 // file is not found
	CodeStrange  CliCode = 10000
)

// Sentinel errors for known CGP CLI error responses.
var (
	ErrAccessDenied      = errors.New("access denied")
	ErrDomainUnknown     = errors.New("unknown domain")
	ErrAccountUnknown    = errors.New("unknown account")
	ErrAuthFailed        = errors.New("authentication failed")
	ErrMailboxNotFound   = errors.New("mailbox not found")
	ErrAccountExists     = errors.New("account already exists")
	ErrAliasExists       = errors.New("alias already exists")
	ErrMailListExists    = errors.New("mailing list already exists")
	ErrGroupExists       = errors.New("group already exists")
	ErrForwarderExists   = errors.New("forwarder already exists")
	ErrMailboxExists     = errors.New("mailbox already exists")
	ErrGroupNotFound     = errors.New("group not found")
	ErrForwarderNotFound = errors.New("forwarder not found")
	ErrFileNotFound      = errors.New("file not found")
)

var cgpErrors = [100]error{
	10: ErrAccessDenied,      // 510
	12: ErrDomainUnknown,     // 512
	13: ErrAccountUnknown,    // 513
	15: ErrAuthFailed,        // 515
	16: ErrMailboxNotFound,   // 516
	20: ErrAccountExists,     // 520
	21: ErrAliasExists,       // 521
	22: ErrMailListExists,    // 522
	23: ErrGroupExists,       // 523
	24: ErrForwarderExists,   // 524
	32: ErrMailboxExists,     // 532
	52: ErrGroupNotFound,     // 552
	53: ErrForwarderNotFound, // 553
	99: ErrFileNotFound,      // 599
}

const defaultPort = "106"

// Cli represents a CommuniGate Pro CLI client session.
// It handles connection persistence, automatic reconnection, and command marshaling.
// Use New() to initialize a session.
//
// Concurrency:
// Cli instances are NOT safe for concurrent use by multiple goroutines.
// If multiple goroutines need to communicate with CommuniGate Pro concurrently,
// each goroutine must maintain its own Cli instance or serialize calls via
// external synchronization.
type Cli struct {
	host        string
	login       string
	password    string
	secure      Secure
	tls         bool
	tlsConfig   *tls.Config
	debug       bool
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	buf         bytes.Buffer
	isConnected bool
	lastAccess  time.Time
	respCode    CliCode
	respData    []byte

	// Configurable timeouts.
	cmdTimeout  time.Duration // read/write deadline for a single command.
	connTimeout time.Duration // TCP dial timeout.
	idleTimeout time.Duration // session idle timeout before automatic reconnect.
}

// Close gracefully terminates the CommuniGate Pro CLI session.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) Close() (err error) {
	if cli.conn == nil {
		return nil
	}

	cli.conn.SetWriteDeadline(time.Now().Add(cli.cmdTimeout))
	cli.writer.Write(cmdQUIT)
	if err = cli.writer.Flush(); err != nil {
		cli.close()
		return
	}

	_ = cli.parseResponse()

	cli.close()
	return
}

// GetRespCode returns the numeric status code from the last server response.
//
// Returns:
//   - CliCode: a numeric status code (e.g., 200 for OK, 201 for OK INLINE).
func (cli *Cli) GetRespCode() CliCode {
	return cli.respCode
}

// GetRespData returns the raw string data (payload) from the last server response.
//
// Returns:
//   - string: the data portion of the CLI response, following the status code.
func (cli *Cli) GetRespData() string {
	return string(cli.respData)
}

// IsSuccess returns true if the last response code indicates success (200 or 201).
func (cli *Cli) IsSuccess() bool {
	return (cli.respCode == CodeOk || cli.respCode == CodeOkInline)
}

// Logout gracefully terminates the CommuniGate Pro CLI session.
// It is an alias for [Cli.Close].
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) Logout() error {
	return cli.Close()
}

// New creates and initializes a new CommuniGate Pro CLI session.
// It establishes a connection, performs authentication, and sets the session
// to INLINE mode for efficient data exchange.
//
// Parameters:
//   - host: the server hostname or IP address.
//   - login: the administrative or user login name.
//   - password: the password for authentication.
//   - secure: the security mode to use (Plain, APOP, or WebUser).
//   - tls: if true, attempts to upgrade the connection via STARTTLS.
//
// Returns:
//   - *Cli: a pointer to the initialized session.
//   - error: an error if the command fails.
func New(host, login, password string, secure Secure, tls bool, opts ...Option) (*Cli, error) {

	cli := &Cli{
		cmdTimeout:  defaultCmdTimeout,
		connTimeout: defaultConnTimeout,
		idleTimeout: defaultIdleTimeout,
	}

	cli.host = net.JoinHostPort(host, defaultPort)
	cli.login = login
	cli.password = password
	cli.secure = secure
	cli.tls = tls

	for _, opt := range opts {
		opt(cli)
	}

	err := cli.connect()

	return cli, err
}

// Query sends a command to CommuniGate Pro CLI and unmarshals the response into a Go type.
// It automatically handles connection health and session state.
//
// Parameters:
//   - cmd: the CLI command string (e.g., "GETACCOUNTSETTINGS").
//   - args: variadic arguments that will be marshaled into the CLI command format.
//
// Returns:
//   - any: the unmarshaled result from the server.
//   - error: an error if the command fails.
func (cli *Cli) Query(cmd string, args ...any) (res any, err error) {
	if err = cli.Send(cmd, args...); err != nil {
		return
	}
	if err = cli.parseResponse(); err != nil {
		return
	}
	if !cli.IsSuccess() {
		err = cli.makeCgpError()
		return
	}

	if len(cli.respData) == 0 {
		return nil, nil
	}

	err = Unmarshal(cli.respData, &res)
	return
}

// QueryNV (No Value) sends a command where only the success status code is needed.
// It is useful for operations like updates where the server doesn't return data.
//
// Parameters:
//   - cmd: the CLI command string.
//   - args: variadic arguments for the command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) QueryNV(cmd string, args ...any) (err error) {
	if err = cli.Send(cmd, args...); err != nil {
		return
	}
	if err = cli.parseResponse(); err != nil {
		return
	}
	if !cli.IsSuccess() {
		err = cli.makeCgpError()
	}
	return
}

// Send marshals and sends a raw command to the CommuniGate Pro CLI.
// It handles automatic reconnection if the session has timed out.
//
// Parameters:
//   - cmd: the CLI command string.
//   - args: variadic arguments to be appended to the command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) Send(cmd string, args ...any) (err error) {
	if time.Since(cli.lastAccess) > cli.idleTimeout || cli.conn == nil {
		cli.close()
		if err = cli.connect(); err != nil {
			return
		}
	}

	cli.lastAccess = time.Now()
	cli.conn.SetWriteDeadline(time.Now().Add(cli.cmdTimeout))
	cli.buf.Reset()

	cli.buf.WriteString(cmd)

	for _, arg := range args {
		cli.buf.WriteByte(' ')
		if err = marshal(&cli.buf, arg, ""); err != nil {
			return
		}
	}
	cli.buf.WriteByte('\n')

	if cli.debug {
		os.Stderr.Write([]byte(">>> SEND: "))
		os.Stderr.Write(cli.buf.Bytes())
	}

	_, err = cli.writer.Write(cli.buf.Bytes())
	if err != nil {
		cli.close()
		return fmt.Errorf("write error: %w", err)
	}

	err = cli.writer.Flush()
	if err != nil {
		cli.close()
		return fmt.Errorf("flush error: %w", err)
	}

	return nil
}

// SetDebug enables or disables raw protocol logging to os.Stderr.
//
// Parameters:
//   - debug: if true, all sent (>>>) and received (<<<) raw data will
//     be printed to the standard error stream.
func (cli *Cli) SetDebug(debug bool) {
	cli.debug = debug
}

// SetCommandTimeout dynamically changes the per-command deadline. It is
// intended for use before issuing long-running commands (e.g., large exports).
// Non-positive values are ignored.
//
// Parameters:
//   - d: the new per-command read/write deadline.
func (cli *Cli) SetCommandTimeout(d time.Duration) {
	if d > 0 {
		cli.cmdTimeout = d
	}
}

// SetConnectTimeout dynamically changes the TCP dial timeout used on (re)connect.
// Non-positive values are ignored.
//
// Parameters:
//   - d: the new TCP dial timeout.
func (cli *Cli) SetConnectTimeout(d time.Duration) {
	if d > 0 {
		cli.connTimeout = d
	}
}

func (cli *Cli) close() {
	if cli.conn != nil {
		_ = cli.conn.Close()
		cli.conn = nil
	}
	cli.isConnected = false
}

func (cli *Cli) connect() (err error) {

	// 1. Устанавливаем базовое TCP соединение
	cli.conn, err = net.DialTimeout("tcp", cli.host, cli.connTimeout)
	if err != nil {
		return
	}

	if cli.reader == nil {
		cli.reader = bufio.NewReaderSize(cli.conn, cgpLineSize)
	} else {
		cli.reader.Reset(cli.conn)
	}

	if cli.writer == nil {
		cli.writer = bufio.NewWriter(cli.conn)
	} else {
		cli.writer.Reset(cli.conn)
	}

	cli.buf.Grow(cgpLineSize)

	err = cli.parseResponse()
	if err != nil {
		err = fmt.Errorf("can't read CGPro Server prompt: %w", err)
		cli.close()
		return
	}

	cli.conn.SetWriteDeadline(time.Now().Add(cli.cmdTimeout))

	// 2. Обработка STARTTLS (STLS)
	if cli.tls {
		cli.writer.Write(cmdSTLS)
		if err = cli.writer.Flush(); err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil || cli.respCode != CodeOk {
			if err == nil {
				err = fmt.Errorf("STLS failed: %d %s", cli.respCode, string(cli.respData))
			}
			cli.close()
			return
		}

		if cli.tlsConfig == nil {
			cli.tlsConfig = &tls.Config{}
		}
		tlsConn := tls.Client(cli.conn, cli.tlsConfig)
		if err = tlsConn.Handshake(); err != nil {
			cli.close()
			return
		}

		cli.conn = net.Conn(tlsConn)
		cli.reader.Reset(cli.conn)
		cli.writer.Reset(cli.conn)
	}

	// 3. Аутентификация
	switch cli.secure {
	case Plain:
		cli.writer.WriteString("USER ")
		cli.writer.WriteString(cli.login)
		cli.writer.WriteByte('\n')
		if err = cli.writer.Flush(); err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil || cli.respCode != CodePassword {
			if err == nil {
				err = fmt.Errorf("USER command failed: %d %s", cli.respCode, string(cli.respData))
			}
			cli.close()
			return
		}

		cli.writer.WriteString("PASS ")
		cli.writer.WriteString(cli.password)
		cli.writer.WriteByte('\n')
		if err = cli.writer.Flush(); err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil {
			cli.close()
			return
		}

	case APOP:
		idStart := bytes.LastIndexByte(cli.respData, '<')
		if idStart == -1 {
			cli.close()
			err = fmt.Errorf("APOP ID not found in server prompt")
			return
		}
		id := cli.respData[idStart:]
		if idEnd := bytes.IndexByte(id, '>'); idEnd != -1 {
			id = id[:idEnd+1]
		}

		hash := md5.Sum(append(id, cli.password...))
		digest := hex.EncodeToString(hash[:])

		cli.writer.WriteString("APOP ")
		cli.writer.WriteString(cli.login)
		cli.writer.WriteByte(' ')
		cli.writer.WriteString(digest)
		cli.writer.WriteByte('\n')
		if err = cli.writer.Flush(); err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil {
			cli.close()
			return
		}

	case WebUser:
		cli.writer.WriteString("AUTH WEBUSER ")
		cli.writer.WriteString(cli.login)
		cli.writer.WriteByte(' ')
		cli.writer.WriteString(cli.password)
		cli.writer.WriteByte('\n')
		if err = cli.writer.Flush(); err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil {
			cli.close()
			return
		}

	default:
		cli.close()
		err = fmt.Errorf("unknown secure type")
		return
	}

	// Проверка результата входа
	if cli.respCode != CodeOk {
		err = fmt.Errorf("login failed: %d %s", cli.respCode, string(cli.respData))
		cli.close()
		return
	}

	// 4. Переход в режим INLINE
	cli.writer.Write(cmdINLINE)
	if err = cli.writer.Flush(); err != nil {
		cli.close()
		return
	}
	err = cli.parseResponse()
	if err != nil || cli.respCode != CodeOk {
		if err == nil {
			err = fmt.Errorf("failed to set INLINE mode: %d %s", cli.respCode, string(cli.respData))
		}
		cli.close()
		return
	}

	cli.isConnected = true
	cli.lastAccess = time.Now()
	return
}

// makeCgpError builds an error from the last CGP server response.
// It attaches a sentinel error for known error codes and preserves
// the original server message for unknown or generic codes.
func (cli *Cli) makeCgpError() error {
	idx := int(cli.respCode) - 500
	if idx >= 0 && idx < len(cgpErrors) && cgpErrors[idx] != nil {
		return fmt.Errorf("(%d) %s (%w)", cli.respCode, string(cli.respData), cgpErrors[idx])
	}
	return fmt.Errorf("(%d) %s", cli.respCode, string(cli.respData))
}

func parseCode(b []byte) CliCode {
	code := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			break
		}
		code = code*10 + int(c-'0')
	}
	return CliCode(code)
}

func (cli *Cli) parseResponse() (err error) {
	cli.conn.SetReadDeadline(time.Now().Add(cli.cmdTimeout))

	line, err := cli.reader.ReadSlice('\n')
	if err != nil {
		cli.close()
		return fmt.Errorf("can't read line: %w", err)
	}

	line = bytes.TrimRight(line, "\r\n")

	if cli.debug {
		os.Stderr.Write([]byte("<<< RECV: "))
		os.Stderr.Write(line)
		os.Stderr.Write([]byte("\n"))
	}

	cli.lastAccess = time.Now()

	idx := bytes.IndexByte(line, ' ')
	if idx == -1 {
		cli.respCode = parseCode(line)
		cli.respData = nil
		return nil
	}

	code := parseCode(line)
	if code < CodeOk || code > CodeErr599 {
		cli.respCode = CodeStrange
		cli.respData = line
		return fmt.Errorf("STRANGE CODE: %d : %s", code, line)
	}

	cli.respCode = CliCode(code)
	cli.respData = line[idx+1:]

	return nil
}
