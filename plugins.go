// # Plugin Administration
//
// The Plugin Administration section provides functionality for managing plugins
// in CommuniGate Pro. It allows uploading, removing plugins, and managing
// Domain-level plugin settings.
//
// Key capabilities include:
//   - Upload: using [Cli.UploadPluginFile] to upload a plugin binary to the
//     server. The plugin is installed unless the DRYRUN flag is passed.
//   - Removal: using [Cli.RemovePlugin] to remove an installed plugin by its ID.
//   - Retrieval: using [Cli.GetDomainPluginsSettings] to fetch the current
//     Domain Plugins settings dictionary.
//   - Configuration: using [Cli.SetDomainPluginsSettings] to replace the entire
//     Domain Plugins settings dictionary with a new one.
package cgpcli

import (
	"fmt"
)

// GetDomainPluginsSettings retrieves the Domain Plugins settings.
//
// Parameters:
//   - domain: an optional Domain name.
//
// This method executes the GETDOMAINPLUGINSSETTINGS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Domain Plugins settings.
//   - error: an error if the command fails.
func (cli *Cli) GetDomainPluginsSettings(domain string) (map[string]any, error) {
	return cli.getMapAny("GETDOMAINPLUGINSSETTINGS", Atom(domain))
}

// RemovePlugin removes an installed plugin with the specified ID.
//
// Parameters:
//   - id: ID of the plugin to remove.
//
// This method executes the REMOVEPLUGIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RemovePlugin(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return cli.QueryNV("REMOVEPLUGIN", id)
}

// SetDomainPluginsSettings replaces the Domain Plugins settings.
//
// Parameters:
//   - domain: an optional Domain name.
//   - settings: a dictionary used to replace the Domain Plugins settings.
//
// This method executes the SETDOMAINPLUGINSSETTINGS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDomainPluginsSettings(domain string, settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDOMAINPLUGINSSETTINGS", Atom(domain), settings)
}

// UploadPluginFile uploads a plugin binary to the server.
//
// Parameters:
//   - content: data block containing the plugin binary representation.
//   - dryrun: if true, no installation will proceed. Only plugin validation
//     will be performed.
//
// This method executes the UPLOADPLUGINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UploadPluginFile(content []byte, dryrun bool) error {
	if len(content) == 0 {
		return fmt.Errorf("content is required")
	}

	const cmd = "UPLOADPLUGINFILE"
	if dryrun {
		return cli.QueryNV(cmd, content, "DRYRUN")
	} else {
		return cli.QueryNV(cmd, content)
	}
}
