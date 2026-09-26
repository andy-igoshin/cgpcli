// # Miscellaneous Commands
//
// The Miscellaneous section provides a diverse set of administrative methods for
// system diagnostics, real-time monitoring, and low-level server control.
//
// Key capabilities include:
//   - System Inspection: retrieving server version, OS details, start time,
//     and cluster instances via [Cli.GetSystemInfo] and predefined [SystemInfoKey] constants.
//   - Dynamic IP Protection: managing temporary blacklists and unblockable IP lists
//     to mitigate brute-force attacks in real-time.
//   - Queue Management: monitoring and managing the message delivery pipeline,
//     including the ability to reject or release messages from SMTP queues.
//   - Diagnostics & Logging: controlling trace facilities (FileIO/FileOp),
//     toggling "Log Everything" mode, and writing custom entries to the system log.
//   - Cluster Operations: force-reconnecting cluster administration channels and
//     identifying the current Dynamic Cluster Controller.
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
//   - file: if true, the objects_dump.txt file is created in the Server base directory.
//     If false, the list is written into the OS syslog.
//
// Note: This list may contain millions of objects and can overload the OS syslog.
// It also effectively suspends CommuniGate Pro Server activities till all objects are listed.
//
// This method executes the DUMPALLOBJECTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DumpAllObjects(file bool) error {
	const cmd = "DUMPALLOBJECTS"
	if file {
		return cli.QueryNV(cmd, "FILE")
	} else {
		return cli.QueryNV(cmd)
	}
}

// Echo produces an output which is a copy of the command parameter object.
//
// Parameters:
//   - obj: the object (string, dictionary, array, etc.) to be echoed back by the server.
//
// This method executes the ECHO CLI command.
//
// Returns:
//   - any: the exact copy of the input object as returned by the server.
//   - error: an error if the command fails.
func (cli *Cli) Echo(obj any) (any, error) {
	return cli.Query("ECHO", obj)
}

// GetCurrentController retrieves the IP address of the current Dynamic Cluster Controller.
//
// This method executes the GETCURRENTCONTROLLER CLI command.
//
// Returns:
//   - net.IP: the Cluster Controller IP Address.
//   - error: an error if the command fails.
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
// This method executes the GETCURRENTTIME CLI command.
//
// Returns:
//   - time.Time: a timestamp with the Server internal timer value.
//   - error: an error if the command fails.
func (cli *Cli) GetCurrentTime() (time.Time, error) {
	const cmd = "GETCURRENTTIME"
	res, err := cli.Query(cmd)
	if err != nil {
		return time.Time{}, err
	}
	t, ok := res.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("%s: %w: expected time.Time, got %T", cmd, ErrUnexpectedType, res)
	}
	return t, nil
}

// GetMessageQueueInfo reads information about a module message Queue.
//
// Parameters:
//   - module: the module name.
//   - queue: the module queue name.
//
// This method executes the GETMESSAGEQUEUEINFO CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the specified queue information (nTotal, size, delayedTill, lastError).
//   - error: an error if the command fails.
func (cli *Cli) GetMessageQueueInfo(module, queue string) (map[string]any, error) {
	if module == "" {
		return nil, ErrModuleRequired
	}
	if queue == "" {
		return nil, ErrQueueRequired
	}
	return cli.getMapAny("GETMESSAGEQUEUEINFO", module, "QUEUE", queue)
}

// GetSystemInfo retrieves information returned by the CG/PL SystemInfo function called with the what parameter.
//
// Parameters:
//   - what: the system information key to retrieve.
//
// This method executes the GETSYSTEMINFO CLI command.
//
// Returns:
//   - any: an object returned by the SystemInfo function.
//   - error: an error if the command fails.
func (cli *Cli) GetSystemInfo(what SystemInfoKey) (any, error) {
	if what == "" {
		return nil, ErrSystemInfoKeyRequired
	}
	return cli.Query("GETSYSTEMINFO", string(what))
}

// GetTempBlacklistedIPs retrieves the set of Temporarily Blocked Addresses.
//
// This method executes the GETTEMPBLACKLISTEDIPS CLI command.
//
// Returns:
//   - *[TempIPList]: a list of Temporary Blocked IP addresses. Each address may have a suffix indicating the block duration.
//   - error: an error if the command fails.
func (cli *Cli) GetTempBlacklistedIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPBLACKLISTEDIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetTempClientIPs retrieves the set of temporary Client IP Addresses.
//
// This method executes the GETTEMPCLIENTIPS CLI command.
//
// Returns:
//   - *[TempIPList]: a list of Temporary Client IP addresses.
//   - error: an error if the command fails.
func (cli *Cli) GetTempClientIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPCLIENTIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetTempUnblockableIPs retrieves the set of Temporary UnBlockable IP Addresses.
//
// This method executes the GETTEMPUNBLOCKABLEIPS CLI command.
//
// Returns:
//   - *[TempIPList]: a list of Temporary UnBlockable IP addresses.
//   - error: an error if the command fails.
func (cli *Cli) GetTempUnblockableIPs() (*TempIPList, error) {
	raw, err := cli.getString("GETTEMPUNBLOCKABLEIPS")
	if err != nil {
		return nil, err
	}
	return &TempIPList{items: parseTempIPsToMap(raw)}, nil
}

// GetVersion retrieves the current CommuniGate Pro Server version.
//
// This method executes the GETVERSION CLI command.
//
// Returns:
//   - string: a string with the server version.
//   - error: an error if the command fails.
func (cli *Cli) GetVersion() (string, error) {
	return cli.getString("GETVERSION")
}

// ListCLICommands retrieves the list of all CLI commands supported by this version of CommuniGate Pro Server.
//
// This method executes the LISTCLICOMMANDS CLI command.
//
// Returns:
//   - []string: an array of strings, where each string is a supported command name.
//   - error: an error if the command fails.
func (cli *Cli) ListCLICommands() ([]string, error) {
	return cli.getSliceString("LISTCLICOMMANDS")
}

// Noop is a command that always completes successfully.
//
// This method executes the NOOP CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) Noop() error {
	return cli.QueryNV("NOOP")
}

// ReconnectClusterAdmin forces a Dynamic Cluster member to re-open all its inter-cluster Administrative connections.
//
// This method executes the RECONNECTCLUSTERADMIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ReconnectClusterAdmin() error {
	return cli.QueryNV("RECONNECTCLUSTERADMIN")
}

// RejectQueueMessage rejects a message from the Server Queue.
//
// Parameters:
//   - id: the message ID.
//   - report: an optional text to be included in the bounce report.
//     If the parameter is "NONDN", no DSN report message is generated.
//
// This method executes the REJECTQUEUEMESSAGE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RejectQueueMessage(id int64, report string) error {
	if report != "" {
		return cli.QueryNV("REJECTQUEUEMESSAGE", id, "REPORT", report)
	} else {
		return cli.QueryNV("REJECTQUEUEMESSAGE", id)
	}
}

// RejectQueueMessages rejects all messages sent by the specified sender from the Server Queue.
//
// Parameters:
//   - sender: the authenticated sender's name.
//   - report: an optional text to be included in the bounce report.
//     If the parameter is "NONDN", no DSN report message is generated.
//
// Note: In a Dynamic Cluster environment this command rejects messages from all server queues.
//
// This method executes the REJECTQUEUEMESSAGES CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RejectQueueMessages(sender, report string) error {
	if sender == "" {
		return ErrSenderRequired
	}

	if report != "" {
		return cli.QueryNV("REJECTQUEUEMESSAGES", "SENDER", sender, "REPORT", report)
	} else {
		return cli.QueryNV("REJECTQUEUEMESSAGES", "SENDER", sender)
	}
}

// ReleaseSMTPQueue releases an SMTP queue.
//
// Parameters:
//   - queue: the queue (domain) name to release.
//
// Note: In a Dynamic Cluster environment this command releases the specified SMTP queue on all servers.
//
// This method executes the RELEASESMTPQUEUE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ReleaseSMTPQueue(queue string) error {
	if queue == "" {
		return ErrQueueRequired
	}
	return cli.QueryNV("RELEASESMTPQUEUE", queue)
}

// ReportFailedLoginAddress increments the counter of failed Login attempts from the specified IP address.
//
// Parameters:
//   - ip: the Network IP Address to report.
//
// This method executes the REPORTFAILEDLOGINADDRESS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ReportFailedLoginAddress(ip net.IP) error {
	if ip == nil {
		return ErrIPAddressRequired
	}
	return cli.QueryNV("REPORTFAILEDLOGINADDRESS", ip.String())
}

// SetLogAll switches on and off the "Log Everything" mode.
//
// Parameters:
//   - on: boolean to switch the mode ON or OFF.
//
// This method executes the SETLOGALL CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
//   - list: a list of IP addresses in the format of the GetTempBlacklistedIPs command.
//
// This method executes the SETTEMPBLACKLISTEDIPS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetTempBlacklistedIPs(list *TempIPList) error {
	if list == nil || list.Len() == 0 {
		return cli.QueryNV("SETTEMPBLACKLISTEDIPS", "")
	}
	return cli.QueryNV("SETTEMPBLACKLISTEDIPS", list)
}

// SetTempUnblockableIPs adds addresses to the Temporary UnBlockable IP Addresses set.
//
// Parameters:
//   - list: a list of IP addresses in the format of the GetTempUnBlockableIPs command.
//
// This method executes the SETTEMPBLACKLISTEDIPS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetTempUnblockableIPs(list *TempIPList) error {
	if list == nil || list.Len() == 0 {
		return cli.QueryNV("SETTEMPBLACKLISTEDIPS", "")
	}
	return cli.QueryNV("SETTEMPUNBLOCKABLEIPS", list)
}

// SetTrace switches on and off internal logging facilities that write to OS syslog.
//
// Parameters:
//   - facility: the facility name.
//   - on: boolean to switch the facility ON or OFF.
//
// This method executes the SETTRACE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetTrace(facility TraceFacility, on bool) error {
	if facility == "" {
		return ErrTraceFacilityRequired
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
//   - ip: the Network IP Address to add or remove.
//   - timeout: the time period the address should be blocked for. Use 0 or DELETE to remove.
//
// This method executes the TEMPBLACKLISTIP CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) TempBlacklistIP(ip net.IP, timeout time.Duration) error {
	if ip == nil {
		return ErrIPAddressRequired
	}

	if timeout <= 0 {
		return cli.QueryNV("TEMPBLACKLISTIP", ip.String(), "DELETE")
	} else {
		return cli.QueryNV("TEMPBLACKLISTIP", ip.String(), "TIMEOUT", int(timeout.Seconds()))
	}
}

// TempUnblockIP adds or removes an address from the Temporary UnBlockable IP Addresses set.
//
// Parameters:
//   - ip: the Network IP Address to add or remove.
//   - timeout: the time period the address should be in the set. Use 0 or DELETE to remove.
//
// This method executes the TEMPUNBLOCKIP CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) TempUnblockIP(ip net.IP, timeout time.Duration) error {
	if ip == nil {
		return ErrIPAddressRequired
	}

	if timeout <= 0 {
		return cli.QueryNV("TEMPUNBLOCKIP", ip.String(), "DELETE")
	} else {
		return cli.QueryNV("TEMPUNBLOCKIP", ip.String(), "TIMEOUT", int(timeout.Seconds()))
	}
}

// TestLoop tests the server CPU load by executing a calculation loop for the specified number of seconds.
//
// Parameters:
//   - seconds: duration of the test in seconds.
//
// This method executes the TESTLOOP CLI command.
//
// Returns:
//   - int: a number that indicates the average CLI thread CPU performance.
//   - error: an error if the command fails.
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
//   - level: the record log level.
//   - record: the string to be placed into the Server Log.
//
// This method executes the WRITELOG CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) WriteLog(level int, record string) error {
	return cli.QueryNV("WRITELOG", level, record)
}
