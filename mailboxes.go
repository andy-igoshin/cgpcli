// # Section: Mailbox Administration
package cgpcli

import (
	"fmt"
)

// CreateMailbox создает папку.
func (cli *Cli) CreateMailbox(account, mailbox, class, authAccount string) error {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account and mailbox names are required")
	}

	args := make([]any, 0, 7)
	args = append(args, Atom(account), "MAILBOX", mailbox)
	if class != "" {
		args = append(args, "CLASS", class)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	return cli.QueryNV("CREATEMAILBOX", args...)
}

// DeleteMailbox удаляет папку или папки.
func (cli *Cli) DeleteMailbox(account, mailbox string, recursive bool, authAccount string) (err error) {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account and mailbox names are required")
	}

	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}
	if authAccount != "" {
		err = cli.QueryNV("DELETEMAILBOX", Atom(account), keyword, mailbox, "AUTH", authAccount)
	} else {
		err = cli.QueryNV("DELETEMAILBOX", Atom(account), keyword, mailbox)
	}
	return
}

// GetMailboxACL возвращает список прав доступа (ACL) к папке.
func (cli *Cli) GetMailboxACL(account, mailbox, authAccount string) (map[string]string, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}

	const cmd = "GETMAILBOXACL"
	var res map[string]any
	var err error
	if authAccount != "" {
		res, err = cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	} else {
		res, err = cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox)
	}
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

// GetMailboxAliases возвращает map алиасов папок.
func (cli *Cli) GetMailboxAliases(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	const cmd = "GETMAILBOXALIASES"
	res, err := cli.getMapAny(cmd, Atom(account))
	if err != nil {
		return nil, err
	}
	return toStringMap(cmd, res)
}

// GetMailboxInfo возвращает информацию о папке (размер, количество писем и т.д.).
func (cli *Cli) GetMailboxInfo(account, mailbox, authAccount string) (res map[string]any, err error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}

	const cmd = "GETMAILBOXINFO"
	if authAccount != "" {
		res, err = cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	} else {
		res, err = cli.getMapAny(cmd, Atom(account), "MAILBOX", mailbox)
	}
	return
}

// GetMailboxRights возвращает эффективные права доступа к папке.
func (cli *Cli) GetMailboxRights(account, mailbox, authAccount string) (string, error) {
	if account == "" || mailbox == "" || authAccount == "" {
		return "", fmt.Errorf("account, mailbox and authAccount names are required")
	}
	return cli.getString("GETMAILBOXRIGHTS", Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
}

// GetMailboxSubscription возвращает список папок, на которые подписан аккаунт.
func (cli *Cli) GetMailboxSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETMAILBOXSUBSCRIPTION", account)
}

// ListMailboxes возвращает дерево папок аккаунта с метаданными.
func (cli *Cli) ListMailboxes(account, filter, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	args := make([]any, 0, 5)
	args = append(args, Atom(account))
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	return cli.getMapAny("LISTMAILBOXES", args...)
}

// RenameMailbox переименовывает ящик (опционально включая вложенные).
func (cli *Cli) RenameMailbox(account, oldName, newName string, recursive bool, authAccount string) (err error) {
	if account == "" || oldName == "" || newName == "" {
		return fmt.Errorf("account, old and new names are required")
	}

	const cmd = "RENAMEMAILBOX"
	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), keyword, oldName, "INTO", newName, "AUTH", authAccount)
	} else {
		err = cli.QueryNV(cmd, Atom(account), keyword, oldName, "INTO", newName)
	}
	return
}

// SetMailboxACL модифицирует права доступа (поддерживает префиксы + и -).
func (cli *Cli) SetMailboxACL(account, mailbox, authAccount string, acl map[string]string) (err error) {
	if account == "" || mailbox == "" || len(acl) == 0 {
		return fmt.Errorf("account, mailbox and non-empty acl are required")
	}

	const cmd = "SETMAILBOXACL"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount, acl)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, acl)
	}
	return
}

// SetMailboxAliases полностью перезаписывает список алиасов ящиков.
func (cli *Cli) SetMailboxAliases(account string, aliases map[string]string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("SETMAILBOXALIASES", Atom(account), aliases)
}

// SetMailboxClass изменяет класс ящика (например, на IPF.Appointment).
func (cli *Cli) SetMailboxClass(account, mailbox, class, authAccount string) (err error) {
	if account == "" || mailbox == "" || class == "" {
		return fmt.Errorf("account, mailbox and class are required")
	}

	const cmd = "SETMAILBOXCLASS"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "AUTH", authAccount, "CLASS", class)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "MAILBOX", mailbox, "CLASS", class)
	}
	return
}

// SetMailboxSubscription полностью перезаписывает список подписок пользователя.
func (cli *Cli) SetMailboxSubscription(account string, subs []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("SETMAILBOXSUBSCRIPTION", Atom(account), subs)
}
