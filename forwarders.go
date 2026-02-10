// # Forwarder Administration
//
// The Forwarders section provides methods for managing standalone mail redirection
// objects within CommuniGate Pro. Unlike account-level rules, Forwarders act as
// dedicated routing entities that redirect incoming messages to target address.
//
// Key capabilities include:
//   - Lifecycle Management: Creating, renaming, and deleting forwarders.
//   - Discovery: Listing all forwarders in a domain or performing reverse
//     lookups with [Cli.FindForwarders] to identify which objects point
//     to a specific target address.
//   - Inspection: Retrieving the target address of a specific forwarder.
//
// This section is particularly useful for managing system-wide redirections
// and migrating legacy mail aliases to dedicated routing objects.
package cgpcli

import (
	"fmt"
)

func (cli *Cli) CreateForwarder(forwarder, target string) error {
	if forwarder == "" || target == "" {
		return fmt.Errorf("forwarder name and target address are required")
	}
	return cli.QueryNV("CREATEFORWARDER", forwarder, "TO", target)
}

func (cli *Cli) DeleteForwarder(forwarder string) error {
	if forwarder == "" {
		return fmt.Errorf("forwarder name is required")
	}
	return cli.QueryNV("DELETEFORWARDER", forwarder)
}

func (cli *Cli) FindForwarders(domain, target string) ([]string, error) {
	if domain == "" || target == "" {
		return nil, fmt.Errorf("domain and target address are required")
	}
	return cli.getSliceString("FINDFORWARDERS", domain, "TO", target)
}

func (cli *Cli) GetForwarder(forwarder string) (string, error) {
	if forwarder == "" {
		return "", fmt.Errorf("forwarder name is required")
	}
	return cli.getString("GETFORWARDER", forwarder)
}

func (cli *Cli) ListForwarders(domain string) ([]string, error) {
	return cli.getSliceString("LISTFORWARDERS", domain)
}

func (cli *Cli) RenameForwarder(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old name and new name are required")
	}
	return cli.QueryNV("RENAMEFORWARDER", oldName, "INTO", newName)
}
