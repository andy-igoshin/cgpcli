// # Mailing Lists Administration
//
// The Lists section provides comprehensive management of CommuniGate Pro Mailing Lists.
// It covers the entire lifecycle of a list, from creation and configuration to
// advanced subscriber management and bounce processing.
//
// Key capabilities include:
//   - List Management: Creating, renaming, and deleting lists, with automatic
//     owner account resolution.
//   - Configuration: Retrieving and updating list settings via [Cli.GetList]
//     and [Cli.UpdateList], including automated handling of CGP-specific
//     line endings (\e).
//   - Subscriber Operations: Subscribing, unsubscribing, and managing posting
//     modes (moderation) using the [Cli.List] and [Cli.SetPostingMode] methods.
//   - Data Retrieval: Detailed subscriber inspection through [SubscriberInfo],
//     supporting filters and pagination for large lists.
//   - Maintenance: Initiating bounce processing for problematic addresses
//     and tracking delivery statistics (posts, bounces, confirmation IDs).
//
// Use these methods to build custom mailing list managers, automated newsletter
// systems, or integration with external CRM platforms.
package cgpcli

import (
	"fmt"
	"strings"
	"time"
)

// SubscriberInfo contains detailed metadata about a mailing list member.
type SubscriberInfo struct {
	Sub            string    // Subscriber's email address.
	RealName       string    // Subscriber's display name.
	Mode           string    // Posting mode (e.g., "feed", "digest", "moderated").
	SubscribeTime  time.Time // Time when the subscription was initiated.
	TimeSubscribed time.Time // Exact time when the subscription was confirmed.
	Posts          int64     // Number of messages posted by this subscriber.
	Bounces        int64     // Number of delivery errors encountered.
	LastBounceTime time.Time // Timestamp of the most recent delivery error.
	ConfirmationID int64     // ID used for subscription confirmation tracking.
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

	res, err := cli.getMapAny("GETACCOUNTLISTS", account)
	if err != nil {
		return nil, err
	}
	list := make(map[string]int, len(res))
	for k, v := range res {
		list[k] = toInt(v)
	}
	return list, nil
}

// GetDomainLists возвращает map "имя списка рассылки -> количество подписчиков".
func (cli *Cli) GetDomainLists(domain string) (map[string]int, error) {
	res, err := cli.getMapAny("GETDOMAINLISTS", Atom(domain))
	if err != nil {
		return nil, err
	}
	list := make(map[string]int, len(res))
	for k, v := range res {
		list[k] = toInt(v)
	}
	return list, nil
}

// GetList получает настройки списка. Все текстовые поля конвертируются из \e в \n.
func (cli *Cli) GetList(list string) (map[string]any, error) {
	if list == "" {
		return nil, fmt.Errorf("list name is required")
	}
	return cli.getMapAny("GETLIST", list)
}

// GetSubscriberInfo возвращает полную информацию о конкретном подписчике.
func (cli *Cli) GetSubscriberInfo(list, subscriber string) (*SubscriberInfo, error) {
	if list == "" || subscriber == "" {
		return nil, fmt.Errorf("list name and subscriber address are required")
	}

	res, err := cli.getMapAny("GETSUBSCRIBERINFO", list, "NAME", subscriber)
	if err != nil {
		return nil, err
	}
	s := mapToSubscriber(res)
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

	args := make([]any, 0, 4)
	args = append(args, list)
	if filter != "" {
		args = append(args, "FILTER", filter)
		if limit > 0 {
			args = append(args, limit)
		}
	}
	return cli.getSliceString("LISTSUBSCRIBERS", args...)
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

	res, err := cli.getSliceAny("READSUBSCRIBERS", args...)
	if err != nil {
		return nil, err
	}
	if len(res) < 2 {
		return nil, nil
	}

	ls, _ := res[1].([]any)

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
	return cli.QueryNV("UPDATELIST", list, settings)
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
