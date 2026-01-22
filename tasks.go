// # Section: Named Task Administration
package cgpcli

import (
	"fmt"
	"strings"
)

// CreateNamedTask создает новую именованную задачу.
func (cli *Cli) CreateNamedTask(task, account string) (err error) {
	if task == "" || account == "" {
		return fmt.Errorf("task name and account name are required")
	}
	// Убираем домен из account, если он есть
	if idx := strings.IndexByte(account, '@'); idx != -1 {
		account = account[:idx]
	}
	_, err = cli.Query("CREATENAMEDTASK", task, "FOR", account)
	return
}

// DeleteNamedTask удаляет именованную задачу.
func (cli *Cli) DeleteNamedTask(task string) (err error) {
	if task == "" {
		return fmt.Errorf("task name is required")
	}
	_, err = cli.Query("DELETENAMEDTASK", task)
	return
}

// GetNamedTask возвращает словарь с настройками именованной задачи.
func (cli *Cli) GetNamedTask(task string) (map[string]any, error) {
	if task == "" {
		return nil, fmt.Errorf("task name is required")
	}
	const cmd = "GETNAMEDTASK"
	res, err := cli.Query(cmd, task)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// ListAccountNamedTasks возвращает список задач, принадлежащих конкретному аккаунту.
func (cli *Cli) ListAccountNamedTasks(account string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	const cmd = "LISTACCOUNTNAMEDTASKS"
	res, err := cli.Query(cmd, account)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// ListDomainNamedTasks возвращает список всех именованных задач в домене.
func (cli *Cli) ListDomainNamedTasks(domain string) (map[string]any, error) {
	const cmd = "LISTDOMAINNAMEDTASKS"
	var res any
	var err error
	if domain != "" {
		res, err = cli.Query(cmd, domain)
	} else {
		res, err = cli.Query(cmd)
	}

	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// RenameNamedTask переименовывает задачу.
func (cli *Cli) RenameNamedTask(oldTask, newTask string) (err error) {
	if oldTask == "" || newTask == "" {
		return fmt.Errorf("oldName and newName are required")
	}
	_, err = cli.Query("RENAMENAMEDTASK", oldTask, "into", newTask)
	return
}

// UpdateNamedTask обновляет настройки именованной задачи.
func (cli *Cli) UpdateNamedTask(task string, settings map[string]any) (err error) {
	if task == "" || len(settings) == 0 {
		return fmt.Errorf("task name and settings are required")
	}
	_, err = cli.Query("UPDATENAMEDTASK", task, settings)
	return
}
