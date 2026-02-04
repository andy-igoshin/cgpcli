package cgpcli

import (
	"fmt"
	"strings"
	"time"
)

// SubscriberInfo описывает данные подписчика списка рассылки.
type SubscriberInfo struct {
	Bounces        int64
	ConfirmationID int64
	LastBounceTime time.Time
	Mode           string
	Posts          int64
	RealName       string
	Sub            string
	SubscribeTime  time.Time
	TimeSubscribed time.Time
}

// CreateList создает новый список рассылки для указанного аккаунта-владельца.
func (cli *Cli) CreateList(list, account string) error {
	if list == "" || account == "" {
		return fmt.Errorf("list name and account name are required")
	}

	owner := account
	if idx := strings.IndexByte(account, '@'); idx != -1 {
		owner = account[:idx]
	}
	return cli.QueryNV("CREATELIST", list, "FOR", owner)
}

// DeleteList удаляет список рассылки.
func (cli *Cli) DeleteList(list string) error {
	if list == "" {
		return fmt.Errorf("list name is required")
	}
	return cli.QueryNV("DELETELIST", list)
}

// GetAccountLists возвращает списки рассылки, принадлежащие конкретному аккаунту.
func (cli *Cli) GetAccountLists(account string) (map[string]int, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	m, err := cli.getMapAny("GETACCOUNTLISTS", account)
	if err != nil {
		return nil, err
	}
	list := make(map[string]int, len(m))
	for k, v := range m {
		list[k] = toInt(v)
	}
	return list, nil
}

// GetDomainLists возвращает map "имя списка рассылки -> количество подписчиков".
func (cli *Cli) GetDomainLists(domain string) (map[string]int, error) {
	m, err := cli.getMapAny("GETDOMAINLISTS", Atom(domain))
	if err != nil {
		return nil, err
	}
	list := make(map[string]int, len(m))
	for k, v := range m {
		list[k] = toInt(v)
	}
	return list, nil
}

// GetList получает настройки списка. Все текстовые поля конвертируются из \e в \n.
func (cli *Cli) GetList(list string) (map[string]any, error) {
	if list == "" {
		return nil, fmt.Errorf("list name is required")
	}

	m, err := cli.getMapAny("GETLIST", list)
	if err != nil {
		return nil, err
	}
	for k, v := range m {
		if s, ok := v.(string); ok {
			m[k] = toGoText(s)
		}
	}
	return m, nil
}

// GetSubscriberInfo возвращает полную информацию о конкретном подписчике.
func (cli *Cli) GetSubscriberInfo(list, subscriber string) (*SubscriberInfo, error) {
	if list == "" || subscriber == "" {
		return nil, fmt.Errorf("list name and subscriber address are required")
	}

	res, err := cli.Query("GETSUBSCRIBERINFO", list, "NAME", subscriber)
	if err != nil {
		return nil, err
	}

	m, ok := res.(map[string]any)
	if !ok || len(m) == 0 {
		return nil, nil
	}

	s := mapToSubscriber(m)
	return &s, nil
}

// List выполняет операции над подписчиками (subscribe, unsubscribe, feed и др.).
func (cli *Cli) List(list, operation, subscriber string, silently, confirm bool) error {
	if list == "" || operation == "" || subscriber == "" {
		return fmt.Errorf("list name, operation and subscriber are required")
	}

	args := make([]any, 0, 5)
	args = append(args, list, operation)
	if silently {
		args = append(args, "silently")
	}
	if confirm {
		args = append(args, "confirm")
	}
	args = append(args, subscriber)
	return cli.QueryNV("LIST", args...)
}

// ListLists возвращает список имен всех рассылок в домене.
func (cli *Cli) ListLists(domain string) ([]string, error) {
	return cli.getSliceString("LISTLISTS", domain)
}

// ListSubscribers возвращает только список адресов подписчиков.
func (cli *Cli) ListSubscribers(list, filter string, limit int) ([]string, error) {
	if list == "" {
		return nil, fmt.Errorf("list name is required")
	}

	const cmd = "LISTSUBSCRIBERS"
	args := make([]any, 0, 4)
	args = append(args, list)
	if filter != "" {
		args = append(args, "FILTER", filter)
		if limit > 0 {
			args = append(args, limit)
		}
	}
	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ProcessBounce инициирует процедуру обработки ошибки доставки (bounce) для адреса.
func (cli *Cli) ProcessBounce(list, subscriber string, fatal bool) (err error) {
	if list == "" || subscriber == "" {
		return fmt.Errorf("list name and subscriber address are required")
	}

	const cmd = "PROCESSBOUNCE"
	if fatal {
		err = cli.QueryNV(cmd, list, "FATAL", "FOR", subscriber)
	} else {
		err = cli.QueryNV(cmd, list, "FOR", subscriber)
	}
	return
}

// ReadSubscribers возвращает общее число подписчиков и массив структур SubscriberInfo.
func (cli *Cli) ReadSubscribers(list, filter string, limit int) ([]SubscriberInfo, error) {
	if list == "" {
		return nil, fmt.Errorf("list name is required")
	}

	args := make([]any, 0, 4)
	args = append(args, list)
	if filter != "" {
		args = append(args, "FILTER", filter)
		if limit > 0 {
			args = append(args, limit)
		}
	}

	res, err := cli.Query("READSUBSCRIBERS", args...)
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok || len(raw) < 2 {
		return nil, nil
	}

	ls, _ := raw[1].([]any)

	subs := make([]SubscriberInfo, 0, len(ls))
	for _, v := range ls {
		if m, ok := v.(map[string]any); ok {
			subs = append(subs, mapToSubscriber(m))
		}
	}
	return subs, nil
}

// RenameList переименовывает список рассылки (без доменной части в новом имени).
func (cli *Cli) RenameList(list, newName string) error {
	if list == "" || newName == "" {
		return fmt.Errorf("list name and new list name are required")
	}

	new := newName
	if idx := strings.IndexByte(newName, '@'); idx != -1 {
		new = newName[:idx]
	}
	return cli.QueryNV("RENAMELIST", list, "INTO", new)
}

// SetPostingMode устанавливает режим модерации/прав публикации для подписчика.
func (cli *Cli) SetPostingMode(list, subscriber string, mode any) (err error) {
	if list == "" || subscriber == "" {
		return fmt.Errorf("list name and subscriber address are required")
	}

	if mode != nil {
		err = cli.QueryNV("SETPOSTINGMODE", list, "FOR", subscriber, mode)
	} else {
		err = cli.QueryNV("SETPOSTINGMODE", list, "FOR", subscriber)
	}
	return
}

// UpdateList модифицирует настройки списка. Текст с \n автоматически конвертируется в \e.
func (cli *Cli) UpdateList(list string, settings map[string]any) error {
	if list == "" || len(settings) == 0 {
		return fmt.Errorf("list name and settings are required")
	}

	prepared := make(map[string]any, len(settings))
	for k, v := range settings {
		if s, ok := v.(string); ok {
			prepared[k] = toCGPText(s)
		} else {
			prepared[k] = v
		}
	}
	return cli.QueryNV("UPDATELIST", list, prepared)
}

// mapToSubscriber преобразует map[string]any в структуру SubscriberInfo.
func mapToSubscriber(m map[string]any) SubscriberInfo {
	const fn = "mapToSubscriber"

	sub, _ := toString(fn, m["Sub"])
	rName, _ := toString(fn, m["RealName"])
	mode, _ := toString(fn, m["mode"])

	s := SubscriberInfo{
		Sub:            sub,
		RealName:       rName,
		Mode:           mode,
		SubscribeTime:  toTime(m["subscribeTime"]),
		Posts:          toInt64(m["posts"]),
		Bounces:        toInt64(m["bounces"]),
		LastBounceTime: toTime(m["lastBounceTime"]),
		ConfirmationID: toInt64(m["ConfirmationID"]),
	}

	if tsStr, err := toString(fn, m["timeSubscribed"]); err == nil && tsStr != "" {
		if t, err := time.Parse("20060102150405", tsStr); err == nil {
			s.TimeSubscribed = t
		}
	}

	return s
}
