// # Section: Real-Time Application Control
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
