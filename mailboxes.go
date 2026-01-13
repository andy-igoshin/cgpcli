// # Section: Mailbox Administration
package cgpcli

import (
	"fmt"
)

// CreateMailbox создает новый почтовый ящик.
func (cli *Cli) CreateMailbox(account, mailbox, class, authAccount string) error {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account and mailbox names are required")
	}

	args := make([]any, 0, 8)
	args = append(args, "CREATEMAILBOX", Atom(account), "MAILBOX", mailbox)
	if class != "" {
		args = append(args, "CLASS", class)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	_, err := cli.Query(args...)
	return err
}

// DeleteMailbox удаляет почтовый ящик (опционально рекурсивно).
func (cli *Cli) DeleteMailbox(account, mailbox string, recursive bool, authAccount string) error {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account and mailbox names are required")
	}

	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}

	args := make([]any, 0, 6)
	args = append(args, "DELETEMAILBOX", Atom(account), keyword, mailbox)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	_, err := cli.Query(args...)
	return err
}

// GetMailboxACL возвращает список прав доступа (ACL) к ящику.
func (cli *Cli) GetMailboxACL(account, mailbox, authAccount string) (map[string]string, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}

	args := make([]any, 0, 6)
	args = append(args, "GETMAILBOXACL", Atom(account), "MAILBOX", mailbox)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return map[string]string{}, nil
	}

	acl := make(map[string]string, len(dict))
	for k, v := range dict {
		s, err := cli.valToString("GETMAILBOXACL", v)
		if err != nil {
			return nil, err
		}
		acl[k] = s
	}
	return acl, nil
}

// GetMailboxAliases возвращает карту алиасов ящиков для аккаунта.
func (cli *Cli) GetMailboxAliases(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	res, err := cli.Query("GETMAILBOXALIASES", Atom(account))
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return map[string]string{}, nil
	}

	aliases := make(map[string]string, len(dict))
	for k, v := range dict {
		s, err := cli.valToString("GETMAILBOXALIASES", v)
		if err != nil {
			return nil, err
		}
		aliases[k] = s
	}
	return aliases, nil
}

// GetMailboxInfo возвращает техническую информацию о ящике (размер, кол-во писем и т.д.).
func (cli *Cli) GetMailboxInfo(account, mailbox, authAccount string) (map[string]any, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}

	args := make([]any, 0, 6)
	args = append(args, "GETMAILBOXINFO", Atom(account), "MAILBOX", mailbox)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return map[string]any{}, nil
	}
	return dict, nil
}

// GetMailboxRights возвращает итоговые (эффективные) права конкретного пользователя.
func (cli *Cli) GetMailboxRights(account, mailbox, authAccount string) (string, error) {
	if account == "" || mailbox == "" || authAccount == "" {
		return "", fmt.Errorf("account, mailbox and authAccount names are required")
	}

	res, err := cli.Query("GETMAILBOXRIGHTS", Atom(account), "MAILBOX", mailbox, "AUTH", authAccount)
	if err != nil {
		return "", err
	}

	return cli.valToString("GETMAILBOXRIGHTS", res)
}

// GetMailboxSubscription возвращает список ящиков, на которые подписан пользователь.
func (cli *Cli) GetMailboxSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	res, err := cli.Query("GETMAILBOXSUBSCRIPTION", Atom(account))
	if err != nil {
		return nil, err
	}

	return cli.toStringSlice2(res)
}

// ListMailboxes возвращает дерево ящиков аккаунта с метаданными.
func (cli *Cli) ListMailboxes(account, filter, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	args := make([]any, 0, 6)
	args = append(args, "LISTMAILBOXES", Atom(account))
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}

	dict, ok := res.(map[string]any)
	if !ok {
		return map[string]any{}, nil
	}
	return dict, nil
}

// RenameMailbox переименовывает ящик (опционально включая вложенные).
func (cli *Cli) RenameMailbox(account, oldName, newName string, recursive bool, authAccount string) error {
	if account == "" || oldName == "" || newName == "" {
		return fmt.Errorf("account, old and new names are required")
	}

	keyword := "MAILBOX"
	if recursive {
		keyword = "MAILBOXES"
	}

	args := make([]any, 0, 8)
	args = append(args, "RENAMEMAILBOX", Atom(account), keyword, oldName, "INTO", newName)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	_, err := cli.Query(args...)
	return err
}

// SetMailboxACL модифицирует права доступа (поддерживает префиксы + и -).
func (cli *Cli) SetMailboxACL(account, mailbox, authAccount string, newACL map[string]string) error {
	if account == "" || mailbox == "" || len(newACL) == 0 {
		return fmt.Errorf("account, mailbox and non-empty newACL are required")
	}

	args := make([]any, 0, 8)
	args = append(args, "SETMAILBOXACL", Atom(account), "MAILBOX", mailbox)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	args = append(args, newACL)

	_, err := cli.Query(args...)
	return err
}

// SetMailboxAliases полностью перезаписывает список алиасов ящиков.
func (cli *Cli) SetMailboxAliases(account string, newAliases map[string]string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETMAILBOXALIASES", Atom(account), newAliases)
	return err
}

// SetMailboxClass изменяет класс ящика (например, на IPF.Appointment).
func (cli *Cli) SetMailboxClass(account, mailbox, newClass, authAccount string) error {
	if account == "" || mailbox == "" || newClass == "" {
		return fmt.Errorf("account, mailbox and newClass are required")
	}

	args := make([]any, 0, 8)
	args = append(args, "SETMAILBOXCLASS", Atom(account), "MAILBOX", mailbox)
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	args = append(args, "CLASS", newClass)

	_, err := cli.Query(args...)
	return err
}

// SetMailboxSubscription полностью перезаписывает список подписок пользователя.
func (cli *Cli) SetMailboxSubscription(account string, newSub []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	_, err := cli.Query("SETMAILBOXSUBSCRIPTION", Atom(account), newSub)
	return err
}
