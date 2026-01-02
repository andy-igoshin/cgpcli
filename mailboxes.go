package cgpcli

import (
	"fmt"
	"strings"
)

// CreateMailbox создает новый почтовый ящик.
func (cli *Cli) CreateMailbox(account, mailbox, class, authAccount string) error {
	if account == "" || mailbox == "" {
		return fmt.Errorf("account and mailbox names are required")
	}
	parts := []string{"CREATEMAILBOX", Marshal(account), "MAILBOX", Marshal(mailbox)}
	if class != "" {
		parts = append(parts, "CLASS", Marshal(class))
	}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	_, err := cli.query(strings.Join(parts, " "))
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
	parts := []string{"DELETEMAILBOX", Marshal(account), keyword, Marshal(mailbox)}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// GetMailboxACL возвращает список прав доступа (ACL) к ящику.
func (cli *Cli) GetMailboxACL(account, mailbox, authAccount string) (map[string]string, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}
	parts := []string{"GETMAILBOXACL", Marshal(account), "MAILBOX", Marshal(mailbox)}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return map[string]string{}, nil
	}
	acl := make(map[string]string, len(dict))
	for k, v := range dict {
		acl[k] = fmt.Sprint(v)
	}
	return acl, nil
}

// GetMailboxAliases возвращает карту алиасов ящиков для аккаунта.
func (cli *Cli) GetMailboxAliases(account string) (map[string]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETMAILBOXALIASES " + Marshal(account))
	if err != nil {
		return nil, err
	}
	dict, ok := res.(map[string]any)
	if !ok {
		return map[string]string{}, nil
	}
	aliases := make(map[string]string, len(dict))
	for k, v := range dict {
		aliases[k] = fmt.Sprint(v)
	}
	return aliases, nil
}

// GetMailboxInfo возвращает техническую информацию о ящике (размер, кол-во писем и т.д.).
func (cli *Cli) GetMailboxInfo(account, mailbox, authAccount string) (map[string]any, error) {
	if account == "" || mailbox == "" {
		return nil, fmt.Errorf("account and mailbox names are required")
	}
	parts := []string{"GETMAILBOXINFO", Marshal(account), "MAILBOX", Marshal(mailbox)}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	res, err := cli.query(strings.Join(parts, " "))
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
	parts := []string{"GETMAILBOXRIGHTS", Marshal(account), "MAILBOX", Marshal(mailbox), "AUTH", Marshal(authAccount)}
	res, err := cli.query(strings.Join(parts, " "))
	if err != nil {
		return "", err
	}
	return fmt.Sprint(res), nil
}

// GetMailboxSubscription возвращает список ящиков, на которые подписан пользователь.
func (cli *Cli) GetMailboxSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	res, err := cli.query("GETMAILBOXSUBSCRIPTION " + Marshal(account))
	if err != nil {
		return nil, err
	}
	raw, ok := res.([]any)
	if !ok {
		return []string{}, nil
	}
	list := make([]string, len(raw))
	for i, v := range raw {
		list[i] = fmt.Sprint(v)
	}
	return list, nil
}

// ListMailboxes возвращает дерево ящиков аккаунта с метаданными.
func (cli *Cli) ListMailboxes(account, filter, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	parts := []string{"LISTMAILBOXES", Marshal(account)}
	if filter != "" {
		parts = append(parts, "FILTER", Marshal(filter))
	}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	res, err := cli.query(strings.Join(parts, " "))
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
	parts := []string{"RENAMEMAILBOX", Marshal(account), keyword, Marshal(oldName), "INTO", Marshal(newName)}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetMailboxACL модифицирует права доступа (поддерживает префиксы + и -).
func (cli *Cli) SetMailboxACL(account, mailbox, authAccount string, newACL map[string]string) error {
	if account == "" || mailbox == "" || len(newACL) == 0 {
		return fmt.Errorf("account, mailbox and non-empty newACL are required")
	}
	parts := []string{"SETMAILBOXACL", Marshal(account), "MAILBOX", Marshal(mailbox)}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	parts = append(parts, Marshal(newACL))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetMailboxAliases полностью перезаписывает список алиасов ящиков.
func (cli *Cli) SetMailboxAliases(account string, newAliases map[string]string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"SETMAILBOXALIASES", Marshal(account), Marshal(newAliases)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetMailboxClass изменяет класс ящика (например, на IPF.Appointment).
func (cli *Cli) SetMailboxClass(account, mailbox, newClass, authAccount string) error {
	if account == "" || mailbox == "" || newClass == "" {
		return fmt.Errorf("account, mailbox and newClass are required")
	}
	parts := []string{"SETMAILBOXCLASS", Marshal(account), "MAILBOX", Marshal(mailbox)}
	if authAccount != "" {
		parts = append(parts, "AUTH", Marshal(authAccount))
	}
	parts = append(parts, "CLASS", Marshal(newClass))
	_, err := cli.query(strings.Join(parts, " "))
	return err
}

// SetMailboxSubscription полностью перезаписывает список подписок пользователя.
func (cli *Cli) SetMailboxSubscription(account string, newSub []string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	parts := []string{"SETMAILBOXSUBSCRIPTION", Marshal(account), Marshal(newSub)}
	_, err := cli.query(strings.Join(parts, " "))
	return err
}
