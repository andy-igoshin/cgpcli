package cgpcli

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Secure int

const (
	Plain Secure = iota
	APOP
	WebUser
)

var (
	cmdSTLS   = []byte("STLS\n")
	cmdQUIT   = []byte("QUIT\n")
	cmdINLINE = []byte("INLINE\n")
)

const (
	cgpLineSize = 129 * 1024
	cgpTimeout  = 5*time.Minute - 5*time.Second
	cmdTimeout  = 1 * time.Minute
	connTimeout = 7 * time.Second
)

type CliCode int

const (
	CodeOk       CliCode = 200
	CodeOkInline         = 201
	CodePassword         = 300
	CodeGenErr           = 500
	CodeGenErr1          = 501
	CodeStrange          = 10000
)

const defaultPort = "106"

type Cli struct {
	host        string
	login       string
	password    string
	secure      Secure
	tls         bool
	debug       bool
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	buf         bytes.Buffer
	isConnected bool
	lastAccess  time.Time
	respCode    CliCode
	respData    []byte
}

func (cli *Cli) GetRespCode() CliCode {
	return cli.respCode
}

func (cli *Cli) GetRespData() string {
	return string(cli.respData)
}

func (cli *Cli) IsSuccess() bool {
	return (cli.respCode == CodeOk || cli.respCode == CodeOkInline)
}

func (cli *Cli) Logout() (err error) {

	if cli.conn != nil {
		cli.writer.Write(cmdQUIT)
		if err = cli.writer.Flush(); err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil && err != io.EOF {
			cli.close()
			return
		}
		cli.close()
	}

	return
}

func New(host, login, password string, secure Secure, tls bool) (cli *Cli, err error) {

	cli = new(Cli)

	cli.host = net.JoinHostPort(host, defaultPort)
	cli.login = login
	cli.password = password
	cli.secure = secure
	cli.tls = tls

	err = cli.connect()

	return
}

func (cli *Cli) SetDebug(debug bool) {
	cli.debug = debug
}

func (cli *Cli) close() {
	if cli.conn != nil {
		cli.conn.Close()
		cli.conn = nil
	}
}

func (cli *Cli) connect() (err error) {

	// 1. Устанавливаем базовое TCP соединение
	cli.conn, err = net.DialTimeout("tcp", cli.host, connTimeout)
	if err != nil {
		return
	}

	cli.reader = bufio.NewReaderSize(cli.conn, cgpLineSize)
	cli.writer = bufio.NewWriter(cli.conn)
	cli.buf.Grow(cgpLineSize)

	err = cli.parseResponse()
	if err != nil {
		err = fmt.Errorf("can't read CGPro Server prompt: %w", err)
		cli.close()
		return
	}

	cli.conn.SetWriteDeadline(time.Now().Add(cmdTimeout))

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

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		tlsConn := tls.Client(cli.conn, tlsconfig)
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
	cli.conn.SetReadDeadline(time.Now().Add(cmdTimeout))

	// Внутренний буфер CGP имеет размер 128k
	line, err := cli.reader.ReadSlice('\n')
	if err != nil {
		return fmt.Errorf("can't read line: %s", err)
	}

	line = bytes.TrimRight(line, "\r\n")

	idx := bytes.IndexByte(line, ' ')
	if idx == -1 {
		cli.respCode = parseCode(line)
		cli.respData = nil
		return nil
	}

	code := parseCode(line)
	if code < CodeOk || code > CodeGenErr1 {
		cli.respCode = CodeStrange
		cli.respData = line
		return nil
	}

	cli.respCode = CliCode(code)
	cli.respData = line[idx+1:]
	cli.lastAccess = time.Now()

	return nil
}

func (cli *Cli) Query(args ...any) (any, error) {
	if err := cli.Send(args...); err != nil {
		return nil, err
	}
	if err := cli.parseResponse(); err != nil {
		return nil, err
	}
	if !cli.IsSuccess() {
		return nil, fmt.Errorf("CGP Error %d: %s", cli.respCode, string(cli.respData))
	}

	if cli.debug {
		os.Stderr.Write([]byte("<<< RECV: "))
		os.Stderr.Write(cli.respData)
		os.Stderr.Write([]byte("\n"))
	}

	var res any
	err := Unmarshal(cli.respData, &res)
	return res, err
}

func (cli *Cli) Send(args ...any) (err error) {
	if time.Since(cli.lastAccess) > cgpTimeout || cli.conn == nil {
		cli.close()
		if err := cli.connect(); err != nil {
			return err
		}
	}

	cli.lastAccess = time.Now()
	cli.conn.SetWriteDeadline(time.Now().Add(cmdTimeout))

	cli.buf.Reset()

	for i, arg := range args {
		if i > 0 {
			cli.buf.WriteByte(' ')
		}
		if err := marshal(&cli.buf, arg, false); err != nil {
			return err
		}
	}
	cli.buf.WriteByte('\n')

	if cli.debug {
		os.Stderr.Write([]byte(">>> SEND: "))
		os.Stderr.Write(cli.buf.Bytes())
	}

	_, err = cli.writer.Write(cli.buf.Bytes())
	if err != nil {
		return err
	}

	return cli.writer.Flush()
}
