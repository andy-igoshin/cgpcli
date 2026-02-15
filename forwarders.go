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

// CreateForwarder creates new Forwarders.
//
// Parameters:
//   - forwarder: This parameter specifies the name for the new Forwarder. The name can contain the @ symbol followed by the Domain name, in this case the Forwarder is created in the specified Domain. If the Domain name is not specified, the command applies to the administrator Domain.
//   - address: This parameter specifies the E-mail address the Forwarder should reroute E-mail messages and Signals to.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) CreateForwarder(forwarder, address string) error {
	if forwarder == "" || address == "" {
		return fmt.Errorf("forwarder name and address are required")
	}
	return cli.QueryNV("CREATEFORWARDER", forwarder, "TO", address)
}

// DeleteForwarder removes Forwarders.
//
// Parameters:
//   - forwarder: This parameter specifies the name of an existing Forwarder. The name can include the Domain name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteForwarder(forwarder string) error {
	if forwarder == "" {
		return fmt.Errorf("forwarder name is required")
	}
	return cli.QueryNV("DELETEFORWARDER", forwarder)
}

// FindForwarders finds all Forwarders pointing to the specified address.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - address: This parameter specifies an E-mail address to look for.
//
// Returns:
//   - []string: An array with the found Forwarder names.
//   - error: An error if the command fails.
func (cli *Cli) FindForwarders(domain, address string) ([]string, error) {
	if domain == "" || address == "" {
		return nil, fmt.Errorf("domain and target address are required")
	}
	return cli.getSliceString("FINDFORWARDERS", domain, "TO", address)
}

// GetForwarder retrieves the Forwarder address.
//
// Parameters:
//   - forwarder: This parameter specifies the name of an existing Forwarder. The name can include the Domain name.
//
// Returns:
//   - string: A string with the E-mail address this Forwarder reroutes all E-mail messages and Signals to.
//   - error: An error if the command fails.
func (cli *Cli) GetForwarder(forwarder string) (string, error) {
	if forwarder == "" {
		return "", fmt.Errorf("forwarder name is required")
	}
	return cli.getString("GETFORWARDER", forwarder)
}

// ListForwarders retrieves the list of all Forwarders in the Domain.
//
// Parameters:
//   - domain: (Optional) This parameter specifies the Domain name. If the Domain name is not specified, the command applies to the administrator Domain.
//
// Returns:
//   - []string: An array with the names of all Forwarders in the specified (or default) Domain.
//   - error: An error if the command fails.
func (cli *Cli) ListForwarders(domain string) ([]string, error) {
	return cli.getSliceString("LISTFORWARDERS", domain)
}

// RenameForwarder renames Forwarders.
//
// Parameters:
//   - oldName: This parameter specifies the name of an existing Forwarder. The name can include the Domain name.
//   - newName: This parameter specifies the new Forwarder name. The name can include the Domain name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) RenameForwarder(oldName, newName string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new names are required")
	}
	return cli.QueryNV("RENAMEFORWARDER", oldName, "INTO", newName)
}
