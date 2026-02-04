// # Section: Web Skins Administration
package cgpcli

import (
	"fmt"
	"time"
)

// CreateClusterSkin создает кастомный кластерный скин.
func (cli *Cli) CreateClusterSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("CREATECLUSTERSKIN", skin)
}

// CreateDomainSkin создает кастомный скин домена.
func (cli *Cli) CreateDomainSkin(domain, skin string) (err error) {
	const cmd = "CREATEDOMAINSKIN"
	if domain != "" {
		err = cli.QueryNV(cmd, domain, "SKIN", skin)
	} else {
		err = cli.QueryNV(cmd, skin)
	}
	return
}

// CreateServerSkin создает кастомный серверный скин.
func (cli *Cli) CreateServerSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("CREATESERVERSKIN", skin)
}

// DeleteClusterSkin удаляет кастомный кластерный скин.
func (cli *Cli) DeleteClusterSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("DELETECLUSTERSKIN", skin)
}

// DeleteDomainSkin удаляет кастомный скин домена.
func (cli *Cli) DeleteDomainSkin(domain, skin string) (err error) {
	const cmd = "DELETEDOMAINSKIN"
	if domain != "" {
		err = cli.QueryNV(cmd, domain, "SKIN", skin)
	} else {
		err = cli.QueryNV(cmd, skin)
	}
	return
}

// DeleteServerSkin удаляет кастомный серверный скин.
func (cli *Cli) DeleteServerSkin(skin string) error {
	if skin == "" {
		return fmt.Errorf("skin name is required")
	}
	return cli.QueryNV("DELETESERVERSKIN", skin)
}

// ListClusterSkinFiles возвращает словарь с атрибутами файлов в кластерном скине.
func (cli *Cli) ListClusterSkinFiles(skin string) (map[string]any, error) {
	if skin == "" {
		return nil, fmt.Errorf("skin name is required")
	}
	return cli.getMapAny("LISTCLUSTERSKINFILES", skin)
}

// ListClusterSkins возвращает список имен кастомных кластерных скинов.
func (cli *Cli) ListClusterSkins() ([]string, error) {
	return cli.getSliceString("LISTCLUSTERSKINS", "")
}

// ListDomainSkinFiles возвращает словарь с атрибутами файлов в скине домена.
func (cli *Cli) ListDomainSkinFiles(domain, skin string) (res map[string]any, err error) {
	const cmd = "LISTDOMAINSKINFILES"
	if domain != "" {
		res, err = cli.getMapAny(cmd, domain, "SKIN", skin)
	} else {
		res, err = cli.getMapAny(cmd, skin)
	}
	return
}

// ListDomainSkins возвращает список имен кастомных скинов домена.
func (cli *Cli) ListDomainSkins(domain string) ([]string, error) {
	return cli.getSliceString("LISTDOMAINSKINS", domain)
}

// ListServerSkinFiles возвращает словарь с атрибутами файлов в серверном скине.
func (cli *Cli) ListServerSkinFiles(skin string) (map[string]any, error) {
	if skin == "" {
		return nil, fmt.Errorf("skin name is required")
	}
	return cli.getMapAny("LISTSERVERSKINFILES", skin)
}

// ListServerSkins возвращает список имен кастомных серверных скинов.
func (cli *Cli) ListServerSkins() ([]string, error) {
	return cli.getSliceString("LISTSERVERSKINS", "")
}

// ListStockSkinFiles возвращает список файлов встроенного (stock) скина.
func (cli *Cli) ListStockSkinFiles(skin string) (map[string]any, error) {
	return cli.getMapAny("LISTSTOCKSKINFILES", skin)
}

// ReadClusterSkinFile читает файл из кастомного кластерного скина.
func (cli *Cli) ReadClusterSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" || file == "" {
		return nil, time.Time{}, fmt.Errorf("skin name and file name are required")
	}

	const cmd = "READCLUSTERSKINFILE"
	res, err := cli.Query(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadDomainSkinFile читает файл из кастомного скина домена.
func (cli *Cli) ReadDomainSkinFile(domain, skin, file string) ([]byte, time.Time, error) {
	if file == "" {
		return nil, time.Time{}, fmt.Errorf("file name is required")
	}

	const cmd = "READDOMAINSKINFILE"
	var res any
	var err error

	if domain != "" {
		res, err = cli.Query(cmd, domain, "SKIN", skin, "FILE", file)
	} else {
		res, err = cli.Query(cmd, skin, "FILE", file)
	}

	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array [data, time], got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadServerSkinFile читает файл из кастомного серверного скина.
func (cli *Cli) ReadServerSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" || file == "" {
		return nil, time.Time{}, fmt.Errorf("skin name and file name are required")
	}

	const cmd = "READSERVERSKINFILE"
	res, err := cli.Query(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadStockSkinFile читает файл из встроенного (stock) скина.
func (cli *Cli) ReadStockSkinFile(skin, file string) ([]byte, time.Time, error) {
	if skin == "" || file == "" {
		return nil, time.Time{}, fmt.Errorf("skin name and file name are required")
	}

	const cmd = "READSTOCKSKINFILE"
	res, err := cli.Query(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// RenameClusterSkin переименовывает кастомный кластерный скин.
func (cli *Cli) RenameClusterSkin(skin, newSkin string) error {
	if skin == "" || newSkin == "" {
		return fmt.Errorf("skin name and newSkin name are required")
	}
	return cli.QueryNV("RENAMECLUSTERSKIN", skin, "INTO", newSkin)
}

// RenameDomainSkin переименовывает кастомный скин домена.
func (cli *Cli) RenameDomainSkin(domain, skin, newSkin string) (err error) {
	if skin == "" || newSkin == "" {
		return fmt.Errorf("skin name and newSkin name are required")
	}

	const cmd = "RENAMEDOMAINSKIN"
	if domain != "" {
		err = cli.QueryNV(cmd, domain, "SKIN", skin, "INTO", newSkin)
	} else {
		err = cli.QueryNV(cmd, skin, "INTO", newSkin)
	}
	return
}

// RenameServerSkin переименовывает кастомный серверный скин.
func (cli *Cli) RenameServerSkin(skin, newSkin string) error {
	if skin == "" || newSkin == "" {
		return fmt.Errorf("skin name and newSkin name are required")
	}
	return cli.QueryNV("RENAMESERVERSKIN", skin, "INTO", newSkin)
}

// StoreClusterSkinFile сохраняет файл в кастомный кластерный скин или удаляет его.
func (cli *Cli) StoreClusterSkinFile(skin, file string, data []byte) (err error) {
	if skin == "" || file == "" {
		return fmt.Errorf("skin name and file name are required")
	}

	const cmd = "STORECLUSTERSKINFILE"
	if data != nil {
		err = cli.QueryNV(cmd, skin, "FILE", file, "DATA", data)
	} else {
		err = cli.QueryNV(cmd, skin, "FILE", file, "DELETE")
	}
	return
}

// StoreDomainSkinFile сохраняет файл в кастомный скин домена или удаляет его.
func (cli *Cli) StoreDomainSkinFile(domain, skin, file string, data []byte) error {
	if file == "" {
		return fmt.Errorf("file name is required")
	}

	args := make([]any, 0, 6)
	if domain != "" {
		args = append(args, domain, "SKIN")
	}
	args = append(args, skin, "FILE", file)

	if data != nil {
		args = append(args, "DATA", data)
	} else {
		args = append(args, "DELETE")
	}

	return cli.QueryNV("STOREDOMAINSKINFILE", args...)
}

// StoreServerSkinFile сохраняет файл в кастомный серверный скин или удаляет его.
func (cli *Cli) StoreServerSkinFile(skin, file string, data []byte) (err error) {
	if skin == "" || file == "" {
		return fmt.Errorf("skin name and file name are required")
	}

	const cmd = "STORESERVERSKINFILE"
	if data != nil {
		err = cli.QueryNV(cmd, skin, "FILE", file, "DATA", data)
	} else {
		err = cli.QueryNV(cmd, skin, "FILE", file, "DELETE")
	}
	return
}
