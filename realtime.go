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
)

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

// ListClusterPBXFiles lists files in the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]any: A dictionary with file names as keys. Values are dictionaries with file attributes.
//   - error: An error if the command fails.
func (cli *Cli) ListClusterPBXFiles(language string) (map[string]any, error) {
	return cli.getMapAny("LISTCLUSTERPBXFILES", Atom(language))
}

// ListDomainPBXFiles lists files in the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: (Optional) This optional parameter specifies the Domain name. If not specified, the administrator Domain is used.
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]any: A dictionary with file names as keys. Values are dictionaries with file attributes.
//   - error: An error if the command fails.
func (cli *Cli) ListDomainPBXFiles(domain, language string) (map[string]any, error) {
	args := make([]any, 0, 3)
	if domain != "" {
		args = append(args, domain)
	}
	if language != "" {
		args = append(args, "FILE", language)
	}
	return cli.getMapAny("LISTDOMAINPBXFILES", args...)
}

// ListServerPBXFiles lists files in the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]any: A dictionary with file names as keys. Values are dictionaries with file attributes.
//   - error: An error if the command fails.
func (cli *Cli) ListServerPBXFiles(language string) (map[string]any, error) {
	return cli.getMapAny("LISTSERVERPBXFILES", Atom(language))
}

// ListStockPBXFiles lists files in the stock (built-in) Real-Time Application Environment.
//
// Parameters:
//   - language: (Optional) This optional parameter specifies a national subset name.
//
// Returns:
//   - map[string]any: A dictionary with file names as keys. Values are dictionaries with file attributes.
//   - error: An error if the command fails.
func (cli *Cli) ListStockPBXFiles(language string) (map[string]any, error) {
	return cli.getMapAny("LISTSTOCKPBXFILES", Atom(language))
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

	const cmd = "READCLUSTERPBXFILE"
	res, err := cli.Query(cmd, file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
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

	const cmd = "READDOMAINPBXFILE"
	res, err := cli.Query(cmd, domain, "FILE", file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
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

	const cmd = "READSERVERPBXFILE"
	res, err := cli.Query(cmd, file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
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

	const cmd = "READSTOCKPBXFILE"
	res, err := cli.Query(cmd, file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
}

// StoreClusterPBXFile stores a file into or deletes it from the cluster-wide Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//   - data: If not nil, stores the file (DATA). If nil, removes the file (DELETE).
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) StoreClusterPBXFile(file string, data []byte) (err error) {
	if file == "" {
		return fmt.Errorf("file is required")
	}

	const cmd = "STORECLUSTERPBXFILE"
	if data != nil {
		err = cli.QueryNV(cmd, file, "DATA", data)
	} else {
		err = cli.QueryNV(cmd, file, "DELETE")
	}
	return
}

// StoreDomainPBXFile stores a file into or deletes it from the Domain Real-Time Application Environment.
//
// Parameters:
//   - domain: This parameter specifies the Domain name.
//   - file: The file name. For national subsets, specify "language/fileName".
//   - data: If not nil, stores the file (DATA). If nil, removes the file (DELETE).
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) StoreDomainPBXFile(domain, file string, data []byte) (err error) {
	if domain == "" || file == "" {
		return fmt.Errorf("domain and file are required")
	}

	const cmd = "STOREDOMAINPBXFILE"
	if data != nil {
		err = cli.QueryNV(cmd, domain, "FILE", file, "DATA", data)
	} else {
		err = cli.QueryNV(cmd, domain, "FILE", file, "DELETE")
	}
	return
}

// StoreServerPBXFile stores a file into or deletes it from the Server-wide Real-Time Application Environment.
//
// Parameters:
//   - file: The file name. For national subsets, specify "language/fileName".
//   - data: If not nil, stores the file (DATA). If nil, removes the file (DELETE).
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) StoreServerPBXFile(file string, data []byte) (err error) {
	if file == "" {
		return fmt.Errorf("file is required")
	}

	const cmd = "STORESERVERPBXFILE"
	if data != nil {
		err = cli.QueryNV(cmd, file, "DATA", data)
	} else {
		err = cli.QueryNV(cmd, file, "DELETE")
	}
	return
}
