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
	Sub            string    // Subscriber's email address string.
	RealName       string    // An optional string with Real name.
	Mode           string    // A string with subscription mode (index, digest, null, etc.).
	SubscribeTime  time.Time // Timestamp data specifying the moment when this user subscribed.
	TimeSubscribed time.Time // Exact time when the address was subscribed (in the ACAP date/time format).
	Posts          int64     // Number of postings on this list.
	Bounces        int64     // Optional numeric data specifying the number of failed delivery reports received for this user.
	LastBounceTime time.Time // Optional timestamp data specifying the last time when messages sent to this user failed.
	ConfirmationID int64     // This element contains the subscriber's Confirmation ID string.
}

// CreateList creates a mailing list.
//
// Parameters:
//   - list: This parameter specifies the name of a mailing list to create. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - account: This parameter specifies the name of the mailing list owner. It should be the name of an already existing Account in the mailing list Domain.
//
// Returns:
//   - error: An error if the command fails.
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
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteList(list string) error {
	if list == "" {
		return fmt.Errorf("list name is required")
	}
	return cli.QueryNV("DELETELIST", list)
}

// GetAccountLists retrieves the list of all mailing lists belonging to the specified Account.
//
// Parameters:
//   - account: This parameter specifies the list's owner Account name.
//
// Returns:
//   - map[string]int: A dictionary where each key is the name of a mailing list and the value is the number of subscribers.
//   - error: An error if the command fails.
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
//   - domain: (Optional) This optional parameter specifies the Domain name.
//
// Returns:
//   - map[string]int: A dictionary where each key is the name of a mailing list and the value is the number of subscribers.
//   - error: An error if the command fails.
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
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//
// Returns:
//   - map[string]any: A dictionary with the mailing list settings.
//   - error: An error if the command fails.
func (cli *Cli) GetList(list string) (map[string]any, error) {
	if list == "" {
		return nil, fmt.Errorf("list name is required")
	}
	return cli.getMapAny("GETLIST", list)
}

// GetSubscriberInfo retrieves information about a list subscriber.
//
// Parameters:
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - subscriber: This parameter specifies the E-mail address of the list subscriber.
//
// Returns:
//   - *[SubscriberInfo]: A structure with subscriber information.
//   - error: An error if the command fails.
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
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - operation: This parameter specifies the operation (subscribe, feed, digest, index, null, banned, unsubscribe).
//   - subscriber: The subscriber address. It can include the comment part used as the subscriber's real name.
//   - silently: (Optional) This parameter tells the server not to send the Welcome/Bye message to the subscriber.
//   - confirm: (Optional) This parameter tells the server to send a confirmation request to the subscriber.
//
// Returns:
//   - error: An error if the command fails.
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
//   - domain: (Optional) This optional parameter specifies the Domain name.
//
// Returns:
//   - []string: An array of strings where each string is the name of a mailing list in the specified (or default) Domain.
//   - error: An error if the command fails.
func (cli *Cli) ListLists(domain string) ([]string, error) {
	return cli.getSliceString("LISTLISTS", domain)
}

// ListSubscribers retrieves list subscribers.
//
// Parameters:
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - filter: (Optional) If this parameter is specified, only the addresses containing the specified string are returned.
//   - limit: (Optional) This parameter limits the number of subscriber addresses returned.
//
// Returns:
//   - []string: An array with subscribers' E-mail addresses.
//   - error: An error if the command fails.
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
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - subscriber: This parameter specifies the E-mail address of the list subscriber.
//   - fatal: Use the FATAL keyword to emulate a "fatal" bounce. Otherwise the command emulates a non-fatal bounce.
//
// Returns:
//   - error: An error if the command fails.
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

// ReadSubscribers retrieves list subscribers with detailed information.
//
// Parameters:
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - filter: (Optional) If this optional parameter is specified, only subscribers with addresses containing the specified string are returned.
//   - limit: (Optional) This optional parameter limits the number of subscriber dictionaries returned.
//
// Returns:
//   - [][SubscriberInfo]: An array of subscriber descriptor structures.
//   - error: An error if the command fails.
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
//   - oldName: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - newName: This parameter specifies the new name for the mailing list.
//
// Returns:
//   - error: An error if the command fails.
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
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - subscriber: This parameter specifies the E-mail address of the list subscriber.
//   - mode: (Optional) This parameter can be UNMODERATED, MODERATEALL, PROHIBITED, SPECIAL, or a number (numberOfModerated).
//
// Returns:
//   - error: An error if the command fails.
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

// UpdateList modifies list settings.
//
// Parameters:
//   - list: This parameter specifies the name of an existing mailing list. It can include the Domain name. If the Domain name is not specified, the user Domain is used by default.
//   - settings: This dictionary is used to update the mailing list settings dictionary. The omitted settings will be left unmodified.
//
// Returns:
//   - error: An error if the command fails.
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
