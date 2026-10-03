// # Web Skins Administration
//
// The Web Skins section provides comprehensive management for CommuniGate Pro
// WebUser interface customizations. It allows developers to manipulate
// skins at three distinct hierarchy levels: Server, Cluster, and Domain.
//
// Key capabilities include:
//   - Hierarchy Management: creating, renaming, and deleting custom skins
//     at different scopes (Server, Cluster, or Domain-specific).
//   - File Operations: full CRUD support for individual skin files,
//     enabling programmatic updates to HTML templates, CSS, and images via
//     [Cli.StoreServerSkinFile] and related methods.
//   - Stock Skin Access: ability to read and list files from built-in
//     "stock" skins provided by the server for reference or inheritance.
//   - Metadata Tracking: retrieving file content along with its modification
//     timestamp ([time.Time]) for efficient caching and synchronization.
//   - Discovery: listing available skins and inspecting their file contents
//     at each hierarchy level.
package cgpcli

import (
	"fmt"
	"time"
)

// CreateClusterSkin creates a custom Cluster Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the CREATECLUSTERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateClusterSkin(skin string) error {
	if skin == "" {
		return ErrSkinRequired
	}
	return cli.QueryNV("CREATECLUSTERSKIN", skin)
}

// CreateDomainSkin creates a custom Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - skin: the Skin name. To create the unnamed Domain Skin, specify an empty string.
//
// This method executes the CREATEDOMAINSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateDomainSkin(domain, skin string) error {
	const cmd = "CREATEDOMAINSKIN"
	if domain != "" {
		return cli.QueryNV(cmd, domain, "SKIN", skin)
	} else {
		return cli.QueryNV(cmd, skin)
	}
}

// CreateServerSkin creates a custom Server Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the CREATESERVERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateServerSkin(skin string) error {
	if skin == "" {
		return ErrSkinRequired
	}
	return cli.QueryNV("CREATESERVERSKIN", skin)
}

// DeleteClusterSkin deletes a custom Cluster Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the DELETECLUSTERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteClusterSkin(skin string) error {
	if skin == "" {
		return ErrSkinRequired
	}
	return cli.QueryNV("DELETECLUSTERSKIN", skin)
}

// DeleteClusterSkinFile removes a file from a custom Cluster Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the STORECLUSTERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteClusterSkinFile(skin, file string) error {
	if skin == "" {
		return ErrSkinRequired
	}
	if file == "" {
		return ErrFileRequired
	}
	return cli.QueryNV("STORECLUSTERSKINFILE", skin, "FILE", file, "DELETE")
}

// DeleteDomainSkin deletes a custom Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - skin: the Skin name. To delete the unnamed Domain Skin, specify an empty string.
//
// This method executes the DELETEDOMAINSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteDomainSkin(domain, skin string) error {
	const cmd = "DELETEDOMAINSKIN"
	if domain != "" {
		return cli.QueryNV(cmd, domain, "SKIN", skin)
	} else {
		return cli.QueryNV(cmd, skin)
	}
}

// DeleteDomainSkinFile removes a file from a custom Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the STOREDOMAINSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteDomainSkinFile(domain, skin, file string) error {
	if file == "" {
		return ErrFileRequired
	}

	const cmd = "STOREDOMAINSKINFILE"
	if domain != "" {
		return cli.QueryNV(cmd, domain, "SKIN", skin, "FILE", file, "DELETE")
	} else {
		return cli.QueryNV(cmd, skin, "FILE", file, "DELETE")
	}
}

// DeleteServerSkin deletes a custom Server Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the DELETESERVERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteServerSkin(skin string) error {
	if skin == "" {
		return ErrSkinRequired
	}
	return cli.QueryNV("DELETESERVERSKIN", skin)
}

// DeleteServerSkinFile removes a file from a custom Server Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the STORESERVERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteServerSkinFile(skin, file string) error {
	if skin == "" {
		return ErrSkinRequired
	}
	if file == "" {
		return ErrFileRequired
	}
	return cli.QueryNV("STORESERVERSKINFILE", skin, "FILE", file, "DELETE")
}

// ListClusterSkinFiles lists files in a custom Cluster Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the LISTCLUSTERSKINFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with Skin file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListClusterSkinFiles(skin string) (map[string]*FileInfo, error) {
	raw, err := cli.getMapAny("LISTCLUSTERSKINFILES", skin)
	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ListClusterSkins lists custom Cluster Skins.
//
// This method executes the LISTCLUSTERSKINS CLI command.
//
// Returns:
//   - []string: an array with Skin names.
//   - error: an error if the command fails.
func (cli *Cli) ListClusterSkins() ([]string, error) {
	return cli.getSliceString("LISTCLUSTERSKINS")
}

// ListDomainSkinFiles lists files in a custom Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - skin: the Skin name.
//
// This method executes the LISTDOMAINSKINFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with Skin file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListDomainSkinFiles(domain, skin string) (map[string]*FileInfo, error) {
	const cmd = "LISTDOMAINSKINFILES"
	var raw map[string]any
	var err error

	if domain != "" {
		raw, err = cli.getMapAny(cmd, domain, "SKIN", skin)
	} else {
		raw, err = cli.getMapAny(cmd, skin)
	}

	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ListDomainSkins lists custom Domain Skins.
//
// Parameters:
//   - domain: an optional Domain name.
//
// This method executes the LISTDOMAINSKINS CLI command.
//
// Returns:
//   - []string: an array with Skin names.
//   - error: an error if the command fails.
func (cli *Cli) ListDomainSkins(domain string) ([]string, error) {
	return cli.getSliceString("LISTDOMAINSKINS", Atom(domain))
}

// ListServerSkinFiles lists files in a custom Server Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the LISTSERVERSKINFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with Skin file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListServerSkinFiles(skin string) (map[string]*FileInfo, error) {
	raw, err := cli.getMapAny("LISTSERVERSKINFILES", skin)
	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ListServerSkins lists custom Server Skins.
//
// This method executes the LISTSERVERSKINS CLI command.
//
// Returns:
//   - []string: an array with Skin names.
//   - error: an error if the command fails.
func (cli *Cli) ListServerSkins() ([]string, error) {
	return cli.getSliceString("LISTSERVERSKINS")
}

// ListStockSkinFiles lists files in a built-in (stock) Skin.
//
// Parameters:
//   - skin: the Skin name.
//
// This method executes the LISTSTOCKSKINFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with Skin file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListStockSkinFiles(skin string) (map[string]*FileInfo, error) {
	raw, err := cli.getMapAny("LISTSTOCKSKINFILES", skin)
	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ReadClusterSkinFile reads a file from a custom Cluster Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the READCLUSTERSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadClusterSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" {
		return nil, time.Time{}, ErrSkinRequired
	}
	if file == "" {
		return nil, time.Time{}, ErrFileRequired
	}
	return cli.readSkinFile("READCLUSTERSKINFILE", skin, "FILE", file)
}

// ReadDomainSkinFile reads a file from a custom Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the READDOMAINSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadDomainSkinFile(domain, skin, file string) ([]byte, time.Time, error) {
	if file == "" {
		return nil, time.Time{}, ErrFileRequired
	}

	const cmd = "READDOMAINSKINFILE"
	var res []any
	var err error

	if domain != "" {
		res, err = cli.getSliceAny(cmd, domain, "SKIN", skin, "FILE", file)
	} else {
		res, err = cli.getSliceAny(cmd, skin, "FILE", file)
	}

	if err != nil {
		return nil, time.Time{}, err
	}

	if len(res) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedResponseFormat)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedType)
	}

	return data, ts, nil
}

// ReadServerSkinFile reads a file from a custom Server Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the READSERVERSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadServerSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" {
		return nil, time.Time{}, ErrSkinRequired
	}
	if file == "" {
		return nil, time.Time{}, ErrFileRequired
	}

	const cmd = "READSERVERSKINFILE"
	res, err := cli.getSliceAny(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(res) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedResponseFormat)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedType)
	}

	return data, ts, nil
}

// ReadStockSkinFile reads a file from a built-in (stock) Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//
// This method executes the READSTOCKSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadStockSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" {
		return nil, time.Time{}, ErrSkinRequired
	}
	if file == "" {
		return nil, time.Time{}, ErrFileRequired
	}

	const cmd = "READSTOCKSKINFILE"
	res, err := cli.getSliceAny(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(res) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedResponseFormat)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedType)
	}

	return data, ts, nil
}

// ReloadPluginSkins reloads Skins files of all domains, or a specified domain,
// from the installed plugins according to the current settings.
//
// Parameters:
//   - domain: an optional Domain name. If empty, Skins files of all domains are reloaded.
//
// This method executes the RELOADPLUGINSKINS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ReloadPluginSkins(domain string) error {
	const cmd = "RELOADPLUGINSKINS"
	if domain == "" {
		return cli.QueryNV(cmd)
	} else {
		return cli.QueryNV(cmd, "DOMAIN", domain)
	}
}

// RenameClusterSkin renames a custom Cluster Skin.
//
// Parameters:
//   - oldSkin: the Skin name.
//   - newSkin: the new Skin name.
//
// This method executes the RENAMECLUSTERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameClusterSkin(oldSkin, newSkin string) error {
	if oldSkin == "" {
		return ErrOldNameRequired
	}
	if newSkin == "" {
		return ErrNewNameRequired
	}
	return cli.QueryNV("RENAMECLUSTERSKIN", oldSkin, "INTO", newSkin)
}

// RenameDomainSkin renames a custom named Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - oldSkin: the Skin name.
//   - newSkin: the new Skin name.
//
// This method executes the RENAMEDOMAINSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameDomainSkin(domain, oldSkin, newSkin string) error {
	if oldSkin == "" {
		return ErrOldNameRequired
	}
	if newSkin == "" {
		return ErrNewNameRequired
	}

	const cmd = "RENAMEDOMAINSKIN"
	if domain != "" {
		return cli.QueryNV(cmd, domain, "SKIN", oldSkin, "INTO", newSkin)
	} else {
		return cli.QueryNV(cmd, oldSkin, "INTO", newSkin)
	}
}

// RenameServerSkin renames a custom Server Skin.
//
// Parameters:
//   - oldSkin: the Skin name.
//   - newSkin: the new Skin name.
//
// This method executes the RENAMESERVERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameServerSkin(oldSkin, newSkin string) error {
	if oldSkin == "" {
		return ErrOldNameRequired
	}
	if newSkin == "" {
		return ErrNewNameRequired
	}
	return cli.QueryNV("RENAMESERVERSKIN", oldSkin, "INTO", newSkin)
}

// StoreClusterSkinFile stores a file into a custom Cluster Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//   - data: the file content.
//
// This method executes the STORECLUSTERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreClusterSkinFile(skin, file string, data []byte) error {
	if skin == "" {
		return ErrSkinRequired
	}
	if file == "" {
		return ErrFileRequired
	}
	if data == nil {
		return ErrDataRequired
	}
	return cli.QueryNV("STORECLUSTERSKINFILE", skin, "FILE", file, "DATA", data)
}

// StoreDomainSkinFile stores a file into a custom Domain Skin.
//
// Parameters:
//   - domain: an optional Domain name.
//   - skin: the Skin name.
//   - file: the file name.
//   - data: the file content.
//
// This method executes the STOREDOMAINSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreDomainSkinFile(domain, skin, file string, data []byte) error {
	if file == "" {
		return ErrFileRequired
	}
	if data == nil {
		return ErrDataRequired
	}

	const cmd = "STOREDOMAINSKINFILE"
	if domain != "" {
		return cli.QueryNV(cmd, domain, "SKIN", skin, "FILE", file, "DATA", data)
	} else {
		return cli.QueryNV(cmd, skin, "FILE", file, "DATA", data)
	}
}

// StoreServerSkinFile stores a file into a custom Server Skin.
//
// Parameters:
//   - skin: the Skin name.
//   - file: the file name.
//   - data: the file content.
//
// This method executes the STORESERVERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreServerSkinFile(skin, file string, data []byte) error {
	if skin == "" {
		return ErrSkinRequired
	}
	if file == "" {
		return ErrFileRequired
	}
	if data == nil {
		return ErrDataRequired
	}
	return cli.QueryNV("STORESERVERSKINFILE", skin, "FILE", file, "DATA", data)
}

// readSkinFile — внутренний хелпер для чтения файлов скинов.
func (cli *Cli) readSkinFile(cmd string, args ...any) ([]byte, time.Time, error) {
	res, err := cli.getSliceAny(cmd, args...)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(res) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedResponseFormat)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: %w", cmd, ErrUnexpectedType)
	}

	return data, ts, nil
}
