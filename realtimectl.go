// # Real-Time Application Control
//
// The RTA Control section provides methods for managing the execution of
// PBX (Real-Time Application) tasks. While the PBX section handles file
// management, this section focuses on the runtime environment and
// inter-process communication within the CommuniGate Pro Signal engine.
//
// Key capabilities include:
//   - Task Lifecycle: starting new PBX tasks for specific accounts using
//     [Cli.StartPBXTask] with custom entry points and initial parameters.
//   - Process Control: monitoring the current status of running tasks and
//     terminating them (killing nodes) when necessary.
//   - Event Management: sending asynchronous events to active tasks via
//     [Cli.SendTaskEvent], enabling interaction between the CLI and
//     running signal applications.
//   - Task Identification: each task is managed via a unique TaskID returned
//     at startup, allowing precise control over concurrent applications.
package cgpcli

import (
	"fmt"
)

// KillNode kills an existing PBX Task.
//
// Parameters:
//   - taskID: the unique Task ID.
//
// This method executes the KILLNODE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) KillNode(taskID string) error {
	if taskID == "" {
		return fmt.Errorf("taskID is required")
	}
	return cli.QueryNV("KILLNODE", taskID)
}

// ReadNodeStatus reads the current application status of an existing PBX Task.
//
// Parameters:
//   - taskID: the Task ID.
//
// This method executes the READNODESTATUS CLI command.
//
// Returns:
//   - any: the application status object returned by the task.
//   - error: an error if the command fails.
func (cli *Cli) ReadNodeStatus(taskID string) (any, error) {
	if taskID == "" {
		return nil, fmt.Errorf("taskID is required")
	}
	return cli.Query("READNODESTATUS", taskID)
}

// SendTaskEvent sends an Event to an existing PBX Task.
//
// Parameters:
//   - taskID: the Task ID.
//   - event: the name of the Event to send.
//   - param: an optional event parameter.
//
// This method executes the SENDTASKEVENT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SendTaskEvent(taskID, event string, param any) error {
	if taskID == "" || event == "" {
		return fmt.Errorf("taskID and event name are required")
	}

	const cmd = "SENDTASKEVENT"
	if param != nil {
		return cli.QueryNV(cmd, taskID, "EVENT", event, "PARAM", param)
	} else {
		return cli.QueryNV(cmd, taskID, "EVENT", event)
	}
}

// StartPBXTask starts a new PBX Task on behalf of the specified Account.
//
// Parameters:
//   - account: the name of the Account. The name can include the Domain name.
//     If the Domain name is not specified, the current user Domain is used.
//   - program: the name of the program (the .sppr file) to start.
//   - entry: an optional program entry point. If empty, "main" is used.
//   - param: an optional specifies the program parameter.
//     Accessible within the program via Vars().startParameter.
//
// This method executes the PROGRAM CLI command.
//
// Returns:
//   - string: a string with the new Task ID.
//   - error: an error if the command fails.
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
