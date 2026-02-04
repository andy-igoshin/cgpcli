// # Section: Named Task Administration
package cgpcli

import (
	"fmt"
	"strings"
)

// CreateNamedTask создает новую именованную задачу.
func (cli *Cli) CreateNamedTask(task, account string) error {
	if task == "" || account == "" {
		return fmt.Errorf("task name and account name are required")
	}
	if idx := strings.IndexByte(account, '@'); idx != -1 {
		account = account[:idx]
	}
	return cli.QueryNV("CREATENAMEDTASK", task, "FOR", account)
}

// DeleteNamedTask удаляет именованную задачу.
func (cli *Cli) DeleteNamedTask(task string) error {
	if task == "" {
		return fmt.Errorf("task name is required")
	}
	return cli.QueryNV("DELETENAMEDTASK", task)
}

// GetNamedTask возвращает словарь с настройками именованной задачи.
func (cli *Cli) GetNamedTask(task string) (map[string]any, error) {
	if task == "" {
		return nil, fmt.Errorf("task name is required")
	}
	return cli.getMapAny("GETNAMEDTASK", task)
}

// ListAccountNamedTasks возвращает список задач, принадлежащих конкретному аккаунту.
func (cli *Cli) ListAccountNamedTasks(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getMapAny("LISTACCOUNTNAMEDTASKS", account)
}

// ListDomainNamedTasks возвращает список всех именованных задач в домене.
func (cli *Cli) ListDomainNamedTasks(domain string) (map[string]any, error) {
	return cli.getMapAny("LISTDOMAINNAMEDTASKS", Atom(domain))
}

// RenameNamedTask переименовывает задачу.
func (cli *Cli) RenameNamedTask(oldTask, newTask string) error {
	if oldTask == "" || newTask == "" {
		return fmt.Errorf("oldName and newName are required")
	}
	return cli.QueryNV("RENAMENAMEDTASK", oldTask, "INTO", newTask)
}

// UpdateNamedTask обновляет настройки именованной задачи.
func (cli *Cli) UpdateNamedTask(task string, settings map[string]any) error {
	if task == "" || len(settings) == 0 {
		return fmt.Errorf("task name and settings are required")
	}
	return cli.QueryNV("UPDATENAMEDTASK", task, settings)
}
