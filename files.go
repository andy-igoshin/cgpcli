// # Section: File Storage Administration
package cgpcli

import (
	"fmt"
	"time"
)

// StorageFileInfo содержит статистику по количеству файлов и их общему объему.
type StorageFileInfo struct {
	TotalSize  int64
	FilesCount int64
}

// StorageFileContent содержит сырые байты файла и метаданные.
type StorageFileContent struct {
	Data         []byte
	ModifiedTime time.Time
	CurrentSize  int64
}

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

func (cli *Cli) GetFileSubscription(account string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	return cli.getSliceString("GETFILESUBSCRIPTION", Atom(account))
}

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

func (cli *Cli) SetFileSubscription(account string, subs []string) error {
	if account == "" || subs == nil {
		return fmt.Errorf("account name and subscription are required")
	}
	return cli.QueryNV("SETFILESUBSCRIPTION", Atom(account), subs)
}

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
