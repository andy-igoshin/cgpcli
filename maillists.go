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
	owner := account
	if idx := strings.IndexByte(account, '@'); idx != -1 {
		owner = account[:idx]
	}
	_, err := cli.Query("CREATELIST", list, "FOR", owner)
	return err
}

// DeleteList удаляет список рассылки.
func (cli *Cli) DeleteList(list string) error {
	_, err := cli.Query("DELETELIST", list)
	return err
}

// GetAccountLists возвращает списки, принадлежащие конкретному аккаунту.
func (cli *Cli) GetAccountLists(account string) (map[string]int, error) {
	res, err := cli.Query("GETACCOUNTLISTS", account)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETACCOUNTLISTS: expected map, got %T", res)
	}
	out := make(map[string]int, len(m))
	for k, v := range m {
		out[k] = toInt(v)
	}
	return out, nil
}

// GetDomainLists возвращает карту "имя списка -> количество подписчиков".
func (cli *Cli) GetDomainLists(domain string) (map[string]int, error) {
	res, err := cli.Query("GETDOMAINLISTS", domain)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETDOMAINLISTS: expected map, got %T", res)
	}

	out := make(map[string]int, len(m))
	for k, v := range m {
		out[k] = toInt(v)
	}
	return out, nil
}

// GetList получает настройки списка. Все текстовые поля конвертируются из \e в \n.
func (cli *Cli) GetList(list string) (map[string]any, error) {
	res, err := cli.Query("GETLIST", list)
	if err != nil {
		return nil, err
	}

	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETLIST: expected dictionary, got %T", res)
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
	args := make([]any, 0, 6)
	args = append(args, "LIST", list, operation)
	if silently {
		args = append(args, "silently")
	}
	if confirm {
		args = append(args, "confirm")
	}
	args = append(args, subscriber)
	_, err := cli.Query(args...)
	return err
}

// ListLists возвращает список имен всех рассылок в домене.
func (cli *Cli) ListLists(domain string) ([]string, error) {
	const cmd = "LISTLISTS"
	res, err := cli.Query(cmd, domain)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ListSubscribers возвращает только список адресов подписчиков.
func (cli *Cli) ListSubscribers(list, filter string, limit int) ([]string, error) {
	const cmd = "LISTSUBSCRIBERS"
	args := []any{cmd, list}
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if limit > 0 {
		args = append(args, limit)
	}
	res, err := cli.Query(args...)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ProcessBounce инициирует процедуру обработки ошибки доставки (bounce) для адреса.
func (cli *Cli) ProcessBounce(list, subscriber string, fatal bool) error {
	args := []any{"PROCESSBOUNCE", list}
	if fatal {
		args = append(args, "FATAL")
	}
	args = append(args, "FOR", subscriber)
	_, err := cli.Query(args...)
	return err
}

// ReadSubscribers возвращает общее число подписчиков и массив структур SubscriberInfo.
func (cli *Cli) ReadSubscribers(list, filter string, limit int) (int64, []SubscriberInfo, error) {
	args := []any{"READSUBSCRIBERS", list}
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if limit > 0 {
		args = append(args, limit)
	}

	res, err := cli.Query(args...)
	if err != nil {
		return 0, nil, err
	}

	root, ok := res.([]any)
	if !ok || len(root) < 2 {
		return 0, nil, nil
	}

	total := toInt64(root[0])
	listRaw, _ := root[1].([]any)

	subs := make([]SubscriberInfo, 0, len(listRaw))
	for _, v := range listRaw {
		if m, ok := v.(map[string]any); ok {
			subs = append(subs, mapToSubscriber(m))
		}
	}
	return total, subs, nil
}

// RenameList переименовывает список рассылки (без доменной части в новом имени).
func (cli *Cli) RenameList(list, newName string) error {
	_, err := cli.Query("RENAMELIST", list, "into", newName)
	return err
}

// SetPostingMode устанавливает режим модерации/прав публикации для подписчика.
func (cli *Cli) SetPostingMode(list, subscriber string, mode any) error {
	_, err := cli.Query("SETPOSTINGMODE", list, "FOR", subscriber, mode)
	return err
}

// UpdateList модифицирует настройки списка. Текст с \n автоматически конвертируется в \e.
func (cli *Cli) UpdateList(list string, settings map[string]any) error {
	prepared := make(map[string]any, len(settings))
	for k, v := range settings {
		if s, ok := v.(string); ok {
			prepared[k] = toCGPText(s)
		} else {
			prepared[k] = v
		}
	}
	_, err := cli.Query("UPDATELIST", list, prepared)
	return err
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
