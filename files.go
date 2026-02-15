// # File Storage Administration
//
// The Files section provides capabilities for managing the CommuniGate Pro File Storage
// (Web-folder) content. It allows for direct manipulation of files and directories
// within an account's storage area, supporting both binary data transfer and
// metadata management.
//
// Key capabilities include:
//   - File Operations: Reading, writing, renaming, and deleting files using
//     streaming-friendly methods like [Cli.ReadStorageFile] and [Cli.WriteStorageFile].
//   - Partial Access: Support for offsets and specific data sizes, enabling
//     efficient processing of large files.
//   - Metadata & Attributes: Retrieving storage statistics via [StorageFileInfo]
//     and managing extended file attributes.
//   - Access Control: The "AUTH" parameter in most methods allows administrators
//     to perform operations on behalf of other accounts.
//   - Subscriptions: Managing file-folder subscriptions for collaborative
//     work and synchronization.
//
// Use these methods to build web-mail file attachments handlers, document
// management systems, or automated backup solutions within CommuniGate Pro.
package cgpcli

import (
	"fmt"
	"time"
)

// StorageFileContent holds the raw data and metadata of a storage file.
type StorageFileContent struct {
	Data         []byte    // Raw file content retrieved from the storage.
	ModifiedTime time.Time // Last modification timestamp of the file.
	CurrentSize  int64     // Total size of the file on the server in bytes.
}

// StorageFileInfo represents cumulative statistics for a storage path.
type StorageFileInfo struct {
	TotalSize  int64 // Total size in bytes of all files in the specified path.
	FilesCount int64 // Total number of files found in the specified path.
}

// DeleteStorageFile removes a file or a file directory from the Account File Storage.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - file: The name of the existing File Storage file or directory.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - error: An error if mandatory parameters are missing or the server fails to delete the file.
func (cli *Cli) DeleteStorageFile(account, file, authAccount string) (err error) {
	if account == "" || file == "" {
		return fmt.Errorf("account name and file name are required")
	}

	const cmd = "DELETESTORAGEFILE"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "FILE", file, "AUTH", authAccount)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "FILE", file)
	}
	return
}

// GetFileSubscription retrieves the list of Account "subscribed files".
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//
// Returns:
//   - []string: An array containing the list of subscribed file names.
//   - error: An error if the command fails.
func (cli *Cli) GetFileSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETFILESUBSCRIPTION", Atom(account))
}

// GetStorageFileInfo retrieves statistical information about all files in the Account File Storage.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - path: (Optional) The name of the File Storage subdirectory.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - *[StorageFileInfo]: A structure with total size and file count.
//   - error: An error if the command fails or parsing of numeric elements fails.
func (cli *Cli) GetStorageFileInfo(account, path, authAccount string) (*StorageFileInfo, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "GETSTORAGEFILEINFO"
	args := make([]any, 0, 5)
	args = append(args, Atom(account))
	if path != "" {
		args = append(args, "PATH", path)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, fmt.Errorf("%s: expected array of 2 elements, got %T", cmd, res)
	}

	return &StorageFileInfo{
		TotalSize:  toInt64(arr[0]),
		FilesCount: toInt64(arr[1]),
	}, nil
}

// ListStorageFiles lists all files in the File Storage directory or its subdirectories.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - path: (Optional) The name of the File Storage subdirectory.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - map[string]any: A dictionary where keys are file names and values contain file details.
//   - error: An error if the command fails.
func (cli *Cli) ListStorageFiles(account, path, authAccount string) (map[string]any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

	const cmd = "LISTSTORAGEFILES"
	args := make([]any, 0, 5)
	args = append(args, Atom(account))
	if path != "" {
		args = append(args, "PATH", path)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	return cli.getMapAny(cmd, args...)
}

// ReadStorageFile retrieves a file or its slice from the Account File Storage.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - file: The name of the File Storage file to be retrieved.
//   - offset: File position to start reading from.
//   - size: Maximum number of data bytes to return. Use 0 for no limit.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - *[StorageFileContent]: Structure containing data, modification time, and current size.
//   - error: An error if retrieval fails or data format is unexpected.
func (cli *Cli) ReadStorageFile(account, file string, offset, size int64, authAccount string) (*StorageFileContent, error) {
	if account == "" || file == "" {
		return nil, fmt.Errorf("account name and file name are required")
	}

	const cmd = "READSTORAGEFILE"
	args := make([]any, 0, 9)
	args = append(args, Atom(account), "FILE", file)

	// Используем != 0, чтобы учитывались константы OffsetBeg/End (-1, -2)
	if offset != 0 {
		args = append(args, "OFFSET", offset)
	}
	if size > 0 {
		args = append(args, "SIZE", size)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 3 {
		return nil, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	var data []byte
	if arr[0] != nil {
		data, ok = arr[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("%s: expected []byte for DataBlock, got %T", cmd, arr[0])
		}
	}

	return &StorageFileContent{
		Data:         data,
		ModifiedTime: toTime(arr[1]),
		CurrentSize:  toInt64(arr[2]),
	}, nil
}

// ReadStorageFileAttr reads attributes of an Account File Storage file or directory.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - file: The name of the file or directory.
//   - attrs: (Optional) An array of attribute names to retrieve.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - []any: An array of XML elements containing the requested attributes.
//   - error: An error if the command fails.
func (cli *Cli) ReadStorageFileAttr(account, file string, attrs []string, authAccount string) ([]any, error) {
	if account == "" || file == "" {
		return nil, fmt.Errorf("account name and file name are required")
	}

	args := make([]any, 0, 6)
	args = append(args, Atom(account), "FILE", file)
	if len(attrs) > 0 {
		args = append(args, attrs)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	return cli.getSliceAny("READSTORAGEFILEATTR", args...)
}

// RenameStorageFile renames a file or a file directory in the Account File Storage.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - oldName: The name of the existing File Storage file or directory.
//   - newName: The new name for the file or directory.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - error: An error if names are missing or the server fails to rename the object.
func (cli *Cli) RenameStorageFile(account, oldName, newName, authAccount string) (err error) {
	if account == "" || oldName == "" || newName == "" {
		return fmt.Errorf("account name, old and new file names are required")
	}

	const cmd = "RENAMESTORAGEFILE"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "FILE", oldName, "INTO", newName, "AUTH", authAccount)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "FILE", oldName, "INTO", newName)
	}
	return
}

// SetFileSubscription sets the Account "subscribed files" list.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - subs: An array of strings where each element is a file name.
//
// Returns:
//   - error: An error if the command fails or input is nil.
func (cli *Cli) SetFileSubscription(account string, subs []string) error {
	if account == "" || subs == nil {
		return fmt.Errorf("account name and subscription are required")
	}
	return cli.QueryNV("SETFILESUBSCRIPTION", Atom(account), subs)
}

// UpdateStorageFileAttr updates attributes of an Account File Storage file or directory.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - file: The name of an existing file or directory.
//   - attrs: An array of XML elements representing the new attribute values.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) UpdateStorageFileAttr(account, file string, attrs []any, authAccount string) (err error) {
	if account == "" || file == "" || attrs == nil {
		return fmt.Errorf("account name, file name and attributes are required")
	}

	const cmd = "UPDATESTORAGEFILEATTR"
	if authAccount != "" {
		err = cli.QueryNV(cmd, Atom(account), "FILE", file, attrs, "AUTH", authAccount)
	} else {
		err = cli.QueryNV(cmd, Atom(account), "FILE", file, attrs)
	}
	return
}

// WriteStorageFile stores a file or creates a directory in the Account File Storage.
//
// Parameters:
//   - account: The name of an existing Account. Use "*" for current authenticated Account.
//   - file: The name for the File Storage file (ends with "/" to create a directory).
//   - offset: Write position. Can be int (position), or string ("BEG", "END", "NEW").
//   - data: Datablock containing the file content. Must be empty for directory creation.
//   - authAccount: (Optional) The name of an Account on whose behalf the operation is executed.
//
// Returns:
//   - error: An error if mandatory parameters are missing or writing fails.
func (cli *Cli) WriteStorageFile(account, file string, offset any, data []byte, authAccount string) error {
	if account == "" || file == "" || data == nil {
		return fmt.Errorf("account name, file name and file data are required")
	}

	args := make([]any, 0, 9)
	args = append(args, Atom(account), "FILE", file)

	if offset != nil {
		switch v := offset.(type) {
		case string:
			args = append(args, "OFFSET", Atom(v))
		default:
			args = append(args, "OFFSET", v)
		}
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}
	args = append(args, "DATA", data)
	return cli.QueryNV("WRITESTORAGEFILE", args...)
}
