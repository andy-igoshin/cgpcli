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
	const cmd = "DELETESTORAGEFILE"
	if authAccount != "" {
		_, err = cli.Query(cmd, Atom(account), "FILE", file, "AUTH", authAccount)
	} else {
		_, err = cli.Query(cmd, Atom(account), "FILE", file)
	}
	return err
}

func (cli *Cli) GetFileSubscription(account string) ([]string, error) {
	const cmd = "GETFILESUBSCRIPTION"
	res, err := cli.Query(cmd, Atom(account))
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

func (cli *Cli) GetStorageFileInfo(account, path, authAccount string) (*StorageFileInfo, error) {
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
	const cmd = "LISTSTORAGEFILES"
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
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

func (cli *Cli) ReadStorageFile(account, file string, offset, size int64, authAccount string) (*StorageFileContent, error) {
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

	// Ожидаем массив: [datablock, timestamp, currentSize]
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
	const cmd = "READSTORAGEFILEATTR"
	args := make([]any, 0, 6)
	args = append(args, Atom(account), "FILE", file)
	if len(attrs) > 0 {
		args = append(args, attrs)
	}
	if authAccount != "" {
		args = append(args, "AUTH", authAccount)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}

	arr, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}
	return arr, nil
}

func (cli *Cli) RenameStorageFile(account, oldFile, newFile, authAccount string) (err error) {
	const cmd = "RENAMESTORAGEFILE"
	if authAccount != "" {
		_, err = cli.Query(cmd, Atom(account), "FILE", oldFile, "INTO", newFile, "AUTH", authAccount)
	} else {
		_, err = cli.Query(cmd, Atom(account), "FILE", oldFile, "INTO", newFile)
	}
	return err
}

func (cli *Cli) SetFileSubscription(account string, subs []string) error {
	_, err := cli.Query("SETFILESUBSCRIPTION", Atom(account), subs)
	return err
}

func (cli *Cli) UpdateStorageFileAttr(account, file string, attrs []any, authAccount string) (err error) {
	const cmd = "UPDATESTORAGEFILEATTR"
	if authAccount != "" {
		_, err = cli.Query(cmd, Atom(account), "FILE", file, attrs, "AUTH", authAccount)
	} else {
		_, err = cli.Query(cmd, Atom(account), "FILE", file, attrs)
	}
	return err
}

func (cli *Cli) WriteStorageFile(account, file string, offset any, data []byte, authAccount string) (err error) {
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

	_, err = cli.Query("WRITESTORAGEFILE", args...)
	return err
}
