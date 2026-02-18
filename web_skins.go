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
//
package cgpcli

import (
	"fmt"
	"time"
)

// CreateClusterSkin creates a custom Cluster Skin.
//
// Parameters:
//   - skin: specifies the name of the new Skin.
//
// This method executes the CREATECLUSTERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateClusterSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("CREATECLUSTERSKIN", skin)
}

// CreateDomainSkin creates a custom Domain Skin.
//
// Parameters:
//   - domain: an optional the Domain name.
//   - skin: specifies the name of the new Skin. To create the unnamed Domain
//     Skin, specify an empty string.
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
//   - skin: specifies the name of the new Skin.
//
// This method executes the CREATESERVERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateServerSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("CREATESERVERSKIN", skin)
}

// DeleteClusterSkin deletes a custom Cluster Skin.
//
// Parameters:
//   - skin: specifies the name of the Skin to be deleted.
//
// This method executes the DELETECLUSTERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteClusterSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("DELETECLUSTERSKIN", skin)
}

// DeleteClusterSkinFile removes a file from a custom Cluster Skin.
//
// Parameters:
//   - skin: specifies the name of an existing Cluster Skin.
//   - file: specifies the Skin file name to be deleted.
//
// This method executes the STORECLUSTERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteClusterSkinFile(skin, file string) error {
	if skin == "" || file == "" {
		return fmt.Errorf("skin name and file name are required")
	}
	return cli.QueryNV("STORECLUSTERSKINFILE", skin, "FILE", file, "DELETE")
}

// DeleteDomainSkin deletes a custom Domain Skin.
//
// Parameters:
//   - domain: an optional the Domain name.
//   - skin: specifies the name of the Skin to be deleted. To delete the
//     unnamed Domain Skin, specify an empty string.
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
//   - domain: an optional the Domain name.
//   - skin: specifies the name of an existing Domain Skin.
//   - file: specifies the Skin file name to be deleted.
//
// This method executes the STOREDOMAINSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteDomainSkinFile(domain, skin, file string) error {
	if file == "" {
		return fmt.Errorf("file name is required")
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
//   - skin: specifies the name of the Skin to be deleted.
//
// This method executes the DELETESERVERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteServerSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("DELETESERVERSKIN", skin)
}

// DeleteServerSkinFile removes a file from a custom Server Skin.
//
// Parameters:
//   - skin: specifies the name of an existing Server Skin.
//   - file: specifies the Skin file name to be deleted.
//
// This method executes the STORESERVERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteServerSkinFile(skin, file string) error {
	if skin == "" || file == "" {
		return fmt.Errorf("skin name and file name are required")
	}
	return cli.QueryNV("STORESERVERSKINFILE", skin, "FILE", file, "DELETE")
}

// ListClusterSkinFiles lists files in a custom Cluster Skin.
//
// Parameters:
//   - skin: specifies the name of an existing Cluster Skin.
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
//   - domain: an optional the Domain name.
//   - skin: specifies the name of an existing Domain Skin.
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
//   - domain: an optional the Domain name.
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
//   - skin: specifies the name of an existing Server Skin.
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
//   - skin: specifies the name of a built-in Skin.
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
//   - skin: specifies the name of an existing Cluster Skin.
//   - file: specifies the name of an existing file in the specified Cluster Skin.
//
// This method executes the READCLUSTERSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadClusterSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" || file == "" {
		return nil, time.Time{}, fmt.Errorf("skin name and file name are required")
	}
	return cli.readSkinFile("READCLUSTERSKINFILE", skin, "FILE", file)
}

// ReadDomainSkinFile reads a file from a custom Domain Skin.
//
// Parameters:
//   - domain: an optional the Domain name.
//   - skin: specifies the name of an existing Domain Skin.
//   - file: specifies the name of an existing file in the Skin.
//
// This method executes the READDOMAINSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadDomainSkinFile(domain, skin, file string) ([]byte, time.Time, error) {
	if file == "" {
		return nil, time.Time{}, fmt.Errorf("file name is required")
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
		return nil, time.Time{}, fmt.Errorf("%s: expected array [data, time], got %T", cmd, res)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadServerSkinFile reads a file from a custom Server Skin.
//
// Parameters:
//   - skin: specifies the name of an existing Server Skin.
//   - file: specifies the name of an existing file in the Skin.
//
// This method executes the READSERVERSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadServerSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" || file == "" {
		return nil, time.Time{}, fmt.Errorf("skin name and file name are required")
	}

	const cmd = "READSERVERSKINFILE"
	res, err := cli.getSliceAny(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(res) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadStockSkinFile reads a file from a built-in (stock) Skin.
//
// Parameters:
//   - skin: specifies the name of an existing built-in Skin.
//   - file: specifies the name of an existing file in the Skin.
//
// This method executes the READSTOCKSKINFILE CLI command.
//
// Returns:
//   - []byte: the Skin file content.
//   - time.Time: the file modification date.
//   - error: an error if the command fails.
func (cli *Cli) ReadStockSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" || file == "" {
		return nil, time.Time{}, fmt.Errorf("skin name and file name are required")
	}

	const cmd = "READSTOCKSKINFILE"
	res, err := cli.getSliceAny(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(res) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// RenameClusterSkin renames a custom Cluster Skin.
//
// Parameters:
//   - oldSkin: specifies the name of an existing Skin.
//   - newSkin: specifies the new name for the Skin.
//
// This method executes the RENAMECLUSTERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameClusterSkin(oldSkin, newSkin string) error {
	if oldSkin == "" || newSkin == "" {
		return fmt.Errorf("old and new skin names are required")
	}
	return cli.QueryNV("RENAMECLUSTERSKIN", oldSkin, "INTO", newSkin)
}

// RenameDomainSkin renames a custom named Domain Skin.
//
// Parameters:
//   - domain: an optional the Domain name.
//   - oldSkin: specifies the name of an existing named Skin.
//   - newSkin: specifies the new name for the Skin.
//
// This method executes the RENAMEDOMAINSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameDomainSkin(domain, oldSkin, newSkin string) error {
	if oldSkin == "" || newSkin == "" {
		return fmt.Errorf("old and new skin names are required")
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
//   - oldSkin: specifies the name of an existing Skin.
//   - newSkin: specifies the new name for the Skin.
//
// This method executes the RENAMESERVERSKIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameServerSkin(oldSkin, newSkin string) error {
	if oldSkin == "" || newSkin == "" {
		return fmt.Errorf("old and new skin names are required")
	}
	return cli.QueryNV("RENAMESERVERSKIN", oldSkin, "INTO", newSkin)
}

// StoreClusterSkinFile stores a file into a custom Cluster Skin.
//
// Parameters:
//   - skin: specifies the name of an existing Cluster Skin.
//   - file: specifies the Skin file name.
//   - data: the file content.
//
// This method executes the STORECLUSTERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreClusterSkinFile(skin, file string, data []byte) error {
	if skin == "" || file == "" || data == nil {
		return fmt.Errorf("skin name, file name and data are required")
	}
	return cli.QueryNV("STORECLUSTERSKINFILE", skin, "FILE", file, "DATA", data)
}

// StoreDomainSkinFile stores a file into a custom Domain Skin.
//
// Parameters:
//   - domain: an optional the Domain name.
//   - skin: specifies the name of an existing Domain Skin.
//   - file: specifies the Skin file name.
//   - data: the file content.
//
// This method executes the STOREDOMAINSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreDomainSkinFile(domain, skin, file string, data []byte) error {
	if file == "" || data == nil {
		return fmt.Errorf("file name and data are required")
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
//   - skin: specifies the name of an existing Server Skin.
//   - file: specifies the Skin file name.
//   - data: the file content.
//
// This method executes the STORESERVERSKINFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreServerSkinFile(skin, file string, data []byte) error {
	if skin == "" || file == "" || data == nil {
		return fmt.Errorf("skin name, file name and data are required")
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
		return nil, time.Time{}, fmt.Errorf("%s: expected array [data, time], got %T", cmd, res)
	}

	data, ok1 := res[0].([]byte)
	ts, ok2 := res[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}
