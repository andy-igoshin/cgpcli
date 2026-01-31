// # Section: Miscellaneous Commands
package cgpcli

import (
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

// DumpAllObjects пишет список всех объектов данных.
func (cli *Cli) DumpAllObjects(file bool) (err error) {
	const cmd = "DUMPALLOBJECTS"
	if file {
		err = cli.QueryNV(cmd, "FILE")
	} else {
		err = cli.QueryNV(cmd)
	}
	return
}

// Echo возвращает переданный объект без изменений.
func (cli *Cli) Echo(obj any) (any, error) {
	return cli.Query("ECHO", obj)
}

// GetCurrentController возвращает IP-адрес текущего контроллера динамического кластера.
func (cli *Cli) GetCurrentController() (net.IP, error) {
	s, err := cli.getString("GETCURRENTCONTROLLER")
	if err != nil {
		return nil, err
	}
	if s == "" {
		return nil, nil
	}
	return net.ParseIP(s), nil
}

// GetCurrentTime возвращает текущее внутреннее время сервера.
func (cli *Cli) GetCurrentTime() (time.Time, error) {
	const cmd = "GETCURRENTTIME"
	res, err := cli.Query(cmd)
	if err != nil {
		return time.Time{}, err
	}
	t, ok := res.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("%s: expected time.Time, got %T", cmd, res)
	}
	return t, nil
}

// GetMessageQueueInfo возвращает информацию об очереди модуля.
func (cli *Cli) GetMessageQueueInfo(module, queue string) (map[string]any, error) {
	if module == "" || queue == "" {
		return nil, fmt.Errorf("module and queue are required")
	}
	const cmd = "GETMESSAGEQUEUEINFO"
	res, err := cli.Query(cmd, module, "QUEUE", queue)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
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

// GetTempBlacklistedIPs возвращает список временно заблокированных IP.
func (cli *Cli) GetTempBlacklistedIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPBLACKLISTEDIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetTempClientIPs возвращает список временных клиентских IP.
func (cli *Cli) GetTempClientIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPCLIENTIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetTempUnblockableIPs возвращает список адресов, защищенных от блокировки.
func (cli *Cli) GetTempUnblockableIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPUNBLOCKABLEIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetVersion возвращает версию сервера.
func (cli *Cli) GetVersion() (string, error) {
	return cli.getString("GETVERSION")
}

// ListCLICommands возвращает список всех поддерживаемых команд.
func (cli *Cli) ListCLICommands() ([]string, error) {
	return cli.getSliceString("LISTCLICOMMANDS", "")
}

// Noop — пустая команда, всегда успешна.
func (cli *Cli) Noop() error {
	return cli.QueryNV("NOOP")
}

// ReconnectClusterAdmin форсирует переподключение административных соединений кластера.
func (cli *Cli) ReconnectClusterAdmin() error {
	return cli.QueryNV("RECONNECTCLUSTERADMIN")
}

// RejectQueueMessage отклоняет сообщение из очереди.
func (cli *Cli) RejectQueueMessage(id int64, report string) (err error) {
	if report != "" {
		err = cli.QueryNV("REJECTQUEUEMESSAGE", id, "REPORT", report)
	} else {
		err = cli.QueryNV("REJECTQUEUEMESSAGE", id)
	}
	return
}

// RejectQueueMessages отклоняет все сообщения отправителя из очереди.
func (cli *Cli) RejectQueueMessages(sender, report string) (err error) {
	if sender == "" {
		return fmt.Errorf("sender is required")
	}
	if report != "" {
		err = cli.QueryNV("REJECTQUEUEMESSAGES", "SENDER", sender, "REPORT", report)
	} else {
		err = cli.QueryNV("REJECTQUEUEMESSAGES", "SENDER", sender)
	}
	return
}

// ReleaseSMTPQueue освобождает SMTP очередь.
func (cli *Cli) ReleaseSMTPQueue(queue string) error {
	if queue == "" {
		return fmt.Errorf("queue name is required")
	}
	return cli.QueryNV("RELEASESMTPQUEUE", queue)
}

// ReportFailedLoginAddress увеличивает счетчик неудачных попыток входа для IP-адреса.
func (cli *Cli) ReportFailedLoginAddress(ip net.IP) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	return cli.QueryNV("REPORTFAILEDLOGINADDRESS", ip.String())
}

// SetLogAll включает или выключает режим "Log Everything".
func (cli *Cli) SetLogAll(on bool) error {
	val := "OFF"
	if on {
		val = "ON"
	}
	return cli.QueryNV("SETLOGALL", val)
}

// SetTempBlacklistedIPs заменяет текущий список временно заблокированных адресов.
func (cli *Cli) SetTempBlacklistedIPs(ips []net.IP) (err error) {
	return cli.QueryNV("SETTEMPBLACKLISTEDIPS", toIPsString(ips))
}

// SetTempUnblockableIPs заменяет текущий список временно не блокируемых адресов.
func (cli *Cli) SetTempUnblockableIPs(ips []net.IP) (err error) {
	return cli.QueryNV("SETTEMPUNBLOCKABLEIPS", toIPsString(ips))
}

// SetTrace включает или выключает трассировку указанного модуля.
func (cli *Cli) SetTrace(facility TraceFacility, on bool) error {
	if facility == "" {
		return fmt.Errorf("trace facility is required")
	}
	val := "OFF"
	if on {
		val = "ON"
	}
	return cli.QueryNV("SETTRACE", string(facility), val)
}

// TempBlacklistIP добавляет адрес в список или удаляет его.
func (cli *Cli) TempBlacklistIP(ip net.IP, timeout time.Duration) (err error) {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	if timeout <= 0 {
		err = cli.QueryNV("TEMPBLACKLISTIP", ip.String(), "DELETE")
	} else {
		err = cli.QueryNV("TEMPBLACKLISTIP", ip.String(), "TIMEOUT", int(timeout.Seconds()))
	}
	return
}

// TempUnblockIP добавляет адрес в список исключений или удаляет его.
func (cli *Cli) TempUnblockIP(ip net.IP, timeout time.Duration) (err error) {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	if timeout <= 0 {
		err = cli.QueryNV("TEMPUNBLOCKIP", ip.String(), "DELETE")
	} else {
		err = cli.QueryNV("TEMPUNBLOCKIP", ip.String(), "TIMEOUT", int(timeout.Seconds()))
	}
	return
}

// TestLoop запускает цикл тестирования производительности CPU.
func (cli *Cli) TestLoop(seconds int) (int, error) {
	res, err := cli.Query("TESTLOOP", seconds)
	if err != nil {
		return 0, err
	}
	return toInt(res), nil
}

// WriteLog пишет запись в системный лог сервера.
func (cli *Cli) WriteLog(level int, record string) error {
	return cli.QueryNV("WRITELOG", level, record)
}

func toIPsString(ips []net.IP) string {
	if len(ips) > 0 {
		var sb strings.Builder
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
		return sb.String()
	}
	return ""
}
