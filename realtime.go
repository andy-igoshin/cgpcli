// # Real-Time Application Administration
//
// The PBX section provides functionality for managing CommuniGate Pro
// Real-Time Applications (RTA). It allows administrators to upload,
// retrieve, and organize files for telephony services, automated
// attendants, and custom signal-processing applications.
//
// Key capabilities include:
//   - Multi-Level Management: handling PBX files at different hierarchy
//     levels: [Stock] (read-only system files), [Server], [Cluster],
//     and [Domain].
//   - File Operations: storing, reading, and deleting RTA assets such
//     as application scripts and localized media resources (audio prompts).
//   - Localization Support: managing national subsets of the RTA
//     environment through language-specific directories.
//   - Discovery: listing available PBX files at various levels to
//     inspect the current state of custom signaling logic.
package cgpcli

import (
	"fmt"
)

// CreateClusterPBX creates the cluster-wide Real-Time Application Environment or its national subset.
//
// Parameters:
//   - language: a national subset name.
//
// This method executes the CREATECLUSTERPBX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateClusterPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("CREATECLUSTERPBX", language)
}

// CreateDomainPBX creates the Domain Real-Time Application Environment or its national subset.
//
// Parameters:
//   - domain: the Domain name.
//   - language: an optional national subset name.
//
// This method executes the CREATEDOMAINPBX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateDomainPBX(domain, language string) error {
	if domain == "" {
		return fmt.Errorf("domain is required")
	}

	const cmd = "CREATEDOMAINPBX"
	if language != "" {
		return cli.QueryNV(cmd, domain, "FILE", language)
	} else {
		return cli.QueryNV(cmd, domain)
	}
}

// CreateServerPBX creates the Server-wide Real-Time Application Environment or its national subset.
//
// Parameters:
//   - language: a national subset name.
//
// This method executes the CREATESERVERPBX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateServerPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("CREATESERVERPBX", language)
}

// DeleteClusterPBX removes a national subset from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - language: a national subset name.
//
// This method executes the DELETECLUSTERPBX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteClusterPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("DELETECLUSTERPBX", language)
}

// DeleteClusterPBXFile removes a specific file from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the STORECLUSTERPBXFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteClusterPBXFile(file string) error {
	if file == "" {
		return fmt.Errorf("file is required")
	}
	return cli.QueryNV("STORECLUSTERPBXFILE", file, "DELETE")
}

// DeleteDomainPBX removes a national subset from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: the Domain name.
//   - language: a national subset name.
//
// This method executes the DELETEDOMAINPBX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteDomainPBX(domain, language string) error {
	if domain == "" || language == "" {
		return fmt.Errorf("domain and language are required")
	}
	return cli.QueryNV("DELETEDOMAINPBX", domain, "FILE", language)
}

// DeleteDomainPBXFile removes a specific file from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: the Domain name.
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the STOREDOMAINPBXFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteDomainPBXFile(domain, file string) error {
	if domain == "" || file == "" {
		return fmt.Errorf("domain and file are required")
	}
	return cli.QueryNV("STOREDOMAINPBXFILE", domain, "FILE", file, "DELETE")
}

// DeleteServerPBX removes a national subset from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - language: a national subset name.
//
// This method executes the DELETESERVERPBX CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteServerPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("DELETESERVERPBX", language)
}

// DeleteServerPBXFile removes a specific file from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the STORESERVERPBXFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteServerPBXFile(file string) error {
	if file == "" {
		return fmt.Errorf("file is required")
	}
	return cli.QueryNV("STORESERVERPBXFILE", file, "DELETE")
}

// ListClusterPBXFiles lists files in the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - language: an optional national subset name.
//
// This method executes the LISTCLUSTERPBXFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListClusterPBXFiles(language string) (map[string]*FileInfo, error) {
	raw, err := cli.getMapAny("LISTCLUSTERPBXFILES", Atom(language))
	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ListDomainPBXFiles lists files in the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: an optional Domain name. If not specified, the administrator Domain is used.
//   - language: an optional national subset name.
//
// This method executes the LISTDOMAINPBXFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListDomainPBXFiles(domain, language string) (map[string]*FileInfo, error) {
	var raw map[string]any
	var err error
	if language != "" {
		raw, err = cli.getMapAny("LISTDOMAINPBXFILES", Atom(domain), "FILE", language)
	} else {
		raw, err = cli.getMapAny("LISTDOMAINPBXFILES", Atom(domain))
	}

	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ListServerPBXFiles lists files in the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - language: an optional national subset name.
//
// This method executes the LISTSERVERPBXFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListServerPBXFiles(language string) (map[string]*FileInfo, error) {
	raw, err := cli.getMapAny("LISTSERVERPBXFILES", Atom(language))
	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ListStockPBXFiles lists files in the stock (built-in) Real-Time Application Environment.
//
// Parameters:
//   - language: an optional national subset name.
//
// This method executes the LISTSTOCKPBXFILES CLI command.
//
// Returns:
//   - map[string]*[FileInfo]: a dictionary with file names as keys.
//   - error: an error if the command fails.
func (cli *Cli) ListStockPBXFiles(language string) (map[string]*FileInfo, error) {
	raw, err := cli.getMapAny("LISTSTOCKPBXFILES", Atom(language))
	if err != nil {
		return nil, err
	}
	return parseFileInfoMap(raw), nil
}

// ReadClusterPBXFile reads a file from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the READCLUSTERPBXFILE CLI command.
//
// Returns:
//   - []byte: a datablock with the file contents.
//   - error: an error if the command fails.
func (cli *Cli) ReadClusterPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	return cli.getSliceByte("READCLUSTERPBXFILE", file)
}

// ReadDomainPBXFile reads a file from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: the Domain name.
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the READDOMAINPBXFILE CLI command.
//
// Returns:
//   - []byte: a datablock with the file contents.
//   - error: an error if the command fails.
func (cli *Cli) ReadDomainPBXFile(domain, file string) ([]byte, error) {
	if domain == "" || file == "" {
		return nil, fmt.Errorf("domain and file are required")
	}
	return cli.getSliceByte("READDOMAINPBXFILE", domain, "FILE", file)
}

// ReadServerPBXFile reads a file from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the READSERVERPBXFILE CLI command.
//
// Returns:
//   - []byte: a datablock with the file contents.
//   - error: an error if the command fails.
func (cli *Cli) ReadServerPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	return cli.getSliceByte("READSERVERPBXFILE", file)
}

// ReadStockPBXFile reads a file from the stock (built-in) Real-Time Application Environment.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//
// This method executes the READSTOCKPBXFILE CLI command.
//
// Returns:
//   - []byte: a datablock with the file contents.
//   - error: an error if the command fails.
func (cli *Cli) ReadStockPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	return cli.getSliceByte("READSTOCKPBXFILE", file)
}

// StoreClusterPBXFile stores a file into the cluster-wide Real-Time Application Environment.
// Available in the Dynamic Cluster only.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//   - data: a datablock with the file contents.
//
// Note: If the environment contains a file with the specified name, the old file is replaced.
// The file is automatically removed from the Environment cache on all cluster members.
//
// This method executes the STORECLUSTERPBXFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreClusterPBXFile(file string, data []byte) error {
	if file == "" || data == nil {
		return fmt.Errorf("file and data are required")
	}
	return cli.QueryNV("STORECLUSTERPBXFILE", file, "DATA", data)
}

// StoreDomainPBXFile stores a file into the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: the Domain name.
//   - file: the file name. For national subsets, specify "language/fileName".
//   - data: a datablock with the file contents.
//
// Note: If the environment contains a file with the specified name, the old file is replaced.
// The file is automatically removed from the Environment cache.
//
// This method executes the STOREDOMAINPBXFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreDomainPBXFile(domain, file string, data []byte) error {
	if domain == "" || file == "" || data == nil {
		return fmt.Errorf("domain, file and data are required")
	}
	return cli.QueryNV("STOREDOMAINPBXFILE", domain, "FILE", file, "DATA", data)
}

// StoreServerPBXFile stores a file into the Server-wide Real-Time Application Environment.
// Available for System Administrators only.
//
// Parameters:
//   - file: the file name. For national subsets, specify "language/fileName".
//   - data: a datablock with the file contents.
//
// Note: If the environment contains a file with the specified name, the old file is replaced.
// The file is automatically removed from the Environment cache.
//
// This method executes the STORESERVERPBXFILE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreServerPBXFile(file string, data []byte) error {
	if file == "" || data == nil {
		return fmt.Errorf("file and data are required")
	}
	return cli.QueryNV("STORESERVERPBXFILE", file, "DATA", data)
}
