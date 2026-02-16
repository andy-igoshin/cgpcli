// # Real-Time Application Administration
//
// The PBX section provides functionality for managing CommuniGate Pro
// Real-Time Applications (RTA). It allows administrators to upload,
// retrieve, and organize files for telephony services, automated
// attendants, and custom signal-processing applications.
//
// Key capabilities include:
//   - Multi-Level Management: Handling PBX files at different hierarchy
//     levels: [Stock] (read-only system files), [Server], [Cluster],
//     and [Domain].
//   - File Operations: Storing, reading, and deleting RTA assets such
//     as application scripts and localized media resources (audio prompts).
//   - Localization Support: Managing national subsets of the RTA
//     environment through language-specific directories.
//   - Discovery: Listing available PBX files at various levels to
//     inspect the current state of custom signaling logic.
//
// Use these methods to deploy custom IVR solutions, manage corporate
// voice prompts, and automate the provisioning of telephony environments
// across multiple domains.
package cgpcli

import (
	"fmt"
	"time"
)

// PBXFileInfo represents metadata for a file in the Real-Time Application environment.
type PBXFileInfo struct {
	Name     string    // File name (as returned by the server)
	Size     int64     // File size in bytes (STFileSize)
	Created  time.Time // File creation timestamp (STCreated)
	Modified time.Time // Last modification timestamp (STModified)
}

// CreateClusterPBX creates the cluster-wide Real-Time Application Environment or its national subset.
//
// Parameters:
//   - language: This parameter specifies a national subset name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) CreateClusterPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("CREATECLUSTERPBX", language)
}

// CreateDomainPBX creates the Domain Real-Time Application Environment or its national subset.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) CreateDomainPBX(domain, language string) (err error) {
	if domain == "" {
		return fmt.Errorf("domain is required")
	}

	const cmd = "CREATEDOMAINPBX"
	if language != "" {
		err = cli.QueryNV(cmd, domain, "FILE", language)
	} else {
		err = cli.QueryNV(cmd, domain)
	}
	return
}

// CreateServerPBX creates the Server-wide Real-Time Application Environment or its national subset.
//
// Parameters:
//   - language: This parameter specifies a national subset name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) CreateServerPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("CREATESERVERPBX", language)
}

// DeleteClusterPBX removes a national subset from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - language: This parameter specifies a national subset name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteClusterPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("DELETECLUSTERPBX", language)
}

// DeleteClusterPBXFile removes a specific file from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteClusterPBXFile(file string) error {
	if file == "" {
		return fmt.Errorf("file is required")
	}
	return cli.QueryNV("STORECLUSTERPBXFILE", file, "DELETE")
}

// DeleteDomainPBX removes a national subset from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - language: This parameter specifies a national subset name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteDomainPBX(domain, language string) error {
	if domain == "" || language == "" {
		return fmt.Errorf("domain and language are required")
	}
	return cli.QueryNV("DELETEDOMAINPBX", domain, "FILE", language)
}

// DeleteDomainPBXFile removes a specific file from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteDomainPBXFile(domain, file string) error {
	if domain == "" || file == "" {
		return fmt.Errorf("domain and file are required")
	}
	return cli.QueryNV("STOREDOMAINPBXFILE", domain, "FILE", file, "DELETE")
}

// DeleteServerPBX removes a national subset from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - language: This parameter specifies a national subset name.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteServerPBX(language string) error {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	return cli.QueryNV("DELETESERVERPBX", language)
}

// DeleteServerPBXFile removes a specific file from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) DeleteServerPBXFile(file string) error {
	if file == "" {
		return fmt.Errorf("file is required")
	}
	return cli.QueryNV("STORESERVERPBXFILE", file, "DELETE")
}

// ListClusterPBXFiles lists files in the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]*[PBXFileInfo]: A dictionary with file names as keys.
//   - error: An error if the command fails.
func (cli *Cli) ListClusterPBXFiles(language string) (map[string]*PBXFileInfo, error) {
	return cli.listPBXFiles("LISTCLUSTERPBXFILES", Atom(language))
}

// ListDomainPBXFiles lists files in the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: (Optional) This optional parameter specifies the Domain name. If not specified, the administrator Domain is used.
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]*[PBXFileInfo]: A dictionary with file names as keys.
//   - error: An error if the command fails.
func (cli *Cli) ListDomainPBXFiles(domain, language string) (map[string]*PBXFileInfo, error) {
	if language != "" {
		return cli.listPBXFiles("LISTDOMAINPBXFILES", Atom(domain), "FILE", language)
	} else {
		return cli.listPBXFiles("LISTDOMAINPBXFILES", Atom(domain))
	}
}

// ListServerPBXFiles lists files in the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]*[PBXFileInfo]: A dictionary with file names as keys.
//   - error: An error if the command fails.
func (cli *Cli) ListServerPBXFiles(language string) (map[string]*PBXFileInfo, error) {
	return cli.listPBXFiles("LISTSERVERPBXFILES", Atom(language))
}

// ListStockPBXFiles lists files in the stock (built-in) Real-Time Application Environment.
//
// Parameters:
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]*[PBXFileInfo]: A dictionary with file names as keys.
//   - error: An error if the command fails.
func (cli *Cli) ListStockPBXFiles(language string) (map[string]*PBXFileInfo, error) {
	return cli.listPBXFiles("LISTSTOCKPBXFILES", Atom(language))
}

// ReadClusterPBXFile reads a file from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - []byte: A datablock with the file contents.
//   - error: An error if the command fails.
func (cli *Cli) ReadClusterPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	return cli.getSliceByte("READCLUSTERPBXFILE", file)
}

// ReadDomainPBXFile reads a file from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - []byte: A datablock with the file contents.
//   - error: An error if the command fails.
func (cli *Cli) ReadDomainPBXFile(domain, file string) ([]byte, error) {
	if domain == "" || file == "" {
		return nil, fmt.Errorf("domain and file are required")
	}
	return cli.getSliceByte("READDOMAINPBXFILE", domain, "FILE", file)
}

// ReadServerPBXFile reads a file from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - []byte: A datablock with the file contents.
//   - error: An error if the command fails.
func (cli *Cli) ReadServerPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	return cli.getSliceByte("READSERVERPBXFILE", file)
}

// ReadStockPBXFile reads a file from the stock (built-in) Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//
// Returns:
//   - []byte: A datablock with the file contents.
//   - error: An error if the command fails.
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
//   - file: The file name. For national subsets, specify "language/fileName".
//   - data: A datablock with the file contents.
//
// Note: If the environment contains a file with the specified name, the old file is replaced.
// The file is automatically removed from the Environment cache on all cluster members.
//
// Returns:
//   - error: An error if the command fails or data is nil.
func (cli *Cli) StoreClusterPBXFile(file string, data []byte) error {
	if file == "" || data == nil {
		return fmt.Errorf("file and data are required")
	}
	return cli.QueryNV("STORECLUSTERPBXFILE", file, "DATA", data)
}

// StoreDomainPBXFile stores a file into the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - file: The file name. For national subsets, specify "language/fileName".
//   - data: A datablock with the file contents.
//
// Note: If the environment contains a file with the specified name, the old file is replaced.
// The file is automatically removed from the Environment cache.
//
// Returns:
//   - error: An error if the command fails or data is nil.
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
//   - file: The file name. For national subsets, specify "language/fileName".
//   - data: A datablock with the file contents.
//
// Note: If the environment contains a file with the specified name, the old file is replaced.
// The file is automatically removed from the Environment cache.
//
// Returns:
//   - error: An error if the command fails or data is nil.
func (cli *Cli) StoreServerPBXFile(file string, data []byte) error {
	if file == "" || data == nil {
		return fmt.Errorf("file and data are required")
	}
	return cli.QueryNV("STORESERVERPBXFILE", file, "DATA", data)
}

// listPBXFiles — внутренний хелпер для обработки ответов LIST*PBXFILES.
// Преобразует map[string]any в типизированную мапу метаданных.
func (cli *Cli) listPBXFiles(cmd string, args ...any) (map[string]*PBXFileInfo, error) {
	raw, err := cli.getMapAny(cmd, args...)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*PBXFileInfo, len(raw))
	for name, attrs := range raw {
		attrMap, ok := attrs.(map[string]any)
		if !ok {
			continue
		}

		info := &PBXFileInfo{
			Name:     name,
			Size:     toInt64(attrMap["STFileSize"]),
			Created:  toTime(attrMap["STCreated"]),
			Modified: toTime(attrMap["STModified"]),
		}
		res[name] = info
	}
	return res, nil
}
