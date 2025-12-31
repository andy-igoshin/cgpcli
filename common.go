package cgpcli

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type Secure int

const (
	Plain Secure = iota
	APOP
	WebUser
)

const (
	cgpLineSize = 128 * 1024
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

const defaultPort = ":106"

type Cli struct {
	host        string
	login       string
	password    string
	secure      Secure
	tls         bool
	debug       bool
	conn        net.Conn
	reader      *bufio.Reader
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
		_, err = cli.conn.Write([]byte("QUIT\n"))
		if err != nil {
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

	if strings.Contains(host, ":") {
		cli.host = host
	} else {
		cli.host = host + defaultPort
	}

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

	// Инициализируем буферизированный ридер один раз при подключении
	cli.reader = bufio.NewReaderSize(cli.conn, cgpLineSize)

	// Читаем приветствие сервера (Banner)
	err = cli.parseResponse()
	if err != nil {
		err = fmt.Errorf("can't read CGPro Server prompt: %w", err)
		cli.close()
		return
	}

	cli.conn.SetWriteDeadline(time.Now().Add(cmdTimeout))

	// 2. Обработка STARTTLS (STLS)
	if cli.tls {
		_, err = cli.conn.Write([]byte("STLS\n"))
		if err != nil {
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

		// Апгрейд соединения до TLS
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		tlsConn := tls.Client(cli.conn, tlsconfig)
		if err = tlsConn.Handshake(); err != nil {
			cli.close()
			return
		}
		cli.conn = net.Conn(tlsConn)

		// Пересоздаем ридер для нового TLS-слоя
		cli.reader = bufio.NewReaderSize(cli.conn, cgpLineSize)
	}

	// 3. Аутентификация
	switch cli.secure {
	case Plain:
		_, err = cli.conn.Write([]byte("USER " + cli.login + "\n"))
		if err != nil {
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

		_, err = cli.conn.Write([]byte("PASS " + cli.password + "\n"))
		if err != nil {
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
		digest := fmt.Sprintf("%x", md5.Sum(append(id, []byte(cli.password)...)))

		_, err = cli.conn.Write([]byte("APOP " + cli.login + " " + digest + "\n"))
		if err != nil {
			cli.close()
			return
		}
		err = cli.parseResponse()
		if err != nil {
			cli.close()
			return
		}

	case WebUser:
		_, err = cli.conn.Write([]byte("AUTH WEBUSER " + cli.login + " " + cli.password + "\n"))
		if err != nil {
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
	_, err = cli.conn.Write([]byte("INLINE\n"))
	if err != nil {
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

func (cli *Cli) parseResponse() (err error) {
	cli.conn.SetReadDeadline(time.Now().Add(cmdTimeout))

	line, err := cli.reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("can't read line: %s", err)
	}

	line = bytes.TrimRight(line, "\r\n")

	// Ищем первый пробел для отделения кода
	idx := bytes.IndexByte(line, ' ')
	if idx == -1 {
		// Если пробела нет (редко, но бывает), считаем всё кодом или ошибкой
		code, _ := strconv.Atoi(string(line))
		cli.respCode = CliCode(code)
		cli.respData = nil
		return nil
	}

	code, err := strconv.Atoi(string(line[:idx]))
	if err != nil {
		cli.respCode = CodeStrange
		cli.respData = line
		return nil
	}

	cli.respCode = CliCode(code)
	cli.respData = line[idx+1:]
	cli.lastAccess = time.Now()

	return nil
}

func (cli *Cli) query(command string) (any, error) {
	if err := cli.send(command); err != nil {
		return nil, err
	}
	if err := cli.parseResponse(); err != nil {
		return nil, err
	}

	// Проверяем успешность (200 или 201)
	if !cli.IsSuccess() {
		return nil, fmt.Errorf("CGP Error %d: %s", cli.respCode, string(cli.respData))
	}

	var res any
	// Unmarshal теперь работает безопасно
	err := Unmarshal(cli.respData, &res)
	return res, err
}

func (cli *Cli) send(command string) (err error) {

	if time.Now().Sub(cli.lastAccess) > cgpTimeout || cli.conn == nil {
		cli.close()
		err = cli.connect()
		if err != nil {
			err = fmt.Errorf("can't reopen CLI connection: %s", err)
			return
		}
	}

	cli.lastAccess = time.Now()
	cli.conn.SetWriteDeadline(time.Now().Add(cmdTimeout))

	_, err = cli.conn.Write([]byte(command + "\n"))
	if err != nil {
		return
	}

	return
}
