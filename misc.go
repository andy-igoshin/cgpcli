// # Section: Miscellaneous Commands
package cgpcli

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

type SystemInfoKey string

const (
	SysVersion         SystemInfoKey = "serverVersion"
	SysOS              SystemInfoKey = "serverOS"
	SysCPU             SystemInfoKey = "serverCPU"
	SysMainDomain      SystemInfoKey = "mainDomainName"
	SysLicenseDomain   SystemInfoKey = "licenseDomainName"
	SysStartTime       SystemInfoKey = "startTime"
	SysServerInstance  SystemInfoKey = "serverInstance"
	SysClusterInstance SystemInfoKey = "clusterInstance"
	SysCharsets        SystemInfoKey = "knownCharsets"
	SysDigesters       SystemInfoKey = "knownDigesters"
	SysCiphers         SystemInfoKey = "knownCiphers"
)

type TraceFacility string

const (
	TraceFileIO TraceFacility = "FileIO"
	TraceFileOp TraceFacility = "FileOp"
)

type TempIP struct {
	IP        net.IP
	TimeLeft  time.Duration // Интервал до разблокировки
	Permanent bool          // Признак вечной блокировки (*)
}

// DumpAllObjects пишет список всех объектов данных.
func (cli *Cli) DumpAllObjects(file bool) error {
	if file {
		_, err := cli.Query("DUMPALLOBJECTS", "FILE")
		return err
	}
	_, err := cli.Query("DUMPALLOBJECTS")
	return err
}

// Echo возвращает переданный объект без изменений.
func (cli *Cli) Echo(obj any) (any, error) {
	return cli.Query("ECHO", obj)
}

// GetCurrentController возвращает IP-адрес текущего контроллера динамического кластера.
func (cli *Cli) GetCurrentController() (net.IP, error) {
	res, err := cli.Query("GETCURRENTCONTROLLER")
	if err != nil {
		return nil, err
	}
	s, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("GETCURRENTCONTROLLER: expected string, got %T", res)
	}
	if s == "" {
		return nil, nil
	}
	return net.ParseIP(s), nil
}

// GetCurrentTime возвращает текущее внутреннее время сервера.
func (cli *Cli) GetCurrentTime() (time.Time, error) {
	res, err := cli.Query("GETCURRENTTIME")
	if err != nil {
		return time.Time{}, err
	}
	t, ok := res.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("GETCURRENTTIME: expected time.Time, got %T", res)
	}
	return t, nil
}

// GetMessageQueueInfo возвращает информацию об очереди модуля.
func (cli *Cli) GetMessageQueueInfo(module, queue string) (map[string]any, error) {
	if module == "" || queue == "" {
		return nil, fmt.Errorf("module and queue are required")
	}
	res, err := cli.Query("GETMESSAGEQUEUEINFO", module, "QUEUE", queue)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETMESSAGEQUEUEINFO: expected map, got %T", res)
	}
	return m, nil
}

// GetSystemInfo возвращает системную информацию сервера по ключу.
func (cli *Cli) GetSystemInfo(what SystemInfoKey) (any, error) {
	if what == "" {
		return nil, fmt.Errorf("system info key is required")
	}
	return cli.Query("GETSYSTEMINFO", string(what))
}

// GetTempBlacklistedIPs возвращает список временно заблокированных IP с временем их жизни.
func (cli *Cli) GetTempBlacklistedIPs() ([]TempIP, error) {
	res, err := cli.Query("GETTEMPBLACKLISTEDIPS")
	if err != nil {
		return nil, err
	}
	s, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("GETTEMPBLACKLISTEDIPS: expected string, got %T", res)
	}
	return cli.parseTempIPs(s), nil
}

// GetTempClientIPs возвращает список временных клиентских IP.
func (cli *Cli) GetTempClientIPs() ([]TempIP, error) {
	res, err := cli.Query("GETTEMPCLIENTIPS")
	if err != nil {
		return nil, err
	}
	s, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("GETTEMPCLIENTIPS: expected string, got %T", res)
	}
	return cli.parseTempIPs(s), nil
}

// GetTempUnblockableIPs возвращает список адресов, защищенных от блокировки.
func (cli *Cli) GetTempUnblockableIPs() ([]TempIP, error) {
	res, err := cli.Query("GETTEMPUNBLOCKABLEIPS")
	if err != nil {
		return nil, err
	}
	s, ok := res.(string)
	if !ok {
		return nil, fmt.Errorf("GETTEMPUNBLOCKABLEIPS: expected string, got %T", res)
	}
	return cli.parseTempIPs(s), nil
}

// GetVersion возвращает версию сервера.
func (cli *Cli) GetVersion() (string, error) {
	res, err := cli.Query("GETVERSION")
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("GETVERSION: expected string, got %T", res)
	}
	return s, nil
}

// ListCLICommands возвращает список всех поддерживаемых команд.
func (cli *Cli) ListCLICommands() ([]string, error) {
	return cli.getSliceString("LISTCLICOMMANDS", "")
}

// Noop — пустая команда, всегда успешна.
func (cli *Cli) Noop() error {
	_, err := cli.Query("NOOP")
	return err
}

// ReconnectClusterAdmin форсирует переподключение административных соединений кластера.
func (cli *Cli) ReconnectClusterAdmin() error {
	_, err := cli.Query("RECONNECTCLUSTERADMIN")
	return err
}

// RejectQueueMessage отклоняет сообщение из очереди.
func (cli *Cli) RejectQueueMessage(id int64, report string) error {
	args := make([]any, 0, 4)
	args = append(args, "REJECTQUEUEMESSAGE", id)
	if report != "" {
		args = append(args, "REPORT", report)
	}
	_, err := cli.Query(args...)
	return err
}

// RejectQueueMessages отклоняет все сообщения отправителя из очереди.
func (cli *Cli) RejectQueueMessages(sender, report string) error {
	if sender == "" {
		return fmt.Errorf("sender is required")
	}
	args := make([]any, 0, 5)
	args = append(args, "REJECTQUEUEMESSAGES", "SENDER", sender)
	if report != "" {
		args = append(args, "REPORT", report)
	}
	_, err := cli.Query(args...)
	return err
}

// ReleaseSMTPQueue освобождает SMTP очередь.
func (cli *Cli) ReleaseSMTPQueue(queue string) error {
	if queue == "" {
		return fmt.Errorf("queue name is required")
	}
	_, err := cli.Query("RELEASESMTPQUEUE", queue)
	return err
}

// ReportFailedLoginAddress увеличивает счетчик неудачных попыток входа для IP-адреса.
func (cli *Cli) ReportFailedLoginAddress(ip net.IP) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	_, err := cli.Query("REPORTFAILEDLOGINADDRESS", ip.String())
	return err
}

// SetLogAll включает или выключает режим "Log Everything".
func (cli *Cli) SetLogAll(on bool) error {
	val := "OFF"
	if on {
		val = "ON"
	}
	_, err := cli.Query("SETLOGALL", val)
	return err
}

// SetTempBlacklistedIPs заменяет текущий список временно заблокированных адресов.
func (cli *Cli) SetTempBlacklistedIPs(ips []net.IP) error {
	if len(ips) == 0 {
		_, err := cli.Query("SETTEMPBLACKLISTEDIPS", "")
		return err
	}

	var sb strings.Builder
	// Примерный расчет: (длина IPv4 + скобки + запятая) * кол-во
	sb.Grow(len(ips) * 18)

	for i, ip := range ips {
		if ip == nil {
			continue
		}
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('[')
		sb.WriteString(ip.String())
		sb.WriteByte(']')
	}

	_, err := cli.Query("SETTEMPBLACKLISTEDIPS", sb.String())
	return err
}

// SetTempUnblockableIPs аналогично заменяет список исключений.
func (cli *Cli) SetTempUnblockableIPs(ips []net.IP) error {
	if len(ips) == 0 {
		_, err := cli.Query("SETTEMPUNBLOCKABLEIPS", "")
		return err
	}

	var sb strings.Builder
	sb.Grow(len(ips) * 18)

	for i, ip := range ips {
		if ip == nil {
			continue
		}
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('[')
		sb.WriteString(ip.String())
		sb.WriteByte(']')
	}

	_, err := cli.Query("SETTEMPUNBLOCKABLEIPS", sb.String())
	return err
}

// SetTrace включает или выключает трассировку указанного модуля в системный лог (syslog).
func (cli *Cli) SetTrace(facility TraceFacility, on bool) error {
	if facility == "" {
		return fmt.Errorf("trace facility is required")
	}
	val := "OFF"
	if on {
		val = "ON"
	}
	_, err := cli.Query("SETTRACE", string(facility), val)
	return err
}

// TempBlacklistIP добавляет адрес в список временно заблокированных или удаляет его.
// Если timeout <= 0, адрес удаляется (DELETE).
func (cli *Cli) TempBlacklistIP(ip net.IP, timeout time.Duration) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	if timeout <= 0 {
		_, err := cli.Query("TEMPBLACKLISTIP", ip.String(), "DELETE")
		return err
	}
	seconds := int(timeout.Seconds())
	_, err := cli.Query("TEMPBLACKLISTIP", ip.String(), "TIMEOUT", seconds)
	return err
}

// TempUnblockIP добавляет адрес в список исключений или удаляет его.
// Если timeout <= 0, адрес удаляется.
func (cli *Cli) TempUnblockIP(ip net.IP, timeout time.Duration) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	if timeout <= 0 {
		_, err := cli.Query("TEMPUNBLOCKIP", ip.String(), "DELETE")
		return err
	}
	seconds := int(timeout.Seconds())
	_, err := cli.Query("TEMPUNBLOCKIP", ip.String(), "TIMEOUT", seconds)
	return err
}

// TestLoop запускает цикл тестирования производительности CPU.
func (cli *Cli) TestLoop(seconds int) (int, error) {
	res, err := cli.Query("TESTLOOP", seconds)
	if err != nil {
		return 0, err
	}
	return cli.toInt(res), nil
}

// WriteLog пишет запись в системный лог сервера.
func (cli *Cli) WriteLog(level int, record string) error {
	_, err := cli.Query("WRITELOG", level, record)
	return err
}

// хелперы

func (cli *Cli) parseTempIPs(raw string) []TempIP {
	if raw == "" {
		return nil
	}

	b := []byte(raw)

	// Считаем запятые, чтобы выделить память один раз.
	// Кол-во элементов = кол-во запятых + 1
	count := bytes.Count(b, []byte{','}) + 1
	res := make([]TempIP, 0, count)

	for i := 0; i < len(b); {
		start := bytes.IndexByte(b[i:], '[')
		if start == -1 {
			break
		}
		start += i

		end := bytes.IndexByte(b[start:], ']')
		if end == -1 {
			break
		}
		end += start

		item := TempIP{
			IP: net.ParseIP(string(b[start+1 : end])),
		}

		curr := end + 1
		if curr < len(b) && b[curr] == '-' {
			curr++
			if curr < len(b) && b[curr] == '*' {
				item.Permanent = true
				curr++
			} else {
				val := 0
				foundDigit := false
				for curr < len(b) && b[curr] >= '0' && b[curr] <= '9' {
					val = val*10 + int(b[curr]-'0')
					curr++
					foundDigit = true
				}
				if foundDigit {
					item.TimeLeft = time.Duration(val) * time.Second
				}
			}
		}

		res = append(res, item)

		nextComma := bytes.IndexByte(b[curr:], ',')
		if nextComma == -1 {
			break
		}
		i = curr + nextComma + 1
	}

	return res
}
