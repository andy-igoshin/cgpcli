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

// KillNode kills an existing PBX Task.
//
// Parameters:
//   - taskID: This parameter specifies the unique Task ID.
//
// Returns:
//   - error: An error if the task cannot be found or the command fails.
func (cli *Cli) KillNode(taskID string) error {
	if taskID == "" {
		return fmt.Errorf("taskID is required")
	}
	return cli.QueryNV("KILLNODE", taskID)
}

// ReadNodeStatus reads the current application status of an existing PBX Task.
//
// Parameters:
//   - taskID: This parameter specifies the Task ID.
//
// Returns:
//   - any: The application status object returned by the task.
//   - error: An error if the command fails.
func (cli *Cli) ReadNodeStatus(taskID string) (any, error) {
	if taskID == "" {
		return nil, fmt.Errorf("taskID is required")
	}
	return cli.Query("READNODESTATUS", taskID)
}

// SendTaskEvent sends an Event to an existing PBX Task.
//
// Parameters:
//   - taskID: This parameter specifies the Task ID.
//   - event: The name of the Event to send.
//   - param: (Optional) This parameter specifies the Event parameter.
//
// Returns:
//   - error: An error if the command fails.
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
//   - account: The name of the Account. The name can include the Domain name.
//     If the Domain name is not specified, the current user Domain is used.
//   - program: The name of the program (the .sppr file) to start.
//   - entry: (Optional) The program entry point. If empty, "main" is used.
//   - param: (Optional) Specifies the program parameter.
//     Accessible within the program via Vars().startParameter.
//
// Returns:
//   - string: A string with the new Task ID.
//   - error: An error if the command fails.
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
