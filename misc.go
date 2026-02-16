// # Miscellaneous Commands
//
// The Miscellaneous section provides a diverse set of administrative methods for
// system diagnostics, real-time monitoring, and low-level server control.
//
// Key capabilities include:
//   - System Inspection: Retrieving server version, OS details, start time,
//     and cluster instances via [Cli.GetSystemInfo] and predefined [SystemInfoKey] constants.
//   - Dynamic IP Protection: Managing temporary blacklists and unblockable IP lists
//     to mitigate brute-force attacks in real-time.
//   - Queue Management: Monitoring and managing the message delivery pipeline,
//     including the ability to reject or release messages from SMTP queues.
//   - Diagnostics & Logging: Controlling trace facilities (FileIO/FileOp),
//     toggling "Log Everything" mode, and writing custom entries to the system log.
//   - Cluster Operations: Force-reconnecting cluster administration channels and
//     identifying the current Dynamic Cluster Controller.
//
// These methods are essential for building health-check monitors, automated
// security responders, and advanced troubleshooting scripts.
package cgpcli

import (
	"fmt"
	"net"
	"time"
)

// SystemInfoKey represents a recognized key for server metadata retrieval.
type SystemInfoKey string

const (
	SysVersion         SystemInfoKey = "serverVersion" // e.g., "6.3.20"
	SysOS              SystemInfoKey = "serverOS"      // Target Operating System
	SysCPU             SystemInfoKey = "serverCPU"     // CPU Architecture
	SysMainDomain      SystemInfoKey = "mainDomainName"
	SysLicenseDomain   SystemInfoKey = "licenseDomainName"
	SysStartTime       SystemInfoKey = "startTime" // Server boot time
	SysServerInstance  SystemInfoKey = "serverInstance"
	SysClusterInstance SystemInfoKey = "clusterInstance"
	SysCharsets        SystemInfoKey = "knownCharsets"
	SysDigesters       SystemInfoKey = "knownDigesters"
	SysCiphers         SystemInfoKey = "knownCiphers" // Supported SSL/TLS ciphers
)

// TraceFacility defines specific server modules available for real-time tracing.
type TraceFacility string

const (
	TraceFileIO TraceFacility = "FileIO" // Record all file read/write/truncate operations
	TraceFileOp TraceFacility = "FileOp" // Record all file create/rename/remove operations
)

// DumpAllObjects writes the list of all application data objects.
//
// Parameters:
//   - file: If true, the objects_dump.txt file is created in the Server base directory.
//     If false, the list is written into the OS syslog.
//
// Note: This list may contain millions of objects and can overload the OS syslog.
// It also effectively suspends CommuniGate Pro Server activities till all objects are listed.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DumpAllObjects(file bool) (err error) {
	const cmd = "DUMPALLOBJECTS"
	if file {
		err = cli.QueryNV(cmd, "FILE")
	} else {
		err = cli.QueryNV(cmd)
	}
	return
}

// Echo produces an output which is a copy of the command parameter object.
//
// Parameters:
//   - obj: The object (string, dictionary, array, etc.) to be echoed back by the server.
//
// Returns:
//   - any: The exact copy of the input object as returned by the server.
//   - error: An error if the command fails.
func (cli *Cli) Echo(obj any) (any, error) {
	return cli.Query("ECHO", obj)
}

// GetCurrentController retrieves the IP address of the current Dynamic Cluster Controller.
//
// Returns:
//   - net.IP: The Cluster Controller IP Address.
//   - error: An error if the command fails.
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

// GetCurrentTime retrieves the current internal timer value of the CommuniGate Pro Server.
//
// Returns:
//   - time.Time: A timestamp with the Server internal timer value.
//   - error: An error if the command fails.
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

// GetMessageQueueInfo reads information about a module message Queue.
//
// Parameters:
//   - module: This parameter specifies the module name.
//   - queue: This parameter specifies the module queue name.
//
// Returns:
//   - map[string]any: A dictionary with the specified queue information (nTotal, size, delayedTill, lastError).
//   - error: An error if the command fails.
func (cli *Cli) GetMessageQueueInfo(module, queue string) (map[string]any, error) {
	if module == "" || queue == "" {
		return nil, fmt.Errorf("module and queue are required")
	}
	return cli.getMapAny("GETMESSAGEQUEUEINFO", module, "QUEUE", queue)
}

// GetSystemInfo retrieves information returned by the CG/PL SystemInfo function called with the what parameter.
//
// Parameters:
//   - what: This parameter specifies the system information key to retrieve.
//
// Returns:
//   - any: An object returned by the SystemInfo function.
//   - error: An error if the function returns a null-object.
func (cli *Cli) GetSystemInfo(what SystemInfoKey) (any, error) {
	if what == "" {
		return nil, fmt.Errorf("system info key is required")
	}
	return cli.Query("GETSYSTEMINFO", string(what))
}

// GetTempBlacklistedIPs retrieves the set of Temporarily Blocked Addresses.
//
// Returns:
//   - *[TempIPList]: A list of Temporary Blocked IP addresses. Each address may have a suffix indicating the block duration.
//   - error: An error if the command fails.
func (cli *Cli) GetTempBlacklistedIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPBLACKLISTEDIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetTempClientIPs retrieves the set of temporary Client IP Addresses.
//
// Returns:
//   - *[TempIPList]: A list of Temporary Client IP addresses.
//   - error: An error if the command fails.
func (cli *Cli) GetTempClientIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPCLIENTIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetTempUnblockableIPs retrieves the set of Temporary UnBlockable IP Addresses.
//
// Returns:
//   - *[TempIPList]: A list of Temporary UnBlockable IP addresses.
//   - error: An error if the command fails.
func (cli *Cli) GetTempUnblockableIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPUNBLOCKABLEIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetVersion retrieves the current CommuniGate Pro Server version.
//
// Returns:
//   - string: A string with the server version.
//   - error: An error if the command fails.
func (cli *Cli) GetVersion() (string, error) {
	return cli.getString("GETVERSION")
}

// ListCLICommands retrieves the list of all CLI commands supported by this version of CommuniGate Pro Server.
//
// Returns:
//   - []string: An array of strings, where each string is a supported command name.
//   - error: An error if the command fails.
func (cli *Cli) ListCLICommands() ([]string, error) {
	return cli.getSliceString("LISTCLICOMMANDS")
}

// Noop is a command that always completes successfully.
//
// Returns:
//   - error: An error if the connection fails.
func (cli *Cli) Noop() error {
	return cli.QueryNV("NOOP")
}

// ReconnectClusterAdmin forces a Dynamic Cluster member to re-open all its inter-cluster Administrative connections.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) ReconnectClusterAdmin() error {
	return cli.QueryNV("RECONNECTCLUSTERADMIN")
}

// RejectQueueMessage rejects a message from the Server Queue.
//
// Parameters:
//   - id: This parameter specifies the message ID.
//   - report: (Optional) This parameter specifies the text to be included into the bounce report.
//     If the parameter is "NONDN", no DSN report message is generated.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) RejectQueueMessage(id int64, report string) (err error) {
	if report != "" {
		err = cli.QueryNV("REJECTQUEUEMESSAGE", id, "REPORT", report)
	} else {
		err = cli.QueryNV("REJECTQUEUEMESSAGE", id)
	}
	return
}

// RejectQueueMessages rejects all messages sent by the specified sender from the Server Queue.
//
// Parameters:
//   - sender: This parameter specifies the authenticated sender's name.
//   - report: (Optional) This parameter specifies the text to be included into the bounce report.
//     If the parameter is "NONDN", no DSN report message is generated.
//
// Note: In a Dynamic Cluster environment this command rejects messages from all server queues.
//
// Returns:
//   - error: An error if the command fails.
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

// ReleaseSMTPQueue releases an SMTP queue.
//
// Parameters:
//   - queue: This parameter specifies the queue (domain) name to release.
//
// Note: In a Dynamic Cluster environment this command releases the specified SMTP queue on all servers.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) ReleaseSMTPQueue(queue string) error {
	if queue == "" {
		return fmt.Errorf("queue name is required")
	}
	return cli.QueryNV("RELEASESMTPQUEUE", queue)
}

// ReportFailedLoginAddress increments the counter of failed Login attempts from the specified IP address.
//
// Parameters:
//   - ip: The Network IP Address to report.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) ReportFailedLoginAddress(ip net.IP) error {
	if ip == nil {
		return fmt.Errorf("IP address is required")
	}
	return cli.QueryNV("REPORTFAILEDLOGINADDRESS", ip.String())
}

// SetLogAll switches on and off the "Log Everything" mode.
//
// Parameters:
//   - on: Boolean to switch the mode ON or OFF.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetLogAll(on bool) error {
	val := "OFF"
	if on {
		val = "ON"
	}
	return cli.QueryNV("SETLOGALL", val)
}

// SetTempBlacklistedIPs adds addresses to the Temporary Blocked IP Addresses list.
//
// Parameters:
//   - list: A list of IP addresses in the format of the GetTempBlacklistedIPs command.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetTempBlacklistedIPs(list *TempIPList) error {
	if list == nil || list.Len() == 0 {
		return cli.QueryNV("SETTEMPBLACKLISTEDIPS", "")
	}
	return cli.QueryNV("SETTEMPBLACKLISTEDIPS", list)
}

// SetTempUnblockableIPs adds addresses to the Temporary UnBlockable IP Addresses set.
//
// Parameters:
//   - list: A list of IP addresses in the format of the GetTempUnBlockableIPs command.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) SetTempUnblockableIPs(list *TempIPList) error {
	if list == nil || list.Len() == 0 {
		return cli.QueryNV("SETTEMPBLACKLISTEDIPS", "")
	}
	return cli.QueryNV("SETTEMPUNBLOCKABLEIPS", list)
}

// SetTrace switches on and off internal logging facilities that write to OS syslog.
//
// Parameters:
//   - facility: The facility name.
//   - on: Boolean to switch the facility ON or OFF.
//
// Returns:
//   - error: An error if the command fails.
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

// TempBlacklistIP adds or removes an address from the Temporarily Blocked Addresses set.
//
// Parameters:
//   - ip: The Network IP Address to add or remove.
//   - timeout: The time period the address should be blocked for. Use 0 or DELETE to remove.
//
// Returns:
//   - error: An error if the command fails.
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

// TempUnblockIP adds or removes an address from the Temporary UnBlockable IP Addresses set.
//
// Parameters:
//   - ip: The Network IP Address to add or remove.
//   - timeout: The time period the address should be in the set. Use 0 or DELETE to remove.
//
// Returns:
//   - error: An error if the command fails.
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

// TestLoop tests the server CPU load by executing a calculation loop for the specified number of seconds.
//
// Parameters:
//   - seconds: Duration of the test in seconds.
//
// Returns:
//   - int: A number that indicates the average CLI thread CPU performance.
//   - error: An error if the command fails.
func (cli *Cli) TestLoop(seconds int) (int, error) {
	res, err := cli.Query("TESTLOOP", seconds)
	if err != nil {
		return 0, err
	}
	return toInt(res), nil
}

// WriteLog stores a record into the Server Log. Log records have the SYSTEM prefix.
//
// Parameters:
//   - level: This parameter specifies the record log level.
//   - record: This parameter specifies the string to be placed into the Server Log.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) WriteLog(level int, record string) error {
	return cli.QueryNV("WRITELOG", level, record)
}
