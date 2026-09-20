// # Full-Text Search
//
// The FTS section provides comprehensive methods for managing Search Indexes
// across Domains, Accounts, and Groups within CommuniGate Pro. It supports
// triggering asynchronous index operations, tracking indexing states, querying
// data sizes, performing raw index lookups, and clearing indexes.
//
// Key capabilities include:
//   - Domain-wide and Group-based multi-user asynchronous index scheduling.
//   - Account-level immediate full re-indexing and space consumption sizing.
//   - Granular state inspection (none, waiting, indexing, ready).
//   - Strict validation enforcement based on CommuniGate Pro server constraints.
package cgpcli

import (
	"fmt"
	"strings"
)

// BuildAccountSearchIndex triggers the creation or rebuilding of an Account's Search Index.
//
// Parameters:
//   - account: the full Account name.
//
// Note: If the specified User already has an index, it will be completely rebuilt from scratch.
// If indexing is currently disabled in the account or domain settings, this command will explicitly
// force-enable it for this specific account.
//
// This method executes the BUILDACCOUNTSEARCHINDEX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) BuildAccountSearchIndex(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("BUILDACCOUNTSEARCHINDEX", Atom(account))
}

// BuildDomainSearchIndex schedules or forces the creation of a Search Index for all Domain Accounts.
//
// Parameters:
//   - domain: the Domain name.
//   - enableForAll: if true, forces the enablement of Search Index in the settings
//     of each account within the domain (except for users who explicitly disabled it in their profile).
//     If false, indexing is performed only for users who already have it enabled.
//
// Note: Building Search Index for large mailboxes (>1 GB) is a heavy asynchronous operation.
// The server may return code 500 ("Processing domain is already in progress") if an indexing
// task is currently running.
//
// This method executes the BUILDDOMAINSEARCHINDEX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) BuildDomainSearchIndex(domain string, enableForAll bool) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "BUILDDOMAINSEARCHINDEX"
	if enableForAll {
		return cli.QueryNV(cmd, Atom(domain), "ENABLEFORALL")
	} else {
		return cli.QueryNV(cmd, Atom(domain))
	}
}

// BuildGroupSearchIndex processes and builds the Search Index for all members belonging to a Group.
//
// Parameters:
//   - group: the full Group name.
//   - enableForAll: if true, force-enables Search Index inside the profiles of all users
//     contained in the group. If false, indexes only members who already have it active.
//
// This method executes the BUILDGROUPSEARCHINDEX CLI command.
//
// Returns:
//   - error: an error if the group is missing (code 552) or command processing fails.
func (cli *Cli) BuildGroupSearchIndex(group string, enableForAll bool) error {
	if group == "" {
		return fmt.Errorf("group name is required")
	}

	group0, domain, found := strings.Cut(group, "@")
	if !found {
		return fmt.Errorf("group name must contains domain name")
	}

	const cmd = "BUILDGROUPSEARCHINDEX"
	if enableForAll {
		return cli.QueryNV(cmd, Atom(domain), Atom(group0), "ENABLEFORALL")
	} else {
		return cli.QueryNV(cmd, Atom(domain), Atom(group0))
	}
}

// GetAccountSearchIndexSize retrieves the actual disk space allocated for the Account's Search Index.
//
// Parameters:
//   - account: the full Account name.
//
// This method executes the GETACCOUNTSEARCHINDEXSIZE CLI command.
//
// Returns:
//   - int64: the size of the index (typically in bytes).
//   - error: an error if the command fails.
func (cli *Cli) GetAccountSearchIndexSize(account string) (int64, error) {
	if account == "" {
		return 0, fmt.Errorf("account name is required")
	}
	return cli.getInt("GETACCOUNTSEARCHINDEXSIZE", Atom(account))
}

// GetAccountSearchIndexState retrieves the current indexing operational status for a specific Account.
//
// Parameters:
//   - account: the full Account name.
//
// This method executes the GETACCOUNTSEARCHINDEXSTATE CLI command.
//
// Returns:
//   - string: the status value ("none", "waiting", "indexing", or "ready").
//   - error: an error if the command fails.
func (cli *Cli) GetAccountSearchIndexState(account string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}
	return cli.getString("GETACCOUNTSEARCHINDEXSTATE", Atom(account))
}

// GetDomainSearchIndexState retrieves the current status of a Domain's Search Index.
//
// Parameters:
//   - domain: the Domain name.
//
// This method executes the GETDOMAINSEARCHINDEXSTATE CLI command.
//
// Returns:
//   - string: the index state string ("none", "waiting", "indexing", or "ready").
//   - error: an error if the command fails.
func (cli *Cli) GetDomainSearchIndexState(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain name is required")
	}
	return cli.getString("GETDOMAINSEARCHINDEXSTATE", Atom(domain))
}

// RemoveAccountSearchIndex purges the localized Search Index files of a specific Account.
//
// Parameters:
//   - account: the full Account name.
//
// Important: The Search Index feature must be explicitly disabled in either the User's
// individual settings or the inherited Domain settings prior to execution, or the server
// will reject the command with code 500.
//
// This method executes the REMOVEACCOUNTSEARCHINDEX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RemoveAccountSearchIndex(account string) error {
	if account == "" {
		return fmt.Errorf("account name is required")
	}
	return cli.QueryNV("REMOVEACCOUNTSEARCHINDEX", Atom(account))
}

// RemoveDomainSearchIndex clears the Search Index for all Accounts within the specified Domain.
//
// Parameters:
//   - domain: the Domain name.
//
// Important: For this cleanup operation to proceed, Search Index must be turned off
// in the Domain's User Defaults. If Search Index remains explicitly enabled in a particular
// user's profile, the cleanup will be skipped for that specific account.
//
// This method executes the REMOVEDOMAINSEARCHINDEX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RemoveDomainSearchIndex(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}
	return cli.QueryNV("REMOVEDOMAINSEARCHINDEX", Atom(domain))
}

// SearchInIndex executes a direct structural text lookup inside a Account's Search Index.
//
// Parameters:
//   - account: the full Account name.
//   - search: the query criteria adhering to CommuniGate Pro Index Query Syntax.
//   - utf8: if true, forces the server to automatically convert underlying UTF7-IMAP results
//     into native UTF-8 encoding before returning.
//   - options: a key-value dictionary passing optional granular search parameters.
//
// This method executes the SEARCHININDEX CLI command.
//
// Returns:
//   - map[string]any: a structural dictionary containing search hits modeled after XIMSS specifications.
//   - error: an error if the syntax is invalid or processing fails.
func (cli *Cli) SearchInIndex(account string, search string, utf8 bool, options map[string]any) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	if search == "" {
		return nil, fmt.Errorf("search string is required")
	}

	args := make([]any, 0, 6)
	args = append(args, Atom(account), "SEARCHSTRING", search)

	if utf8 {
		args = append(args, "UTF8")
	}
	if len(options) > 0 {
		args = append(args, "OPTIONS", options)
	}

	return cli.getMapAny("SEARCHININDEX", args...)
}
