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
	if report != "" {
		_, err := cli.Query("REJECTQUEUEMESSAGE", id, "REPORT", report)
		return err
	}
	_, err := cli.Query("REJECTQUEUEMESSAGE", id)
	return err
}

// RejectQueueMessages отклоняет все сообщения отправителя из очереди.
func (cli *Cli) RejectQueueMessages(sender, report string) error {
	if sender == "" {
		return fmt.Errorf("sender is required")
	}
	if report != "" {
		_, err := cli.Query("REJECTQUEUEMESSAGES", "SENDER", sender, "REPORT", report)
		return err
	}
	_, err := cli.Query("REJECTQUEUEMESSAGES", "SENDER", sender)
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
	// CGP ожидает список в формате [ip1],[ip2]
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

// SetTempUnblockableIPs аналогично.
func (cli *Cli) SetTempUnblockableIPs(ips []net.IP) error {
	if len(ips) == 0 {
		_, err := cli.Query("SETTEMPUNBLOCKABLEIPS", "")
		return err
	}

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

	_, err := cli.Query("SETTEMPUNBLOCKABLEIPS", sb.String())
	return err
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
	_, err := cli.Query("SETTRACE", string(facility), val)
	return err
}

// TempBlacklistIP добавляет адрес в список или удаляет его.
func (cli *Cli) TempBlacklistIP(ip net.IP, timeout time.Duration) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	if timeout <= 0 {
		_, err := cli.Query("TEMPBLACKLISTIP", ip.String(), "DELETE")
		return err
	}
	_, err := cli.Query("TEMPBLACKLISTIP", ip.String(), "TIMEOUT", int(timeout.Seconds()))
	return err
}

// TempUnblockIP добавляет адрес в список исключений или удаляет его.
func (cli *Cli) TempUnblockIP(ip net.IP, timeout time.Duration) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	if timeout <= 0 {
		_, err := cli.Query("TEMPUNBLOCKIP", ip.String(), "DELETE")
		return err
	}
	_, err := cli.Query("TEMPUNBLOCKIP", ip.String(), "TIMEOUT", int(timeout.Seconds()))
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
