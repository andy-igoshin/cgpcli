// # Real-Time Application Control
//
// The RTA Control section provides methods for managing the execution of
// PBX (Real-Time Application) tasks. While the PBX section handles file
// management, this section focuses on the runtime environment and
// inter-process communication within the CommuniGate Pro Signal engine.
//
// Key capabilities include:
//   - Task Lifecycle: Starting new PBX tasks for specific accounts using
//     [Cli.StartPBXTask] with custom entry points and initial parameters.
//   - Process Control: Monitoring the current status of running tasks and
//     terminating them (killing nodes) when necessary.
//   - Event Management: Sending asynchronous events to active tasks via
//     [Cli.SendTaskEvent], enabling interaction between the CLI and
//     running signal applications.
//   - Task Identification: Each task is managed via a unique TaskID returned
//     at startup, allowing precise control over concurrent applications.
//
// This section is vital for building complex telephony controllers,
// external dialers, and real-time monitoring tools for VoIP services.
package cgpcli

import (
	"fmt"
)

// KillNode принудительно завершает выполнение PBX-задачи.
func (cli *Cli) KillNode(taskID string) error {
	if taskID == "" {
		return fmt.Errorf("taskID is required")
	}
	return cli.QueryNV("KILLNODE", taskID)
}

// ReadNodeStatus возвращает текущий статус приложения для указанной задачи.
func (cli *Cli) ReadNodeStatus(taskID string) (any, error) {
	if taskID == "" {
		return nil, fmt.Errorf("taskID is required")
	}
	return cli.Query("READNODESTATUS", taskID)
}

// SendTaskEvent отправляет событие в существующую PBX-задачу.
func (cli *Cli) SendTaskEvent(taskID, event string, param any) (err error) {
	if taskID == "" || event == "" {
		return fmt.Errorf("taskID and event name are required")
	}

	const cmd = "SENDTASKEVENT"
	if param != nil {
		err = cli.QueryNV(cmd, taskID, "EVENT", event, "PARAM", param)
	} else {
		err = cli.QueryNV(cmd, taskID, "EVENT", event)
	}
	return
}

// StartPBXTask запускает новую PBX-задачу для указанного аккаунта.
func (cli *Cli) StartPBXTask(account, program, entry string, param any) (string, error) {
	if account == "" || program == "" {
		return "", fmt.Errorf("account and program name are required")
	}

	args := make([]any, 0, 7)
	args = append(args, account, "PROGRAM", program)
	if entry != "" {
		args = append(args, "ENTRY", entry)
	}
	if param != nil {
		args = append(args, "PARAM", param)
	}
	return cli.getString("STARTPBXTASK", args...)
}
