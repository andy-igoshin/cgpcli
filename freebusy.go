// # Free/Busy
//
// The Free/Busy section provides functionality for managing Domain-level Free/Busy
// settings in CommuniGate Pro. It allows retrieval and modification of the Free/Busy
// configuration that controls availability information sharing across the domain.
//
// Key capabilities include:
//   - Retrieval: using [Cli.GetDomainFreeBusy] to fetch the current Free/Busy
//     settings dictionary for a domain.
//   - Configuration: using [Cli.SetDomainFreeBusy] to replace the entire Free/Busy
//     settings dictionary with a new one.
package cgpcli

import "fmt"

// GetDomainFreeBusy retrieves the Domain Free/Busy settings.
//
// Parameters:
//   - domain: an optional Domain name.
//
// This method executes the GETDOMAINFREEBUSY CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Domain Free/Busy settings.
//   - error: an error if the command fails.
func (cli *Cli) GetDomainFreeBusy(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINFREEBUSY", Atom(domain))
}

// SetDomainFreeBusy replaces the Domain Free/Busy settings.
//
// Parameters:
//   - domain: an optional Domain name.
//   - settings: a dictionary used to replace the Domain Free/Busy settings.
//
// This method executes the SETDOMAINFREEBUSY CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDomainFreeBusy(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary are required")
	}
	return cli.QueryNV("SETDOMAINFREEBUSY", Atom(domain), settings)
}
