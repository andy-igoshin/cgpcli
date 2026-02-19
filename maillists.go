// # Mailing Lists Administration
//
// The Lists section provides comprehensive management of CommuniGate Pro Mailing Lists.
// It covers the entire lifecycle of a list, from creation and configuration to
// advanced subscriber management and bounce processing.
//
// Key capabilities include:
//   - List Management: creating, renaming, and deleting lists, with automatic
//     owner account resolution.
//   - Configuration: retrieving and updating list settings via [Cli.GetList]
//     and [Cli.UpdateList], including automated handling of CGP-specific
//     line endings (\e).
//   - Subscriber Operations: subscribing, unsubscribing, and managing posting
//     modes (moderation) using the [Cli.List] and [Cli.SetPostingMode] methods.
//   - Data Retrieval: detailed subscriber inspection through [SubscriberInfo],
//     supporting filters and pagination for large lists.
//   - Maintenance: initiating bounce processing for problematic addresses
//     and tracking delivery statistics (posts, bounces, confirmation IDs).
package cgpcli

import (
	"fmt"
	"strings"
	"time"
)

// SubscriberInfo contains detailed metadata about a mailing list member.
type SubscriberInfo struct {
	Sub            string    // subscriber's email address.
	RealName       string    // optional real name.
	Mode           string    // subscription mode (index, digest, null, etc.).
	SubscribeTime  time.Time // time when this subscriber was added.
	TimeSubscribed time.Time // time when the address was subscribed (ACAP format).
	Posts          int64     // number of postings on this list.
	Bounces        int64     // optional count of failed delivery reports received for this subscriber.
	LastBounceTime time.Time // optional time when the last delivery failure occurred for this subscriber.
	ConfirmationID int64     // subscriber's confirmation ID.
}

// CreateList creates a mailing list.
//
// Parameters:
//   - list: the name of a mailing list to create. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - account: the name of the mailing list owner. It should be the name of an already existing Account in the mailing list Domain.
//
// This method executes the CREATELIST CLI command.
//
// Returns:
//   - error: an error if the command fails.
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

// DeleteList removes a mailing list.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//
// This method executes the DELETELIST CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteList(list string) error {
	if list == "" {
		return fmt.Errorf("list name is required")
	}
	return cli.QueryNV("DELETELIST", list)
}

// GetAccountLists retrieves the list of all mailing lists belonging to the specified Account.
//
// Parameters:
//   - account: the list's owner Account name.
//
// This method executes the GETACCOUNTLISTS CLI command.
//
// Returns:
//   - map[string]int: a dictionary where each key is the name of a mailing list and the value is the number of subscribers.
//   - error: an error if the command fails.
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

// GetDomainLists retrieves the list of all mailing lists in the Domain.
//
// Parameters:
//   - domain: an optional Domain name.
//
// This method executes the GETDOMAINLISTS CLI command.
//
// Returns:
//   - map[string]int: a dictionary where each key is the name of a mailing list and the value is the number of subscribers.
//   - error: an error if the command fails.
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

// GetList retrieves list settings.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//
// This method executes the GETLIST CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the mailing list settings.
//   - error: an error if the command fails.
func (cli *Cli) GetList(list string) (map[string]any, error) {
	if list == "" {
		return nil, fmt.Errorf("list name is required")
	}
	return cli.getMapAny("GETLIST", list)
}

// GetSubscriberInfo retrieves information about a list subscriber.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - subscriber: the E-mail address of the list subscriber.
//
// This method executes the GETSUBSCRIBERINFO CLI command.
//
// Returns:
//   - *[SubscriberInfo]: a structure with subscriber information.
//   - error: an error if the command fails.
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

// List updates the subscribers list.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - operation: the operation (subscribe, feed, digest, index, null, banned, unsubscribe).
//   - subscriber: the subscriber address. It can include the comment part used as the subscriber's real name.
//   - silently: tells the server not to send the Welcome/Bye message to the subscriber.
//   - confirm: tells the server to send a confirmation request to the subscriber.
//
// Returns:
//   - error: an error if the command fails.
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

// ListLists retrieves the list of all mailing lists in the Domain.
//
// Parameters:
//   - domain: an optional Domain name.
//
// This method executes the LISTLISTS CLI command.
//
// Returns:
//   - []string: an array of strings where each string is the name of a mailing list in the specified (or default) Domain.
//   - error: an error if the command fails.
func (cli *Cli) ListLists(domain string) ([]string, error) {
	return cli.getSliceString("LISTLISTS", domain)
}

// ListSubscribers retrieves list subscribers.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - filter: an optional substring to filter addresses.
//   - limit: limits the number of subscriber addresses returned.
//
// This method executes the LISTSUBSCRIBERS CLI command.
//
// Returns:
//   - []string: an array with subscribers' E-mail addresses.
//   - error: an error if the command fails.
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

// ProcessBounce performs the same action the List Manager performs when it receives a bounce message for the subscriber address.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - subscriber: the E-mail address of the list subscriber.
//   - fatal: use the FATAL keyword to emulate a "fatal" bounce. Otherwise the command emulates a non-fatal bounce.
//
// This method executes the PROCESSBOUNCE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ProcessBounce(list, subscriber string, fatal bool) error {
	if list == "" || subscriber == "" {
		return fmt.Errorf("list name and subscriber address are required")
	}

	const cmd = "PROCESSBOUNCE"
	if fatal {
		return cli.QueryNV(cmd, list, "FATAL", "FOR", subscriber)
	} else {
		return cli.QueryNV(cmd, list, "FOR", subscriber)
	}
}

// ReadSubscribers retrieves list subscribers with detailed information.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - filter: an optional substring to filter subscriber addresses.
//   - limit: limits the number of subscriber dictionaries returned.
//
// This method executes the READSUBSCRIBERS CLI command.
//
// Returns:
//   - [][SubscriberInfo]: an array of subscriber descriptor structures.
//   - error: an error if the command fails.
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

// RenameList renames a mailing list.
//
// Parameters:
//   - oldName: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - newName: the new name for the mailing list.
//
// This method executes the RENAMELIST CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameList(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new list names are required")
	}

	new := newName
	if idx := strings.IndexByte(newName, '@'); idx != -1 {
		new = newName[:idx]
	}
	return cli.QueryNV("RENAMELIST", oldName, "INTO", new)
}

// SetPostingMode sets the posting mode for the specified subscriber.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - subscriber: the E-mail address of the list subscriber.
//   - mode: can be UNMODERATED, MODERATEALL, PROHIBITED, SPECIAL, or a number (numberOfModerated).
//
// This method executes the SETPOSTINGMODE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetPostingMode(list, subscriber string, mode any) error {
	if list == "" || subscriber == "" {
		return fmt.Errorf("list name and subscriber address are required")
	}

	if mode != nil {
		return cli.QueryNV("SETPOSTINGMODE", list, "FOR", subscriber, mode)
	} else {
		return cli.QueryNV("SETPOSTINGMODE", list, "FOR", subscriber)
	}
}

// UpdateList modifies list settings.
//
// Parameters:
//   - list: the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - settings: this dictionary is used to update the mailing list settings dictionary. The omitted settings will be left unmodified.
//
// This method executes the UPDATELIST CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
